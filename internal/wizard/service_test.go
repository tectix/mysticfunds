package wizard

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Alinoureddine1/mysticfunds/pkg/config"
	"github.com/Alinoureddine1/mysticfunds/pkg/logger"
	pb "github.com/Alinoureddine1/mysticfunds/proto/wizard"
)

func setupTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *WizardServiceImpl) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}

	cfg := &config.Config{
		JWTSecret: "test_secret",
	}
	log := logger.NewLogger("debug")

	return db, mock, NewWizardServiceImpl(db, cfg, log)
}

func TestCreateWizard(t *testing.T) {
	db, mock, service := setupTest(t)
	defer db.Close()

	mock.ExpectQuery("INSERT INTO wizards").
		WithArgs(1, "TestWizard", "TestRealm", "Fire").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery("SELECT (.+) FROM wizards").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "realm", "element", "mana_balance", "created_at", "updated_at", "guild_id", "guild_name"}).
			AddRow(1, 1, "TestWizard", "TestRealm", "Fire", 0, nil, nil, nil, nil))

	resp, err := service.CreateWizard(context.Background(), &pb.CreateWizardRequest{
		UserId:  1,
		Name:    "TestWizard",
		Realm:   "TestRealm",
		Element: "Fire",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Id)
	assert.Equal(t, "TestWizard", resp.Name)
	assert.Equal(t, "TestRealm", resp.Realm)
	assert.Equal(t, "Fire", resp.Element)

	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestGetWizard(t *testing.T) {
	db, mock, service := setupTest(t)
	defer db.Close()

	createdAt := time.Now().Add(-1 * time.Hour)
	updatedAt := time.Now()

	mock.ExpectQuery("SELECT (.+) FROM wizards").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "realm", "element", "mana_balance", "created_at", "updated_at", "guild_id", "guild_name"}).
			AddRow(1, 1, "TestWizard", "TestRealm", "Fire", 100, createdAt, updatedAt, 1, "TestGuild"))

	wizard, err := service.GetWizard(context.Background(), &pb.GetWizardRequest{Id: 1})

	assert.NoError(t, err)
	assert.NotNil(t, wizard)
	assert.Equal(t, int64(1), wizard.Id)
	assert.Equal(t, int64(1), wizard.UserId)
	assert.Equal(t, "TestWizard", wizard.Name)
	assert.Equal(t, "TestRealm", wizard.Realm)
	assert.Equal(t, "Fire", wizard.Element)
	assert.Equal(t, int64(100), wizard.ManaBalance)
	assert.Equal(t, timestamppb.New(createdAt), wizard.CreatedAt)
	assert.Equal(t, timestamppb.New(updatedAt), wizard.UpdatedAt)
	assert.NotNil(t, wizard.Guild)
	assert.Equal(t, int64(1), wizard.Guild.Id)
	assert.Equal(t, "TestGuild", wizard.Guild.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateWizard(t *testing.T) {
	db, mock, service := setupTest(t)
	defer db.Close()

	mock.ExpectExec("UPDATE wizards SET").
		WithArgs("UpdatedWizard", "UpdatedRealm", "Water", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT (.+) FROM wizards").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "realm", "element", "mana_balance", "created_at", "updated_at", "guild_id", "guild_name"}).
			AddRow(1, 1, "UpdatedWizard", "UpdatedRealm", "Water", 100, time.Now(), time.Now(), nil, nil))

	wizard, err := service.UpdateWizard(context.Background(), &pb.UpdateWizardRequest{
		Id:      1,
		Name:    "UpdatedWizard",
		Realm:   "UpdatedRealm",
		Element: "Water",
	})

	assert.NoError(t, err)
	assert.NotNil(t, wizard)
	assert.Equal(t, int64(1), wizard.Id)
	assert.Equal(t, "UpdatedWizard", wizard.Name)
	assert.Equal(t, "UpdatedRealm", wizard.Realm)
	assert.Equal(t, "Water", wizard.Element)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListWizards(t *testing.T) {
	db, mock, service := setupTest(t)
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM wizards").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "realm", "element", "mana_balance", "created_at", "updated_at", "guild_id", "guild_name"}).
			AddRow(1, 1, "Wizard1", "Realm1", "Fire", 100, time.Now(), time.Now(), 1, "Guild1").
			AddRow(2, 2, "Wizard2", "Realm2", "Water", 200, time.Now(), time.Now(), 2, "Guild2"))

	mock.ExpectQuery("SELECT COUNT").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	response, err := service.ListWizards(context.Background(), &pb.ListWizardsRequest{
		PageSize:   10,
		PageNumber: 1,
	})

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Wizards))
	assert.Equal(t, int32(2), response.TotalCount)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteWizard(t *testing.T) {
	db, mock, service := setupTest(t)
	defer db.Close()

	mock.ExpectExec("DELETE FROM wizards").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	response, err := service.DeleteWizard(context.Background(), &pb.DeleteWizardRequest{Id: 1})

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Success)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestJoinGuild(t *testing.T) {
	db, mock, service := setupTest(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id FROM guilds").
		WithArgs("TestGuild").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec("UPDATE wizards SET").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	mock.ExpectQuery("SELECT (.+) FROM wizards").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "realm", "element", "mana_balance", "created_at", "updated_at", "guild_id", "guild_name"}).
			AddRow(1, 1, "TestWizard", "TestRealm", "Fire", 100, time.Now(), time.Now(), 1, "TestGuild"))

	wizard, err := service.JoinGuild(context.Background(), &pb.JoinGuildRequest{
		WizardId:  1,
		GuildName: "TestGuild",
	})

	assert.NoError(t, err)
	assert.NotNil(t, wizard)
	assert.NotNil(t, wizard.Guild)
	assert.Equal(t, int64(1), wizard.Guild.Id)
	assert.Equal(t, "TestGuild", wizard.Guild.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLeaveGuild(t *testing.T) {
	db, mock, service := setupTest(t)
	defer db.Close()

	mock.ExpectExec("UPDATE wizards SET").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT (.+) FROM wizards").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "realm", "element", "mana_balance", "created_at", "updated_at", "guild_id", "guild_name"}).
			AddRow(1, 1, "TestWizard", "TestRealm", "Fire", 100, time.Now(), time.Now(), nil, nil))

	wizard, err := service.LeaveGuild(context.Background(), &pb.LeaveGuildRequest{WizardId: 1})

	assert.NoError(t, err)
	assert.NotNil(t, wizard)
	assert.Nil(t, wizard.Guild)

	assert.NoError(t, mock.ExpectationsWereMet())
}
