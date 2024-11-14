package model

type Response struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message,omitempty"`
	Data    interface{}    `json:"data,omitempty"`
}

type ResponseStatus string

const (
	ResponseStatusSuccess = "success"
	ResponseStatusFail    = "fail"
	ResponseStatusError   = "error"
)
