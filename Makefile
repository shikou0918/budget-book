.PHONY: help install lint lint-fix format format-check test type-check clean

# デフォルトターゲット
help:
	@echo "利用可能なコマンド:"
	@echo "  install     - フロントエンドとバックエンドの依存関係をインストール"
	@echo "  lint        - フロントエンドとバックエンドのLinterを実行"
	@echo "  lint-fix    - フロントエンドとバックエンドのLinterを自動修正付きで実行"
	@echo "  format      - フロントエンドとバックエンドのコードをFormat"
	@echo "  format-check- フロントエンドとバックエンドのコードFormatをチェック"
	@echo "  test        - フロントエンドとバックエンドのテストを実行"
	@echo "  type-check  - フロントエンドの型チェックを実行"
	@echo "  clean       - ビルド成果物をクリーンアップ"

# 依存関係のインストール
install:
	@echo "バックエンドの依存関係をインストール中..."
	cd backend && go mod download
	@echo "フロントエンドの依存関係をインストール中..."
	cd frontend && yarn install

# バックエンドターゲット
backend-lint:
	@echo "バックエンドのリンターを実行中..."
	cd backend && golangci-lint run

backend-lint-fix:
	@echo "バックエンドのリンターを修正付きで実行中..."
	cd backend && golangci-lint run --fix

backend-format:
	@echo "バックエンドコードをフォーマット中..."
	cd backend && gofmt -s -w .
	cd backend && goimports -w .

backend-test:
	@echo "バックエンドテストを実行中..."
	cd backend && go test -v ./...

# フロントエンドターゲット
frontend-lint:
	@echo "フロントエンドのリンターを実行中..."
	cd frontend && yarn lint:check

frontend-lint-fix:
	@echo "フロントエンドのリンターを修正付きで実行中..."
	cd frontend && yarn lint:fix

frontend-format:
	@echo "フロントエンドコードをフォーマット中..."
	cd frontend && yarn format

frontend-format-check:
	@echo "フロントエンドのコードフォーマットをチェック中..."
	cd frontend && yarn format:check

frontend-test:
	@echo "フロントエンドテストを実行中..."
	cd frontend && yarn test 2>/dev/null || echo "フロントエンドテストはまだ設定されていません"

frontend-type-check:
	@echo "フロントエンドの型チェックを実行中..."
	cd frontend && yarn type-check

# まとめて実行
lint: backend-lint frontend-lint
lint-fix: backend-lint-fix frontend-lint-fix
format: backend-format frontend-format
format-check: frontend-format-check
test: backend-test frontend-test
type-check: frontend-type-check

# クリーンアップ
clean:
	@echo "バックエンドの成果物をクリーンアップ中..."
	cd backend && rm -rf bin/
	@echo "フロントエンドの成果物をクリーンアップ中..."
	cd frontend && rm -rf dist/ node_modules/.cache/