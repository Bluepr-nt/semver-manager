package utils

import "strconv"

func IsNumerical(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
