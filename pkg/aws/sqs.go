package aws

import (
	"fmt"
	"heartbeat/pkg/hooks"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func generateAwsSession() *session.Session {
	// Check if region env var was set or use defult
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "eu-west-2"
	}
	// initialize aws session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)
	if err != nil {
		fmt.Printf("Error creating new aws session {%v}\n", err)
	}
	return sess
}

func checkIfSQSQueueExists(queueName string) string {
	var queueURL string
	awsSession := generateAwsSession()
	svc := sqs.New(awsSession)
	// set queue name prefix filter
	queueInput := sqs.ListQueuesInput{
		QueueNamePrefix: aws.String(queueName),
	}
	// Get all queues with prefix
	queueList, _ := svc.ListQueues(&queueInput)

	// Check if create queue
	if len(queueList.QueueUrls) == 0 {
		fmt.Printf("Creating %v queue\n", queueName)
		queueURL = createHeartbeatSQSQueue(queueName)
	} else {
		// Search for exact queue name if there are queues with the same prefix
		for _, q := range queueList.QueueUrls {
			// fmt.Printf("Checking %v\n", *q)
			if strings.HasSuffix(*q, queueName) {
				queueURL = *q
			}
		}
	}
	return queueURL
}

func createHeartbeatSQSQueue(queueName string) string {
	awsSession := generateAwsSession()
	svc := sqs.New(awsSession)
	// create queue with default setting
	createQueueInput := &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
		Attributes: map[string]*string{
			"DelaySeconds":           aws.String("0"),
			"MessageRetentionPeriod": aws.String("86400"),
		},
	}
	result, err := svc.CreateQueue(createQueueInput)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	queueURL := result.QueueUrl
	return *queueURL
}

// SendAliveSQSMessage sends a message to SQS Queue
func SendAliveSQSMessage() {
	// heartbeat sqs queue name
	queueName := "heartbeat"
	sendSQSMessage(queueName)
}

// CheckAliveSQSMessage check heartbeat sqs messages
func CheckAliveSQSMessage() {
	// check queue name
	checkerQueueName := "heartbeatChecker"
	// check if exists and get its url back
	checkerQueueURL := checkIfSQSQueueExists(checkerQueueName)

	// send queue name
	queueName := "heartbeat"
	// check if exists and get its url back
	queueURL := checkIfSQSQueueExists(queueName)
	// count send queue count
	heartbeatCount := countMessagesInQueue(queueURL)

	fmt.Printf("Current messages in queue: %v is %v\n", queueName, heartbeatCount)
	if heartbeatCount > 0 {
		// Heartbeat was sent and all is ok
		fmt.Println("Heartbeat was send purging all queues")
		// empty queue
		deleteMessagesInQueue(queueURL)
		// count check queue
		checkQueueCount := countMessagesInQueue(checkerQueueURL)
		if checkQueueCount > 2 {
			fmt.Println("Heartbeat has recovered")
			hooks.SendTelegramMessage("Heartbeat has recovered")
		} else {
			fmt.Println("Heartbeat is up")
		}
		deleteMessagesInQueue(checkerQueueURL)
	} else {
		sendSQSMessage("heartbeatChecker")
		checkIfHeartbeatDown(checkerQueueURL)
	}
}

func checkIfHeartbeatDown(queueURL string) {
	checkerQueueCount := countMessagesInQueue(queueURL)
	if checkerQueueCount > 2 {
		fmt.Println("Heartbeat is down")
		hooks.SendTelegramMessage("Heartbeat is down")
	}
}

func sendSQSMessage(queueName string) {
	queueURL := checkIfSQSQueueExists(queueName)
	awsSession := generateAwsSession()
	svc := sqs.New(awsSession)
	sqsMessage := &sqs.SendMessageInput{
		MessageBody: aws.String("heartbeat"),
		QueueUrl:    aws.String(queueURL),
	}
	_, err := svc.SendMessage(sqsMessage)

	if err != nil {
		fmt.Printf("Failed to send message %v", err)
		os.Exit(1)
	}
	fmt.Printf("%v message sent\n", queueName)
}

func deleteMessagesInQueue(queueURL string) {
	awsSession := generateAwsSession()
	svc := sqs.New(awsSession)
	// Purge queue is limited by AWS API to 60 seconds
	purgeQueueInput := &sqs.PurgeQueueInput{
		QueueUrl: aws.String(queueURL),
	}
	fmt.Printf("Purging queue %v\n", queueURL)
	_, err := svc.PurgeQueue(purgeQueueInput)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func countMessagesInQueue(queueURL string) int {
	awsSession := generateAwsSession()
	svc := sqs.New(awsSession)
	// set input for messages in queue
	queueAttrib := &sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(queueURL),
		AttributeNames: aws.StringSlice([]string{"ApproximateNumberOfMessages"}),
	}
	// get queue count
	queueCount, err := svc.GetQueueAttributes(queueAttrib)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	count, _ := strconv.Atoi(*queueCount.Attributes["ApproximateNumberOfMessages"])
	return count
}
