#! /bin/sh

test_queue=test_queue
sample_queue=sample_queue
sample_dlq=sample_dlq

awslocal sqs create-queue --queue-name ${test_queue}

awslocal sqs create-queue --queue-name ${sample_dlq}
SQS_DEADLETTERQUEUE_URL=$(awslocal sqs get-queue-url --queue-name ${sample_dlq} --output text) \
                        && echo ${SQS_DEADLETTERQUEUE_URL}
SQS_DEADLETTERQUEUE_ARN=$(awslocal sqs get-queue-attributes --queue-url ${SQS_DEADLETTERQUEUE_URL} \
                          --attribute-names QueueArn --output text | sed -e 's/ATTRIBUTES\s//' ) \
                         && echo ${SQS_DEADLETTERQUEUE_ARN}

awslocal sqs create-queue --queue-name ${sample_queue}
SQS_QUEUE_URL=$(awslocal sqs get-queue-url --queue-name ${sample_queue} --output text) && echo ${SQS_QUEUE_URL}
cat << EOF >> deadletter.json
{
	"RedrivePolicy": "{\"deadLetterTargetArn\":\"${SQS_DEADLETTERQUEUE_ARN}\",\"maxReceiveCount\":\"5\"}"
}
EOF
awslocal sqs set-queue-attributes --queue-url ${SQS_QUEUE_URL} --attributes file://deadletter.json

rm deadletter.json