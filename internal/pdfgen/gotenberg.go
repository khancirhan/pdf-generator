package pdfgen

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type GotenbergPDFGenerator struct {
	baseURL    string
	httpClient *http.Client
}

func NewGotenbergPDFGenerator(baseURL string) *GotenbergPDFGenerator {
	return &GotenbergPDFGenerator{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// GeneratePDF from HTML string
func (c *GotenbergPDFGenerator) GeneratePDF(html string, opts PDFOptions) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// HTML content
	part, err := writer.CreateFormFile("files", "index.html")
	if err != nil {
		return nil, err
	}
	part.Write([]byte(html))

	// Options
	writer.WriteField("paperWidth", fmt.Sprintf("%.2f", opts.PaperWidth))
	writer.WriteField("paperHeight", fmt.Sprintf("%.2f", opts.PaperHeight))
	writer.WriteField("marginTop", fmt.Sprintf("%.2f", opts.MarginTop))
	writer.WriteField("marginBottom", fmt.Sprintf("%.2f", opts.MarginBottom))
	writer.WriteField("marginLeft", fmt.Sprintf("%.2f", opts.MarginLeft))
	writer.WriteField("marginRight", fmt.Sprintf("%.2f", opts.MarginRight))
	writer.WriteField("printBackground", fmt.Sprintf("%t", opts.PrintBackground))
	writer.WriteField("preferCssPageSize", fmt.Sprintf("%t", opts.PreferCSSPageSize))
	writer.WriteField("landscape", fmt.Sprintf("%t", opts.Landscape))

	if opts.WaitDelay != "" {
		writer.WriteField("waitDelay", opts.WaitDelay)
	}

	if opts.WaitForExpression != "" {
		writer.WriteField("waitForExpression", opts.WaitForExpression)
	}

	writer.Close()

	req, err := http.NewRequest("POST", c.baseURL+"/forms/chromium/convert/html", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gotenberg error %d: %s", resp.StatusCode, string(errBody))
	}

	return io.ReadAll(resp.Body)
}
