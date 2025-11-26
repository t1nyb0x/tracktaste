# トラック情報取得 API 仕様書

## 処理フロー

1. リクエストを行う
2. パラメータに入っている URL から ID を取得する（例：22ev7LxXzh9gZ274L5UG9c）
3. [Spotify Get track API](./spotify/fetch-track-api.md)へこの ID を使ってリクエストする
4. SpotifyAPI から返ってきたデータを加工する
5. レスポンスとして返す

## リクエスト

ENDPOINT: /track/fetch?url={url}

METHOD: GET

Param example: https://open.spotify.com/intl-ja/track/22ev7LxXzh9gZ274L5UG9c?si=160b83922bf7467c

## レスポンス

```json
{
  "status": 200,
  "result": {
    "album": {
        "url": string,
        "id": string,
        "images": [
            {
                "url": string,
                "height": int,
                "width": int
            }
        ],
        "name": string,
        "release_date": string,
        "artists": [
            {
                "url": string,
                "name": string,
                "id": string,
            }
        ],
    },
    "artists": [
        {
            "url": string,
            "id": string,
            "name": string,
        }
    ],
    "disc_number": int,
    "popularity": int|null,
    "isrc": string|null,
    "url": string,
    "id": string,
    "track_number": int
  }
}
```

## エラーレスポンス

### パラメータが空だったとき

```json
{
  "status": 400,
  "message": "URLが入力されていません",
  "code": "EMPTY_URL"
}
```

### パラメータが不正だったとき

```json
{
  "status": 400,
  "message": "パラメータが不正です",
  "code": "INVALID_PARAM"
}
```

### Spotify 以外の URL だったとき

```json
{
  "status": 400,
  "message": "SpotifyのURLを入力してください",
  "code": "INVALID_URL"
}
```

### Spotify の TrackURL 以外だったとき

```json
{
  "status": 400,
  "message": "TrackURLを入力してください",
  "code": "DIFFERENT_SPOTIFY_URL"
}
```

### SpotifyAPI に何らかの異常が発生しているとき

```json
{
  "status": 503,
  "message": "Spotify APIで問題が発生しているようです",
  "code": "SOMETHING_SPOTIFY_ERROR"
}
```
