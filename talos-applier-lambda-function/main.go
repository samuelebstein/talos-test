package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// Assuming MyEvent is structured for an AWS Lambda trigger, customize as needed
type MyEvent struct {
	Name string `json:"name"` // Example field, adjust based on actual event data
}

// MyResponse indicates the result of the Lambda function execution
type MyResponse struct {
	Message string `json:"message"`
}

func getSecret(ctx context.Context, secretName string) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config, %v", err)
	}

	client := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := client.GetSecretValue(ctx, input)
	if err != nil {
		return "", fmt.Errorf("error retrieving secret %s: %v", secretName, err)
	}

	return *result.SecretString, nil
}

func HandleRequest(ctx context.Context, event MyEvent) (MyResponse, error) {

	secretName := "sam-ebstein-test-talosconfig"
	_, err := getSecret(ctx, secretName)
	if err != nil {
		return MyResponse{}, err
	}

	// If the secret was retrieved successfully, you can proceed with your logic
	// For this example, we're just indicating success without exposing the secret content
	return MyResponse{Message: "Secret retrieved successfully"}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
