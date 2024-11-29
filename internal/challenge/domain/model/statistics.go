package model

import (
	"log"
)

type Statistics struct {
	totalBalance        float64
	credits             []float64
	debits              []float64
	transactionsByMonth map[int]int
}

func NewStatistics() Statistics {
	return Statistics{
		totalBalance:        0,
		credits:             []float64{},
		debits:              []float64{},
		transactionsByMonth: make(map[int]int),
	}
}

func (s *Statistics) Add(tx Transaction) {
	s.totalBalance = s.totalBalance + tx.Transaction
	if tx.Transaction >= 0 {
		s.credits = append(s.credits, tx.Transaction)
	}
	if tx.Transaction < 0 {
		s.debits = append(s.debits, tx.Transaction)
	}

}

func (s Statistics) TotalBalance() float64 {
	return s.totalBalance
}

func (s Statistics) AverageCredit() float64 {
	return average(s.credits)
}

func (s Statistics) AverageDebit() float64 {
	return average(s.debits)
}

func average(numbers []float64) float64 {
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
