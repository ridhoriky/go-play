package database

import (
	"fmt"
	"strings"
)

func BuildBulkInsert(
	table string,
	columns []string,
	rows [][]any,
) (string, []any, error) {

	if table == "" {
		return "", nil, fmt.Errorf("table name cannot be empty")
	}

	if len(columns) == 0 {
		return "", nil, fmt.Errorf("columns cannot be empty")
	}

	if len(rows) == 0 {
		return "", nil, fmt.Errorf("rows cannot be empty")
	}

	colCount := len(columns)

	args := make([]any, 0, len(rows)*colCount)
	values := make([]string, 0, len(rows))

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

	for j, val := range row {

		placeholders[j] = fmt.Sprintf("$%d", argPos)
		argPos++

		args = append(args, val)
	}

	values = append(
		values,
		"("+strings.Join(placeholders, ",")+")",
	)
}


	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		table,
		strings.Join(columns, ","),
		strings.Join(values, ","),
	)

	return query, args, nil
}