package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"nordshare/pkg/note"
	"time"
)

type Note struct {
	ID            string `json:"id"`
	Content       []byte `json:"content"`
	ReadPassword  []byte `json:"readPassword"`
	WritePassword []byte `json:"writePassword"`
	TTL           int64  `json:"ttl"`
}

type repository struct {
	*dynamodb.DynamoDB
	tableName string
}

func NewRepository(client *dynamodb.DynamoDB, tableName string) *repository {
	return &repository{client, tableName}
}

func (r *repository) SaveNote(ctx context.Context, note note.Note, id string) error {
	storageNote := convert(note, id)
	item, err := dynamodbattribute.MarshalMap(storageNote)
	if err != nil {
		return fmt.Errorf("storage: %w", err)
	}
	input := dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(r.tableName),
	}
	if _, err = r.PutItemWithContext(ctx, &input); err != nil {
		return fmt.Errorf("storage: %w", err)
	}
	return nil
}

func convert(note note.Note, id string) Note {
	n := Note{
		ID:            id,
		Content:       note.Content,
		ReadPassword:  note.ReadPassword,
		WritePassword: note.WritePassword,
	}
	if note.TTL > 0 {
		n.TTL = time.Now().Add(time.Duration(note.TTL) * time.Minute).Unix()
	} else {
		n.TTL = time.Now().Add(24 * time.Hour).Unix()
	}
	return n
}
