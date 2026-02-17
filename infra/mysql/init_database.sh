#!/bin/bash

echo "⏳ Aguardando MySQL subir..."
sleep 30

echo "🚀 Criando banco e tabelas..."

docker exec -i orders-mysql mysql -uroot -proot < create_tables.sql

echo "✅ Banco e tabelas criados com sucesso!"
