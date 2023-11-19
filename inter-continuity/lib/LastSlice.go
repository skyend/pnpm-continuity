package lib

import "math"

func LastSlice[T any](slice []T, count int) []T {
	length := len(slice)
	capable := int(math.Max(float64(length-count), 0))

	return slice[capable:]
}
