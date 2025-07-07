package mana

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tectix/mysticfunds/pkg/logger"
	wizardpb "github.com/tectix/mysticfunds/proto/wizard"
)

func TestInvestmentScheduler(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	log := logger.NewLogger("info")
	wizardMock := &MockWizardServiceClient{}
	scheduler := NewInvestmentScheduler(db, log, wizardMock)

	t.Run("Initialize scheduler", func(t *testing.T) {
		sqlMock.ExpectQuery("SELECT id, end_time FROM wizard_investments").
			WillReturnRows(sqlmock.NewRows([]string{"id", "end_time"}).
				AddRow(1, time.Now().Add(1*time.Hour)).
				AddRow(2, time.Now().Add(2*time.Hour)))

		scheduler.Start()
		assert.NoError(t, sqlMock.ExpectationsWereMet())
		scheduler.Stop()
	})

	t.Run("Process investment completion", func(t *testing.T) {
		sqlMock.ExpectBegin()

		sqlMock.ExpectQuery("SELECT (.+) FROM wizard_investments").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{
				"wizard_id", "amount", "base_return_rate", "risk_level",
			}).AddRow(1, 1000, 5.0, 2))

		sqlMock.ExpectExec("UPDATE wizard_investments").
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Mock wizard service call with flexible matching
		wizardMock.On("UpdateManaBalance",
			mock.Anything, // context
			mock.MatchedBy(func(req *wizardpb.UpdateManaBalanceRequest) bool {
				return req.WizardId == int64(1) && req.Reason == "Investment return"
			})).Return(&wizardpb.UpdateManaBalanceResponse{
			NewBalance: 1100,
			Success:    true,
		}, nil)

		sqlMock.ExpectCommit()

		scheduler.processInvestment(1)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
		wizardMock.AssertExpectations(t)
	})
}

func TestScheduleInvestmentCompletion(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	log := logger.NewLogger("info")
	wizardMock := &MockWizardServiceClient{}
	scheduler := NewInvestmentScheduler(db, log, wizardMock)

	t.Run("Schedule future investment", func(t *testing.T) {
		futureTime := time.Now().Add(1 * time.Hour)
		scheduler.ScheduleInvestmentCompletion(1, futureTime)

		assert.Contains(t, scheduler.active, int64(1))
	})

	t.Run("Schedule past investment", func(t *testing.T) {
		pastTime := time.Now().Add(-1 * time.Hour)

		sqlMock.ExpectBegin()
		sqlMock.ExpectQuery("SELECT (.+) FROM wizard_investments").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{
				"wizard_id", "amount", "base_return_rate", "risk_level",
			}).AddRow(1, 1000, 5.0, 2))
		sqlMock.ExpectExec("UPDATE wizard_investments").
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Mock wizard service call with flexible matching
		wizardMock.On("UpdateManaBalance",
			mock.Anything, // context
			mock.MatchedBy(func(req *wizardpb.UpdateManaBalanceRequest) bool {
				return req.WizardId == int64(1) && req.Reason == "Investment return"
			})).Return(&wizardpb.UpdateManaBalanceResponse{
			NewBalance: 1100,
			Success:    true,
		}, nil)

		sqlMock.ExpectCommit()

		scheduler.ScheduleInvestmentCompletion(1, pastTime)

		// Give some time for the goroutine to complete
		time.Sleep(100 * time.Millisecond)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
		wizardMock.AssertExpectations(t)
	})
}

func TestCalculateReturnRate(t *testing.T) {
	baseRate := 10.0
	riskLevel := int32(3)

	// Calculate return rate multiple times to test variance
	for i := 0; i < 10; i++ {
		actualRate := calculateReturnRate(baseRate, riskLevel)

		// Should be within reasonable bounds (baseRate ± variance)
		variance := float64(riskLevel) * 2.0
		assert.GreaterOrEqual(t, actualRate, baseRate-variance-1.0) // Allow small margin
		assert.LessOrEqual(t, actualRate, baseRate+variance+1.0)    // Allow small margin

		// Should never go below -90%
		assert.GreaterOrEqual(t, actualRate, -90.0)
	}
}
