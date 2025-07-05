package wizard

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/tectix/mysticfunds/pkg/config"
	"github.com/tectix/mysticfunds/pkg/logger"
	pb "github.com/tectix/mysticfunds/proto/wizard"
)

type WizardServiceImpl struct {
	db     *sql.DB
	cfg    *config.Config
	logger logger.Logger
	ticker *JobTicker
	pb.UnimplementedWizardServiceServer
}

func NewWizardServiceImpl(db *sql.DB, cfg *config.Config, logger logger.Logger) *WizardServiceImpl {
	service := &WizardServiceImpl{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
	
	// Initialize the job ticker
	service.ticker = NewJobTicker(db, logger, service)
	
	// Start the ticker automatically
	service.ticker.Start()
	
	return service
}

func (s *WizardServiceImpl) CreateWizard(ctx context.Context, req *pb.CreateWizardRequest) (*pb.Wizard, error) {
	// Check if user already has 2 wizards
	var count int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM wizards WHERE user_id = $1", req.UserId).Scan(&count)
	if err != nil {
		s.logger.Error("Failed to check wizard count", "error", err)
		return nil, status.Error(codes.Internal, "Failed to check wizard count")
	}
	
	if count >= 2 {
		return nil, status.Error(codes.FailedPrecondition, "Users can only create up to 2 wizards")
	}

	var id int64
	err = s.db.QueryRowContext(ctx,
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
                w.created_at, w.updated_at, w.experience_points, w.level, g.id, g.name 
         FROM wizards w 
         LEFT JOIN guilds g ON w.guild_id = g.id 
         WHERE w.id = $1`,
		req.Id).Scan(
		&wizard.Id, &wizard.UserId, &wizard.Name, &wizard.Realm, &wizard.Element,
		&wizard.ManaBalance, &createdAt, &updatedAt, &wizard.ExperiencePoints, &wizard.Level, &guildId, &guildName)

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

	// Build dynamic query based on filters
	query := `SELECT w.id, w.user_id, w.name, w.realm, w.element, w.mana_balance, 
                     w.created_at, w.updated_at, w.experience_points, w.level, g.id, g.name 
              FROM wizards w 
              LEFT JOIN guilds g ON w.guild_id = g.id WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	// Add user filter if provided
	if req.UserId > 0 {
		query += " AND w.user_id = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, req.UserId)
		argIndex++
	}

	// Add realm filter if provided
	if req.Realm != "" {
		query += " AND w.realm = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, req.Realm)
		argIndex++
	}

	query += " ORDER BY w.id LIMIT $" + fmt.Sprintf("%d", argIndex) + " OFFSET $" + fmt.Sprintf("%d", argIndex+1)
	args = append(args, req.PageSize, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
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
			&wizard.ManaBalance, &createdAt, &updatedAt, &wizard.ExperiencePoints, &wizard.Level, &guildId, &guildName); err != nil {
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

// Job methods

func (s *WizardServiceImpl) CreateJob(ctx context.Context, req *pb.CreateJobRequest) (*pb.Job, error) {
	// First, get the realm_id from realm_name
	var realmId int64
	err := s.db.QueryRowContext(ctx, "SELECT id FROM realms WHERE name = $1", req.RealmName).Scan(&realmId)
	if err != nil {
		s.logger.Error("Failed to find realm", "error", err)
		return nil, status.Error(codes.NotFound, "Realm not found")
	}

	var jobId int64
	err = s.db.QueryRowContext(ctx,
		`INSERT INTO jobs (realm_id, title, description, required_element, required_level, 
		 mana_reward_per_hour, exp_reward_per_hour, duration_minutes, max_wizards, 
		 difficulty, job_type, location, special_requirements, created_by_wizard_id) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`,
		realmId, req.Title, req.Description, req.RequiredElement, req.RequiredLevel,
		req.ManaRewardPerHour, req.ExpRewardPerHour, req.DurationMinutes, req.MaxWizards,
		req.Difficulty, req.JobType, req.Location, req.SpecialRequirements, req.CreatedByWizardId).Scan(&jobId)
	if err != nil {
		s.logger.Error("Failed to create job", "error", err)
		return nil, status.Error(codes.Internal, "Failed to create job")
	}

	return s.GetJob(ctx, &pb.GetJobRequest{Id: jobId})
}

func (s *WizardServiceImpl) GetJob(ctx context.Context, req *pb.GetJobRequest) (*pb.Job, error) {
	var job pb.Job
	var createdAt, updatedAt sql.NullTime
	var location, specialRequirements sql.NullString
	var createdByWizardId sql.NullInt64

	err := s.db.QueryRowContext(ctx,
		`SELECT j.id, j.realm_id, r.name as realm_name, j.title, j.description, 
		 j.required_element, j.required_level, j.mana_reward_per_hour, j.exp_reward_per_hour,
		 j.duration_minutes, j.max_wizards, j.currently_assigned, j.difficulty, j.job_type, 
		 j.location, j.special_requirements, j.created_by_wizard_id, j.created_at, j.updated_at, j.is_active
		 FROM jobs j
		 JOIN realms r ON j.realm_id = r.id
		 WHERE j.id = $1`,
		req.Id).Scan(
		&job.Id, &job.RealmId, &job.RealmName, &job.Title, &job.Description,
		&job.RequiredElement, &job.RequiredLevel, &job.ManaRewardPerHour, &job.ExpRewardPerHour,
		&job.DurationMinutes, &job.MaxWizards, &job.CurrentlyAssigned, &job.Difficulty, &job.JobType,
		&location, &specialRequirements, &createdByWizardId, &createdAt, &updatedAt, &job.IsActive)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "Job not found")
		}
		s.logger.Error("Failed to get job", "error", err)
		return nil, status.Error(codes.Internal, "Failed to get job")
	}

	if location.Valid {
		job.Location = location.String
	}
	if specialRequirements.Valid {
		job.SpecialRequirements = specialRequirements.String
	}
	if createdByWizardId.Valid {
		job.CreatedByWizardId = createdByWizardId.Int64
	}
	if createdAt.Valid {
		job.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		job.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	return &job, nil
}

func (s *WizardServiceImpl) ListJobs(ctx context.Context, req *pb.ListJobsRequest) (*pb.ListJobsResponse, error) {
	offset := (req.PageNumber - 1) * req.PageSize

	// Build dynamic query based on filters
	query := `SELECT j.id, j.realm_id, r.name as realm_name, j.title, j.description, 
	          j.required_element, j.required_level, j.mana_reward_per_hour, j.exp_reward_per_hour,
	          j.duration_minutes, j.max_wizards, j.currently_assigned, j.difficulty, j.job_type, 
	          j.location, j.special_requirements, j.created_by_wizard_id, j.created_at, j.updated_at, j.is_active
	          FROM jobs j
	          JOIN realms r ON j.realm_id = r.id WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	// Add realm filter if provided
	if req.Realm != "" {
		query += " AND r.name = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, req.Realm)
		argIndex++
	}

	// Add element filter if provided
	if req.Element != "" {
		query += " AND j.required_element = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, req.Element)
		argIndex++
	}

	// Add difficulty filter if provided
	if req.Difficulty != "" {
		query += " AND j.difficulty = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, req.Difficulty)
		argIndex++
	}

	// Add active filter if requested
	if req.OnlyActive {
		query += " AND j.is_active = true"
	}

	query += " ORDER BY j.created_at DESC LIMIT $" + fmt.Sprintf("%d", argIndex) + " OFFSET $" + fmt.Sprintf("%d", argIndex+1)
	args = append(args, req.PageSize, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Failed to list jobs", "error", err)
		return nil, status.Error(codes.Internal, "Failed to list jobs")
	}
	defer rows.Close()

	var jobs []*pb.Job
	for rows.Next() {
		var job pb.Job
		var createdAt, updatedAt sql.NullTime
		var location, specialRequirements sql.NullString
		var createdByWizardId sql.NullInt64

		if err := rows.Scan(
			&job.Id, &job.RealmId, &job.RealmName, &job.Title, &job.Description,
			&job.RequiredElement, &job.RequiredLevel, &job.ManaRewardPerHour, &job.ExpRewardPerHour,
			&job.DurationMinutes, &job.MaxWizards, &job.CurrentlyAssigned, &job.Difficulty, &job.JobType,
			&location, &specialRequirements, &createdByWizardId, &createdAt, &updatedAt, &job.IsActive); err != nil {
			s.logger.Error("Failed to scan job row", "error", err)
			return nil, status.Error(codes.Internal, "Failed to list jobs")
		}

		if location.Valid {
			job.Location = location.String
		}
		if specialRequirements.Valid {
			job.SpecialRequirements = specialRequirements.String
		}
		if createdByWizardId.Valid {
			job.CreatedByWizardId = createdByWizardId.Int64
		}
		if createdAt.Valid {
			job.CreatedAt = timestamppb.New(createdAt.Time)
		}
		if updatedAt.Valid {
			job.UpdatedAt = timestamppb.New(updatedAt.Time)
		}

		jobs = append(jobs, &job)
	}

	// Get total count for pagination
	countQuery := "SELECT COUNT(*) FROM jobs j JOIN realms r ON j.realm_id = r.id WHERE 1=1"
	countArgs := []interface{}{}
	if req.Realm != "" {
		countQuery += " AND r.name = $1"
		countArgs = append(countArgs, req.Realm)
	}
	if req.OnlyActive {
		countQuery += " AND j.is_active = true"
	}

	var totalCount int32
	err = s.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		s.logger.Error("Failed to get total count of jobs", "error", err)
		return nil, status.Error(codes.Internal, "Failed to list jobs")
	}

	return &pb.ListJobsResponse{
		Jobs:       jobs,
		TotalCount: totalCount,
	}, nil
}

func (s *WizardServiceImpl) UpdateJob(ctx context.Context, req *pb.UpdateJobRequest) (*pb.Job, error) {
	_, err := s.db.ExecContext(ctx,
		`UPDATE jobs SET title = $1, description = $2, mana_reward_per_hour = $3, 
		 exp_reward_per_hour = $4, max_wizards = $5, is_active = $6, updated_at = CURRENT_TIMESTAMP 
		 WHERE id = $7`,
		req.Title, req.Description, req.ManaRewardPerHour, req.ExpRewardPerHour, 
		req.MaxWizards, req.IsActive, req.Id)
	if err != nil {
		s.logger.Error("Failed to update job", "error", err)
		return nil, status.Error(codes.Internal, "Failed to update job")
	}

	return s.GetJob(ctx, &pb.GetJobRequest{Id: req.Id})
}

func (s *WizardServiceImpl) DeleteJob(ctx context.Context, req *pb.DeleteJobRequest) (*pb.DeleteJobResponse, error) {
	result, err := s.db.ExecContext(ctx, "DELETE FROM jobs WHERE id = $1", req.Id)
	if err != nil {
		s.logger.Error("Failed to delete job", "error", err)
		return nil, status.Error(codes.Internal, "Failed to delete job")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("Failed to get rows affected", "error", err)
		return nil, status.Error(codes.Internal, "Failed to delete job")
	}

	if rowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "Job not found")
	}

	return &pb.DeleteJobResponse{Success: true}, nil
}

// Job assignment methods

func (s *WizardServiceImpl) AssignWizardToJob(ctx context.Context, req *pb.AssignWizardToJobRequest) (*pb.JobAssignment, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to start transaction")
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			s.logger.Error("Failed to rollback transaction", "error", err)
		}
	}()

	// Check if job exists and has available slots
	var maxWizards, currentlyAssigned, durationMinutes int32
	err = tx.QueryRowContext(ctx, 
		"SELECT max_wizards, currently_assigned, duration_minutes FROM jobs WHERE id = $1 AND is_active = true", 
		req.JobId).Scan(&maxWizards, &currentlyAssigned, &durationMinutes)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "Job not found or inactive")
		}
		s.logger.Error("Failed to check job availability", "error", err)
		return nil, status.Error(codes.Internal, "Failed to assign wizard to job")
	}

	if currentlyAssigned >= maxWizards {
		return nil, status.Error(codes.FailedPrecondition, "Job is full")
	}

	// Check if wizard meets requirements (element and level)
	var wizardElement string
	var wizardLevel int32
	var requiredElement string
	var requiredLevel int32
	
	err = tx.QueryRowContext(ctx, 
		"SELECT w.element, w.level FROM wizards w WHERE w.id = $1", 
		req.WizardId).Scan(&wizardElement, &wizardLevel)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "Wizard not found")
		}
		s.logger.Error("Failed to get wizard info", "error", err)
		return nil, status.Error(codes.Internal, "Failed to assign wizard to job")
	}

	err = tx.QueryRowContext(ctx, 
		"SELECT required_element, required_level FROM jobs WHERE id = $1", 
		req.JobId).Scan(&requiredElement, &requiredLevel)
	if err != nil {
		s.logger.Error("Failed to get job requirements", "error", err)
		return nil, status.Error(codes.Internal, "Failed to assign wizard to job")
	}

	if wizardElement != requiredElement {
		return nil, status.Error(codes.FailedPrecondition, fmt.Sprintf("Wizard element %s does not match required element %s", wizardElement, requiredElement))
	}

	if wizardLevel < requiredLevel {
		return nil, status.Error(codes.FailedPrecondition, fmt.Sprintf("Wizard level %d is below required level %d", wizardLevel, requiredLevel))
	}

	// Try to create new job assignment 
	// The database constraint will prevent duplicate active assignments
	var assignmentId int64
	err = tx.QueryRowContext(ctx,
		`INSERT INTO job_assignments (job_id, wizard_id, status) 
		 VALUES ($1, $2, 'assigned') RETURNING id`,
		req.JobId, req.WizardId).Scan(&assignmentId)
	if err != nil {
		// Check if it's a constraint violation (wizard already assigned)
		if strings.Contains(err.Error(), "job_assignments_active_unique") {
			return nil, status.Error(codes.FailedPrecondition, "Wizard is already assigned to this job")
		}
		s.logger.Error("Failed to create job assignment", "error", err)
		return nil, status.Error(codes.Internal, "Failed to assign wizard to job")
	}

	// Update job's currently_assigned count
	_, err = tx.ExecContext(ctx,
		"UPDATE jobs SET currently_assigned = currently_assigned + 1 WHERE id = $1",
		req.JobId)
	if err != nil {
		s.logger.Error("Failed to update job assignment count", "error", err)
		return nil, status.Error(codes.Internal, "Failed to assign wizard to job")
	}

	// Create job progress record with proper time tracking
	startTime := time.Now()
	endTime := startTime.Add(time.Duration(durationMinutes) * time.Minute)
	
	_, err = tx.ExecContext(ctx,
		`INSERT INTO job_progress (assignment_id, started_at, actual_start_time, expected_end_time, progress_percentage, time_worked_minutes, is_active, last_tick_time) 
		 VALUES ($1, $2, $2, $3, 0, 0, true, $2)`,
		assignmentId, startTime, endTime)
	if err != nil {
		s.logger.Error("Failed to create job progress record", "error", err)
		return nil, status.Error(codes.Internal, "Failed to assign wizard to job")
	}

	// Create activity log for job assignment
	_, err = tx.ExecContext(ctx,
		`INSERT INTO activity_logs (user_id, wizard_id, activity_type, activity_description, metadata) 
		 SELECT w.user_id, w.id, 'job_assigned', 
		        'Started working on job: ' || j.title,
		        json_build_object('job_id', j.id, 'assignment_id', $1::bigint, 'job_title', j.title)
		 FROM wizards w 
		 JOIN jobs j ON j.id = $2
		 WHERE w.id = $3`,
		assignmentId, req.JobId, req.WizardId)
	if err != nil {
		s.logger.Error("Failed to create activity log", "error", err)
		// Don't fail the transaction for activity log issues
	}

	if err = tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction", "error", err)
		return nil, status.Error(codes.Internal, "Failed to assign wizard to job")
	}

	// Return the created assignment
	return s.getJobAssignmentByID(ctx, assignmentId)
}

func (s *WizardServiceImpl) GetJobAssignments(ctx context.Context, req *pb.GetJobAssignmentsRequest) (*pb.GetJobAssignmentsResponse, error) {
	offset := (req.PageNumber - 1) * req.PageSize

	query := `SELECT ja.id, ja.job_id, ja.wizard_id, w.name as wizard_name, ja.assigned_at, 
	          ja.started_at, ja.completed_at, ja.status, ja.mana_earned, ja.exp_earned, ja.notes,
	          j.title, j.description, j.required_element, j.required_level, j.mana_reward_per_hour,
	          j.exp_reward_per_hour, j.duration_minutes, j.max_wizards, j.currently_assigned,
	          j.difficulty, j.job_type, j.location, j.special_requirements, r.name as realm_name,
	          jp.id as progress_id, jp.started_at as progress_started, jp.progress_percentage,
	          jp.time_worked_minutes, jp.is_active as progress_active, jp.last_updated_at as progress_updated
	          FROM job_assignments ja
	          JOIN wizards w ON ja.wizard_id = w.id
	          JOIN jobs j ON ja.job_id = j.id
	          JOIN realms r ON j.realm_id = r.id
	          LEFT JOIN job_progress jp ON ja.id = jp.assignment_id
	          WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if req.WizardId > 0 {
		query += " AND ja.wizard_id = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, req.WizardId)
		argIndex++
	}

	if req.JobId > 0 {
		query += " AND ja.job_id = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, req.JobId)
		argIndex++
	}

	if req.Status != "" {
		query += " AND ja.status = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, req.Status)
		argIndex++
	}

	query += " ORDER BY ja.assigned_at DESC LIMIT $" + fmt.Sprintf("%d", argIndex) + " OFFSET $" + fmt.Sprintf("%d", argIndex+1)
	args = append(args, req.PageSize, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Failed to get job assignments", "error", err)
		return nil, status.Error(codes.Internal, "Failed to get job assignments")
	}
	defer rows.Close()

	var assignments []*pb.JobAssignment
	for rows.Next() {
		var assignment pb.JobAssignment
		var job pb.Job
		var progress pb.JobProgress
		var assignedAt, startedAt, completedAt sql.NullTime
		var notes sql.NullString
		var progressId sql.NullInt64
		var progressStarted, progressUpdated sql.NullTime
		var progressPercentage, timeWorked sql.NullInt32
		var progressActive sql.NullBool

		if err := rows.Scan(
			&assignment.Id, &assignment.JobId, &assignment.WizardId, &assignment.WizardName,
			&assignedAt, &startedAt, &completedAt, &assignment.Status, 
			&assignment.ManaEarned, &assignment.ExpEarned, &notes,
			&job.Title, &job.Description, &job.RequiredElement, &job.RequiredLevel,
			&job.ManaRewardPerHour, &job.ExpRewardPerHour, &job.DurationMinutes,
			&job.MaxWizards, &job.CurrentlyAssigned, &job.Difficulty, &job.JobType,
			&job.Location, &job.SpecialRequirements, &job.RealmName,
			&progressId, &progressStarted, &progressPercentage, &timeWorked,
			&progressActive, &progressUpdated); err != nil {
			s.logger.Error("Failed to scan job assignment row", "error", err)
			return nil, status.Error(codes.Internal, "Failed to get job assignments")
		}

		// Set job data
		job.Id = assignment.JobId
		assignment.Job = &job

		// Set assignment timestamps
		if assignedAt.Valid {
			assignment.AssignedAt = timestamppb.New(assignedAt.Time)
		}
		if startedAt.Valid {
			assignment.StartedAt = timestamppb.New(startedAt.Time)
		}
		if completedAt.Valid {
			assignment.CompletedAt = timestamppb.New(completedAt.Time)
		}
		if notes.Valid {
			assignment.Notes = notes.String
		}

		// Set progress data if available
		if progressId.Valid {
			progress.Id = progressId.Int64
			progress.AssignmentId = assignment.Id
			if progressStarted.Valid {
				progress.StartedAt = timestamppb.New(progressStarted.Time)
			}
			if progressUpdated.Valid {
				progress.LastUpdatedAt = timestamppb.New(progressUpdated.Time)
			}
			if progressPercentage.Valid {
				progress.ProgressPercentage = progressPercentage.Int32
			}
			if timeWorked.Valid {
				progress.TimeWorkedMinutes = timeWorked.Int32
			}
			if progressActive.Valid {
				progress.IsActive = progressActive.Bool
			}
			assignment.Progress = &progress
		}

		assignments = append(assignments, &assignment)
	}

	return &pb.GetJobAssignmentsResponse{
		Assignments: assignments,
		TotalCount:  int32(len(assignments)), // Simplified for now
	}, nil
}

func (s *WizardServiceImpl) CompleteJobAssignment(ctx context.Context, req *pb.CompleteJobAssignmentRequest) (*pb.JobAssignment, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to start transaction")
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			s.logger.Error("Failed to rollback transaction", "error", err)
		}
	}()

	// Get assignment and job details
	var jobId, wizardId int64
	var manaRewardPerHour, expRewardPerHour, durationMinutes int32
	var currentExp, currentLevel int32
	err = tx.QueryRowContext(ctx,
		`SELECT ja.job_id, ja.wizard_id, j.mana_reward_per_hour, j.exp_reward_per_hour, j.duration_minutes,
		        w.experience_points, w.level
		 FROM job_assignments ja
		 JOIN jobs j ON ja.job_id = j.id
		 JOIN wizards w ON ja.wizard_id = w.id
		 WHERE ja.id = $1 AND ja.status IN ('assigned', 'in_progress')`,
		req.AssignmentId).Scan(&jobId, &wizardId, &manaRewardPerHour, &expRewardPerHour, &durationMinutes, &currentExp, &currentLevel)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "Assignment not found or already completed")
		}
		s.logger.Error("Failed to get assignment details", "error", err)
		return nil, status.Error(codes.Internal, "Failed to complete job assignment")
	}

	// Calculate total rewards - convert per-hour rates to per-minute rates
	manaRewardPerMinute := manaRewardPerHour / 60
	expRewardPerMinute := expRewardPerHour / 60
	totalMana := manaRewardPerMinute * durationMinutes
	totalExp := expRewardPerMinute * durationMinutes

	// Update assignment status and rewards
	_, err = tx.ExecContext(ctx,
		`UPDATE job_assignments SET status = 'completed', completed_at = CURRENT_TIMESTAMP,
		 mana_earned = $1, exp_earned = $2 WHERE id = $3`,
		totalMana, totalExp, req.AssignmentId)
	if err != nil {
		s.logger.Error("Failed to update assignment", "error", err)
		return nil, status.Error(codes.Internal, "Failed to complete job assignment")
	}

	// Calculate new experience and level
	newExp := currentExp + totalExp
	newLevel := s.calculateLevel(newExp)
	
	// Update wizard's mana, experience, and level
	_, err = tx.ExecContext(ctx,
		`UPDATE wizards SET mana_balance = mana_balance + $1, experience_points = $2, level = $3
		 WHERE id = $4`,
		totalMana, newExp, newLevel, wizardId)
	if err != nil {
		s.logger.Error("Failed to update wizard rewards", "error", err)
		return nil, status.Error(codes.Internal, "Failed to complete job assignment")
	}

	// Log level up if it occurred
	if newLevel > currentLevel {
		s.logger.Info("Wizard leveled up!", "wizard_id", wizardId, "old_level", currentLevel, "new_level", newLevel)
		
		// Create activity log for level up
		_, err = tx.ExecContext(ctx,
			`INSERT INTO activity_logs (user_id, wizard_id, activity_type, activity_description, metadata) 
			 SELECT w.user_id, w.id, 'level_up', 
			        'Leveled up from ' || $2::text || ' to ' || $3::text || '!',
			        json_build_object('old_level', $2::integer, 'new_level', $3::integer, 'exp_gained', $4::integer)
			 FROM wizards w 
			 WHERE w.id = $1`,
			wizardId, currentLevel, newLevel, totalExp)
		if err != nil {
			s.logger.Error("Failed to create level up activity log", "error", err)
			// Don't fail the transaction for activity log issues
		}
	}

	// Update job's currently_assigned count
	_, err = tx.ExecContext(ctx,
		"UPDATE jobs SET currently_assigned = currently_assigned - 1 WHERE id = $1",
		jobId)
	if err != nil {
		s.logger.Error("Failed to update job assignment count", "error", err)
		return nil, status.Error(codes.Internal, "Failed to complete job assignment")
	}

	// Mark progress as complete and inactive
	_, err = tx.ExecContext(ctx,
		`UPDATE job_progress SET progress_percentage = 100, is_active = false,
		 last_updated_at = CURRENT_TIMESTAMP WHERE assignment_id = $1`,
		req.AssignmentId)
	if err != nil {
		s.logger.Error("Failed to update job progress", "error", err)
		// Don't fail the transaction for progress update issues
	}

	// Create activity log for job completion
	_, err = tx.ExecContext(ctx,
		`INSERT INTO activity_logs (user_id, wizard_id, activity_type, activity_description, metadata) 
		 SELECT w.user_id, w.id, 'job_completed', 
		        'Completed job: ' || j.title || ' - Earned ' || $2::text || ' mana and ' || $3::text || ' EXP',
		        json_build_object('job_id', j.id, 'assignment_id', $1::bigint, 'job_title', j.title, 'mana_earned', $2::integer, 'exp_earned', $3::integer)
		 FROM job_assignments ja
		 JOIN wizards w ON ja.wizard_id = w.id
		 JOIN jobs j ON ja.job_id = j.id
		 WHERE ja.id = $1::bigint`,
		req.AssignmentId, totalMana, totalExp)
	if err != nil {
		s.logger.Error("Failed to create activity log", "error", err)
		// Don't fail the transaction for activity log issues
	}

	if err = tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction", "error", err)
		return nil, status.Error(codes.Internal, "Failed to complete job assignment")
	}

	return s.getJobAssignmentByID(ctx, req.AssignmentId)
}

// calculateLevel determines the wizard's level based on experience points
// Using a standard RPG leveling formula: Level = floor(sqrt(experience / 100)) + 1
// This means: Level 1 = 0-99 exp, Level 2 = 100-399 exp, Level 3 = 400-899 exp, etc.
func (s *WizardServiceImpl) calculateLevel(experiencePoints int32) int32 {
	if experiencePoints < 0 {
		return 1
	}
	
	// Use a gentler curve: level = floor(experience / 100) + 1
	// This means: Level 1 = 0-99 exp, Level 2 = 100-199 exp, Level 3 = 200-299 exp, etc.
	level := (experiencePoints / 100) + 1
	
	// Cap at level 50 for now
	if level > 50 {
		return 50
	}
	
	return level
}

func (s *WizardServiceImpl) CancelJobAssignment(ctx context.Context, req *pb.CancelJobAssignmentRequest) (*pb.JobAssignment, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to start transaction")
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			s.logger.Error("Failed to rollback transaction", "error", err)
		}
	}()

	// Get job_id for updating count
	var jobId int64
	err = tx.QueryRowContext(ctx,
		"SELECT job_id FROM job_assignments WHERE id = $1 AND status IN ('assigned', 'in_progress')",
		req.AssignmentId).Scan(&jobId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "Assignment not found or cannot be cancelled")
		}
		s.logger.Error("Failed to get assignment job_id", "error", err)
		return nil, status.Error(codes.Internal, "Failed to cancel job assignment")
	}

	// Update assignment status
	_, err = tx.ExecContext(ctx,
		"UPDATE job_assignments SET status = 'cancelled', notes = $1 WHERE id = $2",
		req.Reason, req.AssignmentId)
	if err != nil {
		s.logger.Error("Failed to update assignment status", "error", err)
		return nil, status.Error(codes.Internal, "Failed to cancel job assignment")
	}

	// Update job's currently_assigned count
	_, err = tx.ExecContext(ctx,
		"UPDATE jobs SET currently_assigned = currently_assigned - 1 WHERE id = $1",
		jobId)
	if err != nil {
		s.logger.Error("Failed to update job assignment count", "error", err)
		return nil, status.Error(codes.Internal, "Failed to cancel job assignment")
	}

	if err = tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction", "error", err)
		return nil, status.Error(codes.Internal, "Failed to cancel job assignment")
	}

	return s.getJobAssignmentByID(ctx, req.AssignmentId)
}

// Helper function to get job assignment by ID
func (s *WizardServiceImpl) getJobAssignmentByID(ctx context.Context, id int64) (*pb.JobAssignment, error) {
	var assignment pb.JobAssignment
	var assignedAt, startedAt, completedAt sql.NullTime
	var notes sql.NullString
	var progress pb.JobProgress
	var progressStartedAt, progressLastUpdated, progressCreatedAt sql.NullTime

	err := s.db.QueryRowContext(ctx,
		`SELECT ja.id, ja.job_id, ja.wizard_id, w.name as wizard_name, ja.assigned_at, 
		 ja.started_at, ja.completed_at, ja.status, ja.mana_earned, ja.exp_earned, ja.notes,
		 jp.id, jp.assignment_id, jp.started_at, jp.last_updated_at, jp.progress_percentage,
		 jp.time_worked_minutes, jp.is_active, jp.created_at
		 FROM job_assignments ja
		 JOIN wizards w ON ja.wizard_id = w.id
		 LEFT JOIN job_progress jp ON ja.id = jp.assignment_id
		 WHERE ja.id = $1`,
		id).Scan(
		&assignment.Id, &assignment.JobId, &assignment.WizardId, &assignment.WizardName,
		&assignedAt, &startedAt, &completedAt, &assignment.Status, 
		&assignment.ManaEarned, &assignment.ExpEarned, &notes,
		&progress.Id, &progress.AssignmentId, &progressStartedAt, &progressLastUpdated,
		&progress.ProgressPercentage, &progress.TimeWorkedMinutes, &progress.IsActive, &progressCreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "Job assignment not found")
		}
		s.logger.Error("Failed to get job assignment", "error", err)
		return nil, status.Error(codes.Internal, "Failed to get job assignment")
	}

	if assignedAt.Valid {
		assignment.AssignedAt = timestamppb.New(assignedAt.Time)
	}
	if startedAt.Valid {
		assignment.StartedAt = timestamppb.New(startedAt.Time)
	}
	if completedAt.Valid {
		assignment.CompletedAt = timestamppb.New(completedAt.Time)
	}
	if notes.Valid {
		assignment.Notes = notes.String
	}

	// Add progress data if available
	if progress.Id > 0 {
		if progressStartedAt.Valid {
			progress.StartedAt = timestamppb.New(progressStartedAt.Time)
		}
		if progressLastUpdated.Valid {
			progress.LastUpdatedAt = timestamppb.New(progressLastUpdated.Time)
		}
		if progressCreatedAt.Valid {
			progress.CreatedAt = timestamppb.New(progressCreatedAt.Time)
		}
		assignment.Progress = &progress
	}

	return &assignment, nil
}

// Progress tracking methods
func (s *WizardServiceImpl) UpdateJobProgress(ctx context.Context, req *pb.UpdateJobProgressRequest) (*pb.JobProgress, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to start transaction")
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			s.logger.Error("Failed to rollback transaction", "error", err)
		}
	}()

	// Validate that the assignment exists and is active
	var assignmentStatus string
	err = tx.QueryRowContext(ctx,
		"SELECT status FROM job_assignments WHERE id = $1",
		req.AssignmentId).Scan(&assignmentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "Job assignment not found")
		}
		s.logger.Error("Failed to get assignment status", "error", err)
		return nil, status.Error(codes.Internal, "Failed to update job progress")
	}

	// Only allow progress updates for assigned or in_progress assignments
	if assignmentStatus != "assigned" && assignmentStatus != "in_progress" {
		return nil, status.Error(codes.FailedPrecondition, "Cannot update progress for completed or cancelled assignments")
	}

	// Ensure progress doesn't go backwards and is within valid range
	progressPercentage := req.ProgressPercentage
	if progressPercentage < 0 {
		progressPercentage = 0
	}
	if progressPercentage > 100 {
		progressPercentage = 100
	}

	// Check current progress to prevent going backwards
	var currentProgress int32
	err = tx.QueryRowContext(ctx,
		"SELECT COALESCE(progress_percentage, 0) FROM job_progress WHERE assignment_id = $1",
		req.AssignmentId).Scan(&currentProgress)
	if err != nil && err != sql.ErrNoRows {
		s.logger.Error("Failed to get current progress", "error", err)
		return nil, status.Error(codes.Internal, "Failed to update job progress")
	}

	// Only allow progress to increase
	if progressPercentage < currentProgress {
		s.logger.Warn("Attempted to decrease progress", "assignment_id", req.AssignmentId, "current", currentProgress, "new", progressPercentage)
		progressPercentage = currentProgress
	}

	// Update the progress
	_, err = tx.ExecContext(ctx,
		`UPDATE job_progress SET 
		 progress_percentage = $1, time_worked_minutes = $2, last_updated_at = CURRENT_TIMESTAMP
		 WHERE assignment_id = $3`,
		progressPercentage, req.TimeWorkedMinutes, req.AssignmentId)
	if err != nil {
		s.logger.Error("Failed to update job progress", "error", err)
		return nil, status.Error(codes.Internal, "Failed to update job progress")
	}

	// Update assignment status to 'in_progress' if it's still 'assigned' and progress > 0
	if progressPercentage > 0 && assignmentStatus == "assigned" {
		_, err = tx.ExecContext(ctx,
			`UPDATE job_assignments SET 
			 status = 'in_progress', started_at = COALESCE(started_at, CURRENT_TIMESTAMP)
			 WHERE id = $1`,
			req.AssignmentId)
		if err != nil {
			s.logger.Error("Failed to update assignment status", "error", err)
			return nil, status.Error(codes.Internal, "Failed to update job progress")
		}
	}

	// Auto-complete the job if progress reaches 100%
	if progressPercentage >= 100 && assignmentStatus != "completed" {
		s.logger.Info("Auto-completing job assignment", "assignment_id", req.AssignmentId)
		
		// We'll let the frontend handle completion to maintain the existing flow
		// But we can mark the progress as ready for completion
		_, err = tx.ExecContext(ctx,
			`UPDATE job_progress SET is_active = false WHERE assignment_id = $1`,
			req.AssignmentId)
		if err != nil {
			s.logger.Error("Failed to mark progress as complete", "error", err)
			// Don't fail the transaction for this
		}
	}

	if err = tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction", "error", err)
		return nil, status.Error(codes.Internal, "Failed to update job progress")
	}

	return s.GetJobProgress(ctx, &pb.GetJobProgressRequest{AssignmentId: req.AssignmentId})
}

func (s *WizardServiceImpl) GetJobProgress(ctx context.Context, req *pb.GetJobProgressRequest) (*pb.JobProgress, error) {
	var progress pb.JobProgress
	var startedAt, lastUpdated, createdAt, actualStartTime, expectedEndTime sql.NullTime

	err := s.db.QueryRowContext(ctx,
		`SELECT jp.id, jp.assignment_id, jp.started_at, jp.last_updated_at, jp.progress_percentage,
		 jp.time_worked_minutes, jp.is_active, jp.created_at, jp.actual_start_time, jp.expected_end_time
		 FROM job_progress jp WHERE jp.assignment_id = $1`,
		req.AssignmentId).Scan(
		&progress.Id, &progress.AssignmentId, &startedAt, &lastUpdated,
		&progress.ProgressPercentage, &progress.TimeWorkedMinutes, &progress.IsActive, &createdAt,
		&actualStartTime, &expectedEndTime)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "Job progress not found")
		}
		s.logger.Error("Failed to get job progress", "error", err)
		return nil, status.Error(codes.Internal, "Failed to get job progress")
	}

	// Calculate real-time progress if job is active
	if progress.IsActive && actualStartTime.Valid && expectedEndTime.Valid {
		now := time.Now()
		startTime := actualStartTime.Time
		endTime := expectedEndTime.Time
		
		if now.After(startTime) {
			elapsed := now.Sub(startTime)
			total := endTime.Sub(startTime)
			
			if elapsed >= total {
				progress.ProgressPercentage = 100
			} else {
				progressFloat := float64(elapsed) / float64(total) * 100
				realTimeProgress := int32(progressFloat)
				
				// Use the higher of stored progress or real-time progress
				if realTimeProgress > progress.ProgressPercentage {
					progress.ProgressPercentage = realTimeProgress
				}
			}
			
			// Update time worked
			progress.TimeWorkedMinutes = int32(elapsed.Minutes())
		}
	}

	if startedAt.Valid {
		progress.StartedAt = timestamppb.New(startedAt.Time)
	}
	if lastUpdated.Valid {
		progress.LastUpdatedAt = timestamppb.New(lastUpdated.Time)
	}
	if createdAt.Valid {
		progress.CreatedAt = timestamppb.New(createdAt.Time)
	}

	return &progress, nil
}

// Activity logs methods
func (s *WizardServiceImpl) GetActivities(ctx context.Context, req *pb.GetActivitiesRequest) (*pb.GetActivitiesResponse, error) {
	offset := (req.PageNumber - 1) * req.PageSize

	query := `SELECT id, user_id, wizard_id, activity_type, activity_description, 
	          COALESCE(metadata::text, '{}') as metadata, created_at
	          FROM activity_logs WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if req.UserId > 0 {
		query += " AND user_id = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, req.UserId)
		argIndex++
	}

	if req.WizardId > 0 {
		query += " AND wizard_id = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, req.WizardId)
		argIndex++
	}

	if req.ActivityType != "" {
		query += " AND activity_type = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, req.ActivityType)
		argIndex++
	}

	query += " ORDER BY created_at DESC LIMIT $" + fmt.Sprintf("%d", argIndex) + " OFFSET $" + fmt.Sprintf("%d", argIndex+1)
	args = append(args, req.PageSize, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("Failed to list activities", "error", err)
		return nil, status.Error(codes.Internal, "Failed to list activities")
	}
	defer rows.Close()

	var activities []*pb.ActivityLog
	for rows.Next() {
		var activity pb.ActivityLog
		var createdAt sql.NullTime
		var wizardId sql.NullInt64

		if err := rows.Scan(
			&activity.Id, &activity.UserId, &wizardId, &activity.ActivityType,
			&activity.ActivityDescription, &activity.Metadata, &createdAt); err != nil {
			s.logger.Error("Failed to scan activity row", "error", err)
			return nil, status.Error(codes.Internal, "Failed to list activities")
		}

		if wizardId.Valid {
			activity.WizardId = wizardId.Int64
		}
		if createdAt.Valid {
			activity.CreatedAt = timestamppb.New(createdAt.Time)
		}

		activities = append(activities, &activity)
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM activity_logs WHERE 1=1"
	countArgs := []interface{}{}
	if req.UserId > 0 {
		countQuery += " AND user_id = $1"
		countArgs = append(countArgs, req.UserId)
	}

	var totalCount int32
	err = s.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		s.logger.Error("Failed to get total count of activities", "error", err)
		return nil, status.Error(codes.Internal, "Failed to list activities")
	}

	return &pb.GetActivitiesResponse{
		Activities:  activities,
		TotalCount: totalCount,
	}, nil
}

// Realm methods for convenience (since jobs are realm-based)
func (s *WizardServiceImpl) GetRealms(ctx context.Context, req *pb.GetRealmsRequest) (*pb.GetRealmsResponse, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, description FROM realms ORDER BY name")
	if err != nil {
		s.logger.Error("Failed to get realms", "error", err)
		return nil, status.Error(codes.Internal, "Failed to get realms")
	}
	defer rows.Close()

	var realms []*pb.Realm
	for rows.Next() {
		var realm pb.Realm
		err := rows.Scan(&realm.Id, &realm.Name, &realm.Description)
		if err != nil {
			s.logger.Error("Failed to scan realm", "error", err)
			continue
		}
		realms = append(realms, &realm)
	}

	return &pb.GetRealmsResponse{
		Realms: realms,
	}, nil
}

// Mana management methods
func (s *WizardServiceImpl) GetManaBalance(ctx context.Context, req *pb.GetManaBalanceRequest) (*pb.GetManaBalanceResponse, error) {
	var balance int64
	err := s.db.QueryRowContext(ctx,
		"SELECT mana_balance FROM wizards WHERE id = $1",
		req.WizardId).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "Wizard not found")
		}
		s.logger.Error("Failed to get mana balance", "error", err)
		return nil, status.Error(codes.Internal, "Failed to get mana balance")
	}

	return &pb.GetManaBalanceResponse{
		Balance: balance,
	}, nil
}

func (s *WizardServiceImpl) UpdateManaBalance(ctx context.Context, req *pb.UpdateManaBalanceRequest) (*pb.UpdateManaBalanceResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to start transaction")
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			s.logger.Error("Failed to rollback transaction", "error", err)
		}
	}()

	// Check current balance and validate the operation
	var currentBalance int64
	err = tx.QueryRowContext(ctx,
		"SELECT mana_balance FROM wizards WHERE id = $1",
		req.WizardId).Scan(&currentBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "Wizard not found")
		}
		s.logger.Error("Failed to get current mana balance", "error", err)
		return nil, status.Error(codes.Internal, "Failed to update mana balance")
	}

	// Check for insufficient funds if this is a deduction
	newBalance := currentBalance + req.Amount
	if newBalance < 0 {
		return nil, status.Error(codes.FailedPrecondition, "Insufficient mana balance")
	}

	// Update the balance
	_, err = tx.ExecContext(ctx,
		"UPDATE wizards SET mana_balance = $1 WHERE id = $2",
		newBalance, req.WizardId)
	if err != nil {
		s.logger.Error("Failed to update mana balance", "error", err)
		return nil, status.Error(codes.Internal, "Failed to update mana balance")
	}

	// Create activity log if reason is provided
	if req.Reason != "" {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO activity_logs (user_id, wizard_id, activity_type, activity_description, metadata) 
			 SELECT user_id, id, 'mana_update', $2,
			        json_build_object('amount', $3, 'old_balance', $4, 'new_balance', $5, 'reason', $2)
			 FROM wizards WHERE id = $1`,
			req.WizardId, req.Reason, req.Amount, currentBalance, newBalance)
		if err != nil {
			s.logger.Error("Failed to create activity log", "error", err)
			// Don't fail the transaction for activity log issues
		}
	}

	if err = tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction", "error", err)
		return nil, status.Error(codes.Internal, "Failed to update mana balance")
	}

	return &pb.UpdateManaBalanceResponse{
		NewBalance: newBalance,
		Success:    true,
	}, nil
}

func (s *WizardServiceImpl) TransferMana(ctx context.Context, req *pb.TransferManaRequest) (*pb.TransferManaResponse, error) {
	if req.Amount <= 0 {
		return &pb.TransferManaResponse{
			Success: false,
			Message: "Transfer amount must be positive",
		}, nil
	}

	if req.FromWizardId == req.ToWizardId {
		return &pb.TransferManaResponse{
			Success: false,
			Message: "Cannot transfer mana to yourself",
		}, nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to start transaction")
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			s.logger.Error("Failed to rollback transaction", "error", err)
		}
	}()

	// Check sender's balance
	var senderBalance int64
	err = tx.QueryRowContext(ctx,
		"SELECT mana_balance FROM wizards WHERE id = $1",
		req.FromWizardId).Scan(&senderBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.TransferManaResponse{
				Success: false,
				Message: "Sender wizard not found",
			}, nil
		}
		s.logger.Error("Failed to get sender balance", "error", err)
		return nil, status.Error(codes.Internal, "Failed to transfer mana")
	}

	if senderBalance < req.Amount {
		return &pb.TransferManaResponse{
			Success: false,
			Message: "Insufficient mana balance",
		}, nil
	}

	// Check receiver exists
	var receiverBalance int64
	err = tx.QueryRowContext(ctx,
		"SELECT mana_balance FROM wizards WHERE id = $1",
		req.ToWizardId).Scan(&receiverBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.TransferManaResponse{
				Success: false,
				Message: "Receiver wizard not found",
			}, nil
		}
		s.logger.Error("Failed to get receiver balance", "error", err)
		return nil, status.Error(codes.Internal, "Failed to transfer mana")
	}

	// Update balances
	_, err = tx.ExecContext(ctx,
		"UPDATE wizards SET mana_balance = mana_balance - $1 WHERE id = $2",
		req.Amount, req.FromWizardId)
	if err != nil {
		s.logger.Error("Failed to update sender balance", "error", err)
		return nil, status.Error(codes.Internal, "Failed to transfer mana")
	}

	_, err = tx.ExecContext(ctx,
		"UPDATE wizards SET mana_balance = mana_balance + $1 WHERE id = $2",
		req.Amount, req.ToWizardId)
	if err != nil {
		s.logger.Error("Failed to update receiver balance", "error", err)
		return nil, status.Error(codes.Internal, "Failed to transfer mana")
	}

	// Create activity logs for both wizards
	reason := req.Reason
	if reason == "" {
		reason = "Mana transfer"
	}

	// Log for sender
	_, err = tx.ExecContext(ctx,
		`INSERT INTO activity_logs (user_id, wizard_id, activity_type, activity_description, metadata) 
		 SELECT user_id, id, 'mana_transfer_sent', $2,
		        json_build_object('amount', $3, 'to_wizard_id', $4, 'reason', $2)
		 FROM wizards WHERE id = $1`,
		req.FromWizardId, fmt.Sprintf("Sent %d mana: %s", req.Amount, reason), req.Amount, req.ToWizardId)
	if err != nil {
		s.logger.Error("Failed to create sender activity log", "error", err)
	}

	// Log for receiver
	_, err = tx.ExecContext(ctx,
		`INSERT INTO activity_logs (user_id, wizard_id, activity_type, activity_description, metadata) 
		 SELECT user_id, id, 'mana_transfer_received', $2,
		        json_build_object('amount', $3, 'from_wizard_id', $4, 'reason', $2)
		 FROM wizards WHERE id = $1`,
		req.ToWizardId, fmt.Sprintf("Received %d mana: %s", req.Amount, reason), req.Amount, req.FromWizardId)
	if err != nil {
		s.logger.Error("Failed to create receiver activity log", "error", err)
	}

	if err = tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction", "error", err)
		return nil, status.Error(codes.Internal, "Failed to transfer mana")
	}

	return &pb.TransferManaResponse{
		Success: true,
		Message: "Mana transfer completed successfully",
	}, nil
}
