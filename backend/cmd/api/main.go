package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paulsena/asheville-setlist/internal/config"
	"github.com/paulsena/asheville-setlist/internal/db"
	"github.com/paulsena/asheville-setlist/internal/handlers"
	"github.com/paulsena/asheville-setlist/internal/middleware"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection pool
	pool, err := config.NewDatabasePool(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Create database queries
	queries := db.New(pool)

	// Create handlers
	h := handlers.New(queries)

	// Set Gin mode from configuration
	gin.SetMode(cfg.GinMode)

	// Create Gin router (without default middleware)
	router := gin.New()

	// Apply middleware stack
	router.Use(middleware.Recovery())   // Recover from panics
	router.Use(middleware.Logger())     // Log requests
	router.Use(middleware.CORS())       // Handle CORS
	router.Use(middleware.BodyLimit(1)) // 1MB body limit

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "asheville-setlist-api",
			"version": "0.1.0",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Shows
		api.GET("/shows", h.ListShows)
		api.GET("/shows/:id", h.GetShow)
		api.POST("/shows", h.CreateShow)

		// Venues
		api.GET("/venues", h.ListVenues)
		api.GET("/venues/:slug", h.GetVenue)

		// Bands
		api.GET("/bands", h.ListBands)
		api.GET("/bands/:slug", h.GetBand)
		api.GET("/bands/:slug/similar", h.GetSimilarBands)

		// Genres
		api.GET("/genres", h.ListGenres)

		// Search
		api.GET("/search", h.Search)
	}

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on port %s (mode: %s)", cfg.Port, cfg.GinMode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with 5 second timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
