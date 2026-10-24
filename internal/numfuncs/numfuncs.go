package numfuncs

import (
	"strconv"
)

func ValidNum(x string) bool {
	if _, err := strconv.ParseFloat(x, 64); err != nil {
		return false
	}
	return true
}

func TruncateFloat(x float64) string {
	return strconv.FormatFloat(x, 'g', -1, 64)
}
