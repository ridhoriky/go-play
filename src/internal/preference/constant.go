package preference

const (
	//Error code postgre
	PgErrUniqueViolation     = "23505"
	PgErrForeignKeyViolation = "23503"
	PgErrNotNullViolation    = "23502"
	PgErrCheckViolation      = "23514"

	//Generic Errors
	ErrBadRequest          = "BAD_REQUEST"
	ErrUnauthorized        = "UNAUTHORIZED"
	ErrForbidden           = "FORBIDDEN"
	ErrNotFound            = "NOT_FOUND"
	ErrConflict            = "CONFLICT"
	ErrUnprocessableEntity = "UNPROCESSABLE_ENTITY"
	ErrInternalServer      = "INTERNAL_SERVER_ERROR"

	// Error message global
	ErrInvalidReqBody     = "Invalid request body"
	ErrInvalidQueryParams = "Invalid query parameters"

	// Error message product
	ErrInvalidProductID     = "Invalid product ID"
	ErrNoProductCreated     = "No products to create"
	ErrProductAlreadyExists = "Product name already exists"
	ErrProductEmpty         = "Empty Product to add"
	ErrProductNotFound      = "Product not found"
	ErrProductPriceNegative = "Price cannot be negative"
	ErrProductStockNegative = "Stock cannot be negative"
	ErrProductNameRequied   = "Product name is required"
	ErrProductBatchTooLarge = "Number of products exceeds maximum batch size (max: 500)"

	// Error message category
	ErrInvalidCategoryID     = "Invalid category ID"
	ErrCategoryAlreadyExists = "Category name already exists"
	ErrCategoryNotFound      = "Category not found"
	ErrCategoryNameRequied   = "Category name is required"

	// Error message transaction
	ErrInvalidTranscationID = "Invalid transaction ID"

	// Batch size
	MaxBatchSizeProduct = 500

	// Error message system
	ErrDatabaseConnection = "Database connection error"
)

var ErrorCodeByHTTPStatus = map[int]string{
	400: ErrBadRequest,
	401: ErrUnauthorized,
	403: ErrForbidden,
	404: ErrNotFound,
	409: ErrConflict,
	422: ErrUnprocessableEntity,
	500: ErrInternalServer,
}
