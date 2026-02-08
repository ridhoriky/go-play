package database

import (
	"fmt"
	"strings"
)

// rows = [][]any
// tiap row = 1 record
func BuildBulkInsert(
	table string,
	columns []string,
	rows [][]any,
) (string, []any, error) {

	if len(rows) == 0 {
		return "", nil, fmt.Errorf("rows cannot be empty")
	}

	colCount := len(columns)

	if colCount == 0 {
		return "", nil, fmt.Errorf("columns cannot be empty")
	}

	args := make([]any, 0)
	values := make([]string, 0)

	argPos := 1

	for i, row := range rows {

		if len(row) != colCount {
			return "", nil,
				fmt.Errorf(
					"row %d: expected %d columns, got %d",
					i,
					colCount,
					len(row),
				)
		}

		placeholders := make([]string, colCount)

		for j := 0; j < colCount; j++ {
			placeholders[j] = fmt.Sprintf("$%d", argPos)
			argPos++
		}

		values = append(
			values,
			fmt.Sprintf("(%s)", strings.Join(placeholders, ",")),
		)

		args = append(args, row...)
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		table,
		strings.Join(columns, ","),
		strings.Join(values, ","),
	)

	return query, args, nil
}
