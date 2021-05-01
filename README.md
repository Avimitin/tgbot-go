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

