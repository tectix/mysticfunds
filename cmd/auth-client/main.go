package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Alinoureddine1/mysticfunds/proto/auth"
)

const (
	defaultAddress = "localhost:50051"
)

var (
	addr = flag.String("addr", defaultAddress, "the address to connect to")
)

func main() {
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAuthServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Register
	r, err := c.Register(ctx, &pb.RegisterRequest{Username: "testuser", Email: "test@example.com", Password: "password123"})
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}
	fmt.Printf("Register response: %+v\n", r)

	// Login
	l, err := c.Login(ctx, &pb.LoginRequest{Username: "testuser", Password: "password123"})
	if err != nil {
		log.Fatalf("could not login: %v", err)
	}
	fmt.Printf("Login response: %+v\n", l)

	// Validate Token
	v, err := c.ValidateToken(ctx, &pb.ValidateTokenRequest{Token: l.Token})
	if err != nil {
		log.Fatalf("could not validate token: %v", err)
	}
	fmt.Printf("Validate Token response: %+v\n", v)

	// Refresh Token
	rf, err := c.RefreshToken(ctx, &pb.RefreshTokenRequest{Token: l.Token})
	if err != nil {
		log.Fatalf("could not refresh token: %v", err)
	}
	fmt.Printf("Refresh Token response: %+v\n", rf)

	// Logout
	lo, err := c.Logout(ctx, &pb.LogoutRequest{Token: rf.Token})
	if err != nil {
		log.Fatalf("could not logout: %v", err)
	}
	fmt.Printf("Logout response: %+v\n", lo)
}
