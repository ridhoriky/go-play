package dto

// ─── Health Check ──────────────────────────────────────────────────────────────

type HealthCheckResponse struct {
	Status  string `json:"status" example:"alive"`
	Time    string `json:"time" example:"2024-01-15T10:30:00Z"`
	Service string `json:"service" example:"go-play"`
	Version string `json:"version" example:"1.0.0"`
}

// ─── Ready Check ──────────────────────────────────────────────────────────────

type ReadyCheckResponse struct {
	Status       string            `json:"status" example:"ready"` // ready, not_ready
	Time         string            `json:"time"`
	Service      string            `json:"service"`
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"` // database, cache, etc
	Message      string            `json:"message"`
}

// ─── Pagination ──────────────────────────────────────────────────────────────

type PaginationMeta struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

type APIResponse struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    interface{}     `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
}

// ─── Standard Response Wrapper ───────────────────────────────────────────────

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   *Error `json:"error"`
}

type Error struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
}

// ─── Kode Error Standar ──────────────────────────────────────────────────────

const (
	ErrCodeValidation    = "VALIDATION_ERROR"
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeConflict      = "CONFLICT"
	ErrCodeUnprocessable = "UNPROCESSABLE"
	ErrCodeInternal      = "INTERNAL_SERVER_ERROR"
)

func (e *Error) Error() string {
	return e.Message
}

// ─── Helper Constructor ──────────────────────────────────────────────────────

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func NewSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{Success: true, Data: data}
}

func NewErrorResponse(code int, message string, details ...ErrorDetail) ErrorResponse {
	e := &Error{Code: code, Message: message}
	if len(details) > 0 {
		e.Details = details
	}
	return ErrorResponse{Success: false, Error: e}
}
