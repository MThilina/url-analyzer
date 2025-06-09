// cmd/main.go
package main

import (
	"log"

	"url-analyzer/config"
	_ "url-analyzer/docs"

	"github.com/gin-gonic/gin"

	"url-analyzer/internal/handler"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title URL Analyzer API
// @version 1.0
// @description Analyzes a webpage for structure and link metadata
// @host localhost:8080
// @BasePath /
func main() {
	cfg, err := config.LoadConfig("./config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	r := gin.Default()

	// Serve static CSS/JS
	r.Static("/static", "./static")

	// Load HTML templates
	r.LoadHTMLGlob("template/*")

	// Serve UI
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// Serve backend API
	r.POST("/analyze", handler.AnalyzeHandler)

	// Swagger API documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server on configured port
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
