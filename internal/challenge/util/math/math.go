package math

import "log"

func Average(numbers []float64) float64 {
	if len(numbers) == 0 {
		log.Printf("el slice esta vacio")
		return 0
	}
	var sum float64 = 0
	for _, number := range numbers {
		sum += number
	}
	average := sum / float64(len(numbers))
	return average
}
