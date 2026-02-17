#!/bin/bash

echo "🔍 LISTANDO TODAS AS KEYS DO REDIS"
echo "--------------------------------"

redis-cli -h localhost -p 6379 -a 123 KEYS "*"
