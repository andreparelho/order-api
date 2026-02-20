#!/bin/bash

set -e  # para o script se algum comando falhar

echo "🔧 Configurando AWS CLI para LocalStack..."

PROFILE=localstack
ENDPOINT=http://localhost:4566
REGION=us-east-1

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

echo "✅ AWS CLI configurada com profile '$PROFILE'"

# =========================
# Criar tabela DynamoDB
# =========================
echo "📦 Criando tabela payments-orders..."

aws dynamodb create-table \
  --table-name payments-orders \
  --attribute-definitions \
    AttributeName=order_id,AttributeType=S \
    AttributeName=payment_id,AttributeType=S \
  --key-schema \
    AttributeName=order_id,KeyType=HASH \
    AttributeName=payment_id,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST \
  --endpoint-url=$ENDPOINT \
  --region $REGION \
  --profile $PROFILE

echo "✅ Tabela criada com sucesso!"

PROFILE=localstack
ENDPOINT=http://localhost:4566
REGION=us-east-1

aws dynamodb scan \
  --table-name payments-orders \
  --endpoint-url=http://localhost:4566 \
  --region us-east-1 \
  --profile localstack