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

func (p ProcessTransactionsUsecase) Execute(transactions []model.Transaction) error {
	log.Println("procesando transacciones")
	statistics := model.NewStatistics()
	for _, ts := range transactions {
		statistics.Add(ts)
	}
	log.Println(statistics.TotalBalance())
	log.Println(statistics.AverageCredit())
	log.Println(statistics.AverageDebit())
	log.Println("persistiendo transacciones")
	if err := p.persistencePort.InsertTransactions(transactions); err != nil {
		return err
	}
	// Datos del correo
	emailTo := os.Getenv("EMAIL_TO")
	subject := "Correo con Imagen y Adjunto"
	attachmentPath := "transactions/processed/acc1.csv" // TODO cambiar por el correcto
	variables := map[string]string{
		"Nombre": "Juan",
		"Fecha":  "28/11/2024",
	}
	// Enviar el correo
	log.Println("enviando correo")
	if err := p.emailPort.SendEmail(emailTo, subject, variables, attachmentPath); err != nil {
		return err
	}
	log.Println("transacciones procesadas")
	return nil
}
