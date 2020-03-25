package note

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codedeploy"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	lambdasrv "github.com/aws/aws-sdk-go/service/lambda"
	"log"
)

type Services struct {
	*dynamodb.DynamoDB
	*lambdasrv.Lambda
	*codedeploy.CodeDeploy
}

type Handler struct {
	Services
	Versions
	TableName string
}

type scenario func(*Handler) error

func HandleRequest(input codedeploy.PutLifecycleEventHookExecutionStatusInput, h *Handler) error {
	input.Status = aws.String(codedeploy.LifecycleEventStatusFailed)
	defer putLifecycleEventHook(h.CodeDeploy, &input)
	scenarios := []scenario{FullFlowScenario, ExpirationScenario}
	for _, scenario := range scenarios {
		if err := scenario(h); err != nil {
			log.Print(err)
			return nil
		}
	}
	input.Status = aws.String(codedeploy.LifecycleEventStatusSucceeded)
	return nil
}

func putLifecycleEventHook(deploy *codedeploy.CodeDeploy, input *codedeploy.PutLifecycleEventHookExecutionStatusInput) {
	log.Printf("sending execution status: %+v", input)
	_, err := deploy.PutLifecycleEventHookExecutionStatus(input)
	if err != nil {
		log.Print(err)
	}
}

func NewHandler(services Services, versions Versions, tableName string) *Handler {
	return &Handler{Services: services, Versions: versions, TableName: tableName}
}

type Versions struct {
	Saver    string
	Reader   string
	Modifier string
	Remover  string
}
