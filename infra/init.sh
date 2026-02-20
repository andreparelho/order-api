#!/bin/bash

echo "Inicializando Docker Compose"
cd .. && docker-compose down -v && docker-compose up -d 

echo "Iniciando Fila SQS"
cd ./infra && cd ./sqs && sh sqs-init.sh

echo "Iniciando Database Dynamo"
cd .. && cd ./dynamodb && sh dynamo.sh

echo "Iniciando Base de Dados"
cd .. && cd ./mysql && sh init_database.sh
