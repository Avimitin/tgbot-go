FROM golang:1.15-alpine3.13
WORKDIR ${GOPATH}/src/github.com/avimitin/go-bot
COPY . ${GOPATH}/src/github.com/avimitin/go-bot
RUN go build -o bin/go-bot -ldflags '-s -w' ./cmd/go-bot-v2
ENV GOBOT_CONFIG_PATH=/data
ENTRYPOINT ["bin/go-bot"]
