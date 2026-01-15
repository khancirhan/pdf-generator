package services

import (
	"fmt"
	"os"
	"path/filepath"
	"pdf-generator/internal/domain"
	"pdf-generator/internal/pdfgen"

	"github.com/osteele/liquid"
)

type TemplateService struct {
	templatesDir string
	pdfGenerator pdfgen.PDFGenerator
}

func NewTemplateService(templatesDir string, gotenbergURL string) *TemplateService {
	pdfGenerator := pdfgen.NewGotenbergPDFGenerator(gotenbergURL)

	return &TemplateService{
		templatesDir: templatesDir,
		pdfGenerator: pdfGenerator,
	}
}

func (s *TemplateService) GetAll() ([]domain.GetAllTemplateResponse, *domain.AppError) {
	var templates []domain.GetAllTemplateResponse

	entries, err := os.ReadDir(s.templatesDir)
	if err != nil {
		return nil, domain.InternalServerError("Some error occurred, please try again")
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext == ".html" {
			templates = append(templates, domain.GetAllTemplateResponse{Name: entry.Name()})
		}
	}

	return templates, nil
}

func (s *TemplateService) GetByName(name string) (*domain.Template, *domain.AppError) {
	content, err := s.getTemplateContent(name)
	if err != nil {
		return nil, err
	}

	return &domain.Template{Name: name, Content: content}, nil
}

func (s *TemplateService) RenderHTML(name string, data map[string]any) (string, *domain.AppError) {
	content, appErr := s.getTemplateContent(name)
	if appErr != nil {
		return "", appErr
	}

	engine := liquid.NewEngine()
	out, renderErr := engine.ParseAndRenderString(content, data)
	if renderErr != nil {
		return "", domain.InternalServerError("Failed to render template")
	}

	return out, nil
}

func (s *TemplateService) RenderPDF(name string, data map[string]any, opts pdfgen.PDFOptions) ([]byte, *domain.AppError) {
	content, appErr := s.getTemplateContent(name)
	if appErr != nil {
		return nil, appErr
	}

	engine := liquid.NewEngine()
	html, renderErr := engine.ParseAndRenderString(content, data)
	if renderErr != nil {
		return nil, domain.InternalServerError("Failed to render template")
	}

	pdf, pdfErr := s.pdfGenerator.GeneratePDF(html, opts)
	if pdfErr != nil {
		return nil, domain.InternalServerError("Failed to generate PDF")
	}

	return pdf, nil
}

func (s *TemplateService) getTemplateContent(name string) (string, *domain.AppError) {
	filePath := filepath.Join(s.templatesDir, name)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", domain.NotFoundError(fmt.Sprintf("Template with name %s not found", name))
	}

	return string(content), nil
}
