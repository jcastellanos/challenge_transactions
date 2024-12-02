package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatistics(t *testing.T) {
	statistics := NewStatistics()
	statistics.Add(Transaction{
		Id:          0,
		Month:       7,
		Day:         15,
		Transaction: 60.5,
	})
	statistics.Add(Transaction{
		Id:          1,
		Month:       7,
		Day:         15,
		Transaction: -10.3,
	})
	statistics.Add(Transaction{
		Id:          2,
		Month:       8,
		Day:         2,
		Transaction: -20.46,
	})
	statistics.Add(Transaction{
		Id:          3,
		Month:       8,
		Day:         13,
		Transaction: 10,
	})
	assert.Equal(t, 39.74, statistics.TotalBalance())
	assert.Equal(t, -15.38, statistics.AverageDebit())
	assert.Equal(t, 35.25, statistics.AverageCredit())
	expectedTransactionsByMonth := map[int]int{
		7: 2,
		8: 2,
	}
	assert.Equal(t, expectedTransactionsByMonth, statistics.TransactionsByMonth())
}
