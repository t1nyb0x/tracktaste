# TrackTaste アーキテクチャ

## 概要

TrackTaste は Clean Architecture に基づいて設計されています。
依存関係は外側から内側に向かい、内側のレイヤーは外側のレイヤーを知りません。

## ディレクトリ構造

```
tracktaste/
├── cmd/
│   └── server/
│       └── main.go                 # エントリーポイント、DI（依存性注入）
│
├── sidecar/
│   └── ytmusic/                    # YouTube Music Python sidecar
│       ├── main.py                 # FastAPIサーバー (ytmusicapi)
│       ├── Dockerfile
│       └── requirements.txt
│
└── internal/
    ├── domain/                      # ドメイン層（最も内側）
    │   ├── track.go                # Track, SimpleTrack, SimilarTrack, TrackFeatures
    │   ├── artist.go               # Artist, SimpleArtist, ArtistInfo
    │   ├── album.go                # Album
    │   ├── image.go                # Image
    │   └── errors.go               # ドメインエラー定義
    │
    ├── port/                        # ポート層（インターフェース定義）
    │   ├── repository/
    │   │   └── token.go            # TokenRepository interface
    │   └── external/
    │       ├── spotify.go          # SpotifyAPI interface
    │       ├── kkbox.go            # KKBOXAPI interface
    │       ├── deezer.go           # DeezerAPI interface
    │       ├── musicbrainz.go      # MusicBrainzAPI interface
    │       ├── lastfm.go           # LastFMAPI interface
    │       └── ytmusic.go          # YouTubeMusicAPI interface
    │
    ├── usecase/                     # ユースケース層（ビジネスロジック）
    │   ├── genre_matcher.go        # GenreMatcher (V1/V2共通)
    │   │
    │   ├── v1/                      # V1 ユースケース
    │   │   ├── track.go            # TrackUseCase
    │   │   ├── artist.go           # ArtistUseCase
    │   │   ├── album.go            # AlbumUseCase
    │   │   ├── similar_tracks.go   # SimilarTracksUseCase
    │   │   ├── recommend.go        # RecommendUseCase (Spotify Audio Features)
    │   │   └── similarity.go       # SimilarityCalculator
    │   │
│   │   └── v2/                      # V2 ユースケース
│       ├── recommend.go        # RecommendUseCase (マルチソースレコメンド)
│       │   ├── searchSpotifyWithFallback()  # Spotify検索フォールバック
│       │   ├── sanitizeSearchQuery()        # クエリサニタイズ
│       │   ├── simplifyTrackName()          # 曲名簡素化
│       │   └── fuzzyMatchArtist()           # アーティスト曖昧マッチ
│       └── similarity.go       # SimilarityCalculatorV2
    │
    ├── adapter/                     # アダプター層（最も外側）
    │   ├── gateway/                # Secondary Adapters（外部API実装）
    │   │   ├── spotify/
    │   │   │   ├── gateway.go      # SpotifyAPI 実装
    │   │   │   └── types.go        # Spotify API レスポンス型
    │   │   ├── kkbox/
    │   │   │   └── gateway.go      # KKBOXAPI 実装
    │   │   ├── deezer/
    │   │   │   └── gateway.go      # DeezerAPI 実装 (BPM/Gain)
    │   │   ├── musicbrainz/
    │   │   │   └── gateway.go      # MusicBrainzAPI 実装 (Tags/Relations)
    │   │   ├── lastfm/
    │   │   │   └── gateway.go      # LastFMAPI 実装 (track.getSimilar)
    │   │   ├── ytmusic/
    │   │   │   └── gateway.go      # YouTubeMusicAPI 実装 (sidecar client)
    │   │   ├── cache/
    │   │   │   └── repository.go   # 2層キャッシュ TokenRepository 実装
    │   │   └── redis/
    │   │       └── repository.go   # Redis TokenRepository 実装
    │   ├── handler/                # Primary Adapters（HTTP Handler）
    │   │   ├── track.go            # トラック関連ハンドラー
    │   │   ├── artist.go           # アーティスト関連ハンドラー
    │   │   ├── album.go            # アルバム関連ハンドラー
    │   │   ├── recommend.go        # レコメンドハンドラー (V2)
    │   │   ├── response.go         # レスポンスヘルパー
    │   │   └── extract.go          # URL抽出ユーティリティ
    │   └── server/
    │       └── server.go           # HTTPサーバー・ルーティング
    │
    ├── config/
    │   └── config.go               # 設定
    │
    └── util/
        └── logger/
            └── logger.go           # ロギング
```

## レイヤー図

```
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│   ┌─────────────────────────────────────────────────────────────┐   │
│   │                     Adapter Layer                           │   │
│   │  ┌──────────────────┐          ┌──────────────────────────┐ │   │
│   │  │  Primary Adapter │          │   Secondary Adapter      │ │   │
│   │  │  (handler/)      │          │   (gateway/)             │ │   │
│   │  │                  │          │                          │ │   │
│   │  │  • TrackHandler  │          │  • spotify.Gateway       │ │   │
│   │  │  • ArtistHandler │          │  • kkbox.Gateway         │ │   │
│   │  │  • AlbumHandler  │          │  • deezer.Gateway        │ │   │
│   │  │                  │          │  • musicbrainz.Gateway   │ │   │
│   │  │                  │          │  • lastfm.Gateway        │ │   │
│   │  │                  │          │  • ytmusic.Gateway       │ │   │
│   │  │                  │          │  • redis.TokenRepository │ │   │
│   │  └────────┬─────────┘          └─────────────▲────────────┘ │   │
│   │           │                                  │              │   │
│   └───────────│──────────────────────────────────│──────────────┘   │
│               │                                  │                  │
│   ┌───────────▼──────────────────────────────────│──────────────┐   │
│   │                    UseCase Layer             │              │   │
│   │                                              │              │   │
│   │  • TrackUseCase         依存方向             │              │   │
│   │  • ArtistUseCase        ─────────►          │              │   │
│   │  • AlbumUseCase                              │              │   │
│   │  • SimilarTracksUseCase                      │              │   │
│   │  • RecommendUseCaseV2                        │              │   │
│   │                                              │              │   │
│   └───────────┬──────────────────────────────────│──────────────┘   │
│               │                                  │                  │
│   ┌───────────▼──────────────────────────────────│──────────────┐   │
│   │                     Port Layer               │              │   │
│   │                   (interfaces)               │              │   │
│   │                                              │              │   │
│   │  • SpotifyAPI (interface)  ◄─────────────────┘              │   │
│   │  • KKBOXAPI (interface)                                     │   │
│   │  • DeezerAPI (interface)                                    │   │
│   │  • MusicBrainzAPI (interface)                               │   │
│   │  • LastFMAPI (interface)                                    │   │
│   │  • YouTubeMusicAPI (interface)                              │   │
│   │  • TokenRepository (interface)                              │   │
│   │                                                             │   │
│   └───────────┬─────────────────────────────────────────────────┘   │
│               │                                                     │
│   ┌───────────▼─────────────────────────────────────────────────┐   │
│   │                    Domain Layer                             │   │
│   │                  (entities, errors)                         │   │
│   │                                                             │   │
│   │  • Track, Artist, Album, Image                              │   │
│   │  • SimilarTrack, SimpleTrack, SimpleArtist                  │   │
│   │  • ErrTrackNotFound, ErrISRCNotFound, etc.                  │   │
│   │                                                             │   │
│   │              ★ 依存なし（最も内側）                          │   │
│   └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## 依存性の流れ

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Handler    │────►│   UseCase    │────►│    Port      │◄────│   Gateway    │
│  (adapter)   │     │              │     │ (interface)  │     │  (adapter)   │
└──────────────┘     └──────────────┘     └──────────────┘     └──────────────┘
       │                    │                    │                    │
       │                    │                    │                    │
       ▼                    ▼                    ▼                    ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│                              Domain Layer                                    │
│                         (Track, Artist, Album, etc.)                         │
└──────────────────────────────────────────────────────────────────────────────┘

依存方向: Handler → UseCase → Port (interface) ← Gateway
          すべてのレイヤーは Domain を参照可能
```

## データフロー例: GET /v1/track/fetch

```
HTTP Request
     │
     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. Handler (adapter/handler/track.go)                                       │
│    - URLからSpotify Track IDを抽出                                           │
│    - UseCase を呼び出し                                                      │
│    - domain.Track を HTTP レスポンス形式に変換                                │
└─────────────────────────────────────────────────────────────────────────────┘
     │
     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. UseCase (usecase/track.go)                                               │
│    - バリデーション                                                          │
│    - SpotifyAPI interface を通じてデータ取得                                 │
│    - ビジネスロジック適用                                                    │
└─────────────────────────────────────────────────────────────────────────────┘
     │
     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. Gateway (adapter/gateway/spotify/gateway.go)                             │
│    - TokenRepository で認証トークン取得                                      │
│    - Spotify Web API 呼び出し                                               │
│    - 認証エラー時は自動でトークン再取得してリトライ                           │
│    - レスポンスを domain.Track に変換                                        │
└─────────────────────────────────────────────────────────────────────────────┘
     │
     ▼
HTTP Response (JSON)
```

## データフロー例: GET /v1/track/similar

```
HTTP Request
     │
     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. Handler (adapter/handler/track.go)                                       │
│    - URLからSpotify Track IDを抽出                                           │
└─────────────────────────────────────────────────────────────────────────────┘
     │
     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. SimilarTracksUseCase (usecase/v1/similar_tracks.go)                      │
│    ├── SpotifyAPI.GetTrackByID() → ISRC取得                                  │
│    ├── KKBOXAPI.SearchByISRC() → KKBOXトラックID取得                         │
│    ├── KKBOXAPI.GetRecommendedTracks() → レコメンド取得                      │
│    ├── KKBOXAPI.GetTrackDetail() → 各トラックのISRC取得                      │
│    ├── SpotifyAPI.SearchByISRC() × N（並列処理）                             │
│    └── 重複除去、ソート、上限適用                                            │
└─────────────────────────────────────────────────────────────────────────────┘
     │
     ▼
HTTP Response (JSON)
```

## データフロー例: GET /v2/track/recommend

```
HTTP Request
     │
     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. Handler (adapter/handler/recommend.go)                                   │
│    - URLからSpotify Track IDを抽出                                           │
│    - mode, limit パラメータ解析                                              │
└─────────────────────────────────────────────────────────────────────────────┘
     │
     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. RecommendUseCase (usecase/v2/recommend.go)                               │
│    ├── SpotifyAPI.GetTrackByID() → シードトラック情報取得                    │
│    ├── DeezerAPI.GetTrackByISRC() + MusicBrainzAPI → シード特徴量取得        │
│    ├── 並列候補収集:                                                         │
│    │   ├── KKBOXAPI.GetRecommendedTracks()  (30件)                          │
│    │   ├── LastFMAPI.GetSimilarTracks()     (30件) [optional]               │
│    │   ├── MusicBrainzAPI.GetArtistRecordings() (20件)                      │
│    │   └── YouTubeMusicAPI.SearchTracks()   (25件) [optional, sidecar]      │
│    ├── 重複除去（ISRC/name+artist）                                          │
│    ├── Spotify検索による候補情報補完:                                        │
│    │   ├── ISRC持ち → SearchByISRC()                                        │
│    │   └── 名前のみ → searchSpotifyWithFallback() ※4段階フォールバック      │
│    │       ├── 1. 厳密検索 (track:xxx artist:yyy)                           │
│    │       ├── 2. 簡素化曲名で検索                                           │
│    │       ├── 3. フリーテキスト検索                                         │
│    │       └── 4. 簡素化フリーテキスト検索                                   │
│    ├── 候補の特徴量並列取得（Deezer/MusicBrainz）                            │
│    ├── ジャンルフィルタリング（アニソン保護等）                               │
│    ├── SimilarityCalculatorV2 で類似度計算                                   │
│    │   └── BPM, Duration, Gain, TagSimilarity + アーティスト関係ボーナス     │
│    └── スコア順ソート、上限適用                                              │
└─────────────────────────────────────────────────────────────────────────────┘
     │
     ▼
HTTP Response (JSON)
```

## 依存性注入 (DI)

`cmd/server/main.go` で全ての依存関係を組み立てます：

```go
// 1. Infrastructure
tokenRepo := cache.NewCachedTokenRepository(redisRepo)

// 2. Gateways (port interface を実装)
spotifyGW := spotify.NewGateway(clientID, secret, tokenRepo)
kkboxGW := kkbox.NewGateway(clientID, secret, tokenRepo)
deezerGW := deezer.NewGateway()
musicbrainzGW := musicbrainz.NewGateway(userAgent)
lastfmGW := lastfm.NewGateway(apiKey)        // optional
ytmusicGW := ytmusic.NewGateway(sidecarURL)  // optional

// 3. UseCases (port interface に依存)
trackUC := usecasev1.NewTrackUseCase(spotifyGW)
artistUC := usecasev1.NewArtistUseCase(spotifyGW)
albumUC := usecasev1.NewAlbumUseCase(spotifyGW)
similarUC := usecasev1.NewSimilarTracksUseCase(spotifyGW, kkboxGW)
recommendUC := usecasev2.NewRecommendUseCaseFull(
    spotifyGW, kkboxGW, deezerGW, musicbrainzGW, lastfmGW, ytmusicGW,
)

// 4. Handlers (usecase に依存)
trackHandler := handler.NewTrackHandler(trackUC, similarUC)
recommendHandler := handler.NewRecommendHandler(recommendUC)
healthHandler := handler.NewHealthHandler(enabledServices)

// 5. Server
server.New(config, handlers)
```

## ヘルスチェック

`/healthz` エンドポイントは以下の情報を提供:

| 項目       | 内容                                         |
| ---------- | -------------------------------------------- |
| status     | サーバー状態 (`healthy`)                     |
| version    | ビルド時に設定されたバージョン               |
| build_time | ビルド日時（ISO 8601）                       |
| git_commit | Git コミットハッシュ                         |
| uptime     | サーバー起動からの経過時間                   |
| runtime    | Go バージョン、ゴルーチン数、CPU 数、OS/Arch |
| services   | 各外部サービスの有効/無効状態                |

### ビルド時のバージョン設定

```bash
go build -ldflags="-X main.version=1.0.0 \
  -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -X main.gitCommit=$(git rev-parse --short HEAD)" \
  -o server ./cmd/server/...
```

Docker Compose では環境変数で設定:

```bash
VERSION=1.0.0 BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) GIT_COMMIT=$(git rev-parse --short HEAD) docker compose up -d --build
```

## Clean Architecture の利点

| 利点               | 説明                                                                          |
| ------------------ | ----------------------------------------------------------------------------- |
| **テスタビリティ** | interface によりモック化が容易。UseCase のテストで実際の API を呼ばなくて良い |
| **変更容易性**     | Spotify → 別サービスに変更する場合、Gateway だけ差し替えれば良い              |
| **関心の分離**     | HTTP 処理(Handler)、ビジネスロジック(UseCase)、外部 API(Gateway)が分離        |
| **依存性の逆転**   | UseCase は具体的な実装ではなく interface に依存                               |

## 各レイヤーの責務

### Domain Layer

- 外部に依存しない純粋なビジネスエンティティ
- `Track`, `Artist`, `Album` などの型定義
- ドメインエラーの定義

### Port Layer

- UseCase が必要とする外部機能の interface 定義
- `SpotifyAPI`, `KKBOXAPI`, `TokenRepository`

### UseCase Layer

- アプリケーションのビジネスロジック
- Port interface を通じて外部機能を利用
- Domain エンティティを操作

### Adapter Layer

- 外部との接続を担当
- **Handler**: HTTP ↔ Domain の変換
- **Gateway**: Domain ↔ 外部 API の変換
- **Server**: HTTP サーバー設定

## 起動方法

```bash
go run ./cmd/server/...
```

## API エンドポイント

| Method | Path                | Handler                               | 説明                                       |
| ------ | ------------------- | ------------------------------------- | ------------------------------------------ |
| GET    | /healthz            | HealthHandler.Check                   | ヘルスチェック（バージョン・サービス状態） |
| GET    | /v1/track/fetch     | TrackHandler.FetchByURL               | Spotify URL からトラック情報取得           |
| GET    | /v1/track/search    | TrackHandler.Search                   | キーワードでトラック検索                   |
| GET    | /v1/track/similar   | TrackHandler.FetchSimilar             | KKBOX ベースの類似トラック取得             |
| GET    | /v2/track/recommend | RecommendHandler.FetchRecommendations | マルチソースレコメンド取得                 |
| GET    | /v1/artist/fetch    | ArtistHandler.FetchByURL              | Spotify URL からアーティスト情報取得       |
| GET    | /v1/album/fetch     | AlbumHandler.FetchByURL               | Spotify URL からアルバム情報取得           |
