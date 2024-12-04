package adapter

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jcastellanos/challenge_transactions/internal/challenge/domain/model"
)

type sqlitePersistenceAdapter struct {
	db *sql.DB
}

func NewSqlitePersistenceAdapter(db *sql.DB) sqlitePersistenceAdapter {
	return sqlitePersistenceAdapter{
		db: db,
	}
}

func (sa sqlitePersistenceAdapter) InitializeDatabase() error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS 'transactions' (
		id VARCHAR PRIMARY KEY,
		transaction_id INTEGER NOT NULL,
		date VARCHAR NOT NULL,
		'transaction' REAL NOT NULL
	);`
	_, err := sa.db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}
	return nil
}

func (sa sqlitePersistenceAdapter) InsertTransactions(transactions []model.Transaction) error {
	if len(transactions) == 0 {
		log.Printf("empty transactions")
		return nil
	}
	query := "INSERT INTO 'transactions'(id, transaction_id, date, 'transaction') VALUES(?, ?, ?, ?)"
	stmt, err := sa.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error prepared statement: %w", err)
	}
	defer stmt.Close()
	for _, t := range transactions {
		date := fmt.Sprintf("%d/%d", t.Month, t.Day)
		_, err := stmt.Exec(uuid.New().String(), t.Id, date, t.Transaction)
		if err != nil {
			return fmt.Errorf("error executing query: %w", err)
		}
	}
	return nil
}
