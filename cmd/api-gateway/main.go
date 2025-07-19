package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/tectix/mysticfunds/pkg/config"
	"github.com/tectix/mysticfunds/pkg/logger"
	authpb "github.com/tectix/mysticfunds/proto/auth"
	manapb "github.com/tectix/mysticfunds/proto/mana"
	wizardpb "github.com/tectix/mysticfunds/proto/wizard"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Gateway struct {
	authClient   authpb.AuthServiceClient
	wizardClient wizardpb.WizardServiceClient
	manaClient   manapb.ManaServiceClient
	logger       logger.Logger
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	logger := logger.NewLogger(cfg.LogLevel)

	// Connect to gRPC services
	authAddr := cfg.GetString("AUTH_SERVICE_ADDR", "localhost:50051")
	authConn, err := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Failed to connect to auth service", "error", err, "address", authAddr)
	}
	defer authConn.Close()

	wizardAddr := cfg.GetString("WIZARD_SERVICE_ADDR", "localhost:50052")
	wizardConn, err := grpc.Dial(wizardAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Failed to connect to wizard service", "error", err, "address", wizardAddr)
	}
	defer wizardConn.Close()

	manaAddr := cfg.GetString("MANA_SERVICE_ADDR", "localhost:50053")
	manaConn, err := grpc.Dial(manaAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Failed to connect to mana service", "error", err, "address", manaAddr)
	}
	defer manaConn.Close()

	gateway := &Gateway{
		authClient:   authpb.NewAuthServiceClient(authConn),
		wizardClient: wizardpb.NewWizardServiceClient(wizardConn),
		manaClient:   manapb.NewManaServiceClient(manaConn),
		logger:       logger,
	}

	// Setup routes
	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("/api/auth/register", corsMiddleware(gateway.handleRegister))
	mux.HandleFunc("/api/auth/login", corsMiddleware(gateway.handleLogin))
	mux.HandleFunc("/api/auth/refresh", corsMiddleware(gateway.handleRefreshToken))
	mux.HandleFunc("/api/auth/logout", corsMiddleware(gateway.handleLogout))

	// Wizard routes
	mux.HandleFunc("/api/wizards", corsMiddleware(gateway.authMiddleware(gateway.handleWizards)))
	mux.HandleFunc("/api/wizards/explore", corsMiddleware(gateway.authMiddleware(gateway.handleExploreWizards)))
	mux.HandleFunc("/api/wizards/", corsMiddleware(gateway.authMiddleware(gateway.handleWizardByID)))

	// Mana routes
	mux.HandleFunc("/api/mana/balance/", corsMiddleware(gateway.authMiddleware(gateway.handleManaBalance)))
	mux.HandleFunc("/api/mana/transfer", corsMiddleware(gateway.authMiddleware(gateway.handleManaTransfer)))
	mux.HandleFunc("/api/mana/transactions/", corsMiddleware(gateway.authMiddleware(gateway.handleManaTransactions)))
	mux.HandleFunc("/api/mana/investments", corsMiddleware(gateway.authMiddleware(gateway.handleInvestments)))
	mux.HandleFunc("/api/mana/investment-types", corsMiddleware(gateway.authMiddleware(gateway.handleInvestmentTypes)))

	// Job routes
	mux.HandleFunc("/api/jobs", corsMiddleware(gateway.authMiddleware(gateway.handleJobs)))
	mux.HandleFunc("/api/jobs/", corsMiddleware(gateway.authMiddleware(gateway.handleJobByID)))
	mux.HandleFunc("/api/jobs/assign", corsMiddleware(gateway.authMiddleware(gateway.handleJobAssignment)))
	mux.HandleFunc("/api/jobs/assignments", corsMiddleware(gateway.authMiddleware(gateway.handleJobAssignments)))
	mux.HandleFunc("/api/jobs/assignments/", corsMiddleware(gateway.authMiddleware(gateway.handleJobAssignmentByID)))
	mux.HandleFunc("/api/jobs/assignments/cancel", corsMiddleware(gateway.authMiddleware(gateway.handleJobAssignmentCancel)))

	// Progress and activity routes
	mux.HandleFunc("/api/activities", corsMiddleware(gateway.authMiddleware(gateway.handleActivities)))
	mux.HandleFunc("/api/jobs/progress/", corsMiddleware(gateway.authMiddleware(gateway.handleJobProgress)))

	// Realm routes
	mux.HandleFunc("/api/realms", corsMiddleware(gateway.authMiddleware(gateway.handleRealms)))

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Serve static files - find the web directory
	cwd, _ := os.Getwd()
	logger.Info("Current working directory", "cwd", cwd)

	var webDir string

	// Try different possible paths based on where we're running from
	possiblePaths := []string{
		"./web/",        // If running from project root
		"../../web/",    // If running from cmd/api-gateway/
		"../../../web/", // If running from cmd/api-gateway/bin/
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(filepath.Join(path, "index.html")); err == nil {
			webDir = path
			break
		}
	}

	// If still not found, use absolute path relative to this source file
	if webDir == "" {
		// Get the directory where this main.go file is located
		_, filename, _, _ := runtime.Caller(0)
		sourceDir := filepath.Dir(filename)
		projectRoot := filepath.Dir(filepath.Dir(sourceDir)) // Go up 2 levels from cmd/api-gateway/
		webDir = filepath.Join(projectRoot, "web")
	}

	logger.Info("Using web directory", "path", webDir)

	// Verify the directory exists and has index.html
	if _, err := os.Stat(filepath.Join(webDir, "index.html")); err != nil {
		logger.Error("index.html not found in web directory", "error", err, "path", webDir)
	} else {
		logger.Info("Successfully found index.html", "path", filepath.Join(webDir, "index.html"))
	}

	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	port := cfg.GetString("HTTP_PORT", "8080")
	logger.Info("API Gateway starting", "port", port)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (g *Gateway) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := tokenParts[1]
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		resp, err := g.authClient.ValidateToken(ctx, &authpb.ValidateTokenRequest{
			Token: token,
		})
		if err != nil || !resp.IsValid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user ID to request context
		ctx = context.WithValue(r.Context(), "user_id", resp.UserId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

func (g *Gateway) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req authpb.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.authClient.Register(ctx, &req)
	if err != nil {
		g.logger.Error("Register failed", "error", err)

		// Handle specific gRPC status codes
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.AlreadyExists:
				http.Error(w, st.Message(), http.StatusConflict)
				return
			case codes.InvalidArgument:
				http.Error(w, st.Message(), http.StatusBadRequest)
				return
			}
		}

		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req authpb.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.authClient.Login(ctx, &req)
	if err != nil {
		g.logger.Error("Login failed", "error", err)
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req authpb.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.authClient.RefreshToken(ctx, &req)
	if err != nil {
		g.logger.Error("Token refresh failed", "error", err)
		http.Error(w, "Token refresh failed", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req authpb.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.authClient.Logout(ctx, &req)
	if err != nil {
		g.logger.Error("Logout failed", "error", err)
		http.Error(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleWizards(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch r.Method {
	case http.MethodGet:
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
		pageNumber, _ := strconv.Atoi(r.URL.Query().Get("page_number"))
		realm := r.URL.Query().Get("realm")

		if pageSize <= 0 {
			pageSize = 10
		}
		if pageNumber <= 0 {
			pageNumber = 1
		}

		// Get user ID from context for filtering user's own wizards
		userID := r.Context().Value("user_id").(int64)

		resp, err := g.wizardClient.ListWizards(ctx, &wizardpb.ListWizardsRequest{
			PageSize:   int32(pageSize),
			PageNumber: int32(pageNumber),
			UserId:     userID, // Filter to show only user's wizards
			Realm:      realm,  // Optional realm filter for exploration
		})
		if err != nil {
			g.logger.Error("List wizards failed", "error", err)
			http.Error(w, "Failed to list wizards", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	case http.MethodPost:
		var req wizardpb.CreateWizardRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Get user ID from context
		userID := r.Context().Value("user_id").(int64)
		req.UserId = userID

		resp, err := g.wizardClient.CreateWizard(ctx, &req)
		if err != nil {
			g.logger.Error("Create wizard failed", "error", err)
			http.Error(w, "Failed to create wizard", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (g *Gateway) handleWizardByID(w http.ResponseWriter, r *http.Request) {
	// Extract wizard ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/wizards/")
	wizardID, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid wizard ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch r.Method {
	case http.MethodGet:
		resp, err := g.wizardClient.GetWizard(ctx, &wizardpb.GetWizardRequest{
			Id: wizardID,
		})
		if err != nil {
			g.logger.Error("Get wizard failed", "error", err)
			http.Error(w, "Failed to get wizard", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	case http.MethodPut:
		var req wizardpb.UpdateWizardRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		req.Id = wizardID

		resp, err := g.wizardClient.UpdateWizard(ctx, &req)
		if err != nil {
			g.logger.Error("Update wizard failed", "error", err)
			http.Error(w, "Failed to update wizard", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	case http.MethodDelete:
		resp, err := g.wizardClient.DeleteWizard(ctx, &wizardpb.DeleteWizardRequest{
			Id: wizardID,
		})
		if err != nil {
			g.logger.Error("Delete wizard failed", "error", err)
			http.Error(w, "Failed to delete wizard", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (g *Gateway) handleManaBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract wizard ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/mana/balance/")
	wizardID, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid wizard ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get wizard info (including mana balance) from wizard service
	wizard, err := g.wizardClient.GetWizard(ctx, &wizardpb.GetWizardRequest{
		Id: wizardID,
	})
	if err != nil {
		g.logger.Error("Get mana balance failed", "error", err)
		http.Error(w, "Failed to get mana balance", http.StatusInternalServerError)
		return
	}

	// Create response in mana service format for compatibility
	resp := &manapb.GetManaBalanceResponse{
		Balance: wizard.ManaBalance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleManaTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req manapb.TransferManaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.manaClient.TransferMana(ctx, &req)
	if err != nil {
		g.logger.Error("Transfer mana failed", "error", err)
		http.Error(w, "Failed to transfer mana", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleManaTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract wizard ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/mana/transactions/")
	wizardID, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid wizard ID", http.StatusBadRequest)
		return
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	pageNumber, _ := strconv.Atoi(r.URL.Query().Get("page_number"))

	if pageSize <= 0 {
		pageSize = 10
	}
	if pageNumber <= 0 {
		pageNumber = 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.manaClient.ListTransactions(ctx, &manapb.ListTransactionsRequest{
		WizardId:   wizardID,
		PageSize:   int32(pageSize),
		PageNumber: int32(pageNumber),
	})
	if err != nil {
		g.logger.Error("List transactions failed", "error", err)
		http.Error(w, "Failed to list transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleInvestments(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch r.Method {
	case http.MethodGet:
		wizardIDStr := r.URL.Query().Get("wizard_id")
		wizardID, err := strconv.ParseInt(wizardIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid wizard ID", http.StatusBadRequest)
			return
		}

		status := r.URL.Query().Get("status")

		resp, err := g.manaClient.GetInvestments(ctx, &manapb.GetInvestmentsRequest{
			WizardId: wizardID,
			Status:   status,
		})
		if err != nil {
			g.logger.Error("Get investments failed", "error", err)
			http.Error(w, "Failed to get investments", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	case http.MethodPost:
		var req manapb.CreateInvestmentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		resp, err := g.manaClient.CreateInvestment(ctx, &req)
		if err != nil {
			g.logger.Error("Create investment failed", "error", err)
			http.Error(w, "Failed to create investment", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (g *Gateway) handleInvestmentTypes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	minAmount, _ := strconv.ParseInt(r.URL.Query().Get("min_amount"), 10, 64)
	maxAmount, _ := strconv.ParseInt(r.URL.Query().Get("max_amount"), 10, 64)
	riskLevel, _ := strconv.Atoi(r.URL.Query().Get("risk_level"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.manaClient.GetInvestmentTypes(ctx, &manapb.GetInvestmentTypesRequest{
		MinAmount: minAmount,
		MaxAmount: maxAmount,
		RiskLevel: int32(riskLevel),
	})
	if err != nil {
		g.logger.Error("Get investment types failed", "error", err)
		http.Error(w, "Failed to get investment types", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleExploreWizards(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	pageNumber, _ := strconv.Atoi(r.URL.Query().Get("page_number"))
	realm := r.URL.Query().Get("realm")

	if pageSize <= 0 {
		pageSize = 10
	}
	if pageNumber <= 0 {
		pageNumber = 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// For exploration, don't filter by user ID to show all wizards
	resp, err := g.wizardClient.ListWizards(ctx, &wizardpb.ListWizardsRequest{
		PageSize:   int32(pageSize),
		PageNumber: int32(pageNumber),
		Realm:      realm, // Filter by realm if provided
	})
	if err != nil {
		g.logger.Error("Explore wizards failed", "error", err)
		http.Error(w, "Failed to explore wizards", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleJobs(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch r.Method {
	case http.MethodGet:
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
		pageNumber, _ := strconv.Atoi(r.URL.Query().Get("page_number"))
		realm := r.URL.Query().Get("realm")
		element := r.URL.Query().Get("element")
		difficulty := r.URL.Query().Get("difficulty")
		onlyActive := r.URL.Query().Get("only_active") == "true"

		if pageSize <= 0 {
			pageSize = 10
		}
		if pageNumber <= 0 {
			pageNumber = 1
		}

		resp, err := g.wizardClient.ListJobs(ctx, &wizardpb.ListJobsRequest{
			PageSize:   int32(pageSize),
			PageNumber: int32(pageNumber),
			Realm:      realm,
			Element:    element,
			Difficulty: difficulty,
			OnlyActive: onlyActive,
		})
		if err != nil {
			g.logger.Error("List jobs failed", "error", err)
			http.Error(w, "Failed to list jobs", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	case http.MethodPost:
		var req wizardpb.CreateJobRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		resp, err := g.wizardClient.CreateJob(ctx, &req)
		if err != nil {
			g.logger.Error("Create job failed", "error", err)
			http.Error(w, "Failed to create job", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (g *Gateway) handleJobByID(w http.ResponseWriter, r *http.Request) {
	// Extract job ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/jobs/")
	jobID, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch r.Method {
	case http.MethodGet:
		resp, err := g.wizardClient.GetJob(ctx, &wizardpb.GetJobRequest{
			Id: jobID,
		})
		if err != nil {
			g.logger.Error("Get job failed", "error", err)
			http.Error(w, "Failed to get job", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	case http.MethodPut:
		var req wizardpb.UpdateJobRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		req.Id = jobID

		resp, err := g.wizardClient.UpdateJob(ctx, &req)
		if err != nil {
			g.logger.Error("Update job failed", "error", err)
			http.Error(w, "Failed to update job", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	case http.MethodDelete:
		resp, err := g.wizardClient.DeleteJob(ctx, &wizardpb.DeleteJobRequest{
			Id: jobID,
		})
		if err != nil {
			g.logger.Error("Delete job failed", "error", err)
			http.Error(w, "Failed to delete job", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (g *Gateway) handleJobAssignment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req wizardpb.AssignWizardToJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.wizardClient.AssignWizardToJob(ctx, &req)
	if err != nil {
		g.logger.Error("Assign wizard to job failed", "error", err)
		http.Error(w, "Failed to assign wizard to job", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleJobAssignments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	wizardID, _ := strconv.ParseInt(r.URL.Query().Get("wizard_id"), 10, 64)
	jobID, _ := strconv.ParseInt(r.URL.Query().Get("job_id"), 10, 64)
	status := r.URL.Query().Get("status")
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	pageNumber, _ := strconv.Atoi(r.URL.Query().Get("page_number"))

	if pageSize <= 0 {
		pageSize = 10
	}
	if pageNumber <= 0 {
		pageNumber = 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.wizardClient.GetJobAssignments(ctx, &wizardpb.GetJobAssignmentsRequest{
		WizardId:   wizardID,
		JobId:      jobID,
		Status:     status,
		PageSize:   int32(pageSize),
		PageNumber: int32(pageNumber),
	})
	if err != nil {
		g.logger.Error("Get job assignments failed", "error", err)
		http.Error(w, "Failed to get job assignments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleJobAssignmentByID(w http.ResponseWriter, r *http.Request) {
	// Extract assignment ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/jobs/assignments/")
	assignmentID, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid assignment ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch r.Method {
	case http.MethodPut:
		// Check if this is a complete or cancel action based on request body
		var actionReq struct {
			Action string `json:"action"`
			Reason string `json:"reason,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&actionReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if actionReq.Action == "complete" {
			resp, err := g.wizardClient.CompleteJobAssignment(ctx, &wizardpb.CompleteJobAssignmentRequest{
				AssignmentId: assignmentID,
			})
			if err != nil {
				g.logger.Error("Complete job assignment failed", "error", err)
				http.Error(w, "Failed to complete job assignment", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)

		} else if actionReq.Action == "cancel" {
			resp, err := g.wizardClient.CancelJobAssignment(ctx, &wizardpb.CancelJobAssignmentRequest{
				AssignmentId: assignmentID,
				Reason:       actionReq.Reason,
			})
			if err != nil {
				g.logger.Error("Cancel job assignment failed", "error", err)
				http.Error(w, "Failed to cancel job assignment", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)

		} else {
			http.Error(w, "Invalid action. Use 'complete' or 'cancel'", http.StatusBadRequest)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (g *Gateway) handleActivities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int64)
	wizardIDStr := r.URL.Query().Get("wizard_id")
	activityType := r.URL.Query().Get("activity_type")
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	pageNumber, _ := strconv.Atoi(r.URL.Query().Get("page_number"))

	if pageSize <= 0 {
		pageSize = 20
	}
	if pageNumber <= 0 {
		pageNumber = 1
	}

	var wizardID int64
	if wizardIDStr != "" {
		wizardID, _ = strconv.ParseInt(wizardIDStr, 10, 64)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.wizardClient.GetActivities(ctx, &wizardpb.GetActivitiesRequest{
		UserId:       userID,
		WizardId:     wizardID,
		ActivityType: activityType,
		PageSize:     int32(pageSize),
		PageNumber:   int32(pageNumber),
	})
	if err != nil {
		g.logger.Error("Get activities failed", "error", err)
		http.Error(w, "Failed to get activities", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleJobProgress(w http.ResponseWriter, r *http.Request) {
	// Extract assignment ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/jobs/progress/")
	assignmentID, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid assignment ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch r.Method {
	case http.MethodGet:
		resp, err := g.wizardClient.GetJobProgress(ctx, &wizardpb.GetJobProgressRequest{
			AssignmentId: assignmentID,
		})
		if err != nil {
			g.logger.Error("Get job progress failed", "error", err)
			http.Error(w, "Failed to get job progress", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	case http.MethodPut:
		var req wizardpb.UpdateJobProgressRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		req.AssignmentId = assignmentID

		resp, err := g.wizardClient.UpdateJobProgress(ctx, &req)
		if err != nil {
			g.logger.Error("Update job progress failed", "error", err)
			http.Error(w, "Failed to update job progress", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (g *Gateway) handleJobAssignmentCancel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AssignmentId int64  `json:"assignment_id"`
		Reason       string `json:"reason,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.AssignmentId <= 0 {
		http.Error(w, "Invalid assignment ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.wizardClient.CancelJobAssignment(ctx, &wizardpb.CancelJobAssignmentRequest{
		AssignmentId: req.AssignmentId,
		Reason:       req.Reason,
	})
	if err != nil {
		g.logger.Error("Cancel job assignment failed", "error", err)
		http.Error(w, "Failed to cancel job assignment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleRealms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.wizardClient.GetRealms(ctx, &wizardpb.GetRealmsRequest{})
	if err != nil {
		g.logger.Error("Get realms failed", "error", err)
		http.Error(w, "Failed to get realms", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
