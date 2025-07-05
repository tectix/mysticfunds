// Authentication Management
let currentUser = null;

// Check if user is authenticated
function isAuthenticated() {
    return api.token !== null && api.token !== '' && api.token !== 'null';
}

// Initialize app
function initApp() {
    if (isAuthenticated()) {
        document.getElementById('navbar').style.display = 'block';
        document.getElementById('auth-modal').style.display = 'none';
        document.getElementById('main-content').style.display = 'block';
        showDashboard();
        loadUserData();
    } else {
        showAuthModal();
    }
}

// Show authentication modal
function showAuthModal() {
    document.getElementById('auth-modal').style.display = 'flex';
    document.getElementById('navbar').style.display = 'none';
    document.getElementById('main-content').style.display = 'none';
}

// Hide authentication modal
function hideAuthModal() {
    document.getElementById('auth-modal').style.display = 'none';
    document.getElementById('navbar').style.display = 'block';
    document.getElementById('main-content').style.display = 'block';
    showDashboard();
}

// Switch between login and register forms
function showLogin() {
    document.getElementById('login-form').style.display = 'flex';
    document.getElementById('register-form').style.display = 'none';
    document.querySelector('.tab-button:first-child').classList.add('active');
    document.querySelector('.tab-button:last-child').classList.remove('active');
    document.getElementById('auth-title').textContent = 'Welcome Back!';
}

function showRegister() {
    document.getElementById('login-form').style.display = 'none';
    document.getElementById('register-form').style.display = 'flex';
    document.querySelector('.tab-button:first-child').classList.remove('active');
    document.querySelector('.tab-button:last-child').classList.add('active');
    document.getElementById('auth-title').textContent = 'Join MysticFunds!';
}

// Handle login form submission
async function handleLogin(event) {
    event.preventDefault();
    
    const username = document.getElementById('login-username').value;
    const password = document.getElementById('login-password').value;
    
    if (!username || !password) {
        showToast('Please fill in all fields', 'error');
        return;
    }
    
    showLoading(true);
    
    try {
        const response = await api.login(username, password);
        
        if (response.token) {
            api.setToken(response.token);
            currentUser = {
                id: response.user_id,
                username: username
            };
            
            hideAuthModal();
            showToast('Login successful!', 'success');
            loadUserData();
        } else {
            throw new Error('Invalid response from server');
        }
    } catch (error) {
        console.error('Login error:', error);
        showToast(error.message || 'Login failed', 'error');
    } finally {
        showLoading(false);
    }
}

// Handle register form submission
async function handleRegister(event) {
    event.preventDefault();
    
    const username = document.getElementById('register-username').value;
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;
    
    if (!username || !email || !password) {
        showToast('Please fill in all fields', 'error');
        return;
    }
    
    if (password.length < 6) {
        showToast('Password must be at least 6 characters', 'error');
        return;
    }
    
    showLoading(true);
    
    try {
        const response = await api.register(username, email, password);
        
        if (response.token) {
            api.setToken(response.token);
            currentUser = {
                id: response.user_id,
                username: username
            };
            
            hideAuthModal();
            showToast('Registration successful!', 'success');
            loadUserData();
        } else {
            throw new Error('Invalid response from server');
        }
    } catch (error) {
        console.error('Registration error:', error);
        showToast(error.message || 'Registration failed', 'error');
    } finally {
        showLoading(false);
    }
}

// Handle logout
async function logout() {
    if (confirm('Are you sure you want to logout?')) {
        showLoading(true);
        
        try {
            await api.logout();
            currentUser = null;
            showToast('Logged out successfully', 'success');
            showAuthModal();
        } catch (error) {
            console.error('Logout error:', error);
            // Still logout locally even if server request fails
            api.setToken(null);
            currentUser = null;
            showAuthModal();
        } finally {
            showLoading(false);
        }
    }
}

// Load user data after authentication
async function loadUserData() {
    try {
        // Load initial dashboard data
        await Promise.all([
            loadWizardSelectors(),
            loadDashboardStats(),
            loadInvestmentTypes()
        ]);
    } catch (error) {
        console.error('Error loading user data:', error);
        showToast('Error loading data', 'error');
    }
}

// Event listeners
document.addEventListener('DOMContentLoaded', function() {
    // Initialize app
    initApp();
    
    // Auth form listeners
    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('register-form').addEventListener('submit', handleRegister);
    
    // Navigation toggle for mobile
    const navToggle = document.getElementById('nav-toggle');
    const navMenu = document.getElementById('nav-menu');
    
    if (navToggle) {
        navToggle.addEventListener('click', function() {
            navMenu.classList.toggle('active');
        });
    }
    
    // Other form listeners
    document.getElementById('create-wizard-form').addEventListener('submit', handleCreateWizard);
    document.getElementById('transfer-form').addEventListener('submit', handleTransferMana);
    document.getElementById('investment-form').addEventListener('submit', handleCreateInvestment);
});