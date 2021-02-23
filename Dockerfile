FROM golang:1.15-alpine3.13
WORKDIR ${GOPATH}/src/github.com/avimitin/go-bot
COPY . ${GOPATH}/src/github.com/avimitin/go-bot
RUN go build -o bin/go-bot -ldflags '-s -w' cmd/go-bot/main.go
ENV BOTCFGPATH=/data
ENTRYPOINT ["bin/go-bot"]
