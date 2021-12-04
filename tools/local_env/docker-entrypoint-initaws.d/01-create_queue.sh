#! /bin/sh

awslocal sqs create-queue --queue-name sample_dlQ

SQS_DEADLETTERQUEUE_URL=$(awslocal sqs get-queue-url --queue-name sample_dlQ --output text) \
                        && echo ${SQS_DEADLETTERQUEUE_URL}

SQS_DEADLETTERQUEUE_ARN=$(awslocal sqs get-queue-attributes --queue-url ${SQS_DEADLETTERQUEUE_URL} \
                          --attribute-names QueueArn --output text | sed -e 's/ATTRIBUTES\s//' ) \
                         && echo ${SQS_DEADLETTERQUEUE_ARN}

awslocal sqs create-queue --queue-name sample_stdQ

SQS_QUEUE_URL=$(awslocal sqs get-queue-url --queue-name sample_stdQ --output text) && echo ${SQS_QUEUE_URL}

cat << EOF >> deadletter.json
{
	"RedrivePolicy": "{\"deadLetterTargetArn\":\"${SQS_DEADLETTERQUEUE_ARN}\",\"maxReceiveCount\":\"5\"}"
}
EOF

awslocal sqs set-queue-attributes --queue-url ${SQS_QUEUE_URL} --attributes file://deadletter.json

rm -deadletter.json