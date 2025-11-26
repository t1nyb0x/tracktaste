# TrackTaste

- [TrackTaste](#tracktaste)
  - [概要](#概要)
  - [API エンドポイント](#api-エンドポイント)
  - [URL からの ID 抽出](#url-からの-id-抽出)
    - [対応する URL パターン](#対応する-url-パターン)
      - [トラック](#トラック)
      - [アーティスト](#アーティスト)
      - [アルバム](#アルバム)
    - [エラーケース](#エラーケース)
  - [仕様](#仕様)
    - [機能](#機能)
    - [トラック情報取得 API](#トラック情報取得-api)
    - [アーティスト情報取得 API](#アーティスト情報取得-api)
    - [アルバム情報取得 API](#アルバム情報取得-api)
    - [類似トラック取得 API](#類似トラック取得-api)
  - [API リクエストについて](#api-リクエストについて)
    - [Spotify の BearerToken 取得について](#spotify-の-bearertoken-取得について)
      - [レスポンス](#レスポンス)
    - [KKBOX の BearerToken 取得について](#kkbox-の-bearertoken-取得について)
      - [レスポンス](#レスポンス-1)
  - [Redis 接続情報](#redis-接続情報)
    - [キー設計](#キー設計)
    - [保存形式](#保存形式)
  - [ロギングについて](#ロギングについて)
  - [ロギングフォーマット](#ロギングフォーマット)

## 概要

TrackTaste は SpotifyURL から各種情報を取得する Go 製 API です。

ルーティングは chi を使用し、エラー出力には xerrors を適宜使用します。

## API エンドポイント

トラック情報取得: GET /track/fetch?url={url} 実装済みだが、ここの仕様とは異なるレスポンスが返ってきてる状態。要改修。

アーティスト情報取得: GET /artist/fetch?url={url} 実装済みだが、lastfm を使うようになっている。Spotify で取得するように変更したい。

アルバム情報取得: GET /album/fetch?url={url} 未実装

類似トラック取得: GET /track/similar?url={url} 未実装

## URL からの ID 抽出

internal/util/extractSpotifyTrackId を使用するが、artist と album の ID を抽出する処理ができていないため、extractSpotifyTrackId を参照して同様な実装を行う

### 対応する URL パターン

#### トラック

- `https://open.spotify.com/track/{id}`
- `https://open.spotify.com/intl-ja/track/{id}`
- `https://open.spotify.com/track/{id}?si=xxx`

#### アーティスト

- `https://open.spotify.com/artist/{id}`
- `https://open.spotify.com/intl-ja/artist/{id}`
- `https://open.spotify.com/artist/{id}?si=xxx`

#### アルバム

- `https://open.spotify.com/album/{id}`
- `https://open.spotify.com/intl-ja/album/{id}`
- `https://open.spotify.com/album/{id}?si=xxx`

### エラーケース

- 不正な URL 形式 → `INVALID_URL`
- Spotify 以外の URL → `NOT_SPOTIFY_URL`
- 対象外のリソースタイプ → `INVALID_RESOURCE_TYPE`
- 空 → `EMPTY_PARAM`

## 仕様

### 機能

- トラック情報取得 API
- アーティスト情報取得 API
- アルバム情報取得 API
- 類似トラック取得 API

### トラック情報取得 API

[仕様書](./track-fetch-spec.md)を参照

### アーティスト情報取得 API

[仕様書](./artist-fetch-spec.md)を参照

### アルバム情報取得 API

[仕様書](./album-fetch-spec.md)を参照

### 類似トラック取得 API

[仕様書](./track-similar-spec.md)を参照

## API リクエストについて

API リクエストを行う場合は、Bearer Token の取得が必要です。

Spotify は実装済みですが、KKBOX は未実装のため実装を行ってください。

リクエスト時は、BearerToken を Authorization ヘッダーに入れて行ってください。

BearerToken は Redis で管理を行います。

サービス名、トークン、有効期限 yyyy/MM/dd HH:mm:ss で管理すること。

発行したトークンは Redis に保存し、基本的に Redis からトークンを取得します。

トークンが有効期限切れの場合、トークンを再発行します。

### Spotify の BearerToken 取得について

以下の情報を使って BearerToken を取得してください

ENDPOINT POST https://accounts.spotify.com/api/token

Content-Type application/x-www-form-urlencoded

grant_type: client_credentials

client_id: .env の SPOTIFY_CLIENT_ID
client_secret: .env の SPOTIFY_CLIENT_SECRET

BearerToken は戻り値の access_token に格納されています。

BearerToken は有効期限が expires_in に秒数で記載れています。

#### レスポンス

```json
{
  "access_token": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

### KKBOX の BearerToken 取得について

以下の情報を使って BearerToken を取得してください

ENDPOINT POST https://account.kkbox.com/oauth2/token

Content-Type application/x-www-form-urlencoded

grant_type: client_credentials

client_id: .env の KKBOX_ID
client_secret: .env の KKBOX_SECRET

BearerToken は戻り値の access_token に格納されています。

BearerToken は有効期限が expires_in に秒数で記載れています。

#### レスポンス

```json
{
  "access_token": "xxxxxxxxxxxxxxxxxxxxxxxxxx",
  "token_type": "Bearer",
  "expires_in": 576308
}
```

KKBOX は Spotify と access_token の expires_in が異なります

## Redis 接続情報

- Host: 環境変数 `REDIS_HOST` (デフォルト: localhost)
- Port: 環境変数 `REDIS_PORT` (デフォルト: 6379)
- Password: 環境変数 `REDIS_PASSWORD` (デフォルト: 空)
- DB: 環境変数 `REDIS_DB` (デフォルト: 0)

Docker で利用を想定

### キー設計

- Spotify Token: `token:spotify`
- KKBOX Token: `token:kkbox`

### 保存形式

```json
{
  "access_token": "xxx",
  "expires_at": 1732622400
}
```

## ロギングについて

機能ごとにログを出力するようにしてください。必要であればロギングライブラリまたは、自前でのロギングを実装してください。

ロギングは以下のカテゴリに分けてください

- [FATAL]
- [WARNING]
- [INFO]
- [DEBUG]

## ロギングフォーマット

[LEVEL] YYYY-MM-DD HH:mm:ss [機能名] メッセージ

例:
[INFO] 2025-11-26 10:30:00 [TrackFetch] Spotify API リクエスト開始
[ERROR] 2025-11-26 10:30:01 [TrackFetch] Spotify API エラー: 401 Unauthorized
