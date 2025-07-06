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

// Artifact and Spell Utility Functions
function getArtifactIcon(type) {
    switch (type.toLowerCase()) {
        case 'weapon': return 'fa-sword';
        case 'armor': return 'fa-shield-alt';
        case 'accessory': return 'fa-ring';
        case 'tool': return 'fa-hammer';
        case 'potion': return 'fa-flask';
        case 'tome': return 'fa-book';
        case 'orb': return 'fa-circle';
        case 'staff': return 'fa-magic';
        case 'wand': return 'fa-magic';
        case 'crystal': return 'fa-gem';
        case 'scroll': return 'fa-scroll';
        case 'amulet': return 'fa-medallion';
        case 'cloak': return 'fa-user-shield';
        case 'boots': return 'fa-shoe-prints';
        case 'gloves': return 'fa-hand-paper';
        case 'helmet': return 'fa-hard-hat';
        case 'belt': return 'fa-circle-notch';
        case 'robe': return 'fa-user-ninja';
        case 'hat': return 'fa-hat-wizard';
        case 'mask': return 'fa-mask';
        default: return 'fa-magic';
    }
}

function getSpellIcon(school) {
    switch (school.toLowerCase()) {
        case 'fire': return 'fa-fire';
        case 'water': return 'fa-tint';
        case 'earth': return 'fa-mountain';
        case 'air': return 'fa-wind';
        case 'light': return 'fa-sun';
        case 'dark': return 'fa-moon';
        case 'nature': return 'fa-leaf';
        case 'arcane': return 'fa-star';
        case 'healing': return 'fa-heart';
        case 'necromancy': return 'fa-skull';
        case 'illusion': return 'fa-eye';
        case 'enchantment': return 'fa-magic';
        case 'divination': return 'fa-crystal-ball';
        case 'transmutation': return 'fa-exchange-alt';
        case 'conjuration': return 'fa-portal';
        case 'evocation': return 'fa-bolt';
        case 'abjuration': return 'fa-shield';
        case 'universal': return 'fa-infinity';
        default: return 'fa-magic';
    }
}

// Sample data generators for demonstration
function getSampleArtifacts(wizard) {
    const level = wizard.level || 1;
    const element = wizard.element.toLowerCase();
    const artifacts = [];
    
    // Basic artifacts every wizard has
    if (level >= 1) {
        artifacts.push({
            name: `${wizard.element} Apprentice Robe`,
            type: 'armor',
            rarity: 'common',
            equipped: true
        });
    }
    
    if (level >= 3) {
        artifacts.push({
            name: `${wizard.element} Staff`,
            type: 'weapon',
            rarity: level >= 10 ? 'rare' : 'uncommon',
            equipped: true
        });
    }
    
    if (level >= 5) {
        artifacts.push({
            name: `Amulet of ${wizard.element} Protection`,
            type: 'accessory',
            rarity: level >= 15 ? 'epic' : 'uncommon',
            equipped: true
        });
    }
    
    if (level >= 10) {
        artifacts.push({
            name: `${wizard.element} Crystal Orb`,
            type: 'orb',
            rarity: level >= 20 ? 'legendary' : 'rare',
            equipped: true
        });
    }
    
    return artifacts;
}

function getSampleSpells(wizard) {
    const level = wizard.level || 1;
    const element = wizard.element.toLowerCase();
    const spells = [];
    
    // Basic spells based on element and level
    spells.push({
        name: `${wizard.element} Bolt`,
        school: element,
        level: 1
    });
    
    if (level >= 3) {
        spells.push({
            name: `${wizard.element} Shield`,
            school: element,
            level: 2
        });
    }
    
    if (level >= 5) {
        spells.push({
            name: `${wizard.element} Blast`,
            school: element,
            level: 3
        });
    }
    
    if (level >= 7) {
        spells.push({
            name: 'Mana Restore',
            school: 'arcane',
            level: 2
        });
    }
    
    if (level >= 10) {
        spells.push({
            name: `Greater ${wizard.element} Mastery`,
            school: element,
            level: 5
        });
    }
    
    if (level >= 12) {
        spells.push({
            name: 'Teleport',
            school: 'arcane',
            level: 4
        });
    }
    
    if (level >= 15) {
        spells.push({
            name: `${wizard.element} Storm`,
            school: element,
            level: 6
        });
    }
    
    return spells;
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
    setupJobsEventListeners();
}

function showMarketplace() {
    showPage('marketplace');
    loadMarketplace();
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
        await displayWizards();
    } catch (error) {
        console.error('Error loading wizards:', error);
        showToast('Error loading wizards', 'error');
        wizards = [];
    } finally {
        showLoading(false);
    }
}

async function displayWizards() {
    const wizardsGrid = document.getElementById('wizards-grid');
    
    if (wizards.length === 0) {
        wizardsGrid.innerHTML = '<p style="text-align: center; color: white;">No wizards found. Create your first wizard!</p>';
        return;
    }
    
    // For demo purposes, add sample artifacts and spells based on wizard's element and level
    const wizardsWithCollections = wizards.map(wizard => {
        const sampleArtifacts = getSampleArtifacts(wizard);
        const sampleSpells = getSampleSpells(wizard);
        
        return {
            ...wizard,
            equippedArtifacts: sampleArtifacts,
            knownSpells: sampleSpells
        };
    });
    
    wizardsGrid.innerHTML = wizardsWithCollections.map(wizard => `
        <div class="wizard-card" onclick="openWizardModal(${wizard.id})">
            <div class="wizard-header">
                <div style="display: flex; align-items: flex-start; flex: 1;">
                    <div class="wizard-avatar wizard-element-${wizard.element.toLowerCase()}">
                        <i class="fas fa-hat-wizard"></i>
                    </div>
                    <div class="wizard-info-section">
                        <div class="wizard-name">${wizard.name}</div>
                        <div class="wizard-details">
                            <div class="wizard-detail">
                                <i class="fas fa-globe"></i> ${wizard.realm}
                            </div>
                            <div class="wizard-detail">
                                <i class="fas fa-magic"></i> ${wizard.element}
                            </div>
                        </div>
                    </div>
                </div>
                <div class="wizard-id">#${wizard.id}</div>
            </div>
            
            <div class="wizard-stats">
                <div class="wizard-stat wizard-level">
                    <div class="wizard-stat-label">Level</div>
                    <div class="wizard-stat-value">${wizard.level || 1}</div>
                </div>
                <div class="wizard-stat wizard-mana">
                    <div class="wizard-stat-label">Mana</div>
                    <div class="wizard-stat-value">${formatNumber(wizard.mana_balance || 0)}</div>
                </div>
                <div class="wizard-stat wizard-exp">
                    <div class="wizard-stat-label">Exp</div>
                    <div class="wizard-stat-value">${formatNumber(wizard.experience_points || 0)}</div>
                </div>
            </div>
            
            <div class="wizard-preview">
                ${wizard.equippedArtifacts.length > 0 || wizard.knownSpells.length > 0 ? `
                    <div class="wizard-quick-info">
                        ${wizard.equippedArtifacts.length > 0 ? `
                            <div class="quick-stat">
                                <i class="fas fa-gem"></i>
                                <span>${wizard.equippedArtifacts.length} artifact${wizard.equippedArtifacts.length !== 1 ? 's' : ''}</span>
                            </div>
                        ` : ''}
                        ${wizard.knownSpells.length > 0 ? `
                            <div class="quick-stat">
                                <i class="fas fa-scroll"></i>
                                <span>${wizard.knownSpells.length} spell${wizard.knownSpells.length !== 1 ? 's' : ''}</span>
                            </div>
                        ` : ''}
                    </div>
                ` : `
                    <div class="wizard-empty-state">
                        <i class="fas fa-store"></i>
                        <span>Visit Marketplace</span>
                    </div>
                `}
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
        let exploreWizards = response.wizards || [];
        
        // If no wizards from backend, create sample lore-friendly wizards
        if (exploreWizards.length === 0) {
            exploreWizards = generateSampleWizards(realm);
        }
        
        displayExploreWizards(exploreWizards);
    } catch (error) {
        console.error('Error loading explore wizards:', error);
        // Fallback to sample wizards if API fails
        const exploreWizards = generateSampleWizards(realm);
        displayExploreWizards(exploreWizards);
        showToast('Showing sample wizards - API unavailable', 'warning');
    } finally {
        showLoading(false);
    }
}

function displayExploreWizards(exploreWizards) {
    const exploreWizardsGrid = document.getElementById('explore-wizards-grid');
    
    if (exploreWizards.length === 0) {
        exploreWizardsGrid.innerHTML = '<p style="text-align: center; color: #666; padding: 40px;">No wizards found in this realm.</p>';
        return;
    }
    
    // Add artifact and spell collections to explore wizards (mock data for display)
    const exploreWizardsWithCollections = exploreWizards.map(wizard => ({
        ...wizard,
        equippedArtifacts: wizard.artifacts || [],
        knownSpells: wizard.spells || []
    }));
    
    exploreWizardsGrid.innerHTML = exploreWizardsWithCollections.map(wizard => `
        <div class="wizard-card" onclick="openWizardModal(${wizard.id})">
            <div class="wizard-header">
                <div style="display: flex; align-items: flex-start; flex: 1;">
                    <div class="wizard-avatar wizard-element-${wizard.element.toLowerCase()}">
                        <i class="fas fa-hat-wizard"></i>
                    </div>
                    <div class="wizard-info-section">
                        <div class="wizard-name">${wizard.name}</div>
                        <div class="wizard-details">
                            <div class="wizard-detail">
                                <i class="fas fa-globe"></i> ${wizard.realm}
                            </div>
                            <div class="wizard-detail">
                                <i class="fas fa-magic"></i> ${wizard.element}
                            </div>
                        </div>
                    </div>
                </div>
                <div class="wizard-id">#${wizard.id}</div>
            </div>
            
            <div class="wizard-stats">
                <div class="wizard-stat wizard-level">
                    <div class="wizard-stat-label">Level</div>
                    <div class="wizard-stat-value">${wizard.level || 1}</div>
                </div>
                <div class="wizard-stat wizard-mana">
                    <div class="wizard-stat-label">Mana</div>
                    <div class="wizard-stat-value">${formatNumber(wizard.mana_balance || 0)}</div>
                </div>
                <div class="wizard-stat wizard-exp">
                    <div class="wizard-stat-label">Exp</div>
                    <div class="wizard-stat-value">${formatNumber(wizard.experience_points || 0)}</div>
                </div>
            </div>
            
            <div class="wizard-preview">
                ${wizard.equippedArtifacts.length > 0 || wizard.knownSpells.length > 0 ? `
                    <div class="wizard-quick-info">
                        ${wizard.equippedArtifacts.length > 0 ? `
                            <div class="quick-stat">
                                <i class="fas fa-gem"></i>
                                <span>${wizard.equippedArtifacts.length} artifact${wizard.equippedArtifacts.length !== 1 ? 's' : ''}</span>
                            </div>
                        ` : ''}
                        ${wizard.knownSpells.length > 0 ? `
                            <div class="quick-stat">
                                <i class="fas fa-scroll"></i>
                                <span>${wizard.knownSpells.length} spell${wizard.knownSpells.length !== 1 ? 's' : ''}</span>
                            </div>
                        ` : ''}
                    </div>
                ` : `
                    <div class="wizard-empty-state">
                        <i class="fas fa-eye"></i>
                        <span>Click to view</span>
                    </div>
                `}
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

function generateSampleWizards(filterRealm = '') {
    const sampleWizards = [
        {
            id: 101,
            name: "Ignis Pyroclast",
            element: "Fire",
            realm: "Pyrrhian Flame",
            level: 15,
            mana_balance: 8250,
            experience_points: 12400,
            artifacts: [
                {name: "Salamander's Fury Gauntlets", type: "weapon", rarity: "epic"},
                {name: "Molten Core Pendant", type: "accessory", rarity: "rare"}
            ],
            spells: [
                {name: "Inferno Cascade", school: "fire", level: 4},
                {name: "Phoenix Rebirth", school: "fire", level: 3},
                {name: "Volcanic Eruption", school: "fire", level: 5}
            ]
        },
        {
            id: 102,
            name: "Luna Tidewhisper",
            element: "Water",
            realm: "Thalorion Depths",
            level: 12,
            mana_balance: 6800,
            experience_points: 8900,
            artifacts: [
                {name: "Moonbound Trident", type: "weapon", rarity: "rare"},
                {name: "Depths Walker Robes", type: "armor", rarity: "uncommon"}
            ],
            spells: [
                {name: "Tsunami Call", school: "water", level: 3},
                {name: "Healing Springs", school: "water", level: 2},
                {name: "Frost Prison", school: "water", level: 4}
            ]
        },
        {
            id: 103,
            name: "Vex Shadowbane",
            element: "Shadow",
            realm: "Umbros",
            level: 18,
            mana_balance: 15500,
            experience_points: 24600,
            artifacts: [
                {name: "Void Ripper Daggers", type: "weapon", rarity: "legendary"},
                {name: "Umbral Cloak", type: "armor", rarity: "epic"},
                {name: "Shadow Walker Boots", type: "armor", rarity: "rare"}
            ],
            spells: [
                {name: "Shadow Step", school: "shadow", level: 4},
                {name: "Darkness Veil", school: "shadow", level: 3},
                {name: "Soul Drain", school: "shadow", level: 5},
                {name: "Umbral Blast", school: "shadow", level: 2}
            ]
        },
        {
            id: 104,
            name: "Zephyr Stormcaller",
            element: "Air",
            realm: "Zepharion Heights",
            level: 14,
            mana_balance: 9200,
            experience_points: 15300,
            artifacts: [
                {name: "Cyclone Staff", type: "weapon", rarity: "epic"},
                {name: "Wind Walker Amulet", type: "accessory", rarity: "rare"}
            ],
            spells: [
                {name: "Lightning Storm", school: "air", level: 4},
                {name: "Wind Shield", school: "air", level: 2},
                {name: "Thunder Clap", school: "air", level: 3}
            ]
        },
        {
            id: 105,
            name: "Terra Stoneforge",
            element: "Earth",
            realm: "Terravine Hollow",
            level: 16,
            mana_balance: 7400,
            experience_points: 18700,
            artifacts: [
                {name: "World Tree Staff", type: "weapon", rarity: "epic"},
                {name: "Living Stone Armor", type: "armor", rarity: "rare"},
                {name: "Root Network Ring", type: "accessory", rarity: "uncommon"}
            ],
            spells: [
                {name: "Earthquake", school: "earth", level: 5},
                {name: "Nature's Blessing", school: "earth", level: 3},
                {name: "Stone Skin", school: "earth", level: 2}
            ]
        },
        {
            id: 106,
            name: "Solaris Dawnbringer",
            element: "Light",
            realm: "Virelya",
            level: 13,
            mana_balance: 11600,
            experience_points: 10800,
            artifacts: [
                {name: "Radiant Blade", type: "weapon", rarity: "rare"},
                {name: "Prism Robes", type: "armor", rarity: "epic"}
            ],
            spells: [
                {name: "Solar Flare", school: "light", level: 4},
                {name: "Healing Light", school: "light", level: 3},
                {name: "Blinding Flash", school: "light", level: 2}
            ]
        },
        {
            id: 107,
            name: "Chronos Timekeeper",
            element: "Time",
            realm: "Chronarxis",
            level: 20,
            mana_balance: 18900,
            experience_points: 32100,
            artifacts: [
                {name: "Temporal Orb", type: "weapon", rarity: "legendary"},
                {name: "Time Walker Robes", type: "armor", rarity: "epic"},
                {name: "Chronometer Pendant", type: "accessory", rarity: "rare"}
            ],
            spells: [
                {name: "Time Freeze", school: "time", level: 5},
                {name: "Temporal Shift", school: "time", level: 4},
                {name: "Age Reversal", school: "time", level: 6}
            ]
        },
        {
            id: 108,
            name: "Null the Voidwalker",
            element: "Void",
            realm: "Nyxthar",
            level: 17,
            mana_balance: 13200,
            experience_points: 21800,
            artifacts: [
                {name: "Entropy Scythe", type: "weapon", rarity: "legendary"},
                {name: "Void Cloak", type: "armor", rarity: "epic"}
            ],
            spells: [
                {name: "Reality Tear", school: "void", level: 5},
                {name: "Existence Drain", school: "void", level: 4},
                {name: "Null Zone", school: "void", level: 3}
            ]
        },
        {
            id: 109,
            name: "Ethereal Soulweaver",
            element: "Spirit",
            realm: "Aetherion",
            level: 11,
            mana_balance: 9800,
            experience_points: 7200,
            artifacts: [
                {name: "Soul Prism", type: "weapon", rarity: "rare"},
                {name: "Ethereal Bindings", type: "armor", rarity: "uncommon"}
            ],
            spells: [
                {name: "Spirit Walk", school: "spirit", level: 3},
                {name: "Soul Heal", school: "spirit", level: 2},
                {name: "Astral Projection", school: "spirit", level: 4}
            ]
        },
        {
            id: 110,
            name: "Ferros Gearwright",
            element: "Metal",
            realm: "Technarok",
            level: 19,
            mana_balance: 16700,
            experience_points: 28400,
            artifacts: [
                {name: "Mechanical Arm", type: "weapon", rarity: "epic"},
                {name: "Steel Plate Armor", type: "armor", rarity: "rare"},
                {name: "Gear Heart", type: "accessory", rarity: "legendary"}
            ],
            spells: [
                {name: "Metal Storm", school: "metal", level: 4},
                {name: "Construct Summon", school: "metal", level: 5},
                {name: "Magnetic Pull", school: "metal", level: 3}
            ]
        }
    ];
    
    // Filter by realm if specified
    if (filterRealm) {
        return sampleWizards.filter(wizard => wizard.realm === filterRealm);
    }
    
    return sampleWizards;
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
        
        // Store data for filtering
        currentJobsData = [...jobs];
        
        // Apply current filters
        applyJobFilters();
    } catch (error) {
        console.error('Error loading jobs:', error);
        showToast('Error loading jobs', 'error');
        jobs = [];
        currentJobsData = [];
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
    
    // Use jobs as fallback if filteredJobs is empty
    const jobsToShow = filteredJobs.length > 0 ? filteredJobs : jobs;
    
    if (jobsToShow.length === 0) {
        jobsGrid.innerHTML = `
            <div class="empty-state">
                <i class="fas fa-briefcase"></i>
                <h3>No Jobs Available</h3>
                <p>Check back later for new opportunities across the realms!</p>
            </div>
        `;
        return;
    }
    
    jobsGrid.innerHTML = jobsToShow.map(job => `
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

// Global jobs data storage
let currentJobsData = [];

function setupJobsEventListeners() {
    // Search input
    const searchInput = document.getElementById('jobs-search');
    if (searchInput) {
        searchInput.addEventListener('input', applyJobFilters);
    }
    
    // Sort dropdown
    const sortSelect = document.getElementById('jobs-sort');
    if (sortSelect) {
        sortSelect.addEventListener('change', applyJobFilters);
    }
    
    
    // Element dropdown
    const elementSelect = document.getElementById('job-element-filter');
    if (elementSelect) {
        elementSelect.addEventListener('change', applyJobFilters);
    }
    
    // Difficulty checkboxes
    const difficultyCheckboxes = document.querySelectorAll('.jobs-filters input[type="checkbox"]');
    difficultyCheckboxes.forEach(checkbox => {
        if (['Easy', 'Medium', 'Hard', 'Expert', 'Legendary'].includes(checkbox.value)) {
            checkbox.addEventListener('change', applyJobFilters);
        }
    });
    
    // Realm dropdown
    const realmSelect = document.getElementById('job-realm-filter');
    if (realmSelect) {
        realmSelect.addEventListener('change', applyJobFilters);
    }
}

function applyJobFilters() {
    try {
        // Use currentJobsData or fallback to jobs
        const dataToFilter = currentJobsData.length > 0 ? currentJobsData : jobs;
        
        // If no data at all, just show empty
        if (dataToFilter.length === 0) {
            filteredJobs = [];
            displayJobs();
            return;
        }
        
        // Get filter values with safety checks
        const searchInput = document.getElementById('jobs-search');
        const searchTerm = searchInput ? searchInput.value.toLowerCase() : '';
        
        const elementSelect = document.getElementById('job-element-filter');
        const selectedElement = elementSelect ? elementSelect.value : '';
        
        const selectedDifficulties = Array.from(document.querySelectorAll('.jobs-filters input[type="checkbox"]:checked'))
            .filter(cb => ['Easy', 'Medium', 'Hard', 'Expert', 'Legendary'].includes(cb.value))
            .map(cb => cb.value);
        
        const sortSelect = document.getElementById('jobs-sort');
        const sortBy = sortSelect ? sortSelect.value : 'difficulty';
        
        const realmSelect = document.getElementById('job-realm-filter');
        const realmFilter = realmSelect ? realmSelect.value : '';
        
        // Filter the data
        let filteredData = dataToFilter.filter(job => {
            // Search filter
            const matchesSearch = !searchTerm || 
                (job.title && job.title.toLowerCase().includes(searchTerm)) ||
                (job.description && job.description.toLowerCase().includes(searchTerm)) ||
                (job.realm_name && job.realm_name.toLowerCase().includes(searchTerm));
            
            // Element filter
            const matchesElement = !selectedElement || job.required_element === selectedElement;
            
            // Difficulty filter
            const matchesDifficulty = selectedDifficulties.length === 0 || 
                selectedDifficulties.includes(job.difficulty);
            
            
            // Realm filter
            const matchesRealm = !realmFilter || job.realm_name === realmFilter;
            
            return matchesSearch && matchesElement && matchesDifficulty && matchesRealm;
        });
        
        // Sort the data
        filteredData.sort((a, b) => {
            switch(sortBy) {
                case 'difficulty':
                    const difficultyOrder = { 'Easy': 1, 'Medium': 2, 'Hard': 3, 'Expert': 4, 'Legendary': 5 };
                    return (difficultyOrder[b.difficulty] || 0) - (difficultyOrder[a.difficulty] || 0);
                case 'duration':
                    return (a.duration_hours || 0) - (b.duration_hours || 0);
                case 'name':
                    return (a.title || '').localeCompare(b.title || '');
                default:
                    return 0;
            }
        });
        
        // Update filteredJobs and display
        filteredJobs = filteredData;
        displayJobs();
    } catch (error) {
        console.error('Error in applyJobFilters:', error);
        // Fallback: just show all jobs
        filteredJobs = jobs;
        displayJobs();
    }
}

function filterJobsByRealm() {
    applyJobFilters();
}

function filterJobsByElement() {
    applyJobFilters();
}

function filterJobsByDifficulty() {
    applyJobFilters();
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
        const result = await api.completeJobAssignment(assignmentId);
        
        // Show success message with reward details if available
        if (result && result.rewards) {
            const manaGained = result.rewards.mana || 0;
            const expGained = result.rewards.experience || 0;
            showToast(`Job completed! Earned ${formatNumber(manaGained)} mana and ${formatNumber(expGained)} experience.`, 'success');
        } else {
            showToast('Job completed successfully!', 'success');
        }
        
        // Refresh wizard data to show updated mana and experience
        await loadWizards();
        
    } catch (error) {
        console.error('Error completing job assignment:', error);
        showToast(error.message || 'Error completing job assignment', 'error');
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

// Marketplace Functions
async function loadMarketplace() {
    showMarketplaceSection('artifacts');
    
    // Add event listeners for marketplace controls
    setupMarketplaceEventListeners();
}

function setupMarketplaceEventListeners() {
    // Search input
    const searchInput = document.getElementById('marketplace-search');
    if (searchInput) {
        searchInput.addEventListener('input', applyFilters);
    }
    
    // Sort dropdown
    const sortSelect = document.getElementById('marketplace-sort');
    if (sortSelect) {
        sortSelect.addEventListener('change', applyFilters);
    }
    
    // Price range inputs
    const priceMin = document.getElementById('price-min');
    const priceMax = document.getElementById('price-max');
    if (priceMin) {
        priceMin.addEventListener('input', applyFilters);
    }
    if (priceMax) {
        priceMax.addEventListener('input', applyFilters);
    }
    
    // Element dropdown for spells
    const elementSelect = document.getElementById('marketplace-element-filter');
    if (elementSelect) {
        elementSelect.addEventListener('change', applyFilters);
    }
    
    // Rarity checkboxes for artifacts/scrolls
    const rarityCheckboxes = document.querySelectorAll('.filter-options input[type="checkbox"]');
    rarityCheckboxes.forEach(checkbox => {
        checkbox.addEventListener('change', applyFilters);
    });
}

function showMarketplaceSection(sectionName) {
    // Update navigation active state
    const navItems = document.querySelectorAll('.marketplace-nav-item');
    navItems.forEach(item => item.classList.remove('active'));
    
    // Find and activate the correct nav item
    const activeNavItem = Array.from(navItems).find(item => 
        item.onclick && item.onclick.toString().includes(`'${sectionName}'`)
    );
    if (activeNavItem) {
        activeNavItem.classList.add('active');
    }
    
    // Update section title
    const sectionTitle = document.getElementById('marketplace-section-title');
    if (sectionTitle) {
        sectionTitle.textContent = sectionName.charAt(0).toUpperCase() + sectionName.slice(1);
    }
    
    // Update sidebar filters based on section
    updateMarketplaceSidebar(sectionName);
    
    // Load data for the section
    switch(sectionName) {
        case 'artifacts':
            loadMarketplaceArtifacts();
            break;
        case 'scrolls':
            loadMarketplaceScrolls();
            break;
        case 'spells':
            loadMarketplaceSpells();
            break;
    }
}

function updateMarketplaceSidebar(sectionName) {
    const filterSection = document.querySelector('.marketplace-filters .filter-section');
    if (!filterSection) return;
    
    if (sectionName === 'spells') {
        // Show element dropdown for spells
        filterSection.innerHTML = `
            <label>Element</label>
            <select id="marketplace-element-filter">
                <option value="">All Elements</option>
                <option value="fire">Fire</option>
                <option value="water">Water</option>
                <option value="earth">Earth</option>
                <option value="air">Air</option>
                <option value="light">Light</option>
                <option value="dark">Dark</option>
                <option value="shadow">Shadow</option>
                <option value="spirit">Spirit</option>
                <option value="metal">Metal</option>
                <option value="time">Time</option>
                <option value="void">Void</option>
            </select>
        `;
    } else {
        // Show rarity filters for artifacts and scrolls
        filterSection.innerHTML = `
            <label>Rarity</label>
            <div class="filter-options">
                <label class="filter-option">
                    <input type="checkbox" value="common">
                    <span class="rarity-common">Common</span>
                </label>
                <label class="filter-option">
                    <input type="checkbox" value="uncommon">
                    <span class="rarity-uncommon">Uncommon</span>
                </label>
                <label class="filter-option">
                    <input type="checkbox" value="rare">
                    <span class="rarity-rare">Rare</span>
                </label>
                <label class="filter-option">
                    <input type="checkbox" value="epic">
                    <span class="rarity-epic">Epic</span>
                </label>
                <label class="filter-option">
                    <input type="checkbox" value="legendary">
                    <span class="rarity-legendary">Legendary</span>
                </label>
            </div>
        `;
    }
    
    // Re-attach event listeners
    setupMarketplaceEventListeners();
}

// Global marketplace data storage
let currentMarketplaceData = [];
let currentMarketplaceSection = 'artifacts';

function applyFilters() {
    // Get filter values
    const searchTerm = document.getElementById('marketplace-search').value.toLowerCase();
    const minPrice = parseInt(document.getElementById('price-min').value) || 0;
    const maxPrice = parseInt(document.getElementById('price-max').value) || Infinity;
    const sortBy = document.getElementById('marketplace-sort').value;
    
    // Get filter values based on current section
    let selectedFilter = '';
    if (currentMarketplaceSection === 'spells') {
        const elementSelect = document.getElementById('marketplace-element-filter');
        selectedFilter = elementSelect ? elementSelect.value : '';
    } else {
        const selectedRarities = Array.from(document.querySelectorAll('.filter-options input[type="checkbox"]:checked'))
            .map(cb => cb.value);
        selectedFilter = selectedRarities;
    }
    
    // Filter the data
    let filteredData = currentMarketplaceData.filter(item => {
        // Search filter
        const matchesSearch = !searchTerm || 
            item.name.toLowerCase().includes(searchTerm) ||
            item.description.toLowerCase().includes(searchTerm) ||
            (item.seller && item.seller.toLowerCase().includes(searchTerm));
        
        // Filter based on current section
        let matchesFilter = true;
        if (currentMarketplaceSection === 'spells') {
            // For spells, filter by school/element dropdown
            matchesFilter = !selectedFilter || 
                           item.school?.toLowerCase() === selectedFilter ||
                           item.element?.toLowerCase() === selectedFilter;
        } else {
            // For artifacts and scrolls, filter by rarity checkboxes
            matchesFilter = selectedFilter.length === 0 || selectedFilter.includes(item.rarity);
        }
        
        // Price filter
        const matchesPrice = item.price >= minPrice && item.price <= maxPrice;
        
        return matchesSearch && matchesFilter && matchesPrice;
    });
    
    // Sort the data
    filteredData.sort((a, b) => {
        switch(sortBy) {
            case 'name':
                return a.name.localeCompare(b.name);
            case 'price-low':
                return a.price - b.price;
            case 'price-high':
                return b.price - a.price;
            case 'rarity':
                const rarityOrder = { 'common': 1, 'uncommon': 2, 'rare': 3, 'epic': 4, 'legendary': 5 };
                return (rarityOrder[b.rarity] || 0) - (rarityOrder[a.rarity] || 0);
            default:
                return 0;
        }
    });
    
    // Display filtered data
    displayMarketplaceItems(filteredData, currentMarketplaceSection);
}

function displayMarketplaceItems(items, section) {
    const grid = document.getElementById('marketplace-items-grid');
    
    if (items.length === 0) {
        grid.innerHTML = '<div style="grid-column: 1 / -1; text-align: center; color: #666; padding: 40px;">No items match your filters.</div>';
        return;
    }
    
    if (section === 'artifacts') {
        grid.innerHTML = items.map(artifact => createArtifactHTML(artifact)).join('');
    } else if (section === 'scrolls') {
        grid.innerHTML = items.map(scroll => createScrollHTML(scroll)).join('');
    } else if (section === 'spells') {
        grid.innerHTML = items.map(spell => createSpellHTML(spell)).join('');
    }
}

// HTML generation functions
function createArtifactHTML(artifact) {
    return `
        <div class="marketplace-item">
            <div class="marketplace-item-header">
                <div class="marketplace-item-icon artifact-${artifact.rarity}">
                    <i class="fas ${getArtifactIcon(artifact.type)}"></i>
                </div>
                <div class="marketplace-item-info">
                    <div class="marketplace-item-name">${artifact.name}</div>
                    <div class="marketplace-item-meta">
                        <span class="marketplace-rarity rarity-${artifact.rarity}">${artifact.rarity}</span>
                        <span class="marketplace-type">${artifact.type}</span>
                    </div>
                </div>
            </div>
            
            <div class="marketplace-item-description">
                ${artifact.description}
            </div>
            
            <div class="marketplace-item-stats">
                ${Object.entries(artifact.stats).map(([stat, value]) => 
                    `<div class="marketplace-stat">
                        <span class="marketplace-stat-label">${stat.replace('_', ' ')}</span>
                        <span class="marketplace-stat-value">${value}</span>
                    </div>`
                ).join('')}
                ${artifact.passive_effect ? `
                    <div class="marketplace-stat passive-effect">
                        <span class="marketplace-stat-label">Passive Effect</span>
                        <span class="marketplace-stat-value">${artifact.mana_per_hour} mana/hr</span>
                    </div>
                ` : ''}
            </div>
            
            ${artifact.passive_effect ? `
                <div class="marketplace-passive-description">
                    <i class="fas fa-magic"></i>
                    <span>${artifact.passive_effect}</span>
                </div>
            ` : ''}
            
            <div class="marketplace-item-details">
                <div class="marketplace-detail">
                    <i class="fas fa-globe"></i>
                    <span>${artifact.realm}</span>
                </div>
                <div class="marketplace-detail">
                    <i class="fas fa-user"></i>
                    <span>${artifact.seller}</span>
                </div>
            </div>
            
            <div class="marketplace-item-footer">
                <div class="marketplace-price">
                    <i class="fas fa-coins"></i>
                    <span>${formatNumber(artifact.price)}</span>
                </div>
                <button class="btn btn-primary marketplace-buy-btn" onclick="purchaseArtifact(${artifact.id})">
                    <i class="fas fa-shopping-cart"></i>
                    Buy
                </button>
            </div>
        </div>
    `;
}

function createScrollHTML(scroll) {
    return `
        <div class="marketplace-item">
            <div class="marketplace-item-header">
                <div class="marketplace-item-icon scroll-item">
                    <i class="fas fa-scroll"></i>
                </div>
                <div class="marketplace-item-info">
                    <div class="marketplace-item-name">${scroll.name}</div>
                    <div class="marketplace-item-meta">
                        <span class="marketplace-type">Enhancement Scroll</span>
                    </div>
                </div>
            </div>
            
            <div class="marketplace-item-description">
                ${scroll.description}
            </div>
            
            <div class="marketplace-item-stats">
                <div class="marketplace-stat skill-bonus">
                    <span class="marketplace-stat-label">${scroll.skill.replace('_', ' ')} Bonus</span>
                    <span class="marketplace-stat-value">+${scroll.bonus}%</span>
                </div>
            </div>
            
            <div class="marketplace-item-details">
                <div class="marketplace-detail">
                    <i class="fas fa-user"></i>
                    <span>${scroll.seller}</span>
                </div>
            </div>
            
            <div class="marketplace-item-footer">
                <div class="marketplace-price">
                    <i class="fas fa-coins"></i>
                    <span>${formatNumber(scroll.price)}</span>
                </div>
                <button class="btn btn-primary marketplace-buy-btn" onclick="purchaseScroll(${scroll.id})">
                    <i class="fas fa-shopping-cart"></i>
                    Buy
                </button>
            </div>
        </div>
    `;
}

function createSpellHTML(spell) {
    return `
        <div class="marketplace-item">
            <div class="marketplace-item-header">
                <div class="marketplace-item-icon spell-${spell.school}">
                    <i class="fas ${getSpellIcon(spell.school)}"></i>
                </div>
                <div class="marketplace-item-info">
                    <div class="marketplace-item-name">${spell.name}</div>
                    <div class="marketplace-item-meta">
                        <span class="marketplace-type">${spell.school} Magic</span>
                        <span class="marketplace-level">Level ${spell.level}</span>
                    </div>
                </div>
            </div>
            
            <div class="marketplace-item-description">
                ${spell.description}
            </div>
            
            <div class="marketplace-item-details">
                <div class="marketplace-detail">
                    <i class="fas fa-magic"></i>
                    <span>${spell.school} School</span>
                </div>
                <div class="marketplace-detail">
                    <i class="fas fa-user-graduate"></i>
                    <span>${spell.teacher}</span>
                </div>
            </div>
            
            <div class="marketplace-item-footer">
                <div class="marketplace-price">
                    <i class="fas fa-coins"></i>
                    <span>${formatNumber(spell.price)}</span>
                </div>
                <button class="btn btn-primary marketplace-buy-btn" onclick="purchaseSpell(${spell.id})">
                    <i class="fas fa-shopping-cart"></i>
                    Buy
                </button>
            </div>
        </div>
    `;
}

async function loadMarketplaceArtifacts() {
    currentMarketplaceSection = 'artifacts';
    
    // Sample marketplace artifacts with passive effects
    const sampleArtifacts = [
        {
            id: 1,
            name: "Pyrrhian Ember Blade",
            type: "weapon",
            rarity: "rare",
            realm: "Pyrrhian Flame",
            price: 2500,
            seller: "Salamandrine Forgemaster Ignis",
            stats: { power: 85, fire_damage: 120 },
            passive_effect: "Generates 15 mana/hour from absorbed heat energy",
            mana_per_hour: 15,
            description: "A blade forged in the eternal flames of Pyrrhia, it burns with inner fire and slowly accumulates magical energy from heat."
        },
        {
            id: 2,
            name: "Virelya Light Robes",
            type: "armor",
            rarity: "epic",
            realm: "Virelya",
            price: 4200,
            seller: "Radiant Weaver Lumina",
            stats: { defense: 75, light_power: 140, mana_regeneration: 25 },
            passive_effect: "Absorbs ambient light to generate 35 mana/hour",
            mana_per_hour: 35,
            description: "Robes woven from crystallized light itself, they shimmer with pure radiance and slowly channel photonic energy into mana."
        },
        {
            id: 3,
            name: "Umbros Shadow Ring",
            type: "accessory",
            rarity: "legendary",
            realm: "Umbros",
            price: 8500,
            seller: "Voidwhisper Nyx",
            stats: { shadow_power: 200, stealth: 95, void_resistance: 85 },
            passive_effect: "Feeds on darkness to generate 50 mana/hour",
            mana_per_hour: 50,
            description: "A ring carved from crystallized void, it hungers for shadows and converts dark energy into raw magical power."
        },
        {
            id: 4,
            name: "Zepharion Storm Conductor",
            type: "weapon",
            rarity: "epic",
            realm: "Zepharion Heights",
            price: 5300,
            seller: "Cyclone Sage Tempest",
            stats: { power: 120, air_mastery: 110, storm_calling: 95 },
            passive_effect: "Harvests atmospheric energy for 28 mana/hour",
            mana_per_hour: 28,
            description: "A staff that resonates with the eternal cyclone, it draws power from wind currents and atmospheric disturbances."
        },
        {
            id: 5,
            name: "Terravine Growth Amulet",
            type: "accessory",
            rarity: "uncommon",
            realm: "Terravine Hollow",
            price: 1800,
            seller: "Stone Root Thornwick",
            stats: { earth_power: 60, nature_bond: 75, growth_speed: 40 },
            passive_effect: "Channels life force to generate 12 mana/hour",
            mana_per_hour: 12,
            description: "An amulet grown from petrified world-tree roots, it slowly converts life energy from the earth into magical essence."
        },
        {
            id: 6,
            name: "Thalorion Depth Crown",
            type: "armor",
            rarity: "legendary",
            realm: "Thalorion Depths",
            price: 9200,
            seller: "Moonbound Archon Tidal",
            stats: { water_mastery: 180, wisdom: 120, time_dilation: 60 },
            passive_effect: "Draws from ocean's timeless wisdom for 45 mana/hour",
            mana_per_hour: 45,
            description: "A crown from the drowned courts, it channels the eternal patience of deep waters into continuous magical energy."
        }
    ];
    
    // Store data for filtering
    currentMarketplaceData = sampleArtifacts;
    
    // Apply current filters
    applyFilters();
}

async function loadMarketplaceScrolls() {
    currentMarketplaceSection = 'scrolls';
    
    const sampleScrolls = [
        {
            id: 1,
            name: "Scroll of Mana Efficiency",
            skill: "mana_efficiency",
            bonus: 15,
            price: 1500,
            seller: "Scholar Aethen",
            description: "Reduces mana consumption by 15% for all spells"
        },
        {
            id: 2,
            name: "Scroll of Spell Power",
            skill: "spell_power",
            bonus: 20,
            price: 2200,
            seller: "Arcane Master Zara",
            description: "Increases the power of all spells by 20%"
        },
        {
            id: 3,
            name: "Scroll of Artifact Mastery",
            skill: "artifact_mastery",
            bonus: 25,
            price: 3000,
            seller: "Artificer Gorim",
            description: "Improves effectiveness of equipped artifacts by 25%"
        }
    ];
    
    // Store data for filtering (add rarity field for consistency)
    // Store data for filtering (add rarity field for consistency)
    currentMarketplaceData = sampleScrolls.map(scroll => ({
        ...scroll,
        rarity: 'common' // Default rarity for scrolls
    }));
    
    // Apply current filters
    applyFilters();
}

async function loadMarketplaceSpells() {
    currentMarketplaceSection = 'spells';
    
    const sampleSpells = [
        {
            id: 1,
            name: "Greater Fireball",
            school: "fire",
            level: 4,
            teacher: "Pyromancer Blaze",
            price: 3500,
            description: "A devastating fire spell that deals massive damage"
        },
        {
            id: 2,
            name: "Healing Light",
            school: "light",
            level: 3,
            teacher: "Cleric Lumina",
            price: 2800,
            description: "Restores health and removes negative effects"
        },
        {
            id: 3,
            name: "Shadow Step",
            school: "dark",
            level: 2,
            teacher: "Shadow Mage Nyx",
            price: 2000,
            description: "Allows instant teleportation through shadows"
        },
        {
            id: 4,
            name: "Lightning Storm",
            school: "air",
            level: 5,
            teacher: "Storm Caller Zephyr",
            price: 4500,
            description: "Summons a powerful lightning storm"
        }
    ];
    
    // Store data for filtering (add rarity field for consistency)
    // Store data for filtering (add rarity field for consistency)
    currentMarketplaceData = sampleSpells.map(spell => ({
        ...spell,
        rarity: 'uncommon' // Default rarity for spells
    }));
    
    // Apply current filters
    applyFilters();
}

async function loadMyInventory() {
    showToast('Inventory feature coming soon!', 'info');
}

// Purchase functions
async function purchaseArtifact(artifactId) {
    if (wizards.length === 0) {
        showToast('Create a wizard first!', 'error');
        return;
    }
    
    // Simple wizard selection for demo
    const wizardId = wizards[0].id;
    showToast(`Artifact ${artifactId} purchased for wizard ${wizardId}!`, 'success');
}

async function purchaseScroll(scrollId) {
    if (wizards.length === 0) {
        showToast('Create a wizard first!', 'error');
        return;
    }
    
    showToast(`Scroll ${scrollId} purchased!`, 'success');
}

async function purchaseSpell(spellId) {
    if (wizards.length === 0) {
        showToast('Create a wizard first!', 'error');
        return;
    }
    
    const wizardId = wizards[0].id;
    showToast(`Spell ${spellId} purchased for wizard ${wizardId}!`, 'success');
}

async function learnSpell(spellId) {
    if (wizards.length === 0) {
        showToast('Create a wizard first!', 'error');
        return;
    }
    
    const wizardId = wizards[0].id;
    showToast(`Spell ${spellId} learned by wizard ${wizardId}!`, 'success');
}

// Filter functions
function filterArtifacts() {
    showToast('Artifact filtering coming soon!', 'info');
}

function filterScrolls() {
    showToast('Scroll filtering coming soon!', 'info');
}

function filterSpells() {
    showToast('Spell filtering coming soon!', 'info');
}

// Wizard Detail Modal Functions
let currentWizardId = null;

function openWizardModal(wizardId) {
    currentWizardId = wizardId;
    const wizard = wizards.find(w => w.id === wizardId);
    
    if (!wizard) {
        showToast('Wizard not found!', 'error');
        return;
    }
    
    // Update modal content
    document.getElementById('modal-wizard-name').textContent = wizard.name;
    document.getElementById('modal-wizard-realm').innerHTML = `<i class="fas fa-globe"></i> ${wizard.realm}`;
    document.getElementById('modal-wizard-element').innerHTML = `<i class="fas fa-magic"></i> ${wizard.element}`;
    document.getElementById('modal-wizard-level').innerHTML = `<i class="fas fa-star"></i> Level ${wizard.level || 1}`;
    document.getElementById('modal-wizard-mana').textContent = formatNumber(wizard.mana_balance || 0);
    document.getElementById('modal-wizard-exp').textContent = formatNumber(wizard.experience_points || 0);
    
    // Load wizard's collections
    loadWizardInventory(wizardId);
    
    // Show modal
    document.getElementById('wizard-detail-modal').classList.add('show');
}

function closeWizardModal() {
    document.getElementById('wizard-detail-modal').classList.remove('show');
    currentWizardId = null;
}

async function loadWizardInventory(wizardId) {
    try {
        // Get sample artifacts and spells for this wizard
        const wizard = wizards.find(w => w.id === wizardId);
        const sampleArtifacts = getSampleArtifacts(wizard);
        const sampleSpells = getSampleSpells(wizard);
        
        // Update total items count
        document.getElementById('modal-wizard-items').textContent = sampleArtifacts.length + sampleSpells.length;
        
        // Render artifacts
        const artifactsGrid = document.getElementById('modal-wizard-artifacts');
        if (sampleArtifacts.length > 0) {
            artifactsGrid.innerHTML = sampleArtifacts.map(artifact => `
                <div class="inventory-item artifact-${artifact.rarity.toLowerCase()}">
                    <div class="inventory-item-icon" style="background: ${getArtifactRarityColor(artifact.rarity)}">
                        <i class="fas ${getArtifactIcon(artifact.type)}"></i>
                    </div>
                    <div class="inventory-item-info">
                        <div class="inventory-item-name">${artifact.name}</div>
                        <div class="inventory-item-type">${artifact.type} • ${artifact.rarity}</div>
                    </div>
                </div>
            `).join('');
        } else {
            artifactsGrid.innerHTML = `
                <div class="inventory-empty">
                    <i class="fas fa-gem"></i>
                    <p>No artifacts equipped</p>
                </div>
            `;
        }
        
        // Render spells
        const spellsGrid = document.getElementById('modal-wizard-spells');
        if (sampleSpells.length > 0) {
            spellsGrid.innerHTML = sampleSpells.map(spell => `
                <div class="inventory-item spell-${spell.school.toLowerCase()}">
                    <div class="inventory-item-icon" style="background: ${getSpellSchoolColor(spell.school)}">
                        <i class="fas ${getSpellIcon(spell.school)}"></i>
                    </div>
                    <div class="inventory-item-info">
                        <div class="inventory-item-name">${spell.name}</div>
                        <div class="inventory-item-type">${spell.school} • Level ${spell.level}</div>
                    </div>
                </div>
            `).join('');
        } else {
            spellsGrid.innerHTML = `
                <div class="inventory-empty">
                    <i class="fas fa-scroll"></i>
                    <p>No spells learned</p>
                </div>
            `;
        }
        
    } catch (error) {
        console.error('Error loading wizard inventory:', error);
        showToast('Error loading wizard inventory', 'error');
    }
}

function getArtifactRarityColor(rarity) {
    switch (rarity.toLowerCase()) {
        case 'common': return 'linear-gradient(135deg, #999 0%, #777 100%)';
        case 'uncommon': return 'linear-gradient(135deg, #32cd32 0%, #228b22 100%)';
        case 'rare': return 'linear-gradient(135deg, #4169e1 0%, #1e90ff 100%)';
        case 'epic': return 'linear-gradient(135deg, #9370db 0%, #8a2be2 100%)';
        case 'legendary': return 'linear-gradient(135deg, #ffa500 0%, #ff8c00 100%)';
        case 'mythic': return 'linear-gradient(135deg, #ff1493 0%, #dc143c 100%)';
        default: return 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)';
    }
}

function getSpellSchoolColor(school) {
    switch (school.toLowerCase()) {
        case 'fire': return 'linear-gradient(135deg, #ff4500 0%, #ff6347 100%)';
        case 'water': return 'linear-gradient(135deg, #00bfff 0%, #1e90ff 100%)';
        case 'earth': return 'linear-gradient(135deg, #8b4513 0%, #a0522d 100%)';
        case 'air': return 'linear-gradient(135deg, #4682b4 0%, #87ceeb 100%)';
        case 'light': return 'linear-gradient(135deg, #ffd700 0%, #b8860b 100%)';
        case 'dark': return 'linear-gradient(135deg, #4b0082 0%, #8a2be2 100%)';
        case 'nature': return 'linear-gradient(135deg, #32cd32 0%, #228b22 100%)';
        case 'arcane': return 'linear-gradient(135deg, #9370db 0%, #8a2be2 100%)';
        default: return 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)';
    }
}

// Simple page initialization
document.addEventListener('DOMContentLoaded', function() {
    console.log('MysticFunds app loaded');
    startActiveJobsRefresh();
});