package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jcastellanos/challenge_transactions/internal/challenge/adapter"
	"github.com/jcastellanos/challenge_transactions/internal/challenge/domain/usecase"
	"github.com/jcastellanos/challenge_transactions/internal/challenge/handler"
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
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
		})
		if err != nil {
			log.Fatalf("AWS session error: %v", err)
		}
		svc := dynamodb.New(sess)
		persistencePort := adapter.NewDynamoPersistenceAdapter(svc)
		processTransactionsUsecase := usecase.NewProcessTransactionUsecase(emailPort, persistencePort)
		lambdaHadler := handler.NewLambdaHandler(processTransactionsUsecase)
		lambda.Start(lambdaHadler.Handle)
	} else {
		folder := os.Getenv("TRANSACTIONS_FOLDER")
		db, err := sql.Open("sqlite", "challenge.db")
		if err != nil {
			log.Fatalf("error reading de embedded database: %v", err)
		}
		defer db.Close()
		persistencePort := adapter.NewSqlitePersistenceAdapter(db)
		persistencePort.InitializeDatabase()
		processTransactionsUsecase := usecase.NewProcessTransactionUsecase(emailPort, persistencePort)
		standaloneListener := handler.NewStandaloneHandler(folder, processTransactionsUsecase)
		standaloneListener.Run()
	}
}
