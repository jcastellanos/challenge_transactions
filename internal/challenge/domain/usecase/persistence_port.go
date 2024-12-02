package usecase

import "github.com/jcastellanos/challenge_transactions/internal/challenge/domain/model"

type PersistencePort interface {
	InsertTransactions(transactions []model.Transaction) error
}
