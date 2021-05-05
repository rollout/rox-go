package utils

import "strconv"

func ToFloat(value interface{}) (float64, bool) {

	if value, ok := value.(string); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return float64(intValue), true
		} else if doubleValue, err := strconv.ParseFloat(value, 64); err == nil {
			return doubleValue, true
		}
	}

	if value, ok := value.(float64); ok {
		return value, true
	}
	if value, ok := value.(int); ok {
		return float64(value), true
	}
	return 0, false
}
