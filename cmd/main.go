package main

import (
	"log"
	"os"

	"TA072025/internal/auth"
	"TA072025/internal/database"
	"TA072025/internal/handlers"
	"TA072025/internal/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Validate required environment variables
	requiredEnvVars := []string{
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE",
	}
	for _, env := range requiredEnvVars {
		if os.Getenv(env) == "" {
			log.Fatalf("Variabel lingkungan %s tidak diatur", env)
		}
	}

	// Set Gin mode
	gin.SetMode(utils.GetEnvWithDefault("GIN_MODE", "debug"))

	// Initialize database connection
	database.Initialize()

	// Initialize auth service (includes both user and student repositories)
	auth.Initialize()

	// Create admin user
	err := auth.CreateAdminUser()
	if err != nil {
		log.Fatalf("Gagal membuat pengguna admin: %v", err)
	}

	// Create a new Gin router
	router := gin.Default()

	config := cors.Config{
	    AllowOrigins: []string{
	        "http://localhost:3000",
	        "https://eye-disease-detection25.vercel.app",
	    },
	    AllowMethods: []string{
	        "GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
	    },
	    AllowHeaders: []string{
	        "Origin", "Content-Type", "Authorization",
	    },
	    AllowCredentials: true,
	}

	router.Use(cors.New(config))

	// Register authentication routes
	router.POST("/api/auth/login", handlers.Login)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090" // untuk local
	}
	
	log.Printf("Server berjalan di port %s", port)
	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("Gagal memulai server: %v", err)
	}
}
