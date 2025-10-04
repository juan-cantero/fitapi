package main

import (
	"log"

	"github.com/juan-cantero/fitapi/config"
	"github.com/juan-cantero/fitapi/internal/database"

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

	// Initialize Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "ok",
			"database": "connected",
			"supabase": supabaseClient != nil,
		})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
