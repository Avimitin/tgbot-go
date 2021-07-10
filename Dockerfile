# BUILD
FROM golang:alpine AS build
WORKDIR ${GOPATH}/src/github.com/avimitin/go-bot
COPY . ${GOPATH}/src/github.com/avimitin/go-bot
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o /bin/go-bot -ldflags '-s -w' ./cmd/go-bot

# RUN
FROM alpine:3
COPY --from=build /bin/go-bot /bin/go-bot
ENV GOBOT_CONFIG_PATH=/data/config.toml
ENTRYPOINT ["bin/go-bot"]
