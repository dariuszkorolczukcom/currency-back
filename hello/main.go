package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	fmt.Println(request.Body)
	var buf bytes.Buffer

	var forExVals map[string]float64
	forExVals = make(map[string]float64)
	forExVals["GBP_USD"] = 1.37
	forExVals["GBP_AUD"] = 1.79
	forExVals["GBP_EUR"] = 1.13

	var exchange Exchange

	err := json.Unmarshal([]byte(request.Body), &exchange)
	if err != nil {
		fmt.Println("error:", err)
	}

	exchange.count(forExVals[exchange.Currencies])

	fmt.Printf("%+v", exchange)

	body, err := json.Marshal(exchange)
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}

type Exchange struct {
	Currencies string
	ValueFrom  float64
	ValueTo    float64
}

func (e *Exchange) count(ev float64) {
	e.ValueTo = math.Floor(e.ValueFrom*100*ev) / 100
}
