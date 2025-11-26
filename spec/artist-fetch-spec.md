# アーティスト情報取得 API 仕様書

## 処理フロー

1. リクエストを行う
2. パラメータに入っている URL から ID を取得する（例：76pJHTMyTukyWJNR6yRrZS）
3. [Spotify Get artist API](./spotify/fetch-artist-api.md)へこの ID を使ってリクエストする
4. SpotifyAP いから返ってきたデータを加工する
5. レスポンスとして返す

## リクエスト

ENDPOINT: /artist/fetch?url={url}

METHOD: GET

Param example: https://open.spotify.com/intl-ja/artist/76pJHTMyTukyWJNR6yRrZS?si=RngryHkbS0aggUNjeL55PQ

## レスポンス

```json
{
  "status": 200,
  "result": {
    "url": string,
    "followers": string,
    "genres": string[],
    "id": string,
    "images": [
        {
            "url": string,
            "height": int,
            "width": int
        },
    ],
    "name": string,
    "popularity": int|null,
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

### Spotify の ArtistURL 以外だったとき

```json
{
  "status": 400,
  "message": "ArtistURLを入力してください",
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
