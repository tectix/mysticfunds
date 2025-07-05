// Global app state
let wizards = [];
let investmentTypes = [];

// Utility Functions
function showLoading(show) {
    document.getElementById('loading').style.display = show ? 'flex' : 'none';
}

function showToast(message, type = 'info') {
    const toastContainer = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.innerHTML = `
        <div style="display: flex; align-items: center; gap: 10px;">
            <i class="fas fa-${getToastIcon(type)}"></i>
            <span>${message}</span>
        </div>
    `;
    
    toastContainer.appendChild(toast);
    
    // Auto remove after 5 seconds
    setTimeout(() => {
        if (toast.parentNode) {
            toast.parentNode.removeChild(toast);
        }
    }, 5000);
}

function getToastIcon(type) {
    switch (type) {
        case 'success': return 'check-circle';
        case 'error': return 'exclamation-circle';
        case 'warning': return 'exclamation-triangle';
        default: return 'info-circle';
    }
}

function formatNumber(num) {
    return new Intl.NumberFormat().format(num);
}

function formatDate(timestamp) {
    if (!timestamp) return 'Unknown';
    // Handle both seconds and milliseconds timestamps
    const date = typeof timestamp === 'string' ? new Date(timestamp) : new Date(timestamp * 1000);
    return date.toLocaleDateString();
}

function formatDuration(minutes) {
    if (!minutes || minutes === 0) return 'Unknown';
    
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    
    if (hours > 0 && mins > 0) {
        return `${hours}h ${mins}m`;
    } else if (hours > 0) {
        return `${hours}h`;
    } else if (mins > 0) {
        return `${mins}m`;
    } else {
        return 'Unknown';
    }
}

// Navigation Functions
function showPage(pageId) {
    // Hide all pages
    const pages = document.querySelectorAll('.page');
    pages.forEach(page => page.style.display = 'none');
    
    // Show selected page
    document.getElementById(pageId).style.display = 'block';
    
    // Update active nav link
    const navLinks = document.querySelectorAll('.nav-link');
    navLinks.forEach(link => link.classList.remove('active'));
}

function showDashboard() {
    showPage('dashboard');
    loadDashboardStats();
    loadRecentActivity();
}

function showWizards() {
    showPage('wizards');
    loadWizards();
    loadActiveJobs();
}

function showMana() {
    showPage('mana');
    loadWizardSelectors();
}

function showInvestments() {
    showPage('investments');
    loadWizardSelectors();
    loadInvestmentTypes();
}

function showJobs() {
    showPage('jobs');
    loadJobs();
    loadJobFilters();
    checkJobCreationAccess();
}

function checkJobCreationAccess() {
    const createJobBtn = document.getElementById('create-job-btn');
    if (createJobBtn) {
        // Check if user has any level 10+ wizards
        const hasHighLevelWizard = wizards.some(wizard => (wizard.level || 1) >= 10);
        createJobBtn.style.display = hasHighLevelWizard ? 'block' : 'none';
    }
}

// Dashboard Functions
async function loadDashboardStats() {
    try {
        if (wizards.length === 0) {
            await loadWizards();
        }

        // Calculate stats
        let totalMana = 0;
        let activeInvestments = 0;
        let totalReturns = 0;

        for (const wizard of wizards) {
            // Get mana balance
            try {
                const balanceResponse = await api.getManaBalance(wizard.id);
                totalMana += balanceResponse.balance || 0;
            } catch (error) {
                console.warn(`Failed to get balance for wizard ${wizard.id}`);
            }

            // Get investments
            try {
                const investmentResponse = await api.getInvestments(wizard.id);
                if (investmentResponse.investments) {
                    const active = investmentResponse.investments.filter(inv => inv.status === 'active').length;
                    activeInvestments += active;
                    
                    const completed = investmentResponse.investments.filter(inv => inv.status === 'completed');
                    totalReturns += completed.reduce((sum, inv) => sum + (inv.returned_amount - inv.amount), 0);
                }
            } catch (error) {
                console.warn(`Failed to get investments for wizard ${wizard.id}`);
            }
        }

        // Update UI
        document.getElementById('wizard-count').textContent = wizards.length;
        document.getElementById('total-mana').textContent = formatNumber(totalMana);
        document.getElementById('active-investments').textContent = activeInvestments;
        document.getElementById('total-returns').textContent = formatNumber(totalReturns);

    } catch (error) {
        console.error('Error loading dashboard stats:', error);
        showToast('Error loading dashboard stats', 'error');
    }
}

async function loadRecentActivity() {
    const recentTransactionsDiv = document.getElementById('recent-transactions');
    
    try {
        if (wizards.length === 0) return;

        let allActivities = [];
        
        // Get recent activities (including job events)
        try {
            const activitiesResponse = await api.getActivities('', '', 10, 1);
            if (activitiesResponse.activities) {
                allActivities.push(...activitiesResponse.activities.map(activity => ({
                    type: 'activity',
                    activity_type: activity.activity_type,
                    description: activity.activity_description,
                    created_at: activity.created_at,
                    wizard_id: activity.wizard_id
                })));
            }
        } catch (error) {
            console.warn('Failed to get activities:', error);
        }

        // Get recent transactions for each wizard
        for (const wizard of wizards.slice(0, 3)) { // Limit to first 3 wizards for performance
            try {
                const response = await api.getTransactions(wizard.id, 5, 1);
                if (response.transactions) {
                    allActivities.push(...response.transactions.map(t => ({
                        type: 'transaction',
                        ...t,
                        wizard_name: wizard.name
                    })));
                }
            } catch (error) {
                console.warn(`Failed to get transactions for wizard ${wizard.id}`);
            }
        }

        // Sort by date and take most recent
        allActivities.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
        allActivities = allActivities.slice(0, 8);

        if (allActivities.length === 0) {
            recentTransactionsDiv.innerHTML = '<p>No recent activity</p>';
            return;
        }

        recentTransactionsDiv.innerHTML = allActivities.map(item => {
            if (item.type === 'activity') {
                const wizardName = wizards.find(w => w.id === item.wizard_id)?.name || 'Unknown Wizard';
                return `
                    <div class="transaction-item activity-item">
                        <div class="transaction-info">
                            <div><i class="fas fa-star"></i> ${item.description}</div>
                            <div class="transaction-date">${formatDate(item.created_at)}</div>
                        </div>
                        <div class="activity-type">${item.activity_type}</div>
                    </div>
                `;
            } else {
                return `
                    <div class="transaction-item">
                        <div class="transaction-info">
                            <div><i class="fas fa-coins"></i> Transfer: ${item.wizard_name}</div>
                            <div class="transaction-date">${formatDate(item.created_at)}</div>
                        </div>
                        <div class="transaction-amount">${formatNumber(item.amount)} Mana</div>
                    </div>
                `;
            }
        }).join('');

    } catch (error) {
        console.error('Error loading recent activity:', error);
        recentTransactionsDiv.innerHTML = '<p>Error loading activity</p>';
    }
}

// Wizard Functions
async function loadWizards() {
    try {
        showLoading(true);
        const response = await api.getWizards(20, 1);
        wizards = response.wizards || [];
        displayWizards();
    } catch (error) {
        console.error('Error loading wizards:', error);
        showToast('Error loading wizards', 'error');
        wizards = [];
    } finally {
        showLoading(false);
    }
}

function displayWizards() {
    const wizardsGrid = document.getElementById('wizards-grid');
    
    if (wizards.length === 0) {
        wizardsGrid.innerHTML = '<p style="text-align: center; color: white;">No wizards found. Create your first wizard!</p>';
        return;
    }
    
    wizardsGrid.innerHTML = wizards.map(wizard => `
        <div class="wizard-card">
            <div class="wizard-header">
                <div class="wizard-id">ID: ${wizard.id}</div>
                <div class="wizard-avatar">
                    <i class="fas fa-hat-wizard wizard-element-${wizard.element.toLowerCase()}"></i>
                </div>
            </div>
            <div class="wizard-name">${wizard.name}</div>
            <div class="wizard-info">
                <span><i class="fas fa-globe"></i> ${wizard.realm}</span>
                <span><i class="fas fa-magic"></i> ${wizard.element}</span>
            </div>
            <div class="wizard-stats">
                <div class="wizard-level">
                    <i class="fas fa-star"></i> Level ${wizard.level || 1}
                </div>
                <div class="wizard-mana">
                    <i class="fas fa-coins"></i> ${formatNumber(wizard.mana_balance || 0)} Mana
                </div>
                <div class="wizard-exp">
                    <i class="fas fa-trophy"></i> ${formatNumber(wizard.experience_points || 0)} EXP
                </div>
            </div>
            <div class="wizard-actions">
                <button class="btn btn-primary" onclick="viewWizardDetails(${wizard.id})">
                    <i class="fas fa-eye"></i> View
                </button>
                <button class="btn btn-secondary" onclick="transferManaTo(${wizard.id})">
                    <i class="fas fa-paper-plane"></i> Send
                </button>
            </div>
        </div>
    `).join('');
}

async function loadWizardSelectors() {
    if (wizards.length === 0) {
        await loadWizards();
    }

    const selectors = [
        'from-wizard',
        'transaction-wizard',
        'investment-wizard',
        'investments-wizard'
    ];

    selectors.forEach(selectorId => {
        const select = document.getElementById(selectorId);
        if (select) {
            select.innerHTML = '<option value="">Select Wizard</option>' +
                wizards.map(wizard => 
                    `<option value="${wizard.id}">${wizard.name} (${formatNumber(wizard.mana_balance || 0)} Mana)</option>`
                ).join('');
        }
    });
}

function showCreateWizard() {
    document.getElementById('create-wizard-modal').classList.add('show');
}

function closeCreateWizard() {
    document.getElementById('create-wizard-modal').classList.remove('show');
    document.getElementById('create-wizard-form').reset();
}

async function handleCreateWizard(event) {
    event.preventDefault();
    
    const name = document.getElementById('wizard-name').value;
    const realm = document.getElementById('wizard-realm').value;
    const element = document.getElementById('wizard-element').value;
    
    if (!name || !realm || !element) {
        showToast('Please fill in all fields', 'error');
        return;
    }
    
    showLoading(true);
    
    try {
        const newWizard = await api.createWizard(name, realm, element);
        wizards.push(newWizard);
        displayWizards();
        loadWizardSelectors();
        closeCreateWizard();
        showToast('Wizard created successfully!', 'success');
    } catch (error) {
        console.error('Error creating wizard:', error);
        let errorMessage = error.message || 'Failed to create wizard';
        
        // Provide specific error messages
        if (errorMessage.includes('already exists')) {
            errorMessage = 'A wizard with this name already exists. Please choose a different name.';
        } else if (errorMessage.includes('limit')) {
            errorMessage = 'You have reached the maximum number of wizards (2). Delete a wizard to create a new one.';
        } else if (errorMessage.includes('invalid')) {
            errorMessage = 'Please check that all fields are filled correctly.';
        }
        
        showToast(errorMessage, 'error');
    } finally {
        showLoading(false);
    }
}

async function viewWizardDetails(wizardId) {
    try {
        const wizard = await api.getWizard(wizardId);
        const balance = await api.getManaBalance(wizardId);
        
        // Add mana balance to wizard object for modal display
        wizard.mana_balance = balance.balance || 0;
        
        displayWizardModal(wizard);
    } catch (error) {
        console.error('Error loading wizard details:', error);
        showToast('Error loading wizard details', 'error');
    }
}

// Mana Functions
async function handleTransferMana(event) {
    event.preventDefault();
    
    const fromWizardId = parseInt(document.getElementById('from-wizard').value);
    const toWizardId = parseInt(document.getElementById('to-wizard').value);
    const amount = parseInt(document.getElementById('transfer-amount').value);
    
    if (!fromWizardId || !toWizardId || !amount) {
        showToast('Please fill in all fields', 'error');
        return;
    }
    
    if (fromWizardId === toWizardId) {
        showToast('Cannot transfer to the same wizard', 'error');
        return;
    }
    
    if (amount <= 0) {
        showToast('Amount must be positive', 'error');
        return;
    }
    
    showLoading(true);
    
    try {
        await api.transferMana(fromWizardId, toWizardId, amount);
        showToast('Mana transferred successfully!', 'success');
        document.getElementById('transfer-form').reset();
        loadWizardSelectors(); // Refresh balances
    } catch (error) {
        console.error('Error transferring mana:', error);
        let errorMessage = error.message || 'Failed to transfer mana';
        
        // Provide specific error messages for common issues
        if (errorMessage.includes('not found')) {
            errorMessage = 'One of the selected wizards was not found. Please refresh and try again.';
        } else if (errorMessage.includes('insufficient')) {
            errorMessage = 'Insufficient mana balance for this transfer.';
        } else if (errorMessage.includes('same wizard')) {
            errorMessage = 'Cannot transfer mana to the same wizard.';
        }
        
        showToast(errorMessage, 'error');
    } finally {
        showLoading(false);
    }
}

async function loadTransactions() {
    const wizardId = document.getElementById('transaction-wizard').value;
    const transactionsList = document.getElementById('transactions-list');
    
    if (!wizardId) {
        transactionsList.innerHTML = '<p>Please select a wizard</p>';
        return;
    }
    
    try {
        showLoading(true);
        const response = await api.getTransactions(wizardId, 20, 1);
        const transactions = response.transactions || [];
        
        if (transactions.length === 0) {
            transactionsList.innerHTML = '<p>No transactions found</p>';
            return;
        }
        
        transactionsList.innerHTML = transactions.map(transaction => `
            <div class="transaction-item">
                <div class="transaction-info">
                    <div>From: ${transaction.from_wizard_id} → To: ${transaction.to_wizard_id}</div>
                    <div class="transaction-date">${formatDate(transaction.created_at)}</div>
                </div>
                <div class="transaction-amount">${formatNumber(transaction.amount)} Mana</div>
            </div>
        `).join('');
        
    } catch (error) {
        console.error('Error loading transactions:', error);
        transactionsList.innerHTML = '<p>Error loading transactions</p>';
        showToast('Error loading transactions', 'error');
    } finally {
        showLoading(false);
    }
}

// Investment Functions
async function loadInvestmentTypes() {
    try {
        const response = await api.getInvestmentTypes();
        investmentTypes = response.investment_types || [];
        
        const select = document.getElementById('investment-type');
        if (select) {
            select.innerHTML = '<option value="">Select Investment Type</option>' +
                investmentTypes.map(type => 
                    `<option value="${type.id}">${type.name} - ${type.base_return_rate}% (Risk: ${type.risk_level})</option>`
                ).join('');
        }
    } catch (error) {
        console.error('Error loading investment types:', error);
        showToast('Error loading investment types', 'error');
    }
}

async function handleCreateInvestment(event) {
    event.preventDefault();
    
    const wizardId = parseInt(document.getElementById('investment-wizard').value);
    const investmentTypeId = parseInt(document.getElementById('investment-type').value);
    const amount = parseInt(document.getElementById('investment-amount').value);
    
    if (!wizardId || !investmentTypeId || !amount) {
        showToast('Please fill in all fields', 'error');
        return;
    }
    
    if (amount < 100) {
        showToast('Minimum investment amount is 100 Mana', 'error');
        return;
    }
    
    showLoading(true);
    
    try {
        await api.createInvestment(wizardId, investmentTypeId, amount);
        showToast('Investment created successfully!', 'success');
        document.getElementById('investment-form').reset();
        loadWizardSelectors(); // Refresh balances
    } catch (error) {
        console.error('Error creating investment:', error);
        showToast(error.message || 'Error creating investment', 'error');
    } finally {
        showLoading(false);
    }
}

async function loadInvestments() {
    const wizardId = document.getElementById('investments-wizard').value;
    const investmentsList = document.getElementById('investments-list');
    
    if (!wizardId) {
        investmentsList.innerHTML = '<p>Please select a wizard</p>';
        return;
    }
    
    try {
        showLoading(true);
        const response = await api.getInvestments(wizardId);
        const investments = response.investments || [];
        
        if (investments.length === 0) {
            investmentsList.innerHTML = '<p>No investments found</p>';
            return;
        }
        
        investmentsList.innerHTML = investments.map(investment => `
            <div class="investment-item">
                <div class="investment-info">
                    <div><strong>${investment.investment_type}</strong></div>
                    <div>Amount: ${formatNumber(investment.amount)} Mana</div>
                    <div>Return Rate: ${investment.actual_return_rate}%</div>
                    <div class="investment-date">Started: ${formatDate(investment.start_time)}</div>
                </div>
                <div class="investment-status status-${investment.status}">
                    ${investment.status}
                </div>
            </div>
        `).join('');
        
    } catch (error) {
        console.error('Error loading investments:', error);
        investmentsList.innerHTML = '<p>Error loading investments</p>';
        showToast('Error loading investments', 'error');
    } finally {
        showLoading(false);
    }
}

// Explore Wizards Functions
function showExploreWizards() {
    showPage('explore-wizards');
    loadExploreWizards();
}

async function loadExploreWizards(realm = '') {
    try {
        showLoading(true);
        const response = await api.exploreWizards(20, 1, realm);
        const exploreWizards = response.wizards || [];
        displayExploreWizards(exploreWizards);
    } catch (error) {
        console.error('Error loading explore wizards:', error);
        showToast('Error loading wizards for exploration', 'error');
    } finally {
        showLoading(false);
    }
}

function displayExploreWizards(exploreWizards) {
    const exploreWizardsGrid = document.getElementById('explore-wizards-grid');
    
    if (exploreWizards.length === 0) {
        exploreWizardsGrid.innerHTML = '<p style="text-align: center; color: white;">No wizards found in this realm.</p>';
        return;
    }
    
    exploreWizardsGrid.innerHTML = exploreWizards.map(wizard => `
        <div class="wizard-card">
            <div class="wizard-header">
                <h3>${wizard.name}</h3>
                <span class="wizard-element ${wizard.element.toLowerCase()}">${wizard.element}</span>
            </div>
            <div class="wizard-details">
                <div class="wizard-info">
                    <i class="fas fa-map-marker-alt"></i>
                    <span>Realm: ${wizard.realm}</span>
                </div>
                <div class="wizard-info">
                    <i class="fas fa-coins"></i>
                    <span>Mana: ${formatNumber(wizard.mana_balance)}</span>
                </div>
                ${wizard.guild ? `
                    <div class="wizard-info">
                        <i class="fas fa-shield-alt"></i>
                        <span>Guild: ${wizard.guild.name}</span>
                    </div>
                ` : ''}
            </div>
            <div class="wizard-actions">
                <button class="btn btn-sm btn-secondary" onclick="viewWizardProfile(${wizard.id})">
                    <i class="fas fa-eye"></i> View Profile
                </button>
            </div>
        </div>
    `).join('');
}

function filterWizardsByRealm() {
    const realmFilter = document.getElementById('realm-filter').value;
    loadExploreWizards(realmFilter);
}

function viewWizardProfile(wizardId) {
    showWizardDetailsModal(wizardId);
}

// Wizard Details Modal
function showWizardDetailsModal(wizardId) {
    // Find wizard data
    let wizard = null;
    
    // Check if wizard is in current wizards array or fetch it
    if (wizards && wizards.length > 0) {
        wizard = wizards.find(w => w.id === wizardId);
    }
    
    if (wizard) {
        displayWizardModal(wizard);
    } else {
        // Fetch wizard details
        fetchWizardDetails(wizardId);
    }
}

async function fetchWizardDetails(wizardId) {
    try {
        showLoading(true);
        const wizard = await api.getWizard(wizardId);
        displayWizardModal(wizard);
    } catch (error) {
        console.error('Error fetching wizard details:', error);
        showToast('Error loading wizard details', 'error');
    } finally {
        showLoading(false);
    }
}

function displayWizardModal(wizard) {
    const modalHTML = `
        <div id="wizard-details-modal" class="modal show">
            <div class="modal-content wizard-profile-card">
                <div class="wizard-profile-header">
                    <div class="wizard-avatar">
                        <i class="fas fa-hat-wizard wizard-element-${wizard.element.toLowerCase()}"></i>
                    </div>
                    <div class="wizard-basic-info">
                        <div class="wizard-id-badge">ID: ${wizard.id}</div>
                        <h2 class="wizard-name">${wizard.name}</h2>
                        <div class="wizard-element-badge ${wizard.element.toLowerCase()}">
                            <i class="fas fa-magic"></i> ${wizard.element} Wizard
                        </div>
                        <div class="wizard-location">
                            <i class="fas fa-map-marker-alt"></i> ${wizard.realm}
                        </div>
                    </div>
                    <button class="close-btn" onclick="closeWizardDetailsModal()">
                        <i class="fas fa-times"></i>
                    </button>
                </div>
                
                <div class="wizard-profile-content">
                    <div class="wizard-stats-grid">
                        <div class="stat-card mana-stat">
                            <div class="stat-icon">
                                <i class="fas fa-coins"></i>
                            </div>
                            <div class="stat-info">
                                <div class="stat-value">${formatNumber(wizard.mana_balance)}</div>
                                <div class="stat-label">Mana Balance</div>
                            </div>
                        </div>
                        
                        ${wizard.guild ? `
                            <div class="stat-card guild-stat">
                                <div class="stat-icon">
                                    <i class="fas fa-shield-alt"></i>
                                </div>
                                <div class="stat-info">
                                    <div class="stat-value">${wizard.guild.name}</div>
                                    <div class="stat-label">Guild Member</div>
                                </div>
                            </div>
                        ` : `
                            <div class="stat-card guild-stat">
                                <div class="stat-icon">
                                    <i class="fas fa-user"></i>
                                </div>
                                <div class="stat-info">
                                    <div class="stat-value">Solo</div>
                                    <div class="stat-label">Independent</div>
                                </div>
                            </div>
                        `}
                        
                        <div class="stat-card user-stat">
                            <div class="stat-icon">
                                <i class="fas fa-user-tag"></i>
                            </div>
                            <div class="stat-info">
                                <div class="stat-value">
                                    <span id="user-id-display">${wizard.user_id}</span>
                                    <button class="copy-btn" onclick="copyUserId(${wizard.user_id})" title="Copy User ID">
                                        <i class="fas fa-copy"></i>
                                    </button>
                                </div>
                                <div class="stat-label">User ID</div>
                            </div>
                        </div>
                        
                        <div class="stat-card level-stat">
                            <div class="stat-icon">
                                <i class="fas fa-star"></i>
                            </div>
                            <div class="stat-info">
                                <div class="stat-value">Level ${wizard.level || 1}</div>
                                <div class="stat-label">Experience Level</div>
                            </div>
                        </div>
                        
                        <div class="stat-card exp-stat">
                            <div class="stat-icon">
                                <i class="fas fa-trophy"></i>
                            </div>
                            <div class="stat-info">
                                <div class="stat-value">${formatNumber(wizard.experience_points || 0)}</div>
                                <div class="stat-label">Experience Points</div>
                            </div>
                        </div>
                        
                        <div class="stat-card joined-stat">
                            <div class="stat-icon">
                                <i class="fas fa-calendar-alt"></i>
                            </div>
                            <div class="stat-info">
                                <div class="stat-value">${formatDate(wizard.created_at)}</div>
                                <div class="stat-label">Joined</div>
                            </div>
                        </div>
                    </div>
                    
                    <div class="wizard-actions">
                        <button class="action-btn primary" onclick="sendManaToWizard(${wizard.id})">
                            <i class="fas fa-paper-plane"></i>
                            Send Mana
                        </button>
                        <button class="action-btn secondary" onclick="closeWizardDetailsModal()">
                            <i class="fas fa-arrow-left"></i>
                            Back
                        </button>
                    </div>
                </div>
            </div>
        </div>
    `;
    
    // Remove existing modal if any
    const existingModal = document.getElementById('wizard-details-modal');
    if (existingModal) {
        existingModal.remove();
    }
    
    // Add modal to body
    document.body.insertAdjacentHTML('beforeend', modalHTML);
}

function closeWizardDetailsModal() {
    const modal = document.getElementById('wizard-details-modal');
    if (modal) {
        modal.remove();
    }
}

function copyUserId(userId) {
    navigator.clipboard.writeText(userId.toString()).then(() => {
        showToast('User ID copied to clipboard!', 'success');
    }).catch(() => {
        // Fallback for older browsers
        const textArea = document.createElement('textarea');
        textArea.value = userId.toString();
        document.body.appendChild(textArea);
        textArea.select();
        document.execCommand('copy');
        document.body.removeChild(textArea);
        showToast('User ID copied to clipboard!', 'success');
    });
}

function sendManaToWizard(wizardId) {
    closeWizardDetailsModal();
    showMana();
    // Pre-fill the wizard ID in the transfer form
    setTimeout(() => {
        const toWizardInput = document.getElementById('to-wizard');
        if (toWizardInput) {
            toWizardInput.value = wizardId;
            toWizardInput.setAttribute('value', wizardId);
            // Trigger label movement
            toWizardInput.focus();
            toWizardInput.blur();
        }
    }, 100);
}

function transferManaTo(wizardId) {
    showMana();
    // Pre-fill the wizard ID in the transfer form
    setTimeout(() => {
        const toWizardInput = document.getElementById('to-wizard');
        if (toWizardInput) {
            toWizardInput.value = wizardId;
            toWizardInput.setAttribute('value', wizardId);
            // Trigger label movement
            toWizardInput.focus();
            toWizardInput.blur();
        }
    }, 100);
}

// Realms functionality
function showRealms() {
    showPage('realms');
    loadRealms();
}

async function loadRealms() {
    try {
        showLoading(true);
        // Hardcoded realms data since we don't have a backend endpoint yet
        const realms = [
            {
                id: 1,
                name: "Pyrrhian Flame",
                element: "Fire / Heat / Chaos",
                description: "Realm of eternal fire and volcanic titans",
                lore: "Home to the Salamandrine Lords and the Eternal Forge. Pyrrhian Flame births volcanic titans and flame-bonded warriors. Time moves faster here, aging all who enter.",
                artifact_name: "Heart of Cinder",
                artifact_description: "A molten gem that grants destructive power and burns away lies.",
                icon: "🔥",
                type: "fire"
            },
            {
                id: 2,
                name: "Zepharion Heights",
                element: "Wind / Sky / Sound",
                description: "Floating islands around an eternal cyclone",
                lore: "Floating islands encircle a permanent cyclone known as The Whisper. Skyborn sages ride wind-serpents and wield songs that bend reality.",
                artifact_name: "Aeon Harp",
                artifact_description: "Plays melodies that control storms and memories.",
                icon: "🌪️",
                type: "wind"
            },
            {
                id: 3,
                name: "Terravine Hollow",
                element: "Stone / Growth / Gravity",
                description: "Ancient buried realm of stone titans",
                lore: "An ancient, buried realm where roots grow like veins and sentient stone titans slumber. Once a great civilization, now petrified into time.",
                artifact_name: "Verdant Core",
                artifact_description: "Grants dominion over life, soil, and rebirth.",
                icon: "🌿",
                type: "earth"
            },
            {
                id: 4,
                name: "Thalorion Depths",
                element: "Water / Ice / Depth",
                description: "Submerged empire of the Moonbound Court",
                lore: "A submerged empire ruled by the Moonbound Court. Time slows here, and the ocean whispers ancient truths. Home to leviathans and drowned prophets.",
                artifact_name: "Tideglass Mirror",
                artifact_description: "Sees through illusions and to possible futures.",
                icon: "🌊",
                type: "water"
            },
            {
                id: 5,
                name: "Virelya",
                element: "Light / Purity / Illumination",
                description: "Blinding paradise of pure truth",
                lore: "A blinding paradise where truth manifests as form. Ruled by beings known as Radiants. Mortals must wear veilshades to even look upon it.",
                artifact_name: "Lumen Shard",
                artifact_description: "Reveals the true name of anything it touches.",
                icon: "✨",
                type: "light"
            },
            {
                id: 6,
                name: "Umbros",
                element: "Shadow / Secrets / Corruption",
                description: "Void-split realm where light cannot reach",
                lore: "Light cannot reach this void-split realm. Every whisper is a thought stolen, every step a forgotten path. Shadowmages barter in memories.",
                artifact_name: "Eclipse Fang",
                artifact_description: "Severs light, binding a soul to darkness.",
                icon: "🌑",
                type: "shadow"
            },
            {
                id: 7,
                name: "Nyxthar",
                element: "Null / Anti-Matter / Entropy",
                description: "Realm where reality collapses inward",
                lore: "A realm where reality collapses inward. Voidwalkers and Silence Priests seek ultimate release from being. To enter is to forget existence.",
                artifact_name: "Hollow Crown",
                artifact_description: "Nullifies all magic and erases history.",
                icon: "⚫",
                type: "void"
            },
            {
                id: 8,
                name: "Aetherion",
                element: "Spirit / Soul / Dream",
                description: "Realm between realms of dreaming dead",
                lore: "The realm between realms, where the dreaming dead speak. Time is nonlinear, and the laws of logic bend to desire. Spirits travel as thought.",
                artifact_name: "Soulforge Locket",
                artifact_description: "Binds spirits to bodies or frees them eternally.",
                icon: "👻",
                type: "spirit"
            },
            {
                id: 9,
                name: "Chronarxis",
                element: "Time / Fate / Chronomancy",
                description: "Spiral palace of fractured timelines",
                lore: "A spiral palace where timelines fracture and reform. Timekeepers judge anomalies and anomalies fight back. Accessed only through ancient rituals.",
                artifact_name: "Clockheart Mechanism",
                artifact_description: "Rewinds one moment once, but at a cost.",
                icon: "⏰",
                type: "time"
            },
            {
                id: 10,
                name: "Technarok",
                element: "Metal / Machines / Order",
                description: "Fusion of steel gods and nano-intelligences",
                lore: "A fusion of ancient steel gods and nano-intelligences. Run on logic and decay. Home to sentient forges and recursive codebeasts.",
                artifact_name: "Iron Synapse",
                artifact_description: "Merges user with machine intelligence.",
                icon: "⚙️",
                type: "metal"
            }
        ];
        
        displayRealms(realms);
    } catch (error) {
        console.error('Error loading realms:', error);
        showToast('Error loading realms', 'error');
    } finally {
        showLoading(false);
    }
}

function displayRealms(realms) {
    const realmsGrid = document.getElementById('realms-grid');
    
    realmsGrid.innerHTML = realms.map(realm => `
        <div class="realm-card" onclick="showRealmDetails(${realm.id})">
            <div class="realm-header ${realm.type}">
                <div class="realm-icon">${realm.icon}</div>
            </div>
            <div class="realm-content">
                <div class="realm-name">${realm.name}</div>
                <div class="realm-element">
                    <i class="fas fa-magic"></i>
                    ${realm.element}
                </div>
                <div class="realm-description">${realm.description}</div>
                <div class="realm-artifact">
                    <div class="artifact-name">
                        <i class="fas fa-gem"></i>
                        ${realm.artifact_name}
                    </div>
                    <div class="artifact-description">${realm.artifact_description}</div>
                </div>
            </div>
        </div>
    `).join('');
}

function showRealmDetails(realmId) {
    // Find realm from the hardcoded data
    const realms = [
        {
            id: 1,
            name: "Pyrrhian Flame",
            element: "Fire / Heat / Chaos",
            description: "Realm of eternal fire and volcanic titans",
            lore: "Home to the Salamandrine Lords and the Eternal Forge. Pyrrhian Flame births volcanic titans and flame-bonded warriors. Time moves faster here, aging all who enter.",
            artifact_name: "Heart of Cinder",
            artifact_description: "A molten gem that grants destructive power and burns away lies.",
            icon: "🔥",
            type: "fire"
        },
        {
            id: 2,
            name: "Zepharion Heights",
            element: "Wind / Sky / Sound",
            description: "Floating islands around an eternal cyclone",
            lore: "Floating islands encircle a permanent cyclone known as The Whisper. Skyborn sages ride wind-serpents and wield songs that bend reality.",
            artifact_name: "Aeon Harp",
            artifact_description: "Plays melodies that control storms and memories.",
            icon: "🌪️",
            type: "wind"
        },
        {
            id: 3,
            name: "Terravine Hollow",
            element: "Stone / Growth / Gravity",
            description: "Ancient buried realm of stone titans",
            lore: "An ancient, buried realm where roots grow like veins and sentient stone titans slumber. Once a great civilization, now petrified into time.",
            artifact_name: "Verdant Core",
            artifact_description: "Grants dominion over life, soil, and rebirth.",
            icon: "🌿",
            type: "earth"
        },
        {
            id: 4,
            name: "Thalorion Depths",
            element: "Water / Ice / Depth",
            description: "Submerged empire of the Moonbound Court",
            lore: "A submerged empire ruled by the Moonbound Court. Time slows here, and the ocean whispers ancient truths. Home to leviathans and drowned prophets.",
            artifact_name: "Tideglass Mirror",
            artifact_description: "Sees through illusions and to possible futures.",
            icon: "🌊",
            type: "water"
        },
        {
            id: 5,
            name: "Virelya",
            element: "Light / Purity / Illumination",
            description: "Blinding paradise of pure truth",
            lore: "A blinding paradise where truth manifests as form. Ruled by beings known as Radiants. Mortals must wear veilshades to even look upon it.",
            artifact_name: "Lumen Shard",
            artifact_description: "Reveals the true name of anything it touches.",
            icon: "✨",
            type: "light"
        },
        {
            id: 6,
            name: "Umbros",
            element: "Shadow / Secrets / Corruption",
            description: "Void-split realm where light cannot reach",
            lore: "Light cannot reach this void-split realm. Every whisper is a thought stolen, every step a forgotten path. Shadowmages barter in memories.",
            artifact_name: "Eclipse Fang",
            artifact_description: "Severs light, binding a soul to darkness.",
            icon: "🌑",
            type: "shadow"
        },
        {
            id: 7,
            name: "Nyxthar",
            element: "Null / Anti-Matter / Entropy",
            description: "Realm where reality collapses inward",
            lore: "A realm where reality collapses inward. Voidwalkers and Silence Priests seek ultimate release from being. To enter is to forget existence.",
            artifact_name: "Hollow Crown",
            artifact_description: "Nullifies all magic and erases history.",
            icon: "⚫",
            type: "void"
        },
        {
            id: 8,
            name: "Aetherion",
            element: "Spirit / Soul / Dream",
            description: "Realm between realms of dreaming dead",
            lore: "The realm between realms, where the dreaming dead speak. Time is nonlinear, and the laws of logic bend to desire. Spirits travel as thought.",
            artifact_name: "Soulforge Locket",
            artifact_description: "Binds spirits to bodies or frees them eternally.",
            icon: "👻",
            type: "spirit"
        },
        {
            id: 9,
            name: "Chronarxis",
            element: "Time / Fate / Chronomancy",
            description: "Spiral palace of fractured timelines",
            lore: "A spiral palace where timelines fracture and reform. Timekeepers judge anomalies and anomalies fight back. Accessed only through ancient rituals.",
            artifact_name: "Clockheart Mechanism",
            artifact_description: "Rewinds one moment once, but at a cost.",
            icon: "⏰",
            type: "time"
        },
        {
            id: 10,
            name: "Technarok",
            element: "Metal / Machines / Order",
            description: "Fusion of steel gods and nano-intelligences",
            lore: "A fusion of ancient steel gods and nano-intelligences. Run on logic and decay. Home to sentient forges and recursive codebeasts.",
            artifact_name: "Iron Synapse",
            artifact_description: "Merges user with machine intelligence.",
            icon: "⚙️",
            type: "metal"
        }
    ];
    
    const realm = realms.find(r => r.id === realmId);
    if (!realm) return;
    
    displayRealmModal(realm);
}

function displayRealmModal(realm) {
    const modalHTML = `
        <div id="realm-details-modal" class="modal show">
            <div class="modal-content realm-modal-content">
                <div class="realm-modal-header ${realm.type}">
                    <div class="realm-modal-icon">${realm.icon}</div>
                    <div class="realm-modal-info">
                        <h2 class="realm-modal-name">${realm.name}</h2>
                        <div class="realm-modal-element">${realm.element}</div>
                    </div>
                    <button class="close-btn" onclick="closeRealmModal()">
                        <i class="fas fa-times"></i>
                    </button>
                </div>
                
                <div class="realm-modal-content">
                    <div class="realm-lore-section">
                        <h3><i class="fas fa-scroll"></i> Ancient Lore</h3>
                        <p class="realm-lore">${realm.lore}</p>
                    </div>
                    
                    <div class="realm-artifact-section">
                        <h3><i class="fas fa-gem"></i> Legendary Artifact</h3>
                        <div class="artifact-card">
                            <div class="artifact-header">
                                <div class="artifact-icon">
                                    <i class="fas fa-diamond"></i>
                                </div>
                                <div class="artifact-info">
                                    <div class="artifact-title">${realm.artifact_name}</div>
                                    <div class="artifact-power">${realm.artifact_description}</div>
                                </div>
                            </div>
                        </div>
                    </div>
                    
                    <div class="realm-actions">
                        <button class="action-btn primary" onclick="createWizardInRealm('${realm.name}')">
                            <i class="fas fa-plus"></i>
                            Create Wizard Here
                        </button>
                        <button class="action-btn secondary" onclick="exploreRealmWizards('${realm.name}')">
                            <i class="fas fa-search"></i>
                            Explore Wizards
                        </button>
                        <button class="action-btn secondary" onclick="exploreRealmJobs('${realm.name}')">
                            <i class="fas fa-briefcase"></i>
                            Explore Jobs
                        </button>
                    </div>
                </div>
            </div>
        </div>
    `;
    
    // Remove existing modal if any
    const existingModal = document.getElementById('realm-details-modal');
    if (existingModal) {
        existingModal.remove();
    }
    
    // Add modal to body
    document.body.insertAdjacentHTML('beforeend', modalHTML);
}

function closeRealmModal() {
    const modal = document.getElementById('realm-details-modal');
    if (modal) {
        modal.remove();
    }
}

function createWizardInRealm(realmName) {
    closeRealmModal();
    showWizards();
    showCreateWizard();
    // Pre-select the realm
    const realmSelect = document.getElementById('wizard-realm');
    if (realmSelect) {
        realmSelect.value = realmName;
    }
}

function exploreRealmWizards(realmName) {
    closeRealmModal();
    showExploreWizards();
    // Pre-select the realm filter
    const realmFilter = document.getElementById('realm-filter');
    if (realmFilter) {
        realmFilter.value = realmName;
        filterWizardsByRealm();
    }
}

function exploreRealmJobs(realmName) {
    closeRealmModal();
    showJobs();
    // Pre-select the realm filter
    setTimeout(() => {
        const realmFilter = document.getElementById('job-realm-filter');
        if (realmFilter) {
            realmFilter.value = realmName;
            filterJobsByRealm();
        }
    }, 100);
}

// Jobs Functions
let jobs = [];
let filteredJobs = [];

async function loadJobs() {
    try {
        showLoading(true);
        const response = await api.getJobs();
        console.log('Jobs API response:', response); // Debug log
        jobs = response.jobs || [];
        console.log('Processed jobs:', jobs); // Debug log
        filteredJobs = [...jobs];
        displayJobs();
    } catch (error) {
        console.error('Error loading jobs:', error);
        showToast('Error loading jobs', 'error');
        jobs = [];
        filteredJobs = [];
    } finally {
        showLoading(false);
    }
}

async function loadJobFilters() {
    try {
        // Load realms for filter
        const realmResponse = await api.getRealms();
        const realms = realmResponse.realms || [];
        
        const realmFilter = document.getElementById('job-realm-filter');
        if (realmFilter) {
            realmFilter.innerHTML = '<option value="">All Realms</option>' +
                realms.map(realm => `<option value="${realm.name}">${realm.name}</option>`).join('');
        }
    } catch (error) {
        console.error('Error loading job filters:', error);
    }
}

function displayJobs() {
    const jobsGrid = document.getElementById('jobs-grid');
    
    if (filteredJobs.length === 0) {
        jobsGrid.innerHTML = `
            <div class="empty-state">
                <i class="fas fa-briefcase"></i>
                <h3>No Jobs Available</h3>
                <p>Check back later for new opportunities across the realms!</p>
            </div>
        `;
        return;
    }
    
    jobsGrid.innerHTML = filteredJobs.map(job => `
        <div class="job-card" onclick="showJobDetails(${job.id})">
            <div class="job-header ${job.required_element.toLowerCase()}">
                <div class="job-title">${job.title}</div>
                <div class="job-difficulty">${job.difficulty}</div>
            </div>
            
            <div class="job-content">
                <div class="job-meta">
                    <div class="job-meta-item">
                        <i class="fas fa-magic"></i>
                        <span>${job.required_element}</span>
                    </div>
                    <div class="job-meta-item">
                        <i class="fas fa-map-marker-alt"></i>
                        <span>${job.location || job.realm_name}</span>
                    </div>
                    <div class="job-meta-item">
                        <i class="fas fa-clock"></i>
                        <span>${formatDuration(job.duration_minutes)}</span>
                    </div>
                    <div class="job-meta-item">
                        <i class="fas fa-users"></i>
                        <span>${job.currently_assigned || 0}/${job.max_wizards || 1}</span>
                    </div>
                </div>
                
                <div class="job-description">
                    ${job.description}
                </div>
                
                <div class="job-rewards">
                    <div class="job-reward-item">
                        <i class="fas fa-coins"></i>
                        <span>Mana Reward: <span class="job-reward-value">${formatNumber(Math.round(job.mana_reward_per_hour * ((job.duration_minutes || 0) / 60)))} total</span></span>
                    </div>
                    <div class="job-reward-item">
                        <i class="fas fa-hourglass-half"></i>
                        <span>Rate: <span class="job-reward-value">${formatNumber(job.mana_reward_per_hour)}/hour</span></span>
                    </div>
                </div>
                
                <div class="job-availability ${(job.currently_assigned || 0) >= (job.max_wizards || 1) ? 'full' : ''}">
                    <span class="job-slots ${(job.currently_assigned || 0) >= (job.max_wizards || 1) ? 'full' : ''}">
                        ${(job.currently_assigned || 0) >= (job.max_wizards || 1) ? 'FULL' : 'Available'}
                    </span>
                    <span>${job.currently_assigned || 0}/${job.max_wizards || 1} wizards assigned</span>
                </div>
                
                <div class="job-actions">
                    <button class="job-assign-btn" 
                            onclick="event.stopPropagation(); assignWizardToJob(${job.id})"
                            ${(job.currently_assigned || 0) >= (job.max_wizards || 1) ? 'disabled' : ''}>
                        <i class="fas fa-user-plus"></i>
                        ${(job.currently_assigned || 0) >= (job.max_wizards || 1) ? 'Job Full' : 'Assign Wizard'}
                    </button>
                </div>
            </div>
        </div>
    `).join('');
}

function filterJobsByRealm() {
    const realmFilter = document.getElementById('job-realm-filter').value;
    applyJobFilters();
}

function filterJobsByElement() {
    const elementFilter = document.getElementById('job-element-filter').value;
    applyJobFilters();
}

function filterJobsByDifficulty() {
    const difficultyFilter = document.getElementById('job-difficulty-filter').value;
    applyJobFilters();
}

function applyJobFilters() {
    const realmFilter = document.getElementById('job-realm-filter').value;
    const elementFilter = document.getElementById('job-element-filter').value;
    const difficultyFilter = document.getElementById('job-difficulty-filter').value;
    
    filteredJobs = jobs.filter(job => {
        const matchesRealm = !realmFilter || job.realm_name === realmFilter;
        const matchesElement = !elementFilter || job.required_element === elementFilter;
        const matchesDifficulty = !difficultyFilter || job.difficulty === difficultyFilter;
        
        return matchesRealm && matchesElement && matchesDifficulty;
    });
    
    displayJobs();
}

async function assignWizardToJob(jobId) {
    try {
        // Get user's wizards that match the job requirements
        const job = jobs.find(j => j.id === jobId);
        if (!job) {
            showToast('Job not found', 'error');
            return;
        }
        
        const compatibleWizards = wizards.filter(wizard => 
            wizard.element === job.required_element && 
            (wizard.level || 1) >= job.required_level
        );
        
        if (compatibleWizards.length === 0) {
            if (wizards.some(w => w.element === job.required_element)) {
                showToast(`Your ${job.required_element} wizards need to be level ${job.required_level} or higher for this job`, 'warning');
            } else {
                showToast(`You need a ${job.required_element} wizard for this job`, 'warning');
            }
            return;
        }
        
        // Show wizard selection modal
        showWizardSelectionModal(jobId, compatibleWizards);
        
    } catch (error) {
        console.error('Error assigning wizard to job:', error);
        showToast('Error assigning wizard to job', 'error');
    }
}

function showWizardSelectionModal(jobId, compatibleWizards) {
    const job = jobs.find(j => j.id === jobId);
    
    const modalHTML = `
        <div id="wizard-selection-modal" class="modal show">
            <div class="modal-content">
                <div class="modal-header">
                    <h3><i class="fas fa-user-plus"></i> Assign Wizard to ${job.title}</h3>
                    <span class="close" onclick="closeWizardSelectionModal()">&times;</span>
                </div>
                <div class="modal-body">
                    <p>Select a ${job.required_element} wizard to assign to this job:</p>
                    <div class="wizard-selection-grid">
                        ${compatibleWizards.map(wizard => `
                            <div class="wizard-selection-card" onclick="confirmWizardAssignment(${jobId}, ${wizard.id})">
                                <div class="wizard-avatar">
                                    <i class="fas fa-hat-wizard wizard-element-${wizard.element.toLowerCase()}"></i>
                                </div>
                                <div class="wizard-name">${wizard.name}</div>
                                <div class="wizard-element">${wizard.element} • Level ${wizard.level || 1}</div>
                                <div class="wizard-mana">${formatNumber(wizard.mana_balance)} Mana</div>
                                <div class="wizard-exp">${formatNumber(wizard.experience_points || 0)} EXP</div>
                            </div>
                        `).join('')}
                    </div>
                </div>
            </div>
        </div>
    `;
    
    document.body.insertAdjacentHTML('beforeend', modalHTML);
}

function closeWizardSelectionModal() {
    const modal = document.getElementById('wizard-selection-modal');
    if (modal) {
        modal.remove();
    }
}

async function confirmWizardAssignment(jobId, wizardId) {
    try {
        showLoading(true);
        await api.assignWizardToJob(jobId, wizardId);
        
        // Update job data
        const job = jobs.find(j => j.id === jobId);
        if (job) {
            job.currently_assigned += 1;
        }
        
        closeWizardSelectionModal();
        
        // Refresh all relevant data to reflect the assignment
        await Promise.all([
            loadWizards(),     // Refresh wizard data
            loadJobs(),        // Refresh job data
            loadActiveJobs()   // Refresh active jobs
        ]);
        
        showToast('Wizard assigned to job successfully!', 'success');
    } catch (error) {
        console.error('Error confirming wizard assignment:', error);
        showToast(error.message || 'Error assigning wizard to job', 'error');
    } finally {
        showLoading(false);
    }
}

function showJobDetails(jobId) {
    const job = jobs.find(j => j.id === jobId);
    if (!job) return;
    
    // This could expand to show a detailed modal with assignment history, etc.
    console.log('Job details for:', job);
}

function showCreateJob() {
    const highLevelWizards = wizards.filter(wizard => (wizard.level || 1) >= 10);
    
    if (highLevelWizards.length === 0) {
        showToast('You need a level 10+ wizard to create jobs', 'warning');
        return;
    }
    
    const modalHTML = `
        <div id="create-job-modal" class="modal show">
            <div class="modal-content">
                <div class="modal-header">
                    <h3><i class="fas fa-plus"></i> Create New Job</h3>
                    <span class="close" onclick="closeCreateJobModal()">&times;</span>
                </div>
                <div class="modal-body">
                    <form id="create-job-form" onsubmit="handleCreateJob(event)">
                        <div class="form-group">
                            <label for="job-title">Job Title</label>
                            <input type="text" id="job-title" placeholder="Enter job title" required>
                        </div>
                        
                        <div class="form-group">
                            <label for="job-description">Description</label>
                            <textarea id="job-description" placeholder="Describe the job requirements and duties" rows="4" required></textarea>
                        </div>
                        
                        <div class="form-group select-group">
                            <label for="job-realm">Realm</label>
                            <select id="job-realm" required>
                                <option value="">Select Realm</option>
                            </select>
                        </div>
                        
                        <div class="form-group select-group">
                            <label for="job-element">Required Element</label>
                            <select id="job-element" required>
                                <option value="">Select Element</option>
                                <option value="Fire">Fire</option>
                                <option value="Air">Air</option>
                                <option value="Earth">Earth</option>
                                <option value="Water">Water</option>
                                <option value="Light">Light</option>
                            </select>
                        </div>
                        
                        <div class="form-group">
                            <label for="job-level">Required Level</label>
                            <input type="number" id="job-level" min="1" max="50" value="1" required>
                        </div>
                        
                        <div class="form-group">
                            <label for="job-mana-reward">Mana Reward per Hour</label>
                            <input type="number" id="job-mana-reward" min="10" max="1000" value="100" required>
                        </div>
                        
                        <div class="form-group">
                            <label for="job-exp-reward">EXP Reward per Hour</label>
                            <input type="number" id="job-exp-reward" min="5" max="100" value="10" required>
                        </div>
                        
                        <div class="form-group">
                            <label for="job-duration">Duration (hours)</label>
                            <input type="number" id="job-duration" min="1" max="48" value="8" required>
                        </div>
                        
                        <div class="form-group">
                            <label for="job-max-wizards">Maximum Wizards</label>
                            <input type="number" id="job-max-wizards" min="1" max="10" value="1" required>
                        </div>
                        
                        <div class="form-group select-group">
                            <label for="job-difficulty">Difficulty</label>
                            <select id="job-difficulty" required>
                                <option value="">Select Difficulty</option>
                                <option value="Easy">Easy</option>
                                <option value="Medium">Medium</option>
                                <option value="Hard">Hard</option>
                                <option value="Expert">Expert</option>
                                <option value="Legendary">Legendary</option>
                            </select>
                        </div>
                        
                        <div class="form-group">
                            <label for="job-location">Location (optional)</label>
                            <input type="text" id="job-location" placeholder="Specific location within the realm">
                        </div>
                        
                        <div class="form-actions">
                            <button type="button" class="btn btn-secondary" onclick="closeCreateJobModal()">Cancel</button>
                            <button type="submit" class="btn btn-primary">Create Job</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    `;
    
    document.body.insertAdjacentHTML('beforeend', modalHTML);
    
    // Populate realm dropdown
    loadJobFilters().then(() => {
        const realmSelect = document.getElementById('job-realm');
        const filterRealm = document.getElementById('job-realm-filter');
        if (realmSelect && filterRealm) {
            realmSelect.innerHTML = filterRealm.innerHTML;
        }
    });
}

function closeCreateJobModal() {
    const modal = document.getElementById('create-job-modal');
    if (modal) {
        modal.remove();
    }
}

async function handleCreateJob(event) {
    event.preventDefault();
    
    const formData = {
        title: document.getElementById('job-title').value,
        description: document.getElementById('job-description').value,
        realm_name: document.getElementById('job-realm').value,
        required_element: document.getElementById('job-element').value,
        required_level: parseInt(document.getElementById('job-level').value),
        mana_reward_per_hour: parseInt(document.getElementById('job-mana-reward').value),
        exp_reward_per_hour: parseInt(document.getElementById('job-exp-reward').value),
        duration_minutes: parseInt(document.getElementById('job-duration').value) * 60,
        max_wizards: parseInt(document.getElementById('job-max-wizards').value),
        difficulty: document.getElementById('job-difficulty').value,
        location: document.getElementById('job-location').value || null,
        job_type: 'Player Created'
    };
    
    try {
        showLoading(true);
        await api.createJob(formData);
        
        closeCreateJobModal();
        loadJobs(); // Refresh jobs list
        showToast('Job created successfully!', 'success');
    } catch (error) {
        console.error('Error creating job:', error);
        showToast(error.message || 'Error creating job', 'error');
    } finally {
        showLoading(false);
    }
}

// Active Jobs Functions
async function loadActiveJobs() {
    try {
        // Only load if user has wizards
        if (!wizards || wizards.length === 0) {
            hideActiveJobsSection();
            return;
        }

        showLoading(true);
        
        // Get all active assignments for user's wizards (both assigned and in_progress)
        const allActiveJobs = [];
        for (const wizard of wizards) {
            try {
                // Get both assigned and in_progress assignments
                const [assignedResponse, inProgressResponse] = await Promise.all([
                    api.getJobAssignments(wizard.id, 'assigned'),
                    api.getJobAssignments(wizard.id, 'in_progress')
                ]);
                
                // Combine all active assignments
                const assignments = [
                    ...(assignedResponse.assignments || []),
                    ...(inProgressResponse.assignments || [])
                ];
                
                if (assignments.length > 0) {
                    allActiveJobs.push(...assignments.map(assignment => ({
                        ...assignment,
                        wizard: wizard
                    })));
                }
            } catch (error) {
                console.error(`Error loading jobs for wizard ${wizard.id}:`, error);
            }
        }

        displayActiveJobs(allActiveJobs);
    } catch (error) {
        console.error('Error loading active jobs:', error);
        hideActiveJobsSection();
    } finally {
        showLoading(false);
    }
}

function displayActiveJobs(activeJobs) {
    const section = document.getElementById('active-jobs-section');
    const grid = document.getElementById('active-jobs-grid');
    
    if (!activeJobs || activeJobs.length === 0) {
        hideActiveJobsSection();
        return;
    }

    // Store assignment data in global state for visual progress updates
    window.currentAssignments = activeJobs;

    section.style.display = 'block';
    
    grid.innerHTML = activeJobs.map(assignment => {
        const progress = calculateJobProgress(assignment);
        const earnings = calculateEarnings(assignment, progress);
        
        return `
            <div class="active-job-card" data-assignment-id="${assignment.id}">
                <div class="active-job-header">
                    <div>
                        <h3 class="active-job-title">${assignment.job?.title || 'Job Assignment'}</h3>
                        <div class="active-job-wizard">
                            <i class="fas fa-hat-wizard wizard-element-${assignment.wizard.element.toLowerCase()}"></i>
                            ${assignment.wizard.name}
                        </div>
                    </div>
                    <div class="active-job-status ${assignment.status}">
                        ${assignment.status.replace('_', ' ')}
                    </div>
                </div>
                
                <div class="active-job-progress">
                    <div class="progress-info">
                        <span class="progress-label">Progress</span>
                        <span class="progress-time">${progress.timeRemaining}</span>
                    </div>
                    <div class="progress-bar">
                        <div class="progress-fill" style="width: ${progress.percentage}%"></div>
                    </div>
                </div>
                
                <div class="active-job-rewards">
                    <div class="reward-item">
                        <div class="reward-icon mana">
                            <i class="fas fa-coins"></i>
                        </div>
                        <div class="reward-value">${formatNumber(earnings.mana)}</div>
                        <div class="reward-label">Mana Earned</div>
                    </div>
                    <div class="reward-item">
                        <div class="reward-icon exp">
                            <i class="fas fa-star"></i>
                        </div>
                        <div class="reward-value">${formatNumber(earnings.exp)}</div>
                        <div class="reward-label">EXP Earned</div>
                    </div>
                </div>
                
                <div class="active-job-actions">
                    ${progress.percentage >= 100 ? `
                        <button class="job-action-btn complete" onclick="completeJobAssignment(${assignment.id})">
                            <i class="fas fa-check"></i> Complete Job
                        </button>
                    ` : `
                        <button class="job-action-btn view" onclick="viewJobDetails(${assignment.job_id})">
                            <i class="fas fa-eye"></i> View Job
                        </button>
                    `}
                    <button class="job-action-btn cancel" onclick="cancelJobAssignment(${assignment.id})">
                        <i class="fas fa-times"></i> Cancel
                    </button>
                </div>
            </div>
        `;
    }).join('');
}

function hideActiveJobsSection() {
    const section = document.getElementById('active-jobs-section');
    const grid = document.getElementById('active-jobs-grid');
    
    section.style.display = 'none';
    grid.innerHTML = `
        <div class="no-active-jobs">
            <i class="fas fa-briefcase"></i>
            <h3>No Active Jobs</h3>
            <p>Your wizards aren't currently working on any jobs. Visit the Jobs page to find assignments!</p>
        </div>
    `;
}

function calculateJobProgress(assignment) {
    // The backend now handles real-time progress calculation
    // We just need to get the current progress from the assignment
    const progress = assignment.progress || {};
    const percentage = progress.progress_percentage || 0;
    const durationMinutes = assignment.job?.duration_minutes || 480;
    const timeWorkedMinutes = progress.time_worked_minutes || 0;
    
    let timeRemaining;
    if (percentage >= 100) {
        timeRemaining = 'Ready to complete!';
    } else {
        const remainingMinutes = Math.max(0, durationMinutes - timeWorkedMinutes);
        
        if (remainingMinutes <= 0) {
            timeRemaining = 'Ready to complete!';
        } else {
            const remainingHours = Math.floor(remainingMinutes / 60);
            const remainingMins = remainingMinutes % 60;
            
            if (remainingHours > 0) {
                timeRemaining = `${remainingHours}h ${remainingMins}m left`;
            } else if (remainingMins > 0) {
                timeRemaining = `${remainingMins}m left`;
            } else {
                timeRemaining = 'Almost done!';
            }
        }
    }
    
    return {
        percentage: percentage,
        timeRemaining,
        workedMinutes: timeWorkedMinutes,
        isComplete: percentage >= 100
    };
}

function calculateEarnings(assignment, progress) {
    const job = assignment.job;
    if (!job) return { mana: 0, exp: 0 };
    
    const durationMinutes = job.duration_minutes || 0;
    if (durationMinutes === 0) return { mana: 0, exp: 0 };
    
    const minutesWorked = (progress.percentage / 100) * durationMinutes;
    const hoursWorked = minutesWorked / 60;
    const mana = Math.floor(hoursWorked * (job.mana_reward_per_hour || 0));
    const exp = Math.floor(hoursWorked * (job.exp_reward_per_hour || 0));
    
    return { mana, exp };
}

async function completeJobAssignment(assignmentId) {
    try {
        showLoading(true);
        await api.completeJobAssignment(assignmentId);
        showToast('Job completed successfully! Rewards have been added to your wizard.', 'success');
        
        // Refresh both wizards and active jobs
        loadWizards();
        loadActiveJobs();
    } catch (error) {
        console.error('Error completing job assignment:', error);
        showToast(error.message || 'Error completing job assignment', 'error');
    } finally {
        showLoading(false);
    }
}

async function cancelJobAssignment(assignmentId) {
    if (!confirm('Are you sure you want to cancel this job assignment? Any progress will be lost.')) {
        return;
    }
    
    try {
        showLoading(true);
        await api.cancelJobAssignment(assignmentId, 'Cancelled by user');
        showToast('Job assignment cancelled.', 'success');
        
        // Refresh both wizards and active jobs
        loadWizards();
        loadActiveJobs();
    } catch (error) {
        console.error('Error cancelling job assignment:', error);
        showToast(error.message || 'Error cancelling job assignment', 'error');
    } finally {
        showLoading(false);
    }
}

function viewJobDetails(jobId) {
    // Navigate to jobs page and show specific job
    showJobs();
    // Scroll to or highlight the specific job
    setTimeout(() => {
        const jobCard = document.querySelector(`[onclick*="showJobDetails(${jobId})"]`);
        if (jobCard) {
            jobCard.scrollIntoView({ behavior: 'smooth', block: 'center' });
            jobCard.style.animation = 'highlight 2s ease-in-out';
        }
    }, 500);
}

// Auto-refresh active jobs with real-time data
let jobDataRefreshInterval;

function startActiveJobsRefresh() {
    // Clear any existing intervals
    if (jobDataRefreshInterval) {
        clearInterval(jobDataRefreshInterval);
    }
    
    // Set up real-time data refresh every 3 seconds
    jobDataRefreshInterval = setInterval(() => {
        const wizardsPage = document.getElementById('wizards');
        if (wizardsPage && wizardsPage.style.display !== 'none') {
            loadActiveJobsSilently(); // Fetch fresh data from backend
        }
    }, 3000);
}

// Real-time progress updates are now handled by backend ticker
// Frontend just needs to fetch fresh data periodically

// Progress updates are now handled by the backend ticker system
// No need for frontend progress simulation

async function loadActiveJobsSilently() {
    try {
        // Only load if user has wizards
        if (!wizards || wizards.length === 0) {
            hideActiveJobsSection();
            return;
        }

        // Get all active assignments for user's wizards (both assigned and in_progress)
        // Don't show loading spinner for silent updates
        const allActiveJobs = [];
        for (const wizard of wizards) {
            try {
                // Get both assigned and in_progress assignments
                const [assignedResponse, inProgressResponse] = await Promise.all([
                    api.getJobAssignments(wizard.id, 'assigned'),
                    api.getJobAssignments(wizard.id, 'in_progress')
                ]);
                
                // Combine all active assignments
                const assignments = [
                    ...(assignedResponse.assignments || []),
                    ...(inProgressResponse.assignments || [])
                ];
                
                if (assignments.length > 0) {
                    allActiveJobs.push(...assignments.map(assignment => ({
                        ...assignment,
                        wizard: wizard
                    })));
                }
            } catch (error) {
                console.error(`Error loading jobs for wizard ${wizard.id}:`, error);
            }
        }

        displayActiveJobs(allActiveJobs);
    } catch (error) {
        console.error('Error loading active jobs silently:', error);
        hideActiveJobsSection();
    }
}

// Real-time progress is now handled by backend ticker - no manual updates needed

function stopActiveJobsRefresh() {
    if (jobDataRefreshInterval) {
        clearInterval(jobDataRefreshInterval);
        jobDataRefreshInterval = null;
    }
}

// Simple page initialization
document.addEventListener('DOMContentLoaded', function() {
    console.log('MysticFunds app loaded');
    startActiveJobsRefresh();
});