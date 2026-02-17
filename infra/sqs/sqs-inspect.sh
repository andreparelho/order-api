#!/bin/bash

echo "🔧 Configurando AWS CLI para LocalStack..."

# Credenciais fake (LocalStack aceita qualquer valor)
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
export AWS_REGION=us-east-1
export SQS_ENDPOINT=http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/orders-queue   # LocalStack (trocar em prod)
export QUEUE_NAME=orders-queue 

# Profile dedicado
aws configure set aws_access_key_id test --profile localstack
aws configure set aws_secret_access_key test --profile localstack
aws configure set region us-east-1 --profile localstack
aws configure set output json --profile localstack

QUEUE_URL=$(aws --region $AWS_REGION \
  --endpoint-url=$SQS_ENDPOINT \
  sqs get-queue-url \
  --queue-name $QUEUE_NAME \
  --query "QueueUrl" \
  --output text)
  
aws --region $AWS_REGION \
  --endpoint-url=$SQS_ENDPOINT \
  sqs receive-message \
  --queue-url $QUEUE_URL \
  --max-number-of-messages 10 \
  --wait-time-seconds 1 \
  --visibility-timeout 0
