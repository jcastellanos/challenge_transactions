package handler

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jcastellanos/challenge_transactions/internal/challenge/domain/model"
	"github.com/jcastellanos/challenge_transactions/internal/challenge/domain/usecase"
)

type listener struct {
	folder                     string
	processTransactionsUsecase usecase.ProcessTransactionsUsecase
}

func NewStandaloneListener(folder string, processTransactionsUsecase usecase.ProcessTransactionsUsecase) listener {
	return listener{
		folder:                     folder,
		processTransactionsUsecase: processTransactionsUsecase,
	}
}

func (l listener) Run() {
	pendingFolder := l.folder + "/pending"
	log.Printf("Revisando la carpeta: %s", pendingFolder)

	for {
		files, err := os.ReadDir(pendingFolder)
		if err != nil {
			log.Fatalf("Error leyendo la carpeta: %v", err)
		}

		for _, file := range files {
			if !file.IsDir() {
				filePath := filepath.Join(pendingFolder, file.Name())
				l.processFile(filePath)
			}
		}

		time.Sleep(5 * time.Second) // Revisar cada 5 segundos
	}
}

func (l listener) processFile(filePath string) {
	log.Printf("Leyendo archivo: %s", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error abriendo el archivo: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	transactions := []model.Transaction{}
	numLine := 0
	for scanner.Scan() {
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
	if err = l.processTransactionsUsecase.Execute(transactions, filePath); err != nil {
		log.Printf("Error procesando transacciones: %v", err)
	}
	destPath := filepath.Join(l.folder+"/processed", filepath.Base(filePath))
	err = os.Rename(filePath, destPath)
	if err != nil {
		log.Printf("Error moviendo el archivo: %v", err)
	}
}
