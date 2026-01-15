package routes

import (
	"pdf-generator/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterHealthRoutes(api *gin.RouterGroup) {
	handler := handlers.NewHealthHandler()

	api.GET("/health", handler.CheckHealth)
}
