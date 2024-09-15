package auth

import (
	"context"
	"database/sql"
	"testing"
	"time"

	jwtauth "github.com/Alinoureddine1/mysticfunds/pkg/auth"
	"github.com/Alinoureddine1/mysticfunds/pkg/config"
	"github.com/Alinoureddine1/mysticfunds/pkg/logger"
	pb "github.com/Alinoureddine1/mysticfunds/proto/auth"
	"github.com/DATA-DOG/go-sqlmock"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *AuthServiceImpl) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	cfg := &config.Config{
		JWTSecret: "test_secret",
		LogLevel:  "debug",
	}
	log := logger.NewLogger(cfg.LogLevel)

	service := NewAuthServiceImpl(db, cfg, log)
	t.Logf("Service created: %+v", service)
	return db, mock, service
}

func TestRegister(t *testing.T) {
	db, mock, service := setupTest(t)
	defer db.Close()

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("testuser", "test@example.com", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	resp, err := service.Register(context.Background(), &pb.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, int64(1), resp.UserId)
}

func TestLogin(t *testing.T) {
	db, mock, service := setupTest(t)
	defer db.Close()

	// Generate a bcrypt hash for "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	mock.ExpectQuery("SELECT id, password_hash FROM users WHERE username = \\$1").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "password_hash"}).
			AddRow(1, string(hashedPassword)))

	resp, err := service.Login(context.Background(), &pb.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})

	assert.NoError(t, err)
	if err == nil {
		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, int64(1), resp.UserId)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestValidateToken(t *testing.T) {
	_, _, service := setupTest(t)

	token, err := jwtauth.GenerateToken(1, service.cfg.JWTSecret, time.Hour*24)
	assert.NoError(t, err)

	resp, err := service.ValidateToken(context.Background(), &pb.ValidateTokenRequest{
		Token: token,
	})

	assert.NoError(t, err)
	assert.True(t, resp.IsValid)
	assert.Equal(t, int64(1), resp.UserId)
}

func TestRefreshToken(t *testing.T) {
	_, _, service := setupTest(t)

	initialToken, err := jwtauth.GenerateToken(1, service.cfg.JWTSecret, time.Hour*24)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	resp, err := service.RefreshToken(context.Background(), &pb.RefreshTokenRequest{
		Token: initialToken,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.NotEqual(t, initialToken, resp.Token, "Refreshed token should be different from the initial token")
	assert.Equal(t, int64(1), resp.UserId)

	claims, err := jwtauth.ValidateToken(resp.Token, service.cfg.JWTSecret)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), claims.UserID)

	expirationTime := time.Unix(claims.ExpiresAt, 0)
	assert.True(t, expirationTime.After(time.Now()), "New token should have a future expiration time")
}

func printTokenDetails(t *testing.T, tokenString string, label string) {
	token, _ := jwt.ParseWithClaims(tokenString, &jwtauth.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("test_secret"), nil
	})

	if claims, ok := token.Claims.(*jwtauth.JWTClaims); ok {
		t.Logf("%s details:", label)
		t.Logf("  User ID: %d", claims.UserID)
		t.Logf("  Issued At: %v", time.Unix(claims.IssuedAt, 0))
		t.Logf("  Expires At: %v", time.Unix(claims.ExpiresAt, 0))
	}
}
func TestLogout(t *testing.T) {
	_, _, service := setupTest(t)

	token, err := jwtauth.GenerateToken(1, service.cfg.JWTSecret, time.Hour*24)
	assert.NoError(t, err)

	resp, err := service.Logout(context.Background(), &pb.LogoutRequest{
		Token: token,
	})

	assert.NoError(t, err)
	assert.True(t, resp.Success)
}
