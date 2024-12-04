package adapter

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"github.com/jcastellanos/challenge_transactions/internal/challenge/domain/model"
)

type dynamoPersistenceAdapter struct {
	svc *dynamodb.DynamoDB
}

func NewDynamoPersistenceAdapter(svc *dynamodb.DynamoDB) dynamoPersistenceAdapter {
	return dynamoPersistenceAdapter{
		svc: svc,
	}
}

func (sa dynamoPersistenceAdapter) InsertTransactions(transactions []model.Transaction) error {
	if len(transactions) == 0 {
		log.Printf("empty transactions")
		return nil
	}
	tableName := "transactions"
	const batchSize = 25
	for i := 0; i < len(transactions); i += batchSize {
		end := i + batchSize
		if end > len(transactions) {
			end = len(transactions)
		}
		batch := transactions[i:end]
		var writeRequests []*dynamodb.WriteRequest
		for _, t := range batch {
			date := fmt.Sprintf("%d/%d", t.Month, t.Day)
			item := map[string]*dynamodb.AttributeValue{
				"Id":            {S: aws.String(uuid.New().String())},
				"TransactionId": {S: aws.String(strconv.Itoa(t.Id))},
				"Date":          {S: aws.String(date)},
				"Transaction":   {S: aws.String(fmt.Sprintf("%.2f", t.Transaction))},
			}
			writeRequests = append(writeRequests, &dynamodb.WriteRequest{
				PutRequest: &dynamodb.PutRequest{
					Item: item,
				},
			})
		}
		_, err := sa.svc.BatchWriteItem(&dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				tableName: writeRequests,
			},
		})
		if err != nil {
			return fmt.Errorf("error en BatchWriteItem: %w", err)
		}
	}

	return nil
}
