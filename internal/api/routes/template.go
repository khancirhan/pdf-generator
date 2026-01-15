package routes

import (
	"pdf-generator/internal/api/handlers"
	"pdf-generator/internal/config"
	"pdf-generator/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterTemplateRoutes(api *gin.RouterGroup, cfg *config.Config) {
	service := services.NewTemplateService(cfg.TemplatesDir, cfg.GotenbergURL)
	handler := handlers.NewTemplateHandler(service)

	group := api.Group("/templates")
	{
		group.GET("/", handler.GetAll)
		group.GET("/:name", handler.GetByName)
		group.POST("/html", handler.RenderHTML)
		group.POST("/pdf", handler.RenderPDF)
	}
}
