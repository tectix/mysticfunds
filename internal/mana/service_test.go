package mana

import (
	"context"
	"database/sql"
	"testing"

	"github.com/tectix/mysticfunds/pkg/config"
	"github.com/tectix/mysticfunds/pkg/logger"
	pb "github.com/tectix/mysticfunds/proto/mana"
	wizardpb "github.com/tectix/mysticfunds/proto/wizard"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type testSetup struct {
	db          *sql.DB
	mock        sqlmock.Sqlmock
	wizardMock  *MockWizardServiceClient
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
	
	// Create mock wizard client
	wizardMock := &MockWizardServiceClient{}
	
	// Create mock scheduler
	scheduler := NewInvestmentScheduler(db, log, wizardMock)

	// Create service with all dependencies
	service := &ManaServiceImpl{
		db:           db,
		cfg:          cfg,
		log:          log,
		scheduler:    scheduler,
		wizardClient: wizardMock,
	}

	ctx := context.Background()

	return &testSetup{
		db:          db,
		mock:        mock,
		wizardMock:  wizardMock,
		service:     service,
		ctx:         ctx,
		testWizard1: 1,
		testWizard2: 2,
	}
}

func TestTransferMana(t *testing.T) {
	setup := setupTest(t)
	defer setup.db.Close()

	// Mock wizard service transfer call
	setup.wizardMock.On("TransferMana", setup.ctx, &wizardpb.TransferManaRequest{
		FromWizardId: setup.testWizard1,
		ToWizardId:   setup.testWizard2,
		Amount:       100,
		Reason:       "Mana service transfer",
	}).Return(&wizardpb.TransferManaResponse{
		Success: true,
		Message: "Transfer successful",
	}, nil)

	// Mock database transaction insert
	setup.mock.ExpectExec("INSERT INTO mana_transactions").
		WithArgs(setup.testWizard1, setup.testWizard2, int64(100)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute transfer
	resp, err := setup.service.TransferMana(setup.ctx, &pb.TransferManaRequest{
		FromWizardId: setup.testWizard1,
		ToWizardId:   setup.testWizard2,
		Amount:       100,
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// Verify all expectations
	assert.NoError(t, setup.mock.ExpectationsWereMet())
	setup.wizardMock.AssertExpectations(t)
}

func TestTransferManaFailure(t *testing.T) {
	setup := setupTest(t)
	defer setup.db.Close()

	// Mock wizard service transfer call to fail
	setup.wizardMock.On("TransferMana", setup.ctx, &wizardpb.TransferManaRequest{
		FromWizardId: setup.testWizard1,
		ToWizardId:   setup.testWizard2,
		Amount:       100,
		Reason:       "Mana service transfer",
	}).Return(&wizardpb.TransferManaResponse{
		Success: false,
		Message: "Insufficient balance",
	}, nil)

	// Execute transfer
	resp, err := setup.service.TransferMana(setup.ctx, &pb.TransferManaRequest{
		FromWizardId: setup.testWizard1,
		ToWizardId:   setup.testWizard2,
		Amount:       100,
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Success)

	// Verify all expectations
	setup.wizardMock.AssertExpectations(t)
}

func TestGetManaBalance(t *testing.T) {
	setup := setupTest(t)
	defer setup.db.Close()

	expectedBalance := int64(500)

	// Mock wizard service balance call
	setup.wizardMock.On("GetManaBalance", setup.ctx, &wizardpb.GetManaBalanceRequest{
		WizardId: setup.testWizard1,
	}).Return(&wizardpb.GetManaBalanceResponse{
		Balance: expectedBalance,
	}, nil)

	// Execute get balance
	resp, err := setup.service.GetManaBalance(setup.ctx, &pb.GetManaBalanceRequest{
		WizardId: setup.testWizard1,
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedBalance, resp.Balance)

	// Verify all expectations
	setup.wizardMock.AssertExpectations(t)
}

func TestGetManaBalanceWizardNotFound(t *testing.T) {
	setup := setupTest(t)
	defer setup.db.Close()

	// Mock wizard service balance call to fail
	setup.wizardMock.On("GetManaBalance", setup.ctx, &wizardpb.GetManaBalanceRequest{
		WizardId: setup.testWizard1,
	}).Return((*wizardpb.GetManaBalanceResponse)(nil), status.Error(codes.NotFound, "Wizard not found"))

	// Execute get balance
	resp, err := setup.service.GetManaBalance(setup.ctx, &pb.GetManaBalanceRequest{
		WizardId: setup.testWizard1,
	})

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Wizard not found")

	// Verify all expectations
	setup.wizardMock.AssertExpectations(t)
}

func TestCreateInvestment(t *testing.T) {
	setup := setupTest(t)
	defer setup.db.Close()

	investmentTypeId := int64(1)
	amount := int64(200)
	expectedBalance := int64(1000)

	// Mock investment type query
	setup.mock.ExpectQuery("SELECT min_amount, max_amount, duration_hours FROM investment_types").
		WithArgs(investmentTypeId).
		WillReturnRows(sqlmock.NewRows([]string{"min_amount", "max_amount", "duration_hours"}).
			AddRow(100, 1000, 24))

	// Mock wizard service balance check
	setup.wizardMock.On("GetManaBalance", setup.ctx, &wizardpb.GetManaBalanceRequest{
		WizardId: setup.testWizard1,
	}).Return(&wizardpb.GetManaBalanceResponse{
		Balance: expectedBalance,
	}, nil)

	// Mock wizard service balance update
	setup.wizardMock.On("UpdateManaBalance", setup.ctx, &wizardpb.UpdateManaBalanceRequest{
		WizardId: setup.testWizard1,
		Amount:   -amount,
		Reason:   "Investment creation",
	}).Return(&wizardpb.UpdateManaBalanceResponse{
		NewBalance: expectedBalance - amount,
		Success:    true,
	}, nil)

	// Mock transaction begin/commit
	setup.mock.ExpectBegin()
	setup.mock.ExpectQuery("INSERT INTO wizard_investments").
		WithArgs(setup.testWizard1, investmentTypeId, amount, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	setup.mock.ExpectCommit()

	// Execute create investment
	resp, err := setup.service.CreateInvestment(setup.ctx, &pb.CreateInvestmentRequest{
		WizardId:         setup.testWizard1,
		InvestmentTypeId: investmentTypeId,
		Amount:           amount,
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.InvestmentId)
	assert.Greater(t, resp.EndTime, int64(0))

	// Verify all expectations
	assert.NoError(t, setup.mock.ExpectationsWereMet())
	setup.wizardMock.AssertExpectations(t)
}

func TestCreateInvestmentInsufficientBalance(t *testing.T) {
	setup := setupTest(t)
	defer setup.db.Close()

	investmentTypeId := int64(1)
	amount := int64(200)
	insufficientBalance := int64(100)

	// Mock investment type query
	setup.mock.ExpectQuery("SELECT min_amount, max_amount, duration_hours FROM investment_types").
		WithArgs(investmentTypeId).
		WillReturnRows(sqlmock.NewRows([]string{"min_amount", "max_amount", "duration_hours"}).
			AddRow(100, 1000, 24))

	// Mock transaction begin
	setup.mock.ExpectBegin()

	// Mock wizard service balance check with insufficient balance
	setup.wizardMock.On("GetManaBalance", setup.ctx, &wizardpb.GetManaBalanceRequest{
		WizardId: setup.testWizard1,
	}).Return(&wizardpb.GetManaBalanceResponse{
		Balance: insufficientBalance,
	}, nil)

	// Execute create investment
	resp, err := setup.service.CreateInvestment(setup.ctx, &pb.CreateInvestmentRequest{
		WizardId:         setup.testWizard1,
		InvestmentTypeId: investmentTypeId,
		Amount:           amount,
	})

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Insufficient mana balance")

	// Verify all expectations
	setup.wizardMock.AssertExpectations(t)
}