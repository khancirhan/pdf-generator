package pdfgen

type PDFOptions struct {
	PaperWidth        float64 `json:"paperWidth"`
	PaperHeight       float64 `json:"paperHeight"`
	MarginTop         float64 `json:"marginTop"`
	MarginBottom      float64 `json:"marginBottom"`
	MarginLeft        float64 `json:"marginLeft"`
	MarginRight       float64 `json:"marginRight"`
	PrintBackground   bool    `json:"printBackground"`
	PreferCSSPageSize bool    `json:"preferCSSPageSize"`
	Landscape         bool    `json:"landscape"`
	WaitDelay         string  `json:"waitDelay"`         // e.g., "1s"
	WaitForExpression string  `json:"waitForExpression"` // JS expression to wait for
}

func DefaultPDFOptions() PDFOptions {
	return PDFOptions{
		PaperWidth:        8.5,
		PaperHeight:       11,
		MarginTop:         0.5,
		MarginBottom:      0.5,
		MarginLeft:        0.5,
		MarginRight:       0.5,
		PrintBackground:   true,
		PreferCSSPageSize: true,
	}
}

type PDFGenerator interface {
	GeneratePDF(html string, opts PDFOptions) ([]byte, error)
}
