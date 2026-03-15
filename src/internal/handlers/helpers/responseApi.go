package helpers

import (
	"encoding/json"
	"errors"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"
	"net/http"

	"github.com/lib/pq"
)

func ParsePgError(err error) *dto.Error {
	var pqErr *pq.Error

	if !errors.As(err, &pqErr) {
		return &dto.Error{
			Code:    http.StatusInternalServerError,
			Message: "An unexpected error occurred",
		}
	}

	switch pqErr.Code {
	case preference.PgErrUniqueViolation:
		field := extractFieldFromConstraint(pqErr.Constraint)
		return &dto.Error{
			Code:    http.StatusConflict,
			Message: field + " already exists, please use a different value",
		}

	case preference.PgErrForeignKeyViolation:
		return &dto.Error{
			Code:    http.StatusUnprocessableEntity,
			Message: "Referenced resource does not exist",
		}

	case preference.PgErrNotNullViolation:
		return &dto.Error{
			Code:    http.StatusBadRequest,
			Message: "Field " + pqErr.Column + " must not be empty",
		}

	case preference.PgErrCheckViolation:
		return &dto.Error{
			Code:    http.StatusBadRequest,
			Message: "Value does not satisfy the required constraint",
		}

	default:
		return &dto.Error{
			Code:    http.StatusInternalServerError,
			Message: "An unexpected error occurred",
		}
	}
}

func extractFieldFromConstraint(constraint string) string {
	friendlyNames := map[string]string{
		"products_name_key":   "Product Name",
		"categories_name_key": "Category Name",
		"users_email_key":     "Email",
		// Add other constraints according to DB schema
	}

	if name, ok := friendlyNames[constraint]; ok {
		return name
	}
	return "A field"
}

func ResponseSuccess(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := dto.APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func ResponseError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	appErr, ok := err.(*dto.Error)
	if !ok {
		appErr = ParsePgError(err)
	}

	w.WriteHeader(appErr.Code)

	resp := dto.APIResponse{
		Status:  appErr.Code,
		Message: appErr.Message,
		Error:   preference.ErrorCodeByHTTPStatus[appErr.Code],
	}

	_ = json.NewEncoder(w).Encode(resp)
}
