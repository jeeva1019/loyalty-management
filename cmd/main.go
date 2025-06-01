package main

import (
	"fmt"
	"log"
	"loyality_points/config"
	"loyality_points/controllers"
	"loyality_points/helpers"
	middleware "loyality_points/middlewares"
	"loyality_points/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Set up log file for debugging
	logFile, err := helpers.SetupLogger()
	if err != nil {
		log.Fatalf("Failed to set up logger: %v", err)
	}
	defer logFile.Close()

	// Load configurations from TOML
	config.TomlInit()

	// Connect to the database
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Set up services and controllers
	userService := services.NewUserService(db)
	userController := controllers.NewUserController(userService)
	middlewareController := middleware.NewMiddleWare(db)

	// Create base router and subrouter for authenticated APIs
	r := mux.NewRouter()

	// Unauthenticated routes
	r.HandleFunc("/signup", userController.SignUpHandler).Methods(http.MethodPost)
	r.HandleFunc("/login", userController.LogInHandler).Methods(http.MethodPost)

	// Subrouter for authenticated routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middlewareController.Middleware) // Register Middleware

	// Authenticated endpoints
	api.HandleFunc("/transaction", userController.TransactionHandler).Methods(http.MethodPost)
	api.HandleFunc("/points/balance", userController.PointsBalanceHandler).Methods(http.MethodGet)
	api.HandleFunc("/points/redeem", userController.RedeemPointsHandler).Methods(http.MethodPost)
	api.HandleFunc("/points/history", userController.PointsHistoryHandler).Methods(http.MethodGet)

	// Start background scheduler
	go StartExpirationScheduler(userService)

	// Get port from config, default to 8080
	port := config.GetTomlStr("common", "port")
	if port == "" {
		port = "8080"
	}

	// Start server
	fmt.Printf("🚀 Server running on port %s...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// StartExpirationScheduler runs periodically to expire old points
func StartExpirationScheduler(userService *services.UserService) {
	debug := new(helpers.HelperStruct)
	debug.Init()

	debug.Info("StartExpirationScheduler(+)")

	// Default interval to 24 hours unless overridden
	hourStr := config.GetTomlStr("common", "hour")
	interval := time.Duration(24) * time.Hour

	if parsedHour, err := strconv.Atoi(hourStr); err == nil && parsedHour > 0 {
		interval = time.Duration(parsedHour) * time.Hour
	} else {
		debug.Warn("Invalid or missing 'hour' in config, defaulting to 24h")
	}

	ticker := time.NewTicker(interval)

	for range ticker.C {
		debug.Info("Running points expiration job...")
		userService.ExpireOldPoints(debug)
	}
}
