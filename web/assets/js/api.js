// API Configuration - Environment-aware with enhanced session management
function getApiBaseUrl() {
    // Auto-detect API URL based on environment
    if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
        return 'http://localhost:8080/api';
    }
    
    // For Vercel deployment, API calls are proxied through Vercel
    if (window.location.hostname.includes('vercel.app')) {
        return `${window.location.protocol}//${window.location.host}/api`;
    }
    
    // For production, use same origin
    return `${window.location.protocol}//${window.location.host}/api`;
}

const API_BASE_URL = getApiBaseUrl();

// Enhanced API Client Class with session management
class ApiClient {
    constructor() {
        this.token = localStorage.getItem('auth_token');
        this.tokenExpiry = localStorage.getItem('auth_token_expiry');
        this.refreshAttempted = false;
        
        // Clean up invalid tokens
        if (this.token === 'null' || this.token === 'undefined' || this.token === '') {
            this.clearToken();
        }
        
        // Check token expiry
        this.checkTokenExpiry();
    }

    setToken(token, expiryHours = 24) {
        this.token = token;
        if (token) {
            localStorage.setItem('auth_token', token);
            // Set token expiry (default 24 hours)
            const expiry = new Date();
            expiry.setHours(expiry.getHours() + expiryHours);
            this.tokenExpiry = expiry.toISOString();
            localStorage.setItem('auth_token_expiry', this.tokenExpiry);
        } else {
            this.clearToken();
        }
    }

    clearToken() {
        this.token = null;
        this.tokenExpiry = null;
        this.refreshAttempted = false;
        localStorage.removeItem('auth_token');
        localStorage.removeItem('auth_token_expiry');
    }

    checkTokenExpiry() {
        if (!this.token || !this.tokenExpiry) return;
        
        const now = new Date();
        const expiry = new Date(this.tokenExpiry);
        
        // If token expires in less than 1 hour, try to refresh
        const oneHour = 60 * 60 * 1000;
        if ((expiry.getTime() - now.getTime()) < oneHour && !this.refreshAttempted) {
            this.attemptTokenRefresh();
        }
        
        // If token is expired, clear it
        if (now >= expiry) {
            this.clearToken();
        }
    }

    async attemptTokenRefresh() {
        if (this.refreshAttempted || !this.token) return;
        
        this.refreshAttempted = true;
        try {
            const result = await this.refreshToken();
            if (result && result.token) {
                this.setToken(result.token);
                console.log('Token refreshed successfully');
            }
        } catch (error) {
            console.warn('Token refresh failed:', error.message);
            // Don't clear token here, let it expire naturally
        }
    }

    getHeaders() {
        const headers = {
            'Content-Type': 'application/json',
        };
        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }
        return headers;
    }

    async request(endpoint, options = {}, retryCount = 0) {
        // Check token expiry before making request
        this.checkTokenExpiry();
        
        const url = `${API_BASE_URL}${endpoint}`;
        const config = {
            ...options,
            headers: {
                ...this.getHeaders(),
                ...options.headers,
            },
        };

        try {
            const response = await fetch(url, config);
            
            // Always clone the response to avoid "body stream already read" errors
            const responseClone = response.clone();
            
            if (!response.ok) {
                if (response.status === 401 && retryCount === 0 && this.token) {
                    // Try to refresh token once
                    try {
                        await this.attemptTokenRefresh();
                        if (this.token) {
                            // Retry the request with new token
                            return this.request(endpoint, options, retryCount + 1);
                        }
                    } catch (refreshError) {
                        console.warn('Token refresh failed during request:', refreshError);
                    }
                    // If refresh fails or no token, proceed with auth error
                    this.clearToken();
                    if (typeof showAuthModal === 'function') {
                        showAuthModal();
                    }
                    throw new Error('Please log in to continue');
                }
                
                let errorMessage = 'Unknown error occurred';
                try {
                    const responseText = await responseClone.text();
                    try {
                        const errorData = JSON.parse(responseText);
                        errorMessage = errorData.error || errorData.message || responseText;
                    } catch {
                        errorMessage = responseText || 'Unknown error occurred';
                    }
                } catch {
                    errorMessage = 'Failed to read error response';
                }
                
                // Provide user-friendly error messages
                switch (response.status) {
                    case 400:
                        throw new Error(`Invalid request: ${errorMessage}`);
                    case 403:
                        throw new Error('Access denied. You may not have permission for this action.');
                    case 404:
                        throw new Error('Resource not found. It may have been deleted or moved.');
                    case 500:
                        throw new Error('Server error. Please try again later.');
                    case 503:
                        throw new Error('Service temporarily unavailable. Please try again later.');
                    default:
                        throw new Error(errorMessage || `Request failed (${response.status})`);
                }
            }

            const contentType = response.headers.get('content-type');
            if (contentType && contentType.includes('application/json')) {
                return await response.json();
            }
            return await response.text();
        } catch (error) {
            console.error('API request failed:', error);
            throw error;
        }
    }

    // Auth API calls
    async register(username, email, password) {
        return this.request('/auth/register', {
            method: 'POST',
            body: JSON.stringify({ username, email, password }),
        });
    }

    async login(username, password) {
        return this.request('/auth/login', {
            method: 'POST',
            body: JSON.stringify({ username, password }),
        });
    }

    async logout() {
        const result = await this.request('/auth/logout', {
            method: 'POST',
            body: JSON.stringify({ token: this.token }),
        });
        this.setToken(null);
        return result;
    }

    async refreshToken() {
        return this.request('/auth/refresh', {
            method: 'POST',
            body: JSON.stringify({ token: this.token }),
        });
    }

    // Wizard API calls
    async getWizards(pageSize = 10, pageNumber = 1) {
        return this.request(`/wizards?page_size=${pageSize}&page_number=${pageNumber}`);
    }

    async exploreWizards(pageSize = 10, pageNumber = 1, realm = '') {
        let url = `/wizards/explore?page_size=${pageSize}&page_number=${pageNumber}`;
        if (realm) {
            url += `&realm=${encodeURIComponent(realm)}`;
        }
        return this.request(url);
    }

    async getWizard(id) {
        return this.request(`/wizards/${id}`);
    }

    async createWizard(name, realm, element) {
        return this.request('/wizards', {
            method: 'POST',
            body: JSON.stringify({
                name,
                realm,
                element,
            }),
        });
    }

    async updateWizard(id, name, realm, element) {
        return this.request(`/wizards/${id}`, {
            method: 'PUT',
            body: JSON.stringify({
                name,
                realm,
                element,
            }),
        });
    }

    async deleteWizard(id) {
        return this.request(`/wizards/${id}`, {
            method: 'DELETE',
        });
    }

    // Mana API calls
    async getManaBalance(wizardId) {
        return this.request(`/mana/balance/${wizardId}`);
    }

    async transferMana(fromWizardId, toWizardId, amount) {
        return this.request('/mana/transfer', {
            method: 'POST',
            body: JSON.stringify({
                from_wizard_id: fromWizardId,
                to_wizard_id: toWizardId,
                amount: amount,
            }),
        });
    }

    async getTransactions(wizardId, pageSize = 10, pageNumber = 1) {
        return this.request(`/mana/transactions/${wizardId}?page_size=${pageSize}&page_number=${pageNumber}`);
    }

    // Investment API calls
    async getInvestmentTypes(minAmount, maxAmount, riskLevel) {
        let query = '';
        const params = [];
        if (minAmount) params.push(`min_amount=${minAmount}`);
        if (maxAmount) params.push(`max_amount=${maxAmount}`);
        if (riskLevel) params.push(`risk_level=${riskLevel}`);
        if (params.length > 0) query = '?' + params.join('&');
        
        return this.request(`/mana/investment-types${query}`);
    }

    async createInvestment(wizardId, investmentTypeId, amount) {
        return this.request('/mana/investments', {
            method: 'POST',
            body: JSON.stringify({
                wizard_id: wizardId,
                investment_type_id: investmentTypeId,
                amount: amount,
            }),
        });
    }

    async getInvestments(wizardId, status = '') {
        let query = `wizard_id=${wizardId}`;
        if (status) query += `&status=${status}`;
        
        return this.request(`/mana/investments?${query}`);
    }

    // Jobs API calls
    async getJobs(realm = '', element = '', difficulty = '', pageSize = 20, pageNumber = 1) {
        let query = `page_size=${pageSize}&page_number=${pageNumber}&only_active=true`;
        const params = [];
        if (realm) params.push(`realm=${encodeURIComponent(realm)}`);
        if (element) params.push(`element=${encodeURIComponent(element)}`);
        if (difficulty) params.push(`difficulty=${encodeURIComponent(difficulty)}`);
        
        if (params.length > 0) {
            query += '&' + params.join('&');
        }
        
        return this.request(`/jobs?${query}`);
    }

    async createJob(jobData) {
        return this.request('/jobs', {
            method: 'POST',
            body: JSON.stringify(jobData),
        });
    }

    async assignWizardToJob(jobId, wizardId) {
        return this.request('/jobs/assign', {
            method: 'POST',
            body: JSON.stringify({
                job_id: jobId,
                wizard_id: wizardId,
            }),
        });
    }

    async getJobAssignments(wizardId = '', status = '', pageSize = 20, pageNumber = 1) {
        let query = `page_size=${pageSize}&page_number=${pageNumber}`;
        const params = [];
        if (wizardId) params.push(`wizard_id=${wizardId}`);
        if (status) params.push(`status=${encodeURIComponent(status)}`);
        
        if (params.length > 0) {
            query += '&' + params.join('&');
        }
        
        return this.request(`/jobs/assignments?${query}`);
    }

    async completeJobAssignment(assignmentId) {
        return this.request(`/jobs/assignments/${assignmentId}`, {
            method: 'PUT',
            body: JSON.stringify({
                action: 'complete'
            }),
        });
    }

    async cancelJobAssignment(assignmentId, reason = '') {
        return this.request(`/jobs/assignments/${assignmentId}`, {
            method: 'PUT',
            body: JSON.stringify({
                action: 'cancel',
                reason: reason
            }),
        });
    }

    // Job Progress API calls
    async getJobProgress(assignmentId) {
        return this.request(`/jobs/progress/${assignmentId}`);
    }

    async updateJobProgress(assignmentId, progressData) {
        return this.request(`/jobs/progress/${assignmentId}`, {
            method: 'PUT',
            body: JSON.stringify(progressData),
        });
    }

    // Activities API calls
    async getActivities(wizardId = '', activityType = '', pageSize = 20, pageNumber = 1) {
        let query = `page_size=${pageSize}&page_number=${pageNumber}`;
        const params = [];
        if (wizardId) params.push(`wizard_id=${wizardId}`);
        if (activityType) params.push(`activity_type=${encodeURIComponent(activityType)}`);
        
        if (params.length > 0) {
            query += '&' + params.join('&');
        }
        
        return this.request(`/activities?${query}`);
    }

    // Realms API calls
    async getRealms() {
        return this.request('/realms');
    }

    // Marketplace API calls - Basic methods for wizard collections
    async getWizardArtifacts(wizardId) {
        return this.request(`/marketplace/artifacts/wizard/${wizardId}`);
    }

    async getWizardSpells(wizardId) {
        return this.request(`/marketplace/spells/wizard/${wizardId}`);
    }
}

// Global API client instance
const api = new ApiClient();