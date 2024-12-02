package standaloneListener

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jcastellanos/challenge_transactions/internal/challenge/domain/model"
	"github.com/jcastellanos/challenge_transactions/internal/challenge/domain/usecase"
)

type listener struct {
	folder                     string
	processTransactionsUsecase usecase.ProcessTransactionsUsecase
}

func NewListener(folder string, processTransactionsUsecase usecase.ProcessTransactionsUsecase) listener {
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
			if transaction, err := parseTransaction(strTx); err == nil {
				log.Println(strTx)
				log.Println(transaction)
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

func parseTransaction(transaction string) (*model.Transaction, error) {
	tokens := strings.Split(transaction, ",")
	if len(tokens) != 3 {
		return nil, fmt.Errorf("entrada no válida, se esperaban 3 partes separadas por comas")
	}
	id, err := strconv.Atoi(tokens[0])
	if err != nil {
		return nil, fmt.Errorf("error al convertir Id: %v", err)
	}
	date := tokens[1]
	month, day, err := parseDate(date)
	if err != nil {
		return nil, fmt.Errorf("error al convertir transaction: %v", err)
	}
	tx, err := strconv.ParseFloat(tokens[2], 64)
	if err != nil {
		return nil, fmt.Errorf("error al convertir transaction: %v", err)
	}

	return &model.Transaction{
		Id:          id,
		Month:       month,
		Day:         day,
		Transaction: float64(tx),
	}, nil
}

func parseDate(date string) (int, int, error) {
	const layout = "1/2" // dd/mm en Go es representado por 02/01

	// Intentar analizar la fecha
	if parseDate, err := time.Parse(layout, date); err == nil {
		return int(parseDate.Month()), parseDate.Day(), nil
	}
	return 0, 0, fmt.Errorf("fecha de transaccion no valida: %s", date)
}
