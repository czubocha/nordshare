package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"math"
	"nordshare/pkg/note"
	"time"
)

type Error string

func (e Error) Error() string { return string(e) }

const ErrNoteExpired = Error("storage: note has expired")

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
	storageNote := convertToStorage(note, id)
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

func (r *repository) UpdateNote(ctx context.Context, content []byte, ttl int64, id string) error {
	const (
		contentAttributeName            = "content"
		contentExpressionAttributeValue = ":c"
		ttlAttributeName                = "ttl"
		ttlExpressionAttributeName      = "#ttl"
		ttlExpressionAttributeValue     = ":t"
	)
	ttlAttVal, err := dynamodbattribute.Marshal(nowPlusMinutes(ttl))
	if err != nil {
		return fmt.Errorf("storage: %w", err)
	}
	input := dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			contentExpressionAttributeValue: {B: content},
			ttlExpressionAttributeValue:     ttlAttVal},
		ExpressionAttributeNames: map[string]*string{
			ttlExpressionAttributeName: aws.String(ttlAttributeName),
		},
		Key:       map[string]*dynamodb.AttributeValue{"id": {S: aws.String(id)}},
		TableName: aws.String(r.tableName),
		UpdateExpression: aws.String(fmt.Sprintf("set %v = %v, %v = %v",
			contentAttributeName, contentExpressionAttributeValue, ttlExpressionAttributeName, ttlExpressionAttributeValue)),
	}
	if _, err := r.UpdateItemWithContext(ctx, &input); err != nil {
		return fmt.Errorf("storage: %w", err)
	}
	return nil
}

func (r *repository) ReadNote(ctx context.Context, id string) (note.Note, error) {
	input := dynamodb.GetItemInput{
		Key:       map[string]*dynamodb.AttributeValue{"id": {S: aws.String(id)}},
		TableName: aws.String(r.tableName)}
	output, err := r.GetItemWithContext(ctx, &input)
	if err != nil {
		log.Printf("storage: %v", err)
		return note.Note{}, ErrNoteExpired
	}
	n := Note{}
	if err = dynamodbattribute.UnmarshalMap(output.Item, &n); err != nil {
		return note.Note{}, fmt.Errorf("storage: %w", err)
	}
	if isExpired(n.TTL) {
		return note.Note{}, ErrNoteExpired
	}
	return convertFromStorage(n), nil
}

func (r *repository) DeleteNote(ctx context.Context, id string) error {
	input := dynamodb.DeleteItemInput{
		Key:       map[string]*dynamodb.AttributeValue{"id": {S: aws.String(id)}},
		TableName: aws.String(r.tableName)}
	_, err := r.DeleteItemWithContext(ctx, &input)
	if err != nil {
		return fmt.Errorf("storage: %w", err)
	}
	return nil
}

func isExpired(ttl int64) bool {
	return time.Now().After(time.Unix(ttl, 0))
}

func convertToStorage(note note.Note, id string) Note {
	n := Note{
		ID:            id,
		Content:       note.Content,
		ReadPassword:  note.ReadPassword,
		WritePassword: note.WritePassword,
	}
	if note.TTL > 0 {
		n.TTL = nowPlusMinutes(note.TTL)
	} else {
		n.TTL = tomorrow()
	}
	return n
}

func tomorrow() int64 {
	return time.Now().Add(24 * time.Hour).Unix()
}

func nowPlusMinutes(minutes int64) int64 {
	return time.Now().Add(time.Duration(minutes) * time.Minute).Unix()
}

func convertFromStorage(n Note) note.Note {
	return note.Note{
		Content:       n.Content,
		ReadPassword:  n.ReadPassword,
		WritePassword: n.WritePassword,
		TTL:           getRemainingMinutes(n.TTL),
	}
}

func getRemainingMinutes(unixTime int64) int64 {
	return int64(math.Ceil(time.Until(time.Unix(unixTime, 0)).Minutes()))
}
