package dto

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

func (e *Error) Error() string {
	return e.Message
}

// ─── Kode Error Standar ──────────────────────────────────────────────────────

const (
	ErrCodeValidation    = "VALIDATION_ERROR"
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeConflict      = "CONFLICT"
	ErrCodeUnprocessable = "UNPROCESSABLE"
	ErrCodeInternal      = "INTERNAL_SERVER_ERROR"
)

// ─── Helper Constructor ──────────────────────────────────────────────────────

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
