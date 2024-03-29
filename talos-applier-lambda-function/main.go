package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	machineapi "github.com/siderolabs/talos/pkg/machinery/api/machine"
	"github.com/siderolabs/talos/pkg/machinery/client"
	// clientconfig "github.com/siderolabs/talos/pkg/machinery/client/config"
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
	// TODO: should move client configuration outside of function declaration
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

	// Check if the instance state is not "running" because my eventbridge rule wasn't working
	if ec2Event.Detail.State != "running" {
		// Log and skip execution. Maybe error in future as lambda shouldn't be invoked on other state changes
		fmt.Printf("Skipping execution as the instance state is '%s', not 'running'.\n", ec2Event.Detail.State)
		return MyResponse{Message: fmt.Sprintf("Skipped execution for instance %s as its state is '%s'.", ec2Event.Detail.InstanceID, ec2Event.Detail.State)}, nil
	}

	fmt.Printf("Received EC2 state change event: Instance ID %s, State %s\n", ec2Event.Detail.InstanceID, ec2Event.Detail.State)

	// Retrieve the public IP address of the instance
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

	// Should be able to simplify this found logic with the assumption that we're only querying for details of one instance
	var newWorkerNodeIP string
	found := false
	for _, reservation := range describeInstancesOutput.Reservations {
		for _, instance := range reservation.Instances {
			if instance.PublicIpAddress != nil { // Add nil check here
				newWorkerNodeIP = *instance.PublicIpAddress
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

	fmt.Printf("Public IP for instance %s is %s\n", ec2Event.Detail.InstanceID, newWorkerNodeIP)

	// Retrieve the worker config secret to apply to node
	workerConfigSecretName := "sam-ebstein-test-talos-worker-yaml"
	workerConfigSecretString, err := getSecret(ctx, workerConfigSecretName)
	if err != nil {
		return MyResponse{}, fmt.Errorf("failed to retrieve secret %s: %v", workerConfigSecretName, err)
	}

	fmt.Printf("Retrieved secret %s", workerConfigSecretName)

	// https://www.talos.dev/v1.6/introduction/getting-started/
	// "The --insecure flag is necessary because the PKI infrastructure has not yet been made available to the node. Note: the connection will be encrypted, but not authenticated.
	// When using the --insecure flag, it is not necessary to specify an endpoint."
	// https://github.com/siderolabs/talos/blob/f02aeec922b6327dad6d4fee917987b147abbf2a/pkg/cluster/apply-config.go#L31-L51
	// Apply config code above from siderolabs actually sends the endpoints as the node ip.. which seems opposite to me as usually the node is what is getting targeted.
	talosClient, err := client.New(ctx, client.WithTLSConfig(&tls.Config{
		InsecureSkipVerify: true,
	}), client.WithEndpoints(newWorkerNodeIP))
	if err != nil {
		return MyResponse{}, fmt.Errorf("failed to create talos client", err)
	}
	// Closing the client after function execution. Not sure if it autocloses or issues with not closing explicitly
	defer talosClient.Close()
	fmt.Printf("Created talos client")

	// GetEndpoints returns the client's endpoints from the override set with WithEndpoints or from the configuration.
	// I wouldn't think we need any endpoints but the client initialization checks if there's > 1 endpoints?
	fmt.Printf("Endpoints: %v", talosClient.GetEndpoints())

	// TODO: there should be a check on the instance being succesfully initialized/"healthy" before running apply configuration. Or have retries enabled. Because lambda failed but then retried due to default lambda retries

	// Apply the configuration to the node
	req := &machineapi.ApplyConfigurationRequest{
		Data:   []byte(workerConfigSecretString),
		Mode:   machineapi.ApplyConfigurationRequest_AUTO,
		DryRun: false,
	}
	_, err = talosClient.ApplyConfiguration(ctx, req)
	if err != nil {
		return MyResponse{}, fmt.Errorf("failed to apply configuration to instance %s with public ip %s, %v", ec2Event.Detail.InstanceID, newWorkerNodeIP, err)
	}

	// How to automatically test that applying worked?

	// Is running state the only configuration to check? We might need to check the status of the instance ie "ready"?
	return MyResponse{Message: "Configuration process completed."}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
