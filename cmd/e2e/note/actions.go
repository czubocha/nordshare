package note

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	lambdasrv "github.com/aws/aws-sdk-go/service/lambda"
	"github/czubocha/nordshare/cmd/e2e/lambda"
	"net/http"
	"time"
)

const (
	content       = "test note"
	readPassword  = "abc"
	writePassword = "def"
	ttl           = 60

	modifiedContent = "modified note"
	modifiedTtl     = 30

	idPathParamName    = "id"
	passwordHeaderName = "password"
)

type id struct {
	ID string
}

type CreateNote struct {
	Content       string `json:"content"`
	ReadPassword  string `json:"readPassword"`
	WritePassword string `json:"writePassword"`
	TTL           int64  `json:"ttl"`
}

type ReadModifyNote struct {
	Content string `json:"content"`
	TTL     int64  `json:"ttl"`
}

// check if status is 201
func create(l *lambdasrv.Lambda, version string, note CreateNote) (string, error) {
	noteId := &id{}
	status, err := lambda.InvokeLambda(l, version, nil, nil, note, noteId)
	if err != nil {
		return "", fmt.Errorf("saver lambda: %w", err)
	}
	if status != http.StatusCreated {
		return "", fmt.Errorf("create note: expected status: %v got %v", http.StatusCreated, status)
	}
	return noteId.ID, nil
}

// check if 1) status is 200 2) content and TTL is valid
func get(l *lambdasrv.Lambda, version, id, password string, expectedNote ReadModifyNote) error {
	n := &ReadModifyNote{}
	status, err := lambda.InvokeLambda(l, version, map[string]string{passwordHeaderName: password},
		map[string]string{idPathParamName: id}, nil, n)
	if err != nil {
		return fmt.Errorf("reader lambda: %w", err)
	}
	if status != http.StatusOK {
		return fmt.Errorf("read note: expected status: %v got %v", http.StatusOK, status)
	}
	if n.Content != expectedNote.Content {
		return fmt.Errorf("read note: expected content: %v got %v", expectedNote.Content, n.Content)
	}
	if expectedNote.TTL < n.TTL {
		return fmt.Errorf("read note: ttl greater than original: original: %v got %v", expectedNote.TTL, n.TTL)
	}
	return nil
}

// check if status is 200
func modify(l *lambdasrv.Lambda, version, id, password string, modification ReadModifyNote) error {
	status, err := lambda.InvokeLambda(l, version, map[string]string{passwordHeaderName: password},
		map[string]string{idPathParamName: id}, modification, nil)
	if err != nil {
		return fmt.Errorf("modifier lambda: %w", err)
	}
	if status != http.StatusOK {
		return fmt.Errorf("modify note: expected status: %v got %v", http.StatusOK, status)
	}
	return nil
}

// check if status is 200
func remove(l *lambdasrv.Lambda, version, password, id string) error {
	status, err := lambda.InvokeLambda(l, version, map[string]string{passwordHeaderName: password},
		map[string]string{idPathParamName: id}, nil, nil)
	if err != nil {
		return fmt.Errorf("remover lambda: %w", err)
	}
	if status != http.StatusOK {
		return fmt.Errorf("remove note: expected status: %v got %v", http.StatusOK, status)
	}
	return nil
}

func expire(id, tableName string, db *dynamodb.DynamoDB) error {
	const ttlAttributeName, ttlExpressionAttributeName, ttlExpressionAttributeValue = "ttl", "#ttl", ":t"
	const keyAttributeName = "id"
	passedTTL := time.Now().Add(-10 * time.Minute).Unix()
	passedTTLAtt, err := dynamodbattribute.Marshal(passedTTL)
	if err != nil {
		return fmt.Errorf("unable to marshal ttl: %w", err)
	}
	input := dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{ttlExpressionAttributeValue: passedTTLAtt},
		ExpressionAttributeNames:  map[string]*string{ttlExpressionAttributeName: aws.String(ttlAttributeName)},
		Key:                       map[string]*dynamodb.AttributeValue{keyAttributeName: {S: aws.String(id)}},
		TableName:                 aws.String(tableName),
		UpdateExpression: aws.String(fmt.Sprintf("set %v = %v",
			ttlExpressionAttributeName, ttlExpressionAttributeValue)),
	}
	if _, err = db.UpdateItem(&input); err != nil {
		return fmt.Errorf("unable to update item: %w", err)
	}
	return nil
}
