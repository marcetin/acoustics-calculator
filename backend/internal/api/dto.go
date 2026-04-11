package api

// Response envelopes
type HealthResponse struct {
	OK bool `json:"ok"`
}

type SuccessResponse struct {
	OK       bool        `json:"ok"`
	Data     interface{} `json:"data"`
	Warnings []string    `json:"warnings"`
}

type ErrorResponse struct {
	OK    bool        `json:"ok"`
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}
