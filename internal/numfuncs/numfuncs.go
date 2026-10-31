package numfuncs

import (
	"math"
	"strconv"
)

func ValidNum(x string) bool {
	n, err := strconv.ParseFloat(x, 64)
	return err != nil && math.IsNaN(n) && math.IsInf(n, 0)
}

func TruncateFloat(x float64) string {
	return strconv.FormatFloat(x, 'g', -1, 64)
}
