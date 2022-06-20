GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
# GOTEST=$(GOCMD) test
# GOGET=$(GOCMD) get
GO_MOD=$(GOCMD) mod
GO_ENV=$(GOCMD) env
BINARY_NAME=iobscan-ibc-explorer-backend
export GO111MODULE = on
export GOSUMDB=off
export GIT_TERMINAL_PROMPT=1
export GOPROXY=https://goproxy.cn,direct

all: get_vendor build

get_vendor:
	@rm -rf vendor/
	@echo "--> Downloading dependencies"
	$(GO_MOD) download
	$(GO_MOD) vendor

build:
	$(GOBUILD) -o $(BINARY_NAME) .

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)