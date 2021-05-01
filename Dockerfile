# BUILD
FROM golang:1.15-alpine3.13 AS build
WORKDIR ${GOPATH}/src/github.com/avimitin/go-bot
COPY . ${GOPATH}/src/github.com/avimitin/go-bot
RUN go build -o /bin/go-bot -ldflags '-s -w' ./cmd/go-bot

# RUN
FROM alpine:3
COPY --from=build /bin/go-bot /bin/go-bot
ENV GOBOT_CONFIG_PATH=/data
ENTRYPOINT ["bin/go-bot"]
