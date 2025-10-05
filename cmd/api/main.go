package main

import (
	"log"

	"github.com/juan-cantero/fitapi/config"
	"github.com/juan-cantero/fitapi/internal/database"
	"github.com/juan-cantero/fitapi/internal/handlers"
	"github.com/juan-cantero/fitapi/internal/middleware"
	"github.com/juan-cantero/fitapi/internal/repositories"
	"github.com/juan-cantero/fitapi/internal/services"

	"github.com/gin-gonic/gin"
	supa "github.com/supabase-community/supabase-go"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Initialize database connection
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Supabase client
	supabaseClient, err := supa.NewClient(cfg.SupabaseURL, cfg.SupabaseKey, &supa.ClientOptions{})
	if err != nil {
		log.Fatalf("Failed to initialize Supabase client: %v", err)
	}

	// Initialize repositories
	equipmentRepo := repositories.NewPostgresEquipmentRepository(db.Pool)

	// Initialize services
	equipmentService := services.NewEquipmentService(equipmentRepo)

	// Initialize handlers
	equipmentHandler := handlers.NewEquipmentHandler(equipmentService)

	// Initialize Gin router
	router := gin.Default()

	// Public routes (no authentication required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "ok",
			"database": "connected",
			"supabase": supabaseClient != nil,
		})
	})

	// Protected routes (authentication required)
	api := router.Group("/api")
	api.Use(middleware.AuthRequired())
	{
		// Test endpoint to verify auth is working
		api.GET("/me", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			userEmail, _ := c.Get("user_email")

			c.JSON(200, gin.H{
				"user_id": userID,
				"email":   userEmail,
				"message": "Authentication successful!",
			})
		})

		// Equipment endpoints
		api.POST("/equipment", equipmentHandler.Create)
		api.GET("/equipment", equipmentHandler.List)
		api.GET("/equipment/:id", equipmentHandler.GetByID)
		api.PUT("/equipment/:id", equipmentHandler.Update)
		api.DELETE("/equipment/:id", equipmentHandler.Delete)
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
