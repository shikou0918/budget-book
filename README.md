# Budget Book

家計簿管理アプリケーションです。収入・支出の記録、カテゴリ管理、予算設定、月次サマリーの機能を提供します。

## 技術スタック

### Backend
- **Go** 1.25
- **Echo v4** - Web フレームワーク
- **GORM** - ORM
- **MySQL** 8.0 - データベース
- **SQLite** - テスト用データベース
- **Validator v9** - バリデーション
- **Testify** - テストライブラリ
- **mockgen** - モック生成ツール
- **Air** - ホットリロード（開発環境）

### Frontend
- **Vue 3** - フロントエンドフレームワーク
- **TypeScript** - 型安全性
- **Vite** - ビルドツール
- **Pinia** - 状態管理
- **Vue Router** - ルーティング
- **Vuetify** - UI コンポーネントライブラリ
- **Chart.js** - グラフ表示
- **Axios** - HTTP クライアント
- **Vitest** - テストフレームワーク
- **Vue Test Utils** - Vue コンポーネントテスト
- **jsdom** - ブラウザ環境シミュレーション

### インフラ
- **Docker & Docker Compose** - コンテナ化
- **GitHub Actions** - CI/CD

## プロジェクト構成

```
.
├── backend/                 # Go API サーバー (Clean Architecture)
│   ├── cmd/api/            # アプリケーションエントリーポイント・DI設定
│   ├── entity/             # エンティティ定義・ドメインロジック
│   │   └── *_test.go       # エンティティテスト
│   ├── usecase/            # ビジネスロジック + リポジトリインターフェース定義
│   │   └── *_test.go       # ユースケーステスト（モック使用）
│   ├── infrastructure/     # リポジトリ実装（usecaseインターフェースに依存）
│   │   ├── database/       # データベース接続
│   │   └── repository/     # リポジトリ実装
│   ├── interface/          # HTTP ハンドラー・ミドルウェア
│   │   ├── handler/        # HTTP ハンドラー
│   │   │   └── *_test.go   # ハンドラーテスト・統合テスト
│   │   └── middleware/     # ミドルウェア
│   ├── mocks/              # テスト用モック（mockgen 生成）
│   │   ├── repository/     # リポジトリモック
│   │   └── usecase/        # ユースケースモック
│   ├── config/             # 設定管理
│   └── migrations/         # データベースマイグレーション
├── frontend/               # Vue.js フロントエンド
│   ├── src/
│   │   ├── components/     # Vue コンポーネント
│   │   │   └── **/__tests__/  # コンポーネントテスト
│   │   ├── views/          # ページコンポーネント
│   │   ├── stores/         # Pinia ストア
│   │   │   └── __tests__/     # ストアテスト
│   │   ├── utils/          # ユーティリティ関数
│   │   │   └── __tests__/     # ユーティリティテスト
│   │   ├── router/         # ルーティング設定
│   │   └── types/          # TypeScript 型定義
│   ├── public/             # 静的ファイル
│   ├── vitest.config.ts    # Vitest 設定
│   └── src/test-setup.ts   # テストセットアップ
└── docker-compose.yml      # Docker 設定
```

## アーキテクチャ

### Clean Architecture + 依存性逆転の原則

バックエンドは厳格なレイヤー分離とDI（依存性注入）を採用:

```
entity/              # ドメインエンティティ + ビジネスルール
  ↑
usecase/            # ビジネスロジック + リポジトリインターフェース定義
  ↑
infrastructure/     # リポジトリ実装（usecaseインターフェースに依存）
  repository/
  ↑
interface/          # HTTPハンドラー（usecaseに依存）
  handler/
```

**重要パターン**: リポジトリインターフェースは`usecase/`パッケージで定義され（例: `usecase/transaction.go`の`TransactionRepositoryInterface`）、実装は`infrastructure/repository/`に配置されます。これにより依存性逆転原則を実現しています。

**DI フロー** (`cmd/api/main.go:25-87`参照):
1. DB接続でリポジトリ実装を作成
2. リポジトリをユースケースに注入（ユースケースはインターフェースのみ知っている）
3. ユースケースをハンドラーに注入
4. ハンドラーをEchoルートに登録

## セットアップ

### 前提条件
- Docker & Docker Compose
- Go 1.25+ (ローカル開発時)
- Node.js 20+ & Yarn (ローカル開発時)
- mockgen (モック生成時)
- Air (Goホットリロード、Dockerコンテナ内で自動インストール)

### Docker での起動

```bash
# 全サービス起動
docker-compose up -d

# ログ確認
docker-compose logs -f

# バックエンドのみ再ビルド（コード変更時）
docker-compose up -d --build backend
```

サービスURL:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- MySQL: localhost:3306

**開発時の注意**: バックエンドは [Air](https://github.com/cosmtrek/air) によるホットリロードを使用しています。`backend/`ディレクトリ内のGoファイルを編集すると、自動的に再ビルド・再起動されます。

### ローカル開発

#### Dockerを使った開発（推奨）

```bash
# 全サービス起動（ホットリロード有効）
docker-compose up -d

# ログをリアルタイムで確認
docker-compose logs -f backend

# バックエンドのコードを編集すると、Airが自動で再ビルド・再起動
```

#### ローカル環境での開発

```bash
# 依存関係インストール
make install

# データベース起動
docker-compose up -d mysql

# バックエンド起動（Airでホットリロード）
cd backend
air

# または通常起動
go run cmd/api/main.go

# フロントエンド起動 (別ターミナル)
cd frontend
yarn dev
```

## 開発コマンド

### Make コマンド

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

# 型チェック実行
make type-check

# ビルド成果物削除
make clean
```

### バックエンドテスト

```bash
cd backend

# 全テスト実行
go test ./...

# 詳細表示でテスト実行
go test -v ./...

# カバレッジ付きでテスト実行
go test -cover ./...

# 特定パッケージのテスト実行
go test ./usecase/
go test ./interface/handler/
```

### モック生成（バックエンド）

リポジトリインターフェースを変更した後は、モックを再生成する必要があります:

```bash
cd backend

# Transaction リポジトリモック再生成
mockgen -source=usecase/transaction.go -destination=mocks/repository/transaction_mock.go -package=repository

# Category リポジトリモック再生成
mockgen -source=usecase/category.go -destination=mocks/repository/category_mock.go -package=repository

# Budget リポジトリモック再生成
mockgen -source=usecase/budget.go -destination=mocks/repository/budget_mock.go -package=repository
```

### フロントエンドテスト

```bash
cd frontend

# テスト実行（ウォッチモード）
yarn test

# テスト実行（一回のみ）
yarn test:run

# テスト UI でブラウザ表示
yarn test:ui

# カバレッジ付きでテスト実行
yarn test:coverage
```

### テスト構成

#### フロントエンド
- **ユニットテスト**: Vitest + Vue Test Utils
  - `src/components/**/__tests__/` - Vue コンポーネントテスト
  - `src/stores/__tests__/` - Pinia ストアテスト
  - `src/utils/__tests__/` - ユーティリティ関数テスト
- **モック**: vi.mock() でAPI呼び出しをモック
- **テスト環境**: jsdom でブラウザ環境をシミュレーション

#### バックエンド
- **ユニットテスト**: Go標準のtesting + testify
  - `entity/*_test.go` - エンティティ・ドメインロジックテスト
  - `usecase/*_test.go` - ビジネスロジックテスト（モック使用）
  - `interface/handler/*_test.go` - HTTPハンドラーテスト
- **統合テスト**: 実際のデータベース（SQLite）を使用したE2Eテスト
- **モック**: mockgen で生成されたタイプセーフなモック
- **アーキテクチャ**: クリーンアーキテクチャ + 依存性逆転原則に基づくテスト設計

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

### アプリケーション機能
- ✅ 収入・支出の記録管理
- ✅ カテゴリ管理
- ✅ 予算設定・管理
- ✅ 月次サマリー・統計表示
- ✅ レスポンシブデザイン
- ✅ データのバリデーション
- ✅ REST API

### 開発・品質保証
- ✅ TypeScript による型安全性
- ✅ ESLint + Prettier による コード品質管理
- ✅ Go による静的型付け・コンパイル時チェック
- ✅ Vitest による単体テスト（フロントエンド）
- ✅ Go testing + testify による単体テスト（バックエンド）
- ✅ Vue コンポーネントテスト
- ✅ Pinia ストアテスト
- ✅ mockgen によるタイプセーフなモック生成
- ✅ 統合テスト（HTTPハンドラー + データベース）
- ✅ クリーンアーキテクチャによる保守性の高い設計
- ✅ 依存性逆転原則に基づくテスタブルな構造
- ✅ CI/CD パイプライン (GitHub Actions)
- ✅ Docker によるコンテナ化