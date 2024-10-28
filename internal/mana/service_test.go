package mana

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Alinoureddine1/mysticfunds/pkg/config"
	"github.com/Alinoureddine1/mysticfunds/pkg/logger"
	pb "github.com/Alinoureddine1/mysticfunds/proto/mana"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type testSetup struct {
	db          *sql.DB
	mock        sqlmock.Sqlmock
	service     *ManaServiceImpl
	ctx         context.Context
	testWizard1 int64
	testWizard2 int64
}

func setupTest(t *testing.T) *testSetup {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}

	cfg := &config.Config{
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg.LogLevel)
	scheduler := NewInvestmentScheduler(db, log)

	service := NewManaServiceImpl(db, cfg, log, scheduler)
	ctx := context.Background()

	return &testSetup{
		db:          db,
		mock:        mock,
		service:     service,
		ctx:         ctx,
		testWizard1: 1,
		testWizard2: 2,
	}
}

func TestTransferMana(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(sqlmock.Sqlmock)
		request       *pb.TransferManaRequest
		expectedError codes.Code
	}{
		{
			name: "Successful transfer",
			setupMocks: func(mock sqlmock.Sqlmock) {
				// Begin transaction
				mock.ExpectBegin()

				// Check sender balance
				mock.ExpectQuery("SELECT mana_balance FROM wizards").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"mana_balance"}).AddRow(1000))

				// Update sender balance
				mock.ExpectExec("UPDATE wizards SET mana_balance").
					WithArgs(500, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Update receiver balance
				mock.ExpectExec("UPDATE wizards SET mana_balance").
					WithArgs(500, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Record transaction
				mock.ExpectExec("INSERT INTO mana_transactions").
					WithArgs(1, 2, 500).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Commit transaction
				mock.ExpectCommit()
			},
			request: &pb.TransferManaRequest{
				FromWizardId: 1,
				ToWizardId:   2,
				Amount:       500,
			},
			expectedError: codes.OK,
		},
		{
			name: "Insufficient balance",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT mana_balance FROM wizards").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"mana_balance"}).AddRow(100))
				mock.ExpectRollback()
			},
			request: &pb.TransferManaRequest{
				FromWizardId: 1,
				ToWizardId:   2,
				Amount:       500,
			},
			expectedError: codes.FailedPrecondition,
		},
		{
			name: "Sender not found",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT mana_balance FROM wizards").
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			request: &pb.TransferManaRequest{
				FromWizardId: 1,
				ToWizardId:   2,
				Amount:       500,
			},
			expectedError: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupTest(t)
			defer setup.db.Close()

			tt.setupMocks(setup.mock)

			_, err := setup.service.TransferMana(setup.ctx, tt.request)

			if tt.expectedError == codes.OK {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				status, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedError, status.Code())
			}

			assert.NoError(t, setup.mock.ExpectationsWereMet())
		})
	}
}

func TestCreateInvestment(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(sqlmock.Sqlmock)
		request       *pb.CreateInvestmentRequest
		expectedError codes.Code
	}{
		{
			name: "Successful investment creation",
			setupMocks: func(mock sqlmock.Sqlmock) {
				// Get investment type details
				mock.ExpectQuery("SELECT min_amount, max_amount, duration_hours FROM investment_types").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"min_amount", "max_amount", "duration_hours"}).
						AddRow(100, 1000, 24))

				mock.ExpectBegin()

				// Check wizard balance
				mock.ExpectQuery("SELECT mana_balance FROM wizards").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"mana_balance"}).AddRow(1000))

				// Update wizard balance
				mock.ExpectExec("UPDATE wizards SET mana_balance").
					WithArgs(500, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Create investment record
				mock.ExpectQuery("INSERT INTO wizard_investments").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()
			},
			request: &pb.CreateInvestmentRequest{
				WizardId:         1,
				InvestmentTypeId: 1,
				Amount:           500,
			},
			expectedError: codes.OK,
		},
		{
			name: "Investment amount below minimum",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT min_amount, max_amount, duration_hours FROM investment_types").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"min_amount", "max_amount", "duration_hours"}).
						AddRow(1000, 5000, 24))
			},
			request: &pb.CreateInvestmentRequest{
				WizardId:         1,
				InvestmentTypeId: 1,
				Amount:           500,
			},
			expectedError: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupTest(t)
			defer setup.db.Close()

			tt.setupMocks(setup.mock)

			_, err := setup.service.CreateInvestment(setup.ctx, tt.request)

			if tt.expectedError == codes.OK {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				status, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedError, status.Code())
			}

			assert.NoError(t, setup.mock.ExpectationsWereMet())
		})
	}
}

func TestGetInvestments(t *testing.T) {
	now := time.Now()
	future := now.Add(24 * time.Hour)

	tests := []struct {
		name          string
		setupMocks    func(sqlmock.Sqlmock)
		request       *pb.GetInvestmentsRequest
		expectedCount int
		expectedError codes.Code
	}{
		{
			name: "Successfully retrieve investments",
			setupMocks: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "amount", "start_time", "end_time", "status",
					"actual_return_rate", "returned_amount", "name", "risk_level",
				})
				rows.AddRow(
					1,             // id
					500,           // amount
					now.Unix(),    // start_time
					future.Unix(), // end_time
					"active",      // status
					sql.NullFloat64{Float64: 5.5, Valid: true}, // actual_return_rate
					sql.NullFloat64{Float64: 525, Valid: true}, // returned_amount
					"Test Investment",                          // name
					2,                                          // risk_level
				)
				rows.AddRow(
					2,                             // id
					1000,                          // amount
					now.Unix(),                    // start_time
					future.Unix(),                 // end_time
					"active",                      // status
					sql.NullFloat64{Valid: false}, // actual_return_rate
					sql.NullFloat64{Valid: false}, // returned_amount
					"Test Investment 2",           // name
					3,                             // risk_level
				)

				mock.ExpectQuery(`SELECT i.id, i.amount, i.start_time, i.end_time, i.status, 
					i.actual_return_rate, i.returned_amount, t.name, t.risk_level
					FROM wizard_investments i
					JOIN investment_types t ON i.investment_type_id = t.id
					WHERE i.wizard_id = \$1
					ORDER BY i.created_at DESC`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			request: &pb.GetInvestmentsRequest{
				WizardId: 1,
			},
			expectedCount: 2,
			expectedError: codes.OK,
		},
		{
			name: "No investments found",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT i.id, i.amount, i.start_time, i.end_time, i.status, 
					i.actual_return_rate, i.returned_amount, t.name, t.risk_level
					FROM wizard_investments i
					JOIN investment_types t ON i.investment_type_id = t.id
					WHERE i.wizard_id = \$1
					ORDER BY i.created_at DESC`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "amount", "start_time", "end_time", "status",
						"actual_return_rate", "returned_amount", "name", "risk_level",
					}))
			},
			request: &pb.GetInvestmentsRequest{
				WizardId: 1,
			},
			expectedCount: 0,
			expectedError: codes.OK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupTest(t)
			defer setup.db.Close()

			tt.setupMocks(setup.mock)

			response, err := setup.service.GetInvestments(setup.ctx, tt.request)

			if tt.expectedError == codes.OK {
				assert.NoError(t, err)
				if assert.NotNil(t, response) {
					assert.Equal(t, tt.expectedCount, len(response.Investments))

					if tt.expectedCount > 0 {
						// Verify first investment details
						inv := response.Investments[0]
						assert.Equal(t, int64(1), inv.Id)
						assert.Equal(t, int64(500), inv.Amount)
						assert.Equal(t, "active", inv.Status)
						assert.Equal(t, 5.5, inv.ActualReturnRate)
						assert.Equal(t, int64(525), inv.ReturnedAmount)
						assert.Equal(t, "Test Investment", inv.InvestmentType)
						assert.Equal(t, int32(2), inv.RiskLevel)
					}
				}
			} else {
				assert.Error(t, err)
				status, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedError, status.Code())
			}

			assert.NoError(t, setup.mock.ExpectationsWereMet())
		})
	}
}
