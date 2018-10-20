package ticker

import "strconv"

// ValidateTimestamp validates the timestamp parameter for the ticker endpoint
func ValidateTimestamp(variable string) (bool, interface{}) {
	ts, err := strconv.ParseInt(variable, 10, 64)
	if err != nil {
		return false, nil
	}
	return true, ts
}
