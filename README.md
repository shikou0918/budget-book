# Budget Book

家計簿管理アプリケーションです。収入・支出の記録、カテゴリ管理、予算設定、月次サマリーの機能を提供します。

## 技術スタック

### Backend
- **Go** 1.21
- **Echo v4** - Web フレームワーク
- **GORM** - ORM
- **MySQL** 8.0 - データベース
- **Validator v9** - バリデーション

### Frontend
- **Vue 3** - フロントエンドフレームワーク
- **TypeScript** - 型安全性
- **Vite** - ビルドツール
- **Pinia** - 状態管理
- **Vue Router** - ルーティング
- **Chart.js** - グラフ表示
- **Axios** - HTTP クライアント

### インフラ
- **Docker & Docker Compose** - コンテナ化
- **GitHub Actions** - CI/CD

## プロジェクト構成

```
.
├── backend/                 # Go API サーバー
│   ├── cmd/api/            # アプリケーションエントリーポイント
│   ├── entity/             # エンティティ定義
│   ├── usecase/            # ビジネスロジック
│   ├── infrastructure/     # データベース・リポジトリ
│   ├── interface/          # ハンドラー・ミドルウェア
│   ├── config/             # 設定管理
│   └── migrations/         # データベースマイグレーション
├── frontend/               # Vue.js フロントエンド
│   ├── src/
│   │   ├── components/     # Vue コンポーネント
│   │   ├── views/          # ページコンポーネント
│   │   ├── stores/         # Pinia ストア
│   │   └── router/         # ルーティング設定
│   └── public/             # 静的ファイル
└── docker-compose.yml      # Docker 設定
```

## セットアップ

### 前提条件
- Docker & Docker Compose
- Go 1.21+ (ローカル開発時)
- Node.js 18+ & Yarn (ローカル開発時)

### Docker での起動

```bash
# 全サービス起動
docker-compose up -d

# ログ確認
docker-compose logs -f
```

サービスURL:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- MySQL: localhost:3306

### ローカル開発

```bash
# 依存関係インストール
make install

# データベース起動
docker-compose up -d mysql

# バックエンド起動
cd backend
go run cmd/api/main.go

# フロントエンド起動 (別ターミナル)
cd frontend
yarn dev
```

## 開発コマンド

```bash
# ヘルプ表示
make help

# 依存関係インストール
make install

# リント実行
make lint

# リント自動修正
make lint-fix

# コードフォーマット
make format

# フォーマットチェック
make format-check

# テスト実行
make test

# ビルド成果物削除
make clean
```

## API エンドポイント

### 取引 (Transactions)
- `GET /api/transactions` - 取引一覧取得
- `POST /api/transactions` - 取引作成
- `GET /api/transactions/:id` - 取引詳細取得
- `PUT /api/transactions/:id` - 取引更新
- `DELETE /api/transactions/:id` - 取引削除

### カテゴリ (Categories)
- `GET /api/categories` - カテゴリ一覧取得
- `POST /api/categories` - カテゴリ作成
- `GET /api/categories/:id` - カテゴリ詳細取得
- `PUT /api/categories/:id` - カテゴリ更新
- `DELETE /api/categories/:id` - カテゴリ削除

### 予算 (Budgets)
- `GET /api/budgets` - 予算一覧取得
- `POST /api/budgets` - 予算作成
- `GET /api/budgets/:id` - 予算詳細取得
- `PUT /api/budgets/:id` - 予算更新
- `DELETE /api/budgets/:id` - 予算削除

### サマリー (Summary)
- `GET /api/summary/:year/:month` - 月次サマリー取得

## データベース

### マイグレーション

```bash
# マイグレーション実行
cd backend
atlas migrate apply --url "mysql://root:password@localhost:3306/budget_book"
```

## 機能

- ✅ 収入・支出の記録管理
- ✅ カテゴリ管理
- ✅ 予算設定・管理
- ✅ 月次サマリー・統計表示
- ✅ レスポンシブデザイン
- ✅ データのバリデーション
- ✅ REST API

## ライセンス

このプロジェクトはプライベートプロジェクトです。