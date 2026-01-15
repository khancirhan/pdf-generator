package domain

type Template struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type GetAllTemplateResponse struct {
	Name string `json:"name"`
}
