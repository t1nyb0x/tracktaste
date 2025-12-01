# TrackTaste

[![CI](https://github.com/t1nyb0x/tracktaste/actions/workflows/ci.yml/badge.svg)](https://github.com/t1nyb0x/tracktaste/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/t1nyb0x/tracktaste/branch/main/graph/badge.svg)](https://codecov.io/gh/t1nyb0x/tracktaste)

Spotify と KKBOX を連携した音楽トラック情報取得・類似曲検索 API サーバーです。

## 機能

- **トラック情報取得**: Spotify URL からトラックの詳細情報を取得
- **トラック検索**: キーワードで Spotify のトラックを検索
- **類似トラック検索**: Spotify URL を元に KKBOX のレコメンド機能を活用した類似曲を取得
- **レコメンド V2**: マルチソース候補収集 + Deezer/MusicBrainz 特徴量による高精度レコメンド
  - **候補ソース**: KKBOX, Last.fm, MusicBrainz (アーティスト曲), YouTube Music
  - **特徴量**: Deezer (BPM/Duration/Gain) + MusicBrainz (Tags/Relations)
- **アーティスト情報取得**: Spotify URL からアーティストの詳細情報を取得
- **アルバム情報取得**: Spotify URL からアルバムの詳細情報を取得

## 技術スタック

- **言語**: Go 1.24
- **フレームワーク**: [go-chi/chi](https://github.com/go-chi/chi) v5
- **キャッシュ**: 2 層キャッシュ（L1: インメモリ, L2: Redis）
- **外部 API**: Spotify, KKBOX, Deezer, MusicBrainz, Last.fm, YouTube Music (sidecar)
- **アーキテクチャ**: Clean Architecture

## 必要要件

- Go 1.24 以上
- Docker (YouTube Music sidecar 用)
- Redis（オプション、L2 キャッシュ用。なくてもインメモリキャッシュで動作）
- Spotify Developer アカウント
- KKBOX Developer アカウント
- Last.fm API Key（オプション、https://www.last.fm/api/account/create で無料取得）

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

# Last.fm API (optional - for multi-source candidates)
LASTFM_API_KEY=your_lastfm_api_key

# YouTube Music Sidecar (optional - for multi-source candidates)
YTMUSIC_SIDECAR_URL=http://localhost:8081

# Redis (optional - L2 cache)
REDIS_URL=localhost:6379
REDIS_PASSWORD=
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

## Docker

### Docker Compose で起動（推奨）

```bash
# 開発環境
docker compose up -d

# 本番環境
docker compose -f docker-compose.prod.yml up -d
```

### Docker イメージの取得

```bash
# 最新版
docker pull ghcr.io/t1nyb0x/tracktaste:latest

# 特定バージョン
docker pull ghcr.io/t1nyb0x/tracktaste:v1.0.0
```

### 手動でビルド

```bash
docker build -t tracktaste .
docker run -p 8080:8080 --env-file .env tracktaste
```

## テスト

```bash
# 全テスト実行
go test ./...

# カバレッジ付き
go test ./... -cover

# 詳細出力
go test ./... -v

# 特定パッケージのみ
go test ./internal/usecase/v1/...  # V1ユースケース
go test ./internal/usecase/v2/...  # V2ユースケース
```

### Lint

```bash
# golangci-lint をインストール
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# lint 実行
golangci-lint run
```

## API エンドポイント

### ヘルスチェック

```
GET /healthz
```

### トラック

| Method | Endpoint              | パラメータ             | 説明                                        |
| ------ | --------------------- | ---------------------- | ------------------------------------------- |
| GET    | `/v1/track/fetch`     | `url`                  | Spotify URL からトラック情報を取得          |
| GET    | `/v1/track/search`    | `q`                    | キーワードでトラックを検索                  |
| GET    | `/v1/track/similar`   | `url`                  | 類似トラックを取得（KKBOX レコメンド）      |
| GET    | `/v2/track/recommend` | `url`, `mode`, `limit` | Deezer + MusicBrainz ベースのレコメンド取得 |

#### `/v2/track/recommend` パラメータ詳細

| パラメータ | 必須 | デフォルト | 説明                                                |
| ---------- | ---- | ---------- | --------------------------------------------------- |
| `url`      | ○    | -          | Spotify トラック URL                                |
| `mode`     | -    | `balanced` | レコメンドモード (`similar`, `related`, `balanced`) |
| `limit`    | -    | `20`       | 返却件数（1〜30）                                   |

##### レコメンドモード (`mode`)

| モード     | 説明                                 | 用途                           |
| ---------- | ------------------------------------ | ------------------------------ |
| `similar`  | BPM/Duration/Gain の類似度を最優先   | テンポや雰囲気の似た曲を探す時 |
| `related`  | タグ・ジャンルの関連性を優先         | 同じジャンルの新しい曲を探す時 |
| `balanced` | 両方をバランスよく考慮（デフォルト） | 一般的なレコメンド             |

##### レコメンドエンジン V2 について

従来の Spotify Audio Features (廃止済み) に代わり、以下のデータソースを使用:

**候補収集（並列実行）**
| ソース | 内容 | 候補数 |
|--------|------|--------|
| KKBOX | レコメンドトラック | 30 件 |
| Last.fm | track.getSimilar | 30 件 |
| MusicBrainz | 同一アーティストの他の曲 | 20 件 |
| YouTube Music | ラジオ/類似曲 (sidecar) | 25 件 |

**特徴量取得**

- **Deezer API**: BPM、Duration（秒）、Gain（ReplayGain dB）
- **MusicBrainz API**: タグ（ジャンル、ムード等）、アーティスト関連情報

類似度計算には Jaccard 係数（タグ類似度）と各特徴量の正規化距離を組み合わせ、ジャンルボーナス/ペナルティを適用しています。

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

### レコメンドトラックの取得

```bash
# デフォルト（balanced モード、20件）
curl "http://localhost:8080/v2/track/recommend?url=https://open.spotify.com/track/4uLU6hMCjMI75M1A2tKUQC"

# similar モード、10件
curl "http://localhost:8080/v2/track/recommend?url=https://open.spotify.com/track/4uLU6hMCjMI75M1A2tKUQC&mode=similar&limit=10"
```

#### レコメンドレスポンス例

```json
{
  "status": 200,
  "result": {
    "seed_track": {
      "id": "xxx",
      "name": "入力楽曲名",
      "artists": [{"id": "...", "name": "アーティスト名", "url": "..."}],
      "audio_features": {
        "bpm": 128.0,
        "duration_seconds": 240,
        "gain": -5.2,
        "tags": ["j-pop", "anime", "electronic"]
      },
      "genres": ["anime", "j-pop"]
    },
    "items": [
      {
        "id": "yyy",
        "name": "レコメンド楽曲名",
        "artists": [{"id": "...", "name": "アーティスト名", "url": "..."}],
        "album": {...},
        "url": "https://open.spotify.com/track/yyy",
        "similarity_score": 0.92,
        "genre_bonus": 1.5,
        "final_score": 1.38,
        "match_reasons": ["bpm", "duration", "same_tags"],
        "audio_features": {
          "bpm": 126.0,
          "duration_seconds": 235,
          "gain": -4.8,
          "tags": ["j-pop", "pop", "anime"]
        }
      }
    ],
    "mode": "balanced"
  }
}
```

## プロジェクト構成

```
tracktaste/
├── cmd/server/          # エントリーポイント
├── sidecar/
│   └── ytmusic/         # YouTube Music Python sidecar (ytmusicapi)
│       ├── main.py      # FastAPIサーバー
│       ├── Dockerfile
│       └── requirements.txt
├── internal/
│   ├── domain/          # ドメインモデル
│   ├── port/            # インターフェース定義
│   ├── usecase/         # ビジネスロジック
│   │   ├── recommend_v2.go    # レコメンドロジック
│   │   ├── similarity.go      # 類似度計算
│   │   └── genre_matcher.go   # ジャンルマッチング
│   ├── adapter/         # 外部接続
│   │   ├── gateway/     # 外部API実装
│   │   │   ├── cache/       # 2層キャッシュ（L1:メモリ, L2:Redis）
│   │   │   ├── redis/       # Redisクライアント
│   │   │   ├── spotify/     # Spotify API
│   │   │   ├── kkbox/       # KKBOX API
│   │   │   ├── deezer/      # Deezer API（BPM/Gain取得）
│   │   │   ├── musicbrainz/ # MusicBrainz API（タグ/関連情報）
│   │   │   ├── lastfm/      # Last.fm API（類似曲取得）
│   │   │   └── ytmusic/     # YouTube Music sidecarクライアント
│   │   ├── handler/     # HTTPハンドラー
│   │   └── server/      # サーバー設定
│   ├── config/          # 設定
│   └── util/            # ユーティリティ
├── docs/                # ドキュメント
└── spec/                # API仕様書
```

詳細なアーキテクチャについては [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) を参照してください。

## CI/CD

GitHub Actions によるワークフロー:

| ワークフロー                | トリガー     | 内容                                 |
| --------------------------- | ------------ | ------------------------------------ |
| **CI** (`ci.yml`)           | push/PR      | Lint, Test, Build                    |
| **Publish** (`publish.yml`) | main push    | Docker イメージを ghcr.io に publish |
| **Release** (`release.yml`) | Release 作成 | バージョンタグ付きイメージを publish |

### Docker イメージタグ

- `latest` - 最新の main ブランチ
- `v1.2.3` - リリースバージョン
- `v1.2` - マイナーバージョン
- `<sha>` - コミットハッシュ

## ライセンス

[MIT License](LICENSE)

## Author

shika ([@t1nyb0x](https://github.com/t1nyb0x))
