package util

const maxLimit, defaultLimit = 100, 10

func ValidatePage(page int) int {
	if page < 1 {
		return 1
	}

	return page
}

func ValidatePageSize(pageSize int) int {
	if pageSize < 1 {
		return defaultLimit
	} else if pageSize > maxLimit {
		return maxLimit
	}

	return pageSize
}
