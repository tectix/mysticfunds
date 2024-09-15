package auth

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	jwtauth "github.com/Alinoureddine1/mysticfunds/pkg/auth"
	"github.com/Alinoureddine1/mysticfunds/pkg/config"
	"github.com/Alinoureddine1/mysticfunds/pkg/logger"
	pb "github.com/Alinoureddine1/mysticfunds/proto/auth"
)

type AuthServiceImpl struct {
	db     *sql.DB
	cfg    *config.Config
	logger logger.Logger
	pb.UnimplementedAuthServiceServer
}

func NewAuthServiceImpl(db *sql.DB, cfg *config.Config, logger logger.Logger) *AuthServiceImpl {
	return &AuthServiceImpl{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", "error", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	var userId int64
	err = s.db.QueryRowContext(ctx,
		"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id",
		req.Username, req.Email, string(hashedPassword)).Scan(&userId)
	if err != nil {
		s.logger.Error("Failed to insert user", "error", err)
		return nil, status.Error(codes.Internal, "Failed to register user")
	}

	token, err := jwtauth.GenerateToken(userId, s.cfg.JWTSecret, time.Hour*24)
	if err != nil {
		s.logger.Error("Failed to generate JWT", "error", err)
		return nil, status.Error(codes.Internal, "Failed to generate token")
	}

	return &pb.AuthResponse{Token: token, UserId: userId}, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	var (
		id           int64
		passwordHash string
	)

	err := s.db.QueryRowContext(ctx,
		"SELECT id, password_hash FROM users WHERE username = $1",
		req.Username).Scan(&id, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		s.logger.Error("Failed to query user", "error", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	}

	token, err := jwtauth.GenerateToken(id, s.cfg.JWTSecret, time.Hour*24)
	if err != nil {
		s.logger.Error("Failed to generate JWT", "error", err)
		return nil, status.Error(codes.Internal, "Failed to generate token")
	}

	return &pb.AuthResponse{Token: token, UserId: id}, nil
}

func (s *AuthServiceImpl) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	claims, err := jwtauth.ValidateToken(req.Token, s.cfg.JWTSecret)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	newToken, err := jwtauth.GenerateToken(claims.UserID, s.cfg.JWTSecret, time.Hour*24)
	if err != nil {
		s.logger.Error("Failed to generate new JWT", "error", err)
		return nil, status.Error(codes.Internal, "Failed to generate new token")
	}

	return &pb.RefreshTokenResponse{Token: newToken, UserId: claims.UserID}, nil
}

func (s *AuthServiceImpl) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := jwtauth.ValidateToken(req.Token, s.cfg.JWTSecret)
	if err != nil {
		return &pb.ValidateTokenResponse{IsValid: false}, nil
	}

	return &pb.ValidateTokenResponse{
		IsValid: true,
		UserId:  claims.UserID,
	}, nil
}

func (s *AuthServiceImpl) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	_, err := jwtauth.ValidateToken(req.Token, s.cfg.JWTSecret)
	if err != nil {
		return &pb.LogoutResponse{Success: false}, nil
	}

	return &pb.LogoutResponse{Success: true}, nil
}
