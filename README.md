# go-bot

A telegram bot with basic commands.

## build

```bash
go build -o ./bin/go-bot -ldflags '-s -w' ./cmd/go-bot
```

## run

```bash
mkdir -p ~/.config/go-bot
cat ./fixtures/config.toml > ~/.config/go-bot/config.toml

./bin/go-bot
```

## docker

```bash
mkdir bot
curl -o ./bot/docker-compose.yml \
       -sSL "https://raw.githubusercontent.com/Avimitin/go-bot/master/docker-compose.yml.example"
cd ./bot; docker-compose up -d
```
