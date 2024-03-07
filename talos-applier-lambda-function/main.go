package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// EventDetail represents the detail of the EventBridge event for EC2 state change
type EventDetail struct {
	InstanceID string `json:"instance-id"`
	State      string `json:"state"`
}

// EC2StateChangeEvent represents the structure of an EC2 state change event from EventBridge
type EC2StateChangeEvent struct {
	Detail EventDetail `json:"detail"`
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

func HandleRequest(ctx context.Context, event json.RawMessage) (MyResponse, error) {
	var ec2Event EC2StateChangeEvent
	if err := json.Unmarshal(event, &ec2Event); err != nil {
		return MyResponse{}, fmt.Errorf("failed to unmarshal event: %v", err)
	}

	// Check if the instance state is not "running"
	if ec2Event.Detail.State != "running" {
		// Log and skip execution
		fmt.Printf("Skipping execution as the instance state is '%s', not 'running'.\n", ec2Event.Detail.State)
		return MyResponse{Message: fmt.Sprintf("Skipped execution for instance %s as its state is '%s'.", ec2Event.Detail.InstanceID, ec2Event.Detail.State)}, nil
	}

	fmt.Printf("Received EC2 state change event: Instance ID %s, State %s\n", ec2Event.Detail.InstanceID, ec2Event.Detail.State)

	// Example: Retrieve the public IP address of the instance
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return MyResponse{}, fmt.Errorf("unable to load SDK config, %v", err)
	}
	ec2Client := ec2.NewFromConfig(cfg)

	describeInstancesOutput, err := ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{ec2Event.Detail.InstanceID},
	})
	if err != nil {
		return MyResponse{}, fmt.Errorf("failed to describe instances: %v", err)
	}

	var publicIP string
	found := false
	for _, reservation := range describeInstancesOutput.Reservations {
		for _, instance := range reservation.Instances {
			if instance.PublicIpAddress != nil { // Add nil check here
				publicIP = *instance.PublicIpAddress
				found = true
				break // Assuming only one instance per reservation for simplicity
			}
		}
		if found {
			break // Found the public IP, no need to check further
		}
	}

	if !found {
		msg := fmt.Sprintf("Public IP not found for instance %s. Exiting function.", ec2Event.Detail.InstanceID)
		fmt.Println(msg)
		return MyResponse{Message: msg}, nil
	}

	fmt.Printf("Public IP for instance %s is %s\n", ec2Event.Detail.InstanceID, publicIP)

	// Proceed with applying the configuration using the public IP...
	return MyResponse{Message: "Configuration process can be initiated"}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
