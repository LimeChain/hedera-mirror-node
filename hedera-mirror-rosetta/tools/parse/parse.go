package parse

import "strconv"

func ToInt(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
