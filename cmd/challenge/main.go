package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jcastellanos/challenge_transactions/internal/challenge/adapter"
	"github.com/jcastellanos/challenge_transactions/internal/challenge/domain/usecase"
	handler "github.com/jcastellanos/challenge_transactions/internal/challenge/ports/input"
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
		persistencePort := adapter.NewMockPersistenceAdapter()
		processTransactionsUsecase := usecase.NewProcessTransactionUsecase(emailPort, persistencePort)
		lambdaHadler := handler.NewLambdaHandler(processTransactionsUsecase)
		lambda.Start(lambdaHadler.Handle)
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
		standaloneListener := handler.NewStandaloneListener(folder, processTransactionsUsecase)
		standaloneListener.Run()
	}
}
