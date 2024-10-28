package mana

import (
	"testing"
	"time"

	"github.com/Alinoureddine1/mysticfunds/pkg/logger"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInvestmentScheduler(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	log := logger.NewLogger("info")
	scheduler := NewInvestmentScheduler(db, log)

	t.Run("Initialize scheduler", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, end_time FROM wizard_investments").
			WillReturnRows(sqlmock.NewRows([]string{"id", "end_time"}).
				AddRow(1, time.Now().Add(1*time.Hour)).
				AddRow(2, time.Now().Add(2*time.Hour)))

		scheduler.Start()
		assert.NoError(t, mock.ExpectationsWereMet())
		scheduler.Stop()
	})

	t.Run("Process investment completion", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectQuery("SELECT (.+) FROM wizard_investments").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{
				"wizard_id", "amount", "base_return_rate", "risk_level",
			}).AddRow(1, 1000, 5.0, 2))

		mock.ExpectExec("UPDATE wizard_investments").
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("UPDATE wizards").
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		scheduler.processInvestment(1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestScheduleInvestmentCompletion(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	log := logger.NewLogger("info")
	scheduler := NewInvestmentScheduler(db, log)

	t.Run("Schedule future investment", func(t *testing.T) {
		futureTime := time.Now().Add(1 * time.Hour)
		scheduler.ScheduleInvestmentCompletion(1, futureTime)

		assert.Contains(t, scheduler.active, int64(1))
	})

	t.Run("Schedule past investment", func(t *testing.T) {
		pastTime := time.Now().Add(-1 * time.Hour)

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT (.+) FROM wizard_investments").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{
				"wizard_id", "amount", "base_return_rate", "risk_level",
			}).AddRow(1, 1000, 5.0, 2))
		mock.ExpectExec("UPDATE wizard_investments").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("UPDATE wizards").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		scheduler.ScheduleInvestmentCompletion(1, pastTime)

		// Give some time for the goroutine to complete
		time.Sleep(100 * time.Millisecond)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
