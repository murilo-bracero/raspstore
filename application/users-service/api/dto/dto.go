package dto

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Code    string `json:"code,omitempty"`
}
