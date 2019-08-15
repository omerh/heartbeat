# Heartbeat

This project is designed to monitor internet connectivity/power outages async using AWS SQS Queues.
When there is a connectivity issues, after 2 checks a telegram message will be sent to you
Can be used only with the single heartbeat binary using cobra.

The monitored site run `./heartbeat send`
The checker site run `./heartbeat check`

## Build

Running `make` will build 2 binaries for linux os, one for running on linux host and the second for running in lambda

```bash
make
```

## Create a telegram channel for alerts

1. Start a chat with the BotFather
2. send `/newbot`and follow the instructions
3. once finished, create a new channel and add the new bot as admin and allow it to push messages

### Monitored site

Configure aws credentials and `AWS_REGION` environment variable, if you want it not to send to the default region **`eu-west-2`**
Set a cron for a minimum every 1 min or more that runs `./heartbeat send`

### Checking site

Do the same as the monitored site, just run `./heartbeat check`

### Checking lambda

Create the Lambda rule to have access to 2 sqs queues as follow:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "sqs:GetQueueUrl",
                "sqs:PurgeQueue",
                "sqs:ReceiveMessage",
                "sqs:SendMessage",
                "sqs:GetQueueAttributes",
                "sqs:CreateQueue"
            ],
            "Resource": [
                "arn:aws:sqs:*:*:heartbeatChecker",
                "arn:aws:sqs:*:*:heartbeat"
            ]
        },
        {
            "Sid": "VisualEditor1",
            "Effect": "Allow",
            "Action": "sqs:ListQueues",
            "Resource": "*"
        }
    ]
}
```

>Can be reduces to SQS Account and region specific

The ARN gives access to the following SQS

1. heartbeat
2. heartbeatChecker

Create AWS Lambda in the region you have set (defaults to eu-west-2) with the following settings:

1. runtime is Go1.x
2. memory 128m
3. timeout 5 sec
4. role you've created with lambda-basic-execution and sqs access
5. add environment variables `TELEGRAM_TOKEN` and `TELEGRAM_CHANNEL` that you've created before
6. create trigger minimum rate(1 minute)

Enjoy!
