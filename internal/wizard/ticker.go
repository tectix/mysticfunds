package wizard

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/tectix/mysticfunds/pkg/logger"
	pb "github.com/tectix/mysticfunds/proto/wizard"
)

// JobTicker handles real-time job progress updates
type JobTicker struct {
	db          *sql.DB
	logger      logger.Logger
	tickerMutex sync.RWMutex
	running     bool
	stopCh      chan struct{}
	service     *WizardServiceImpl
}

// NewJobTicker creates a new job ticker instance
func NewJobTicker(db *sql.DB, logger logger.Logger, service *WizardServiceImpl) *JobTicker {
	return &JobTicker{
		db:      db,
		logger:  logger,
		service: service,
		stopCh:  make(chan struct{}),
	}
}

// Start begins the job ticker that updates progress every 5 seconds
func (jt *JobTicker) Start() {
	jt.tickerMutex.Lock()
	defer jt.tickerMutex.Unlock()

	if jt.running {
		jt.logger.Info("Job ticker already running")
		return
	}

	jt.running = true
	jt.logger.Info("Starting job ticker")

	go jt.tickerLoop()
}

// Stop halts the job ticker
func (jt *JobTicker) Stop() {
	jt.tickerMutex.Lock()
	defer jt.tickerMutex.Unlock()

	if !jt.running {
		return
	}

	jt.running = false
	close(jt.stopCh)
	jt.logger.Info("Job ticker stopped")
}

// tickerLoop runs the main ticker loop
func (jt *JobTicker) tickerLoop() {
	ticker := time.NewTicker(5 * time.Second) // Tick every 5 seconds
	defer ticker.Stop()

	// Initial tick
	jt.processTick()

	for {
		select {
		case <-ticker.C:
			jt.processTick()
		case <-jt.stopCh:
			jt.logger.Info("Job ticker loop terminated")
			return
		}
	}
}

// processTick handles a single tick of the job ticker
func (jt *JobTicker) processTick() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Update all active job progress
	err := jt.updateAllJobProgress(ctx)
	if err != nil {
		jt.logger.Error("Failed to update job progress", "error", err)
	}

	// Auto-complete finished jobs
	err = jt.autoCompleteFinishedJobs(ctx)
	if err != nil {
		jt.logger.Error("Failed to auto-complete finished jobs", "error", err)
	}
}

// updateAllJobProgress updates progress for all active jobs based on elapsed time
func (jt *JobTicker) updateAllJobProgress(ctx context.Context) error {
	// Get all active job progress records with their job details
	query := `
		SELECT 
			jp.id,
			jp.assignment_id,
			jp.actual_start_time,
			jp.expected_end_time,
			jp.progress_percentage,
			jp.is_active,
			j.duration_minutes,
			ja.status
		FROM job_progress jp
		JOIN job_assignments ja ON jp.assignment_id = ja.id
		JOIN jobs j ON ja.job_id = j.id
		WHERE jp.is_active = true 
		AND ja.status IN ('assigned', 'in_progress')
		AND jp.actual_start_time IS NOT NULL
		AND jp.expected_end_time IS NOT NULL
	`

	rows, err := jt.db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var updates []struct {
		progressID         int64
		assignmentID       int64
		newProgress        int32
		timeWorkedMinutes  int32
		shouldMarkComplete bool
		status             string
	}

	now := time.Now()

	for rows.Next() {
		var progressID, assignmentID int64
		var actualStartTime, expectedEndTime time.Time
		var currentProgress int32
		var isActive bool
		var durationMinutes int32
		var status string

		err := rows.Scan(
			&progressID,
			&assignmentID,
			&actualStartTime,
			&expectedEndTime,
			&currentProgress,
			&isActive,
			&durationMinutes,
			&status,
		)
		if err != nil {
			jt.logger.Error("Failed to scan progress row", "error", err)
			continue
		}

		// Calculate real-time progress
		elapsed := now.Sub(actualStartTime)
		totalDuration := expectedEndTime.Sub(actualStartTime)

		var newProgress int32
		if elapsed <= 0 {
			newProgress = 0
		} else if elapsed >= totalDuration {
			newProgress = 100
		} else {
			progressFloat := float64(elapsed) / float64(totalDuration) * 100
			newProgress = int32(progressFloat)
		}

		// Calculate time worked in minutes
		timeWorkedMinutes := int32(elapsed.Minutes())
		if timeWorkedMinutes < 0 {
			timeWorkedMinutes = 0
		}

		// Only update if progress has changed and increased
		if newProgress > currentProgress {
			updates = append(updates, struct {
				progressID         int64
				assignmentID       int64
				newProgress        int32
				timeWorkedMinutes  int32
				shouldMarkComplete bool
				status             string
			}{
				progressID:         progressID,
				assignmentID:       assignmentID,
				newProgress:        newProgress,
				timeWorkedMinutes:  timeWorkedMinutes,
				shouldMarkComplete: newProgress >= 100,
				status:             status,
			})
		}
	}

	// Apply all updates in a transaction
	if len(updates) > 0 {
		return jt.applyProgressUpdates(ctx, updates)
	}

	return nil
}

// applyProgressUpdates applies all progress updates in a single transaction
func (jt *JobTicker) applyProgressUpdates(ctx context.Context, updates []struct {
	progressID         int64
	assignmentID       int64
	newProgress        int32
	timeWorkedMinutes  int32
	shouldMarkComplete bool
	status             string
}) error {
	tx, err := jt.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			jt.logger.Error("Failed to rollback progress update transaction", "error", err)
		}
	}()

	for _, update := range updates {
		// Update progress record
		_, err = tx.ExecContext(ctx, `
			UPDATE job_progress 
			SET progress_percentage = $1, 
			    time_worked_minutes = $2, 
			    last_updated_at = CURRENT_TIMESTAMP,
			    last_tick_time = CURRENT_TIMESTAMP
			WHERE id = $3`,
			update.newProgress, update.timeWorkedMinutes, update.progressID)
		if err != nil {
			jt.logger.Error("Failed to update job progress", "progress_id", update.progressID, "error", err)
			continue
		}

		// Update assignment status to 'in_progress' if it's still 'assigned'
		if update.newProgress > 0 && update.status == "assigned" {
			_, err = tx.ExecContext(ctx, `
				UPDATE job_assignments 
				SET status = 'in_progress', 
				    started_at = COALESCE(started_at, CURRENT_TIMESTAMP) 
				WHERE id = $1`,
				update.assignmentID)
			if err != nil {
				jt.logger.Error("Failed to update assignment status", "assignment_id", update.assignmentID, "error", err)
				continue
			}
		}

		jt.logger.Debug("Updated job progress",
			"assignment_id", update.assignmentID,
			"progress", update.newProgress,
			"time_worked", update.timeWorkedMinutes)
	}

	return tx.Commit()
}

// autoCompleteFinishedJobs automatically completes jobs that have reached 100% progress
func (jt *JobTicker) autoCompleteFinishedJobs(ctx context.Context) error {
	// Find jobs that are at 100% progress but not yet completed
	query := `
		SELECT jp.assignment_id
		FROM job_progress jp
		JOIN job_assignments ja ON jp.assignment_id = ja.id
		WHERE jp.progress_percentage >= 100
		AND ja.status IN ('assigned', 'in_progress')
		AND jp.is_active = true
	`

	rows, err := jt.db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var assignmentIDs []int64
	for rows.Next() {
		var assignmentID int64
		if err := rows.Scan(&assignmentID); err != nil {
			jt.logger.Error("Failed to scan assignment ID", "error", err)
			continue
		}
		assignmentIDs = append(assignmentIDs, assignmentID)
	}

	// Auto-complete each finished job
	for _, assignmentID := range assignmentIDs {
		jt.logger.Info("Auto-completing finished job", "assignment_id", assignmentID)

		_, err := jt.service.CompleteJobAssignment(ctx, &pb.CompleteJobAssignmentRequest{
			AssignmentId: assignmentID,
		})
		if err != nil {
			jt.logger.Error("Failed to auto-complete job assignment",
				"assignment_id", assignmentID, "error", err)
		}
	}

	return nil
}

// IsRunning returns whether the ticker is currently running
func (jt *JobTicker) IsRunning() bool {
	jt.tickerMutex.RLock()
	defer jt.tickerMutex.RUnlock()
	return jt.running
}
