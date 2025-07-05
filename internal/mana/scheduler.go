package mana

import (
	"context"
	"database/sql"
	"math/rand"
	"sync"
	"time"

	"github.com/tectix/mysticfunds/pkg/logger"
	wizardpb "github.com/tectix/mysticfunds/proto/wizard"
)

type InvestmentScheduler struct {
	db           *sql.DB
	log          logger.Logger
	wizardClient wizardpb.WizardServiceClient
	done         chan struct{}
	mutex        sync.Mutex
	active       map[int64]*time.Timer
}

func NewInvestmentScheduler(db *sql.DB, log logger.Logger, wizardClient wizardpb.WizardServiceClient) *InvestmentScheduler {
	return &InvestmentScheduler{
		db:           db,
		log:          log,
		wizardClient: wizardClient,
		done:         make(chan struct{}),
		active:       make(map[int64]*time.Timer),
	}
}

func (s *InvestmentScheduler) Start() {
	s.log.Info("Starting investment scheduler")

	// Load and schedule existing active investments
	rows, err := s.db.Query(`
		SELECT id, end_time 
		FROM wizard_investments 
		WHERE status = 'active' AND end_time > NOW()
	`)
	if err != nil {
		s.log.Error("Failed to load active investments", "error", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var endTime time.Time
		if err := rows.Scan(&id, &endTime); err != nil {
			s.log.Error("Failed to scan investment", "error", err)
			continue
		}
		s.ScheduleInvestmentCompletion(id, endTime)
	}

	// Start periodic cleanup of completed investments
	go s.cleanupRoutine()
}

func (s *InvestmentScheduler) Stop() {
	s.log.Info("Stopping investment scheduler")
	close(s.done)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Stop all active timers
	for _, timer := range s.active {
		timer.Stop()
	}
	s.active = make(map[int64]*time.Timer)
}

func (s *InvestmentScheduler) ScheduleInvestmentCompletion(investmentId int64, endTime time.Time) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Cancel existing timer if any
	if timer, exists := s.active[investmentId]; exists {
		timer.Stop()
	}

	duration := time.Until(endTime)
	if duration < 0 {
		// Investment already past due, process immediately
		go s.processInvestment(investmentId)
		return
	}

	timer := time.AfterFunc(duration, func() {
		s.processInvestment(investmentId)
	})
	s.active[investmentId] = timer
}

func (s *InvestmentScheduler) processInvestment(investmentId int64) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.log.Error("Failed to begin transaction", "error", err, "investmentId", investmentId)
		return
	}
	defer tx.Rollback()

	// Get investment details
	var investment struct {
		wizardId       int64
		amount         int64
		baseReturnRate float64
		riskLevel      int32
	}

	err = tx.QueryRowContext(ctx, `
		SELECT i.wizard_id, i.amount, t.base_return_rate, t.risk_level
		FROM wizard_investments i
		JOIN investment_types t ON i.investment_type_id = t.id
		WHERE i.id = $1 AND i.status = 'active'`,
		investmentId).Scan(
		&investment.wizardId,
		&investment.amount,
		&investment.baseReturnRate,
		&investment.riskLevel,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Info("Investment already processed", "investmentId", investmentId)
			return
		}
		s.log.Error("Failed to fetch investment", "error", err, "investmentId", investmentId)
		return
	}

	// Calculate return based on risk level and random variance
	actualReturnRate := calculateReturnRate(investment.baseReturnRate, investment.riskLevel)
	returnedAmount := int64(float64(investment.amount) * (1 + actualReturnRate/100))

	// Update investment status and return
	_, err = tx.ExecContext(ctx, `
		UPDATE wizard_investments 
		SET status = 'completed',
			actual_return_rate = $1,
			returned_amount = $2,
			updated_at = NOW()
		WHERE id = $3`,
		actualReturnRate, returnedAmount, investmentId)
	if err != nil {
		s.log.Error("Failed to update investment", "error", err, "investmentId", investmentId)
		return
	}

	// Credit returned amount to wizard via wizard service
	_, err = s.wizardClient.UpdateManaBalance(ctx, &wizardpb.UpdateManaBalanceRequest{
		WizardId: investment.wizardId,
		Amount:   returnedAmount,
		Reason:   "Investment return",
	})
	if err != nil {
		s.log.Error("Failed to credit return", "error", err, "investmentId", investmentId)
		return
	}

	if err = tx.Commit(); err != nil {
		s.log.Error("Failed to commit transaction", "error", err, "investmentId", investmentId)
		return
	}

	s.mutex.Lock()
	delete(s.active, investmentId)
	s.mutex.Unlock()

	s.log.Info("Investment completed successfully",
		"investmentId", investmentId,
		"returnRate", actualReturnRate,
		"returnedAmount", returnedAmount)
}

func (s *InvestmentScheduler) cleanupRoutine() {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanupExpiredInvestments()
		case <-s.done:
			return
		}
	}
}

func (s *InvestmentScheduler) cleanupExpiredInvestments() {
	ctx := context.Background()

	// Find expired but still active investments
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, end_time 
		FROM wizard_investments 
		WHERE status = 'active' AND end_time < NOW()
	`)
	if err != nil {
		s.log.Error("Failed to query expired investments", "error", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var endTime time.Time
		if err := rows.Scan(&id, &endTime); err != nil {
			s.log.Error("Failed to scan expired investment", "error", err)
			continue
		}
		// Process expired investment
		go s.processInvestment(id)
	}
}

func calculateReturnRate(baseRate float64, riskLevel int32) float64 {
	rand.Seed(time.Now().UnixNano())

	// Calculate variance based on risk level (higher risk = higher variance)
	variance := float64(riskLevel) * 2.0

	// Generate random adjustment within variance range
	adjustment := (rand.Float64()*2 - 1) * variance

	// Calculate actual return rate
	actualRate := baseRate + adjustment

	// Ensure return rate doesn't go below -90% (to prevent total loss)
	if actualRate < -90 {
		actualRate = -90
	}

	return actualRate
}
