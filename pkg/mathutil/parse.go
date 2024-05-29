package mathutil

import (
	"strconv"
)

func ParseInt(size string) (int, error) {
	if size == "" {
		return 0, nil
	}

	value, err := strconv.ParseInt(size, 10, 32)

	return int(value), err
}

func ParseInt64(size string) (int64, error) {
	if size == "" {
		return 0, nil
	}

	return strconv.ParseInt(size, 10, 64)

}
