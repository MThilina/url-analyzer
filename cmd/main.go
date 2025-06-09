// cmd/main.go
package main

import (
	"log"

	_ "url-analyzer/docs"
	"url-analyzer/internal/config"
	"url-analyzer/internal/handler"

	"github.com/gin-gonic/gin"

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

	r.POST("/analyze", handler.AnalyzeHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":" + cfg.Server.Port)
}
