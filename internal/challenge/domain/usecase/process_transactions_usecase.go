package usecase

import (
	"log"
	"os"

	"github.com/jcastellanos/challenge_transactions/internal/challenge/domain/model"
)

type ProcessTransactionsUsecase struct {
	emailPort       EmailPort
	persistencePort PersistencePort
}

func NewProcessTransactionUsecase(emailPort EmailPort, persistencePort PersistencePort) ProcessTransactionsUsecase {
	return ProcessTransactionsUsecase{
		emailPort:       emailPort,
		persistencePort: persistencePort,
	}
}

func (p ProcessTransactionsUsecase) Execute(transactions []model.Transaction, transactionsPath string) error {
	log.Println("processing transactions ...")
	statistics := model.NewStatistics()
	for _, ts := range transactions {
		statistics.Add(ts)
	}
	log.Println("saving transactions to database")
	if err := p.persistencePort.InsertTransactions(transactions); err != nil {
		return err
	}
	emailTo := os.Getenv("EMAIL_TO")
	subject := "Your transactions"

	log.Println("sending email")
	if err := p.emailPort.SendEmail(emailTo, subject, statistics, transactionsPath); err != nil {
		return err
	}
	log.Println("transactions processed successfully")
	return nil
}
