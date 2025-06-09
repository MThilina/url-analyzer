package model

type AnalyzeRequest struct {
	URL string `json:"url" binding:"required,url"`
}

type AnalyzeResponse struct {
	HTMLVersion  string         `json:"htmlVersion"`
	Title        string         `json:"title"`
	Headings     map[string]int `json:"headings"`
	Links        LinkSummary    `json:"links"`
	HasLoginForm bool           `json:"hasLoginForm"`
}

type LinkSummary struct {
	Internal     int `json:"internal"`
	External     int `json:"external"`
	Inaccessible int `json:"inaccessible"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
