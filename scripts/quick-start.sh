#!/bin/bash

echo "Blog serverni doimiy ishga tushirish..."

# 1. Binary yaratish
go build -o blog-server

# 2. Screen sessiyasida ishga tushirish
screen -dmS blog-server bash -c './blog-server'

echo "✅ Server ishga tushdi!"
echo "📱 Sayt: http://localhost:8080"
echo "🔍 Status ko'rish: screen -r blog-server"
echo "🛑 To'xtatish: screen -r blog-server, keyin Ctrl+C"
echo "📋 Barcha screen sessiyalar: screen -ls"
