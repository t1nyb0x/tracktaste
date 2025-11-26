# TrackTaste

Spotify と KKBOX を連携した音楽トラック情報取得・類似曲検索 API サーバーです。

## 機能

- **トラック情報取得**: Spotify URL からトラックの詳細情報を取得
- **トラック検索**: キーワードで Spotify のトラックを検索
- **類似トラック検索**: Spotify URL を元に KKBOX のレコメンド機能を活用した類似曲を取得
- **アーティスト情報取得**: Spotify URL からアーティストの詳細情報を取得
- **アルバム情報取得**: Spotify URL からアルバムの詳細情報を取得

## 技術スタック

- **言語**: Go 1.24
- **フレームワーク**: [go-chi/chi](https://github.com/go-chi/chi) v5
- **キャッシュ**: Redis（トークンキャッシュ用）
- **外部 API**: Spotify Web API, KKBOX Open API
- **アーキテクチャ**: Clean Architecture

## 必要要件

- Go 1.24 以上
- Redis（オプション、トークンキャッシュ用）
- Spotify Developer アカウント
- KKBOX Developer アカウント

## セットアップ

### 1. リポジトリのクローン

```bash
git clone https://github.com/t1nyb0x/tracktaste.git
cd tracktaste
```

### 2. 環境変数の設定

プロジェクトルートに `.env` ファイルを作成:

```env
# Server
HTTP_ADDR=:8080

# Spotify API
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret

# KKBOX API
KKBOX_ID=your_kkbox_client_id
KKBOX_SECRET=your_kkbox_client_secret

# Redis (optional)
REDIS_ADDR=localhost:6379
```

### 3. 依存関係のインストール

```bash
go mod download
```

### 4. サーバーの起動

```bash
go run ./cmd/server/...
```

開発時はホットリロードに [air](https://github.com/cosmtrek/air) を使用できます:

```bash
cd cmd/server
air
```

## API エンドポイント

### ヘルスチェック

```
GET /healthz
```

### トラック

| Method | Endpoint            | パラメータ | 説明                               |
| ------ | ------------------- | ---------- | ---------------------------------- |
| GET    | `/v1/track/fetch`   | `url`      | Spotify URL からトラック情報を取得 |
| GET    | `/v1/track/search`  | `q`        | キーワードでトラックを検索         |
| GET    | `/v1/track/similar` | `url`      | 類似トラックを取得                 |

### アーティスト

| Method | Endpoint           | パラメータ | 説明                                   |
| ------ | ------------------ | ---------- | -------------------------------------- |
| GET    | `/v1/artist/fetch` | `url`      | Spotify URL からアーティスト情報を取得 |

### アルバム

| Method | Endpoint          | パラメータ | 説明                               |
| ------ | ----------------- | ---------- | ---------------------------------- |
| GET    | `/v1/album/fetch` | `url`      | Spotify URL からアルバム情報を取得 |

## 使用例

### トラック情報の取得

```bash
curl "http://localhost:8080/v1/track/fetch?url=https://open.spotify.com/track/4uLU6hMCjMI75M1A2tKUQC"
```

### トラックの検索

```bash
curl "http://localhost:8080/v1/track/search?q=米津玄師%20Lemon"
```

### 類似トラックの取得

```bash
curl "http://localhost:8080/v1/track/similar?url=https://open.spotify.com/track/4uLU6hMCjMI75M1A2tKUQC"
```

## プロジェクト構成

```
tracktaste/
├── cmd/server/          # エントリーポイント
├── internal/
│   ├── domain/          # ドメインモデル
│   ├── port/            # インターフェース定義
│   ├── usecase/         # ビジネスロジック
│   ├── adapter/         # 外部接続
│   │   ├── gateway/     # 外部API実装
│   │   ├── handler/     # HTTPハンドラー
│   │   └── server/      # サーバー設定
│   ├── config/          # 設定
│   └── util/            # ユーティリティ
├── docs/                # ドキュメント
└── spec/                # API仕様書
```

詳細なアーキテクチャについては [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) を参照してください。

## ライセンス

[MIT License](LICENSE)

## Author

shika ([@t1nyb0x](https://github.com/t1nyb0x))
