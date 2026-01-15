package main

import (
	"log"
	"pdf-generator/internal/api/routes"
	"pdf-generator/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// pdfGenerator, err := pdfgen.NewPDFGenerator()
	// if err != nil {
	// 	log.Fatalf("Failed to initialize PDF generator: %v", err)
	// 	panic(err)
	// }
	// defer pdfGenerator.Close()

	r := gin.Default()

	// Equivalent to:
	// r := gin.New()
	// r.Use(gin.Logger())
	// r.Use(gin.Recovery())

	// gin.SetMode(gin.DebugMode)

	routes.RegisterRoutes(r, cfg)

	log.Printf("Starting server on :%s", cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
