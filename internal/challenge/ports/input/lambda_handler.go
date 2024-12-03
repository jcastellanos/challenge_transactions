package handler

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jcastellanos/challenge_transactions/internal/challenge/domain/model"
	"github.com/jcastellanos/challenge_transactions/internal/challenge/domain/usecase"
)

type lambdaHandler struct {
	processTransactionsUsecase usecase.ProcessTransactionsUsecase
}

func NewLambdaHandler(processTransactionsUsecase usecase.ProcessTransactionsUsecase) lambdaHandler {
	return lambdaHandler{
		processTransactionsUsecase: processTransactionsUsecase,
	}
}

func (lh lambdaHandler) Handle(ctx context.Context, event events.S3Event) {
	sess := session.Must(session.NewSession())
	s3Client := s3.New(sess)

	for _, record := range event.Records {
		s3Bucket := record.S3.Bucket.Name
		s3Key := record.S3.Object.Key

		log.Printf("Archivo subido al bucket: %s, clave: %s", s3Bucket, s3Key)

		output, err := s3Client.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    aws.String(s3Key),
		})
		if err != nil {
			log.Printf("Error obteniendo el archivo de S3: %v", err)
			continue
		}

		defer output.Body.Close()
		// Save file
		localFilePath := filepath.Join("/tmp", filepath.Base(s3Key))
		localFile, err := os.Create(localFilePath)
		if err != nil {
			log.Printf("Error creando archivo en /tmp: %v", err)
			continue
		}
		defer localFile.Close()
		log.Printf("Guardando una copia del archivo en: %s", localFilePath)
		if _, err := io.Copy(localFile, output.Body); err != nil {
			log.Printf("Error guardando el archivo en /tmp: %v", err)
		}
		// Volver a abrir el archivo para leer su contenido línea por línea
		localFile.Seek(0, io.SeekStart)
		scanner := bufio.NewScanner(localFile)
		transactions := []model.Transaction{}
		numLine := 0
		for scanner.Scan() {
			log.Println(scanner.Text())
			strTx := scanner.Text()
			if numLine > 0 {
				if transaction, err := model.ParseTransaction(strTx); err == nil {
					transactions = append(transactions, *transaction)
				} else {
					log.Printf("Error con la estructura del registro: %s", strTx)
				}
			}
			numLine++
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Error leyendo el archivo: %v", err)
		}
		if err = lh.processTransactionsUsecase.Execute(transactions, localFilePath); err != nil {
			log.Printf("Error procesando transacciones: %v", err)
		}
	}
}
