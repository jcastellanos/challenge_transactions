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

type standaloneHandler struct {
	folder                     string
	processTransactionsUsecase usecase.ProcessTransactionsUsecase
}

func NewStandaloneHandler(folder string, processTransactionsUsecase usecase.ProcessTransactionsUsecase) standaloneHandler {
	return standaloneHandler{
		folder:                     folder,
		processTransactionsUsecase: processTransactionsUsecase,
	}
}

func (l standaloneHandler) Run() {
	pendingFolder := l.folder + "/pending"
	log.Printf("scanning the folder: %s", pendingFolder)

	for {
		files, err := os.ReadDir(pendingFolder)
		if err != nil {
			log.Fatalf("error reading the folder: %v", err)
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

func (l standaloneHandler) processFile(filePath string) {
	log.Printf("reading file: %s", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("error reading file: %v", err)
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
				log.Printf("error with the structure of the row: %s", strTx)
			}
		}
		numLine++
	}
	if err := scanner.Err(); err != nil {
		log.Printf("error reading the file: %v", err)
	}
	if err = l.processTransactionsUsecase.Execute(transactions, filePath); err != nil {
		log.Printf("error processing transactions: %v", err)
	}
	destPath := filepath.Join(l.folder+"/processed", filepath.Base(filePath))
	err = os.Rename(filePath, destPath)
	if err != nil {
		log.Printf("error moving the file: %v", err)
	}
}
