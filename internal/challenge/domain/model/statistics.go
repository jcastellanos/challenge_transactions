package model

import "github.com/jcastellanos/challenge_transactions/internal/challenge/util/math"

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
	s.transactionsByMonth[tx.Month] = s.transactionsByMonth[tx.Month] + 1

}

func (s Statistics) TotalBalance() float64 {
	return s.totalBalance
}

func (s Statistics) AverageCredit() float64 {
	return math.Average(s.credits)
}

func (s Statistics) AverageDebit() float64 {
	return math.Average(s.debits)
}

// Return the number of transactions by month.
// It returns a map, the key of the map is the month of the year (1 -> January), the
// map only contains the month than has transactions.
func (s Statistics) TransactionsByMonth() map[int]int {
	return s.transactionsByMonth
}
