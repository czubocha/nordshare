package api

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"strings"
)

func NewResponse(status int, body ...interface{}) (response events.APIGatewayProxyResponse, err error) {
	response = events.APIGatewayProxyResponse{
		StatusCode: status,
	}
	if len(body) == 0 {
		return
	}
	bytes, err := json.Marshal(body[0])
	if err != nil {
		log.Printf("api: %v", err)
		return
	}
	response.Body = string(bytes)
	return
}

func GetHeaderIncasesensible(headers map[string]string, key string) string {
	for k, v := range headers {
		if strings.ToLower(k) == strings.ToLower(key) {
			return v
		}
	}
	return ""
}
