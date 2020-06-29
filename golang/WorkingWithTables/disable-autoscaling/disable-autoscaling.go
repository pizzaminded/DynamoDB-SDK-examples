package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/applicationautoscaling"
    "fmt"
)

var tableName      = "Music"
var roleARN        = "arn:aws:iam::618326157558:role/Music_TableScalingRole"
var policyName     = fmt.Sprintf("%s_%s", tableName, "TableScalingPolicy")
var resourceID     = fmt.Sprintf("%s/%s", "table", tableName)
var readDimension  = "dynamodb:table:ReadCapacityUnits"
var writeDimension = "dynamodb:table:WriteCapacityUnits"

func getSession() (*session.Session) {
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
        // Provide SDK Config options, such as Region and Endpoint
        Config: aws.Config{
            Region: aws.String("us-west-2"),
	    },
    }))

    return sess
}

func registerScalableTarget(
    autoscalingClient *applicationautoscaling.ApplicationAutoScaling,
    dimension string,
    resourceID string,
    roleARN string,
) {
    input := &applicationautoscaling.RegisterScalableTargetInput{
        MaxCapacity:       aws.Int64(500),
        MinCapacity:       aws.Int64(1),
        ResourceId:        aws.String(resourceID),
        RoleARN:           aws.String(roleARN),
        ScalableDimension: aws.String(dimension),
        ServiceNamespace:  aws.String("dynamodb"),
    }
    autoscalingClient.RegisterScalableTarget(input)
}

func deleteScalingPolicy(
	autoscalingClient *applicationautoscaling.ApplicationAutoScaling,
	dimension string,
) {
	input := &applicationautoscaling.DeleteScalingPolicyInput{
		PolicyName:        aws.String(policyName),
		ResourceId:        aws.String(resourceID),
		ServiceNamespace:  aws.String("dynamodb"),
		ScalableDimension: aws.String(dimension),
	}
	autoscalingClient.DeleteScalingPolicy(input)
}

func deregisterScalableTarget(
	autoscalingClient *applicationautoscaling.ApplicationAutoScaling,
	dimension string,
) {
	input := &applicationautoscaling.DeregisterScalableTargetInput{
		ResourceId:        aws.String(resourceID),
		ServiceNamespace:  aws.String("dynamodb"),
		ScalableDimension: aws.String(dimension),
	}
	autoscalingClient.DeregisterScalableTarget(input)
}

func disableAutoscaling() {
    autoscalingClient := applicationautoscaling.New(getSession())

    deleteScalingPolicy(autoscalingClient, readDimension)
    fmt.Println("Read scaling policy deleted ...")

    deleteScalingPolicy(autoscalingClient, writeDimension)
    fmt.Println("Write scaling policy deleted ...")

    deregisterScalableTarget(autoscalingClient, readDimension)
    fmt.Println("Write scalable target registered ...")

    deregisterScalableTarget(autoscalingClient, writeDimension)
    fmt.Println("Write scalable target registered ...")
}

func main() {
    fmt.Println("Updating table to enable autoscaling ...")
    disableAutoscaling()
    fmt.Println("Finished ...")
}
