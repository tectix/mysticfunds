package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Alinoureddine1/mysticfunds/internal/mana"
	"github.com/Alinoureddine1/mysticfunds/pkg/config"
	"github.com/Alinoureddine1/mysticfunds/pkg/database"
	"github.com/Alinoureddine1/mysticfunds/pkg/logger"
	pb "github.com/Alinoureddine1/mysticfunds/proto/mana"
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

	// Initialize investment scheduler
	scheduler := mana.NewInvestmentScheduler(db, log)
	defer scheduler.Stop()

	// Initialize mana service with investment handling
	manaService := mana.NewManaServiceImpl(db, cfg, log, scheduler)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterManaServiceServer(grpcServer, manaService)

	// Start listening on configured port
	address := fmt.Sprintf(":%d", cfg.GRPCPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Failed to listen", "error", err)
	}

	// Start server in a goroutine
	go func() {
		log.Info("Starting Mana Service", "address", address)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve", "error", err)
		}
	}()

	// Start investment scheduler in a goroutine
	go scheduler.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Mana Service")
	grpcServer.GracefulStop()
}
