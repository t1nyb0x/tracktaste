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
└── internal/
    ├── domain/                      # ドメイン層（最も内側）
    │   ├── track.go                # Track, SimpleTrack, SimilarTrack
    │   ├── artist.go               # Artist, SimpleArtist
    │   ├── album.go                # Album
    │   ├── image.go                # Image
    │   └── errors.go               # ドメインエラー定義
    │
    ├── port/                        # ポート層（インターフェース定義）
    │   ├── repository/
    │   │   └── token.go            # TokenRepository interface
    │   └── external/
    │       ├── spotify.go          # SpotifyAPI interface
    │       └── kkbox.go            # KKBOXAPI interface
    │
    ├── usecase/                     # ユースケース層（ビジネスロジック）
    │   ├── track.go                # TrackUseCase
    │   ├── artist.go               # ArtistUseCase
    │   ├── album.go                # AlbumUseCase
    │   └── similar_tracks.go       # SimilarTracksUseCase
    │
    ├── adapter/                     # アダプター層（最も外側）
    │   ├── gateway/                # Secondary Adapters（外部API実装）
    │   │   ├── spotify/
    │   │   │   ├── gateway.go      # SpotifyAPI 実装
    │   │   │   └── types.go        # Spotify API レスポンス型
    │   │   ├── kkbox/
    │   │   │   └── gateway.go      # KKBOXAPI 実装
    │   │   └── redis/
    │   │       └── repository.go   # TokenRepository 実装
    │   ├── handler/                # Primary Adapters（HTTP Handler）
    │   │   ├── track.go            # トラック関連ハンドラー
    │   │   ├── artist.go           # アーティスト関連ハンドラー
    │   │   ├── album.go            # アルバム関連ハンドラー
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
│   │  │  • AlbumHandler  │          │  • redis.TokenRepository │ │   │
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
│   │                                              │              │   │
│   └───────────┬──────────────────────────────────│──────────────┘   │
│               │                                  │                  │
│   ┌───────────▼──────────────────────────────────│──────────────┐   │
│   │                     Port Layer               │              │   │
│   │                   (interfaces)               │              │   │
│   │                                              │              │   │
│   │  • SpotifyAPI (interface)  ◄─────────────────┘              │   │
│   │  • KKBOXAPI (interface)                                     │   │
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
│ 2. SimilarTracksUseCase (usecase/similar_tracks.go)                         │
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

## 依存性注入 (DI)

`cmd/server/main.go` で全ての依存関係を組み立てます：

```go
// 1. Infrastructure
tokenRepo := redis.NewTokenRepository()

// 2. Gateways (port interface を実装)
spotifyGW := spotify.NewGateway(clientID, secret, tokenRepo)
kkboxGW := kkbox.NewGateway(clientID, secret, tokenRepo)

// 3. UseCases (port interface に依存)
trackUC := usecase.NewTrackUseCase(spotifyGW)     // SpotifyAPI interface
similarUC := usecase.NewSimilarTracksUseCase(spotifyGW, kkboxGW)

// 4. Handlers (usecase に依存)
trackHandler := handler.NewTrackHandler(trackUC, similarUC)

// 5. Server
server.New(config, handlers)
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

| Method | Path              | Handler                   |
| ------ | ----------------- | ------------------------- |
| GET    | /healthz          | Health check              |
| GET    | /v1/track/fetch   | TrackHandler.FetchByURL   |
| GET    | /v1/track/search  | TrackHandler.Search       |
| GET    | /v1/track/similar | TrackHandler.FetchSimilar |
| GET    | /v1/artist/fetch  | ArtistHandler.FetchByURL  |
| GET    | /v1/album/fetch   | AlbumHandler.FetchByURL   |
