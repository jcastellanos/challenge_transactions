package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/jcastellanos/challenge_transactions/internal/challenge/adapter"
	"github.com/jcastellanos/challenge_transactions/internal/challenge/domain/usecase"
	standaloneListener "github.com/jcastellanos/challenge_transactions/internal/challenge/ports/input"
	_ "modernc.org/sqlite"
)

func main() {

	runtime := os.Getenv("RUNTIME")
	emailUsername := os.Getenv("EMAIL_USERNAME")
	emailPassword := os.Getenv("EMAIL_PASSWORD")

	emailConfig := adapter.EmailConfig{
		SMTPServer: "smtp.gmail.com",
		Port:       "587",
		Username:   emailUsername,
		Password:   emailPassword,
	}
	emailPort := adapter.NewGmailEmailAdapter(emailConfig)
	if runtime == "lambda" {

	} else {
		folder := os.Getenv("TRANSACTIONS_FOLDER")
		db, err := sql.Open("sqlite", "challenge.db")
		if err != nil {
			log.Fatalf("Error conectando a la base de datos: %v", err)
		}
		defer db.Close()
		persistencePort := adapter.NewSqlitePersistenceAdapter(db)
		persistencePort.InitializeDatabase()
		processTransactionsUsecase := usecase.NewProcessTransactionUsecase(emailPort, persistencePort)
		standaloneListener := standaloneListener.NewListener(folder, processTransactionsUsecase)
		standaloneListener.Run()
	}
}
