.PHONY: help build clean deps

# 变量
MODULE := freecharge/go-freecharge
GO := go
GOFMT := gofmt
GOVET := go vet
BINARY := bin/freecharge-server
BUILD_DIR := $(dir $(BINARY))

# 默认目标
.DEFAULT_GOAL := build

deps: ## 下载并验证依赖
	$(GO) mod download
	$(GO) mod verify
	$(GO) mod tidy

build: ## 构建可执行文件
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BINARY) ./cmd/server

clean: ## 清理构建产物
	$(GO) clean
	rm -f $(BINARY)

