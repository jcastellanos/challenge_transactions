package adapter

import (
	"log"

	"github.com/jcastellanos/challenge_transactions/internal/challenge/domain/model"
)

type mockPersistenceAdapter struct {
}

func NewMockPersistenceAdapter() mockPersistenceAdapter {
	return mockPersistenceAdapter{}
}

func (sa mockPersistenceAdapter) InsertTransactions(transactions []model.Transaction) error {
	log.Println("mock InsertTransactions")
	return nil
}
