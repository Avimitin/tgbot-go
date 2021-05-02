GO=$(shell which go)

TAG=$(shell cat .tag)

GOPATH=$(shell $(GO) env GOPATH)
GOARCH=$(shell $(GO) env GOARCH)
GOOS=$(shell $(GO) env GOOS)

PROJECT_PATH=$(GOPATH)/src/github.com/Avimitin/go-bot

export GOBOT_CONFIG_PATH=$(PROJECT_PATH)/fixtures/config.toml
export GO111MODULE=on

GO_TEST_ARGS='-v'

.PHONY: build
build:
	$(GO) build -o ./bin/go-bot-$(TAG)-$(GOARCH)-$(GOOS) -ldflags '-s -w' ./cmd/go-bot

.PHONY: test
test:
	$(GO) test $(GO_TEST_ARGS) $(PROJECT_PATH)/...
