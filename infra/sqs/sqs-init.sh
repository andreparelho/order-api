#!/bin/bash

echo "Criando Fila SQS"
aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name orders-queue
