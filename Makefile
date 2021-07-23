PKG := "gitlab.com/beewar/beewar-be"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

default: build

build: clean
	@go build -o cmd/bin/main cmd/main/main.go
clean: ## clean binaries
	rm -rf cmd/bin
run: ## run server for local test
	@go run cmd/main/main.go
runregress: ## run regression test
	@go run cmd/regression/main.go
runseeder: ## run seeder
	@go run cmd/seeder/main.go
lint:
	@golint -set_exit_status ${PKG_LIST}
fmt:
	@go fmt ${PKG_LIST}
tests:
	@go test -short ${PKG_LIST}
bench:
	@go test -bench=. ${PKG_LIST}
covertests:
	@go test -short ${PKG_LIST} -coverprofile=coverage.cov
covershow:
	@go tool cover -func=coverage.cov
cover: covertests covershow

check: lint fmt tests runregress ## check before push to remote