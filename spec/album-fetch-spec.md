# アルバム情報取得 API 仕様書

## 処理フロー

1. リクエストを行う
2. パラメータに入っている URL から ID を取得する（例：4ft2GMEQ8itLEL66WX9lfi）
3. [Spotify Get album API](./spotify/fetch-album-api.md)へこの ID を使ってリクエストする
4. SpotifyApi から返ってきたデータを加工する
5. レスポンスとして返す

## リクエスト

ENDPOINT: /album/fetch?url={url}

METHOD: GET

Param example: https://open.spotify.com/intl-ja/album/4ft2GMEQ8itLEL66WX9lfi?si=wmWLw4O1Qpak0Fs0_ANMaQ

## レスポンス

```json
{
  "status": 200,
  "result": {
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
        }
    ],
    "tracks": {
        "items": [
            {
                "artists": [
                    {
                        "url": string,
                        "name": string,
                    }
                ],
                "url": string,
                "id": string,
                "name": string,
                "track_number": int,
            }
        ]
    },
    "popularity": int|null,
    "upc": string|null,
    "genres": string[],
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

### Spotify の AlbumURL 以外だったとき

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
