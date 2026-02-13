package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type PgInt4Array []int

// TODO: testing
func (s *PgInt4Array) Scan(src interface{}) error {
	str, ok := src.(string)
	if !ok {
		return fmt.Errorf("[PgInt4Array.Scan]: expected src to be string, got %T", src)
	}

	// Parse Postgres array string, e.g. "{1, 2,3, 4}"
	var arr []int

	str = strings.TrimSpace(str)
	if len(str) >= 2 && str[0] == '{' && str[len(str)-1] == '}' {
		str = str[1 : len(str)-1]
	}
	if strings.TrimSpace(str) == "" {
		*s = arr
		return nil
	}
	elems := strings.Split(str, ",")
	for _, elem := range elems {
		v := strings.TrimSpace(elem)
		n, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("failed parsing int in Postgres array: %w", err)
		}
		arr = append(arr, n)
	}
	*s = arr
	return nil
}
