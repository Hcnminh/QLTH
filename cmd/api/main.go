package main

import (
	"log"
	"os"

	"qlth/internal/api"
)

func main() {
	router := api.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting EduPrimary-KPI Core Server on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
