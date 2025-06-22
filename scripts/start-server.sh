#!/bin/bash

# Screen sessiya yaratish va serverni ishga tushirish
screen -dmS blog-server bash -c 'cd /path/to/your/project && go run *.go'

echo "Server screen sessiyasida ishga tushdi"
echo "Ko'rish uchun: screen -r blog-server"
echo "Chiqish uchun: Ctrl+A, D"
