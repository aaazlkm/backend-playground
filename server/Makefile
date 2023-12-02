# ---------------------------------------------------------------------------- #
#                                      タスク                                   s#
# ---------------------------------------------------------------------------- #

# フォーマット
fmt:
	go fmt ./...
	go mod tidy
	golangci-lint run --fix
.PHONY: fmt

# リントタスク
lint: fmt
	golangci-lint run
.PHONY: lint

# テストタスク
test:
	go test -race -timeout 30m $(option) ./...
.PHONY: test

# テストカバレッジを取得
test-coverage:
	go tool cover -html=cover.out -o cover.html
	echo "create cover.html"
.PHONY: test-coverage

# ビルドタスク
build: fmt lint test
	go build main.go
.PHONY: build

# 実行
run: fmt lint test
	go run main.go
.PHONY: build

# キャッシュ削除
clean:
	@go clean
.PHONY: clean

docker-build:
	docker-compose -f build/docker-compose.yaml up --build --force-recreate --no-deps
.PHONY: docker-build
