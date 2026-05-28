package utils

// HTTPError represents an HTTP error.
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}
