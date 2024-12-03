package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	Id          int
	Month       int
	Day         int
	Transaction float64
}

func ParseTransaction(transaction string) (*Transaction, error) {
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

	return &Transaction{
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
