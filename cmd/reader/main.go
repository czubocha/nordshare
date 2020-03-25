package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/kelseyhightower/envconfig"
	"github/czubocha/nordshare/cmd/reader/api"
	"github/czubocha/nordshare/pkg/encryption"
	"github/czubocha/nordshare/pkg/storage"
)

var h *api.Handler

type envVars struct {
	TableName string `required:"true" envconfig:"TABLE_NAME"`
	KeyID     string `required:"true" envconfig:"KEY_ID"`
}

func init() {
	var env envVars
	envconfig.MustProcess("", &env)
	h = resolveHandler(env)
}

func resolveHandler(env envVars) *api.Handler {
	sess := session.Must(session.NewSession())
	xray.AWSSession(sess)
	db := dynamodb.New(sess)
	repository := storage.NewRepository(db, env.TableName)
	k := kms.New(sess)
	encryptionSrv := encryption.NewService(env.KeyID, k)
	return api.NewHandler(repository, encryptionSrv)
}

func handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return api.HandleRequest(ctx, request, h)
}

func main() {
	lambda.Start(handle)
}
