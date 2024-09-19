package wizard

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Alinoureddine1/mysticfunds/pkg/config"
	"github.com/Alinoureddine1/mysticfunds/pkg/logger"
	pb "github.com/Alinoureddine1/mysticfunds/proto/wizard"
)

type WizardServiceImpl struct {
	db     *sql.DB
	cfg    *config.Config
	logger logger.Logger
	pb.UnimplementedWizardServiceServer
}

func NewWizardServiceImpl(db *sql.DB, cfg *config.Config, logger logger.Logger) *WizardServiceImpl {
	return &WizardServiceImpl{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *WizardServiceImpl) CreateWizard(ctx context.Context, req *pb.CreateWizardRequest) (*pb.Wizard, error) {
	var id int64
	err := s.db.QueryRowContext(ctx,
		"INSERT INTO wizards (user_id, name, realm, element) VALUES ($1, $2, $3, $4) RETURNING id",
		req.UserId, req.Name, req.Realm, req.Element).Scan(&id)
	if err != nil {
		s.logger.Error("Failed to create wizard", "error", err)
		return nil, status.Error(codes.Internal, "Failed to create wizard")
	}

	return s.GetWizard(ctx, &pb.GetWizardRequest{Id: id})
}

func (s *WizardServiceImpl) GetWizard(ctx context.Context, req *pb.GetWizardRequest) (*pb.Wizard, error) {
	var wizard pb.Wizard
	var guildId sql.NullInt64
	var guildName sql.NullString
	var createdAt, updatedAt sql.NullTime

	err := s.db.QueryRowContext(ctx,
		`SELECT w.id, w.user_id, w.name, w.realm, w.element, w.mana_balance, 
                w.created_at, w.updated_at, g.id, g.name 
         FROM wizards w 
         LEFT JOIN guilds g ON w.guild_id = g.id 
         WHERE w.id = $1`,
		req.Id).Scan(
		&wizard.Id, &wizard.UserId, &wizard.Name, &wizard.Realm, &wizard.Element,
		&wizard.ManaBalance, &createdAt, &updatedAt, &guildId, &guildName)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "Wizard not found")
		}
		s.logger.Error("Failed to get wizard", "error", err)
		return nil, status.Error(codes.Internal, "Failed to get wizard")
	}

	if guildId.Valid && guildName.Valid {
		wizard.Guild = &pb.Guild{
			Id:   guildId.Int64,
			Name: guildName.String,
		}
	}

	if createdAt.Valid {
		wizard.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		wizard.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	return &wizard, nil
}

func (s *WizardServiceImpl) UpdateWizard(ctx context.Context, req *pb.UpdateWizardRequest) (*pb.Wizard, error) {
	_, err := s.db.ExecContext(ctx,
		"UPDATE wizards SET name = $1, realm = $2, element = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4",
		req.Name, req.Realm, req.Element, req.Id)
	if err != nil {
		s.logger.Error("Failed to update wizard", "error", err)
		return nil, status.Error(codes.Internal, "Failed to update wizard")
	}

	return s.GetWizard(ctx, &pb.GetWizardRequest{Id: req.Id})
}

func (s *WizardServiceImpl) ListWizards(ctx context.Context, req *pb.ListWizardsRequest) (*pb.ListWizardsResponse, error) {
	offset := (req.PageNumber - 1) * req.PageSize

	rows, err := s.db.QueryContext(ctx,
		`SELECT w.id, w.user_id, w.name, w.realm, w.element, w.mana_balance, 
                w.created_at, w.updated_at, g.id, g.name 
         FROM wizards w 
         LEFT JOIN guilds g ON w.guild_id = g.id 
         ORDER BY w.id LIMIT $1 OFFSET $2`,
		req.PageSize, offset)
	if err != nil {
		s.logger.Error("Failed to list wizards", "error", err)
		return nil, status.Error(codes.Internal, "Failed to list wizards")
	}
	defer rows.Close()

	var wizards []*pb.Wizard
	for rows.Next() {
		var wizard pb.Wizard
		var guildId sql.NullInt64
		var guildName sql.NullString
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(
			&wizard.Id, &wizard.UserId, &wizard.Name, &wizard.Realm, &wizard.Element,
			&wizard.ManaBalance, &createdAt, &updatedAt, &guildId, &guildName); err != nil {
			s.logger.Error("Failed to scan wizard row", "error", err)
			return nil, status.Error(codes.Internal, "Failed to list wizards")
		}
		if guildId.Valid && guildName.Valid {
			wizard.Guild = &pb.Guild{
				Id:   guildId.Int64,
				Name: guildName.String,
			}
		}
		if createdAt.Valid {
			wizard.CreatedAt = timestamppb.New(createdAt.Time)
		}
		if updatedAt.Valid {
			wizard.UpdatedAt = timestamppb.New(updatedAt.Time)
		}
		wizards = append(wizards, &wizard)
	}

	var totalCount int32
	err = s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM wizards").Scan(&totalCount)
	if err != nil {
		s.logger.Error("Failed to get total count of wizards", "error", err)
		return nil, status.Error(codes.Internal, "Failed to list wizards")
	}

	return &pb.ListWizardsResponse{
		Wizards:    wizards,
		TotalCount: totalCount,
	}, nil
}

func (s *WizardServiceImpl) DeleteWizard(ctx context.Context, req *pb.DeleteWizardRequest) (*pb.DeleteWizardResponse, error) {
	result, err := s.db.ExecContext(ctx, "DELETE FROM wizards WHERE id = $1", req.Id)
	if err != nil {
		s.logger.Error("Failed to delete wizard", "error", err)
		return nil, status.Error(codes.Internal, "Failed to delete wizard")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("Failed to get rows affected", "error", err)
		return nil, status.Error(codes.Internal, "Failed to delete wizard")
	}

	if rowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "Wizard not found")
	}

	return &pb.DeleteWizardResponse{Success: true}, nil
}

func (s *WizardServiceImpl) JoinGuild(ctx context.Context, req *pb.JoinGuildRequest) (*pb.Wizard, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to start transaction")
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			s.logger.Error("Failed to rollback transaction", "error", err)
		}
	}()

	var guildId int64
	err = tx.QueryRowContext(ctx, "SELECT id FROM guilds WHERE name = $1", req.GuildName).Scan(&guildId)
	if err == sql.ErrNoRows {
		// Create new guild if it doesn't exist
		err = tx.QueryRowContext(ctx, "INSERT INTO guilds (name) VALUES ($1) RETURNING id", req.GuildName).Scan(&guildId)
		if err != nil {
			s.logger.Error("Failed to create guild", "error", err)
			return nil, status.Error(codes.Internal, "Failed to create guild")
		}
	} else if err != nil {
		s.logger.Error("Failed to query guild", "error", err)
		return nil, status.Error(codes.Internal, "Failed to join guild")
	}

	_, err = tx.ExecContext(ctx, "UPDATE wizards SET guild_id = $1 WHERE id = $2", guildId, req.WizardId)
	if err != nil {
		s.logger.Error("Failed to update wizard's guild", "error", err)
		return nil, status.Error(codes.Internal, "Failed to join guild")
	}

	if err = tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction", "error", err)
		return nil, status.Error(codes.Internal, "Failed to join guild")
	}

	return s.GetWizard(ctx, &pb.GetWizardRequest{Id: req.WizardId})
}

func (s *WizardServiceImpl) LeaveGuild(ctx context.Context, req *pb.LeaveGuildRequest) (*pb.Wizard, error) {
	_, err := s.db.ExecContext(ctx, "UPDATE wizards SET guild_id = NULL WHERE id = $1", req.WizardId)
	if err != nil {
		s.logger.Error("Failed to leave guild", "error", err)
		return nil, status.Error(codes.Internal, "Failed to leave guild")
	}

	return s.GetWizard(ctx, &pb.GetWizardRequest{Id: req.WizardId})
}
