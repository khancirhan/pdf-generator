package routes

import (
	"pdf-generator/internal/api/middlewares"
	"pdf-generator/internal/config"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, cfg *config.Config) {
	api := r.Group("/api/v1")

	api.Use(middlewares.Errorhandler())

	RegisterHealthRoutes(api)
	RegisterTemplateRoutes(api, cfg)
}
