package mana

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/tectix/mysticfunds/pkg/config"
	"github.com/tectix/mysticfunds/pkg/logger"
	pb "github.com/tectix/mysticfunds/proto/mana"
	wizardpb "github.com/tectix/mysticfunds/proto/wizard"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ManaServiceImpl struct {
	db           *sql.DB
	cfg          *config.Config
	log          logger.Logger
	scheduler    *InvestmentScheduler
	wizardClient wizardpb.WizardServiceClient
	pb.UnimplementedManaServiceServer
}

func NewManaServiceImpl(db *sql.DB, cfg *config.Config, log logger.Logger) *ManaServiceImpl {
	// Create wizard service client
	wizardConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Failed to connect to wizard service", "error", err)
		panic(err)
	}
	wizardClient := wizardpb.NewWizardServiceClient(wizardConn)

	// Create scheduler with wizard client
	scheduler := NewInvestmentScheduler(db, log, wizardClient)
	scheduler.Start()

	return &ManaServiceImpl{
		db:           db,
		cfg:          cfg,
		log:          log,
		scheduler:    scheduler,
		wizardClient: wizardClient,
	}
}

func (s *ManaServiceImpl) TransferMana(ctx context.Context, req *pb.TransferManaRequest) (*pb.TransferManaResponse, error) {
	// Use wizard service to handle the transfer
	transferResp, err := s.wizardClient.TransferMana(ctx, &wizardpb.TransferManaRequest{
		FromWizardId: req.FromWizardId,
		ToWizardId:   req.ToWizardId,
		Amount:       req.Amount,
		Reason:       "Mana service transfer",
	})
	if err != nil {
		s.log.Error("Failed to transfer mana via wizard service", "error", err)
		return nil, status.Errorf(codes.Internal, "Failed to transfer mana: %v", err)
	}

	if !transferResp.Success {
		return &pb.TransferManaResponse{Success: false}, nil
	}

	// Record transaction in mana service database
	_, err = s.db.ExecContext(ctx,
		"INSERT INTO mana_transactions (from_wizard_id, to_wizard_id, amount) VALUES ($1, $2, $3)",
		req.FromWizardId, req.ToWizardId, req.Amount)
	if err != nil {
		s.log.Error("Failed to record mana transaction", "error", err)
		// Don't fail the entire operation since the transfer succeeded
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

	// Check wizard's balance via wizard service
	balanceResp, err := s.wizardClient.GetManaBalance(ctx, &wizardpb.GetManaBalanceRequest{
		WizardId: req.WizardId,
	})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Wizard not found: %v", err)
	}

	if balanceResp.Balance < req.Amount {
		return nil, status.Errorf(codes.FailedPrecondition, "Insufficient mana balance")
	}

	// Deduct investment amount via wizard service
	_, err = s.wizardClient.UpdateManaBalance(ctx, &wizardpb.UpdateManaBalanceRequest{
		WizardId: req.WizardId,
		Amount:   -req.Amount, // Negative amount to deduct
		Reason:   "Investment creation",
	})
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

func (s *ManaServiceImpl) GetManaBalance(ctx context.Context, req *pb.GetManaBalanceRequest) (*pb.GetManaBalanceResponse, error) {
	// Get balance from wizard service
	balanceResp, err := s.wizardClient.GetManaBalance(ctx, &wizardpb.GetManaBalanceRequest{
		WizardId: req.WizardId,
	})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Wizard not found: %v", err)
	}

	return &pb.GetManaBalanceResponse{
		Balance: balanceResp.Balance,
	}, nil
}

func (s *ManaServiceImpl) ListTransactions(ctx context.Context, req *pb.ListTransactionsRequest) (*pb.ListTransactionsResponse, error) {
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	pageNumber := req.PageNumber
	if pageNumber <= 0 {
		pageNumber = 1
	}

	offset := (pageNumber - 1) * pageSize

	// Get total count
	var totalCount int32
	err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM mana_transactions WHERE from_wizard_id = $1 OR to_wizard_id = $1",
		req.WizardId).Scan(&totalCount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to count transactions: %v", err)
	}

	// Get transactions
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, from_wizard_id, to_wizard_id, amount, created_at 
		 FROM mana_transactions 
		 WHERE from_wizard_id = $1 OR to_wizard_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`,
		req.WizardId, pageSize, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to query transactions: %v", err)
	}
	defer rows.Close()

	var transactions []*pb.ManaTransaction
	for rows.Next() {
		var tx pb.ManaTransaction
		var createdAt time.Time

		err := rows.Scan(
			&tx.Id,
			&tx.FromWizardId,
			&tx.ToWizardId,
			&tx.Amount,
			&createdAt,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to scan transaction: %v", err)
		}

		tx.CreatedAt = timestamppb.New(createdAt)
		transactions = append(transactions, &tx)
	}

	return &pb.ListTransactionsResponse{
		Transactions: transactions,
		TotalCount:   totalCount,
	}, nil
}

func (s *ManaServiceImpl) GetInvestmentTypes(ctx context.Context, req *pb.GetInvestmentTypesRequest) (*pb.GetInvestmentTypesResponse, error) {
	query := `SELECT id, name, description, min_amount, max_amount, duration_hours, base_return_rate, risk_level
	          FROM investment_types WHERE is_active = true`
	var args []interface{}
	argCount := 0

	// Add filters if provided
	if req.MinAmount > 0 {
		argCount++
		query += fmt.Sprintf(" AND min_amount >= $%d", argCount)
		args = append(args, req.MinAmount)
	}

	if req.MaxAmount > 0 {
		argCount++
		query += fmt.Sprintf(" AND (max_amount IS NULL OR max_amount <= $%d)", argCount)
		args = append(args, req.MaxAmount)
	}

	if req.RiskLevel > 0 {
		argCount++
		query += fmt.Sprintf(" AND risk_level = $%d", argCount)
		args = append(args, req.RiskLevel)
	}

	query += " ORDER BY risk_level, min_amount"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to query investment types: %v", err)
	}
	defer rows.Close()

	var investmentTypes []*pb.InvestmentType
	for rows.Next() {
		var it pb.InvestmentType
		var maxAmount sql.NullInt64

		err := rows.Scan(
			&it.Id,
			&it.Name,
			&it.Description,
			&it.MinAmount,
			&maxAmount,
			&it.DurationHours,
			&it.BaseReturnRate,
			&it.RiskLevel,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to scan investment type: %v", err)
		}

		if maxAmount.Valid {
			it.MaxAmount = maxAmount.Int64
		}

		investmentTypes = append(investmentTypes, &it)
	}

	return &pb.GetInvestmentTypesResponse{
		InvestmentTypes: investmentTypes,
	}, nil
}
