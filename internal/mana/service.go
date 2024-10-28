package mana

import (
	"context"
	"database/sql"
	"time"

	"github.com/Alinoureddine1/mysticfunds/pkg/config"
	"github.com/Alinoureddine1/mysticfunds/pkg/logger"
	pb "github.com/Alinoureddine1/mysticfunds/proto/mana"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ManaServiceImpl struct {
	db        *sql.DB
	cfg       *config.Config
	log       logger.Logger
	scheduler *InvestmentScheduler
	pb.UnimplementedManaServiceServer
}

func NewManaServiceImpl(db *sql.DB, cfg *config.Config, log logger.Logger, scheduler *InvestmentScheduler) *ManaServiceImpl {
	return &ManaServiceImpl{
		db:        db,
		cfg:       cfg,
		log:       log,
		scheduler: scheduler,
	}
}

func (s *ManaServiceImpl) TransferMana(ctx context.Context, req *pb.TransferManaRequest) (*pb.TransferManaResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Check sender's balance
	var senderBalance int64
	err = tx.QueryRowContext(ctx,
		"SELECT mana_balance FROM wizards WHERE id = $1",
		req.FromWizardId).Scan(&senderBalance)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Sender wizard not found: %v", err)
	}

	if senderBalance < req.Amount {
		return nil, status.Errorf(codes.FailedPrecondition, "Insufficient mana balance")
	}

	// Update balances
	_, err = tx.ExecContext(ctx,
		"UPDATE wizards SET mana_balance = mana_balance - $1 WHERE id = $2",
		req.Amount, req.FromWizardId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update sender balance: %v", err)
	}

	_, err = tx.ExecContext(ctx,
		"UPDATE wizards SET mana_balance = mana_balance + $1 WHERE id = $2",
		req.Amount, req.ToWizardId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update receiver balance: %v", err)
	}

	// Record transaction
	_, err = tx.ExecContext(ctx,
		"INSERT INTO mana_transactions (from_wizard_id, to_wizard_id, amount) VALUES ($1, $2, $3)",
		req.FromWizardId, req.ToWizardId, req.Amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to record transaction: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to commit transaction: %v", err)
	}

	return &pb.TransferManaResponse{Success: true}, nil
}

func (s *ManaServiceImpl) CreateInvestment(ctx context.Context, req *pb.CreateInvestmentRequest) (*pb.CreateInvestmentResponse, error) {
	// Validate investment type exists
	var minAmount, maxAmount int64
	var duration int32
	err := s.db.QueryRowContext(ctx,
		"SELECT min_amount, max_amount, duration_hours FROM investment_types WHERE id = $1",
		req.InvestmentTypeId).Scan(&minAmount, &maxAmount, &duration)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Investment type not found: %v", err)
	}

	// Validate investment amount
	if req.Amount < minAmount || (maxAmount > 0 && req.Amount > maxAmount) {
		return nil, status.Errorf(codes.InvalidArgument, "Investment amount outside allowed range")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Check wizard's balance
	var balance int64
	err = tx.QueryRowContext(ctx,
		"SELECT mana_balance FROM wizards WHERE id = $1",
		req.WizardId).Scan(&balance)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Wizard not found: %v", err)
	}

	if balance < req.Amount {
		return nil, status.Errorf(codes.FailedPrecondition, "Insufficient mana balance")
	}

	// Deduct investment amount
	_, err = tx.ExecContext(ctx,
		"UPDATE wizards SET mana_balance = mana_balance - $1 WHERE id = $2",
		req.Amount, req.WizardId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update balance: %v", err)
	}

	// Create investment record
	endTime := time.Now().Add(time.Duration(duration) * time.Hour)
	var investmentId int64
	err = tx.QueryRowContext(ctx,
		`INSERT INTO wizard_investments 
		(wizard_id, investment_type_id, amount, end_time, status) 
		VALUES ($1, $2, $3, $4, 'active') 
		RETURNING id`,
		req.WizardId, req.InvestmentTypeId, req.Amount, endTime).Scan(&investmentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create investment: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to commit transaction: %v", err)
	}

	// Schedule investment completion
	s.scheduler.ScheduleInvestmentCompletion(investmentId, endTime)

	return &pb.CreateInvestmentResponse{
		InvestmentId: investmentId,
		EndTime:      endTime.Unix(),
	}, nil
}

func (s *ManaServiceImpl) GetInvestments(ctx context.Context, req *pb.GetInvestmentsRequest) (*pb.GetInvestmentsResponse, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT i.id, i.amount, i.start_time, i.end_time, i.status, 
		        i.actual_return_rate, i.returned_amount, t.name, t.risk_level
		 FROM wizard_investments i
		 JOIN investment_types t ON i.investment_type_id = t.id
		 WHERE i.wizard_id = $1
		 ORDER BY i.created_at DESC`,
		req.WizardId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to fetch investments: %v", err)
	}
	defer rows.Close()

	var investments []*pb.Investment
	for rows.Next() {
		var inv pb.Investment
		var returnRate, returnedAmount sql.NullFloat64
		if err := rows.Scan(
			&inv.Id,
			&inv.Amount,
			&inv.StartTime,
			&inv.EndTime,
			&inv.Status,
			&returnRate,
			&returnedAmount,
			&inv.InvestmentType,
			&inv.RiskLevel,
		); err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to scan investment: %v", err)
		}

		if returnRate.Valid {
			inv.ActualReturnRate = returnRate.Float64
		}
		if returnedAmount.Valid {
			inv.ReturnedAmount = int64(returnedAmount.Float64)
		}

		investments = append(investments, &inv)
	}

	return &pb.GetInvestmentsResponse{
		Investments: investments,
	}, nil
}
