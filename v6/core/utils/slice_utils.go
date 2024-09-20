package utils

func ContainsString(items []string, item string) bool {
	for _, it := range items {
		if it == item {
			return true
		}
	}
	return false
}

func ContainsInt(items []int, item int) bool {
	for _, it := range items {
		if it == item {
			return true
		}
	}
	return false
}

func ContainsDouble(items []float64, item float64) bool {
	for _, it := range items {
		if it == item {
			return true
		}
	}
	return false
}
