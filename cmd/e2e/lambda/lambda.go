package lambda

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/google/uuid"
	"log"
	"text/tabwriter"
)

func InvokeLambda(l *lambda.Lambda, version string, headers, pathParameters map[string]string, body, lambdaOut interface{}) (status float64, err error) {
	var marshalledBody []byte
	if body != nil {
		if marshalledBody, err = json.Marshal(body); err != nil {
			return
		}
	}
	payload := events.APIGatewayProxyRequest{
		RequestContext: events.APIGatewayProxyRequestContext{RequestID: uuid.New().String()},
		Headers:        headers,
		PathParameters: pathParameters,
		Body:           string(marshalledBody)}
	marshalledBody, err = json.Marshal(payload)
	if err != nil {
		return
	}
	invokeInput := lambda.InvokeInput{
		FunctionName: &version,
		Payload:      marshalledBody}
	req, out := l.InvokeRequest(&invokeInput)
	if err = req.Send(); err != nil {
		return
	}
	var resp map[string]interface{}
	if err = json.Unmarshal(out.Payload, &resp); err != nil {
		return
	}
	status = resp["statusCode"].(float64)
	defer printInvocationResult(version, headers, pathParameters, status, lambdaOut)
	if lambdaOut == nil {
		return
	}
	if resp["body"] == "" {
		err = fmt.Errorf("invoke lambda: body expected but empty")
		return
	}
	err = json.Unmarshal([]byte(resp["body"].(string)), lambdaOut)
	return
}

func printInvocationResult(lambda string, headers, pathParameters map[string]string, status float64, output interface{}) {
	buf := &bytes.Buffer{}
	w := tabwriter.NewWriter(buf, 0, 0, 3, ' ', tabwriter.TabIndent)
	if _, err := fmt.Fprintf(w, "\nlambda\t%v\nheaders\t%v\npath parameters\t%v\nstatus\t%v\noutput\t%+v\n\n",
		lambda, headers, pathParameters, status, output); err != nil {
		log.Println(err)
	}
	if err := w.Flush(); err != nil {
		log.Println(err)
	}
	log.Print(buf.String())
}
