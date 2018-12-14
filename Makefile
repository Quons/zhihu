.PHONY: build clean tool lint help

all: build
#每行的分隔符为一个tab键，不能用多个空格
build:
    # @关闭回声
	swag init
	@go build -v .

dev:
    # @关闭回声
	swag init
	@go build -v .
	./go-gin-example

run:
    # @关闭回声
	@go build -v .
	sudo systemctl restart gin

tool:
	go vet ./...; true
	gofmt -w .

lint:
	golint ./...

clean:
	rm -rf go-gin-example
	go clean -i .

help:
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"
	@echo "make lint: golint ./..."
	@echo "make clean: remove object files and cached files"
