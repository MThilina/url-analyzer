package handler

import (
	"net/http"
	"url-analyzer/internal/model"
	"url-analyzer/internal/service"
	"url-analyzer/internal/utils"

	"github.com/gin-gonic/gin"
)

// AnalyzeHandler godoc
// @Summary Analyze a webpage
// @Description Extracts HTML metadata from a given URL
// @Accept json
// @Produce json
// @Param request body model.AnalyzeRequest true "URL to analyze"
// @Success 200 {object} model.AnalyzeResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /analyze [post]
func AnalyzeHandler(c *gin.Context) {
	var req model.AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Invalid request: " + err.Error()})
		return
	}

	normalizedURL, err := utils.NormalizeAndValidateURL(req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: err.Error()})
		return
	}

	resp, err := service.AnalyzeURL(normalizedURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
