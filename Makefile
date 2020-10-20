PKG := "gitlab.com/otqee/otqee-be"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

default: build

build: clean
	@go build -o bin/main cmd/main/main.go
clean:
	rm -rf bin/main
run: ## run server for local test
	@go run cmd/main/main.go
lint:
	@golint -set_exit_status ${PKG_LIST}
fmt:
	@go fmt ${PKG_LIST}
tests:
	@go test -short ${PKG_LIST}
covertests:
	@go test -short ${PKG_LIST} -coverprofile=coverage.cov
covershow:
	@go tool cover -func=coverage.cov
cover: covertests covershow