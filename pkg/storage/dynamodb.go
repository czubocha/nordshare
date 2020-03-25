package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github/czubocha/nordshare"
	"log"
	"math"
	"time"
)

const keyAttributeName = "id"

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

func (r *repository) SaveNote(ctx context.Context, note nordshare.Note, id string) error {
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
	ttlAttVal, err := dynamodbattribute.Marshal(nowPlusMinutes(ttl, time.Now()))
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
		Key:       map[string]*dynamodb.AttributeValue{keyAttributeName: {S: aws.String(id)}},
		TableName: aws.String(r.tableName),
		UpdateExpression: aws.String(fmt.Sprintf("set %v = %v, %v = %v",
			contentAttributeName, contentExpressionAttributeValue, ttlExpressionAttributeName, ttlExpressionAttributeValue)),
	}
	if _, err := r.UpdateItemWithContext(ctx, &input); err != nil {
		return fmt.Errorf("storage: %w", err)
	}
	return nil
}

func (r *repository) ReadNote(ctx context.Context, id string) (nordshare.Note, error) {
	input := dynamodb.GetItemInput{
		Key:       map[string]*dynamodb.AttributeValue{keyAttributeName: {S: aws.String(id)}},
		TableName: aws.String(r.tableName)}
	output, err := r.GetItemWithContext(ctx, &input)
	if err != nil {
		log.Printf("storage: %v", err)
		return nordshare.Note{}, ErrNoteExpired
	}
	n := Note{}
	if err = dynamodbattribute.UnmarshalMap(output.Item, &n); err != nil {
		return nordshare.Note{}, fmt.Errorf("storage: %w", err)
	}
	if isExpired(n.TTL, time.Now()) {
		return nordshare.Note{}, ErrNoteExpired
	}
	return convertFromStorage(n), nil
}

func (r *repository) DeleteNote(ctx context.Context, id string) error {
	input := dynamodb.DeleteItemInput{
		Key:       map[string]*dynamodb.AttributeValue{keyAttributeName: {S: aws.String(id)}},
		TableName: aws.String(r.tableName)}
	_, err := r.DeleteItemWithContext(ctx, &input)
	if err != nil {
		return fmt.Errorf("storage: %w", err)
	}
	return nil
}

func isExpired(ttl int64, now time.Time) bool {
	return now.After(time.Unix(ttl, 0))
}

func convertToStorage(note nordshare.Note, id string) Note {
	return Note{
		ID:            id,
		Content:       note.Content,
		ReadPassword:  note.ReadPassword,
		WritePassword: note.WritePassword,
		TTL:           calculateTTL(note.TTL, time.Now())}
}

func calculateTTL(minutesToExpire int64, now time.Time) int64 {
	if minutesToExpire > 0 {
		return nowPlusMinutes(minutesToExpire, now)
	}
	return tomorrow(now)
}

func tomorrow(now time.Time) int64 {
	return now.Add(24 * time.Hour).Unix()
}

func nowPlusMinutes(minutes int64, now time.Time) int64 {
	return now.Add(time.Duration(minutes) * time.Minute).Unix()
}

func convertFromStorage(n Note) nordshare.Note {
	return nordshare.Note{
		Content:       n.Content,
		ReadPassword:  n.ReadPassword,
		WritePassword: n.WritePassword,
		TTL:           getRemainingMinutes(n.TTL, time.Now()),
	}
}

func getRemainingMinutes(date int64, now time.Time) int64 {
	durationTillDate := time.Unix(date, 0).Sub(now)
	minutesRoundedUp := math.Ceil(durationTillDate.Minutes())
	if minutesRoundedUp < 0 {
		return 0
	}
	return int64(minutesRoundedUp)
}
