package usecase

import (
	"log"

	"github.com/jcastellanos/challenge_transactions/internal/challenge/domain/model"
)

type ProcessTransactionsUsecase struct {
}

func NewProcessTransactionUsecase() ProcessTransactionsUsecase {
	return ProcessTransactionsUsecase{}
}

func (p ProcessTransactionsUsecase) Execute(transactions []model.Transaction) {
	log.Println(transactions)
	statistics := model.NewStatistics()
	for _, ts := range transactions {
		statistics.Add(ts)
	}
	log.Println(statistics.TotalBalance())
	log.Println(statistics.AverageCredit())
	log.Println(statistics.AverageDebit())
}
