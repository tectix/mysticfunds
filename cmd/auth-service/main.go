package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Alinoureddine1/mysticfunds/internal/auth"
	"github.com/Alinoureddine1/mysticfunds/pkg/config"
	"github.com/Alinoureddine1/mysticfunds/pkg/database"
	"github.com/Alinoureddine1/mysticfunds/pkg/logger"
	pb "github.com/Alinoureddine1/mysticfunds/proto/auth"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	log := logger.NewLogger(cfg.LogLevel)

	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	authService := auth.NewAuthServiceImpl(db, cfg, log)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authService)

	address := fmt.Sprintf(":%d", cfg.GRPCPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Failed to listen", "error", err)
	}

	go func() {
		log.Info("Starting Auth Service", "address", address)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Auth Service")
	grpcServer.GracefulStop()
}
