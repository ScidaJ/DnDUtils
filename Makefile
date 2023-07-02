include .env

build-bot:
	go build -o bin/bot.exe bot/main.go

run-bot:
	go run bot/main.go -t $(BOT_TOKEN) -g $(GUILD_ID) -e $(API_ENDPOINT)