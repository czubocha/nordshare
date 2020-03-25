package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	codedeploysrv "github.com/aws/aws-sdk-go/service/codedeploy"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	lambdasrv "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/kelseyhightower/envconfig"
	"github/czubocha/nordshare/cmd/e2e/note"
)

var h *note.Handler

type envVars struct {
	Saver     string `required:"true"`
	Reader    string `required:"true"`
	Modifier  string `required:"true"`
	Remover   string `required:"true"`
	TableName string `required:"true" envconfig:"TABLE_NAME"`
}

func init() {
	var env envVars
	envconfig.MustProcess("version", &env)
	h = resolveHandler(env)
}

func resolveHandler(env envVars) *note.Handler {
	sess := session.Must(session.NewSession())
	db := dynamodb.New(sess)
	codeDeployService := codedeploysrv.New(sess)
	lambdaService := lambdasrv.New(sess)
	versions := convert(env)
	return note.NewHandler(
		note.Services{DynamoDB: db, Lambda: lambdaService, CodeDeploy: codeDeployService},
		versions, env.TableName)
}

func convert(env envVars) note.Versions {
	return note.Versions{
		Saver:    env.Saver,
		Reader:   env.Reader,
		Modifier: env.Modifier,
		Remover:  env.Remover,
	}
}

func handle(event codedeploysrv.PutLifecycleEventHookExecutionStatusInput) error {
	return note.HandleRequest(event, h)
}

func main() {
	lambda.Start(handle)
}
