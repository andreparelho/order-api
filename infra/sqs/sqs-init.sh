#!/bin/bash

echo "🔧 Configurando AWS CLI para LocalStack..."

# Credenciais fake (LocalStack aceita qualquer valor)
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
export AWS_REGION=us-east-1

# Profile dedicado
aws configure set aws_access_key_id test --profile localstack
aws configure set aws_secret_access_key test --profile localstack
aws configure set region us-east-1 --profile localstack
aws configure set output json --profile localstack

echo "Criando Fila SQS"
aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name orders-queue
