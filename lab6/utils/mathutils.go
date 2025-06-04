package utils

func GetMax(data []float64) float64 {
	max := data[0]
	for _, value := range data {
		if value > max {
			max = value
		}
	}
	return max
}

func GetMin(data []float64) float64 {
	min := data[0]
	for _, value := range data {
		if value < min {
			min = value
		}
	}
	return min
}
