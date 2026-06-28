package helpers

import (
	"errors"
	"net/http"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func ParsePgError(err error) *dto.Error {
	var pqErr *pq.Error

	if !errors.As(err, &pqErr) {
		return dto.NewError(http.StatusInternalServerError, "An unexpected error occurred")
	}

	switch pqErr.Code {
	case preference.PgErrUniqueViolation:
		field := extractFieldFromConstraint(pqErr.Constraint)
		return dto.NewError(http.StatusConflict, field+" already exists, please use a different value")

	case preference.PgErrForeignKeyViolation:
		return dto.NewError(http.StatusUnprocessableEntity, "Referenced resource does not exist")

	case preference.PgErrNotNullViolation:
		return dto.NewError(http.StatusBadRequest, "Field "+pqErr.Column+" must not be empty")

	case preference.PgErrCheckViolation:
		return dto.NewError(http.StatusBadRequest, "Value does not satisfy the required constraint")

	default:
		return dto.NewError(http.StatusInternalServerError, "An unexpected error occurred")
	}
}

func extractFieldFromConstraint(constraint string) string {
	friendlyNames := map[string]string{
		"products_name_key":   "Product Name",
		"categories_name_key": "Category Name",
		"users_email_key":     "Email",
	}

	if name, ok := friendlyNames[constraint]; ok {
		return name
	}
	return "A field"
}

func ResponseSuccess(c *gin.Context, status int, message string, data any) {
	resp := dto.APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
	c.JSON(status, resp)
}

func ResponseSuccessWithMeta(c *gin.Context, status int, message string, data any, meta *dto.PaginationMeta) {
	resp := dto.APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
	c.JSON(status, resp)
}

func ResponseError(c *gin.Context, err error) {
	appErr := &dto.Error{}
	ok := errors.As(err, &appErr)
	if !ok {
		appErr = ParsePgError(err)
	}

	_ = c.Error(err)

	resp := dto.APIResponse{
		Status:  appErr.Code,
		Message: appErr.Message,
		Error:   preference.ErrorCodeByHTTPStatus[appErr.Code],
	}

	c.JSON(appErr.Code, resp)
}
