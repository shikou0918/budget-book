# Budget Book API Documentation

このディレクトリには、Budget Book API のドキュメントが含まれています。

## 📖 API 仕様書

API 仕様書は Swagger UI を使用して GitHub Pages で公開されています。

- **Live Documentation**: [GitHub Pages で確認](https://[your-username].github.io/budget-book/)
- **OpenAPI 仕様ファイル**: [openapi.yml](../openapi.yml)

## 🚀 機能

- **インタラクティブ API 探索**: Swagger UI で API エンドポイントを直接テスト可能
- **日本語対応**: API 仕様と UI が日本語で記述
- **自動デプロイ**: `openapi.yml`の更新時に自動でドキュメントが更新

## 📋 API エンドポイント概要

### 取引 (Transactions)

- `GET /api/transactions` - 取引一覧取得
- `POST /api/transactions` - 取引作成
- `GET /api/transactions/{id}` - 取引詳細取得
- `PUT /api/transactions/{id}` - 取引更新
- `DELETE /api/transactions/{id}` - 取引削除

### カテゴリ (Categories)

- `GET /api/categories` - カテゴリ一覧取得
- `POST /api/categories` - カテゴリ作成
- `GET /api/categories/{id}` - カテゴリ詳細取得
- `PUT /api/categories/{id}` - カテゴリ更新
- `DELETE /api/categories/{id}` - カテゴリ削除

### 予算 (Budgets)

- `GET /api/budgets` - 予算一覧取得
- `POST /api/budgets` - 予算作成
- `GET /api/budgets/{id}` - 予算詳細取得
- `PUT /api/budgets/{id}` - 予算更新
- `DELETE /api/budgets/{id}` - 予算削除

### サマリー (Summary)

- `GET /api/summary/{year}/{month}` - 月次サマリー取得

## 🔧 開発者向け

### ローカルでの確認

1. Swagger UI をローカルで起動:

```bash
# Swagger UIのDockerイメージを使用
docker run -p 8080:8080 -e SWAGGER_JSON=/openapi.yml -v $(pwd)/openapi.yml:/openapi.yml swaggerapi/swagger-ui
```

2. ブラウザで http://localhost:8080 にアクセス

### API 仕様の更新

1. `openapi.yml`を編集
2. 変更をコミット・プッシュ
3. GitHub Actions が自動でドキュメントを更新

## 📝 注意事項

- GitHub Pages の設定でソースを「GitHub Actions」に設定する必要があります
- 初回デプロイ時は、リポジトリの設定で GitHub Pages を有効にしてください
