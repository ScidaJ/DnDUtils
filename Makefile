include .env

build-bot:
	go build -o bin/bot.exe bot/main.go

run-bot:
	go run bot/main.go

run-api:
	go run api/main.go