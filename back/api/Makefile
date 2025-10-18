.PHONY: start generate migrate migrate-host test test-coverage lint fmt vet build clean install-tools ci

start:
	air -c .air.toml

# SQLBoiler モデル生成
generate:
	sqlboiler mysql

# マイグレーション実行 (DBコンテナから実行)
migrate:
	@echo "マイグレーションを実行します..."
	docker compose exec db mysql -u collage -pcollagepass collage -e "$(shell cat migrations/001_create_users_table.sql)"
	@echo "マイグレーション完了"

# マイグレーション実行 (ホストから実行)
migrate-host:
	@echo "マイグレーションを実行します..."
	@for file in migrations/*.sql; do \
		echo "実行中: $$file"; \
		docker compose exec -T db mysql -u collage -pcollagepass collage < $$file; \
	done
	@echo "マイグレーション完了"

# テスト実行
test:
	@echo "Running tests..."
	GOWORK=off go test -v -race ./...

# テストカバレッジ
test-coverage:
	@echo "Running tests with coverage..."
	GOWORK=off go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	GOWORK=off go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# リント実行
lint:
	@echo "Running linter..."
	golangci-lint run --timeout=5m

# コードフォーマット
fmt:
	@echo "Formatting code..."
	GOWORK=off go fmt ./...
	gofmt -s -w .

# go vet実行
vet:
	@echo "Running go vet..."
	GOWORK=off go vet ./...

# ビルド
build:
	@echo "Building application..."
	GOWORK=off go build -o bin/api ./cmd/main.go

# クリーンアップ
clean:
	@echo "Cleaning up..."
	rm -rf bin/ coverage.out coverage.html tmp/

# 開発ツールのインストール
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/volatiletech/sqlboiler/v4@latest
	go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest

# CI用のコマンド（lint + test + build）
ci: fmt vet lint test build
	@echo "CI checks passed!"
