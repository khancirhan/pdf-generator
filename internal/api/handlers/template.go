package handlers

import (
	"net/http"
	"pdf-generator/internal/domain"
	"pdf-generator/internal/pdfgen"
	"pdf-generator/internal/services"

	"github.com/gin-gonic/gin"
)

type RenderHTMLRequest struct {
	Template string         `json:"template" binding:"required"`
	Data     map[string]any `json:"data"`
}

type RenderPDFRequest struct {
	Template string            `json:"template" binding:"required"`
	Data     map[string]any    `json:"data"`
	Options  pdfgen.PDFOptions `json:"options"`
}

type TemplateHandler struct {
	templateService *services.TemplateService
}

func NewTemplateHandler(service *services.TemplateService) *TemplateHandler {
	return &TemplateHandler{
		templateService: service,
	}
}

func (h *TemplateHandler) GetAll(ctx *gin.Context) {
	templates, err := h.templateService.GetAll()

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, templates)
}

func (h *TemplateHandler) GetByName(ctx *gin.Context) {
	name := ctx.Param("name")

	template, err := h.templateService.GetByName(name)

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, template)

}

func (h *TemplateHandler) RenderHTML(ctx *gin.Context) {
	var req RenderHTMLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(domain.BadRequestError("Invalid JSON body"))
		return
	}

	html, err := h.templateService.RenderHTML(req.Template, req.Data)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func (h *TemplateHandler) RenderPDF(ctx *gin.Context) {
	var req RenderPDFRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(domain.BadRequestError("Invalid JSON body"))
		return
	}

	pdf, err := h.templateService.RenderPDF(req.Template, req.Data, req.Options)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Data(http.StatusOK, "application/pdf", pdf)
}
