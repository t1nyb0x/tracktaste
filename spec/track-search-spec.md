# トラック検索 API 仕様書

## 処理フロー

1. リクエストを行う
2. クエリパラメータ `q` から検索キーワードを取得する
3. [Spotify Search API](./spotify/search-api.md)へキーワードを使ってリクエストする
4. SpotifyAPI から返ってきたデータを加工する
5. レスポンスとして返す

## リクエスト

ENDPOINT: /track/search?q={query}

METHOD: GET

Param example: 米津玄師 Lemon

## レスポンス

```json
{
  "status": 200,
  "result": {
    "items": [
      {
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
              "id": string
            }
          ]
        },
        "artists": [
          {
            "url": string,
            "id": string,
            "name": string
          }
        ],
        "disc_number": int,
        "popularity": int|null,
        "isrc": string|null,
        "url": string,
        "id": string,
        "name": string,
        "track_number": int,
        "duration_ms": int,
        "explicit": bool
      }
    ]
  }
}
```

## エラーレスポンス

### 検索クエリが空だったとき

```json
{
  "status": 400,
  "message": "検索クエリが入力されていません",
  "code": "EMPTY_QUERY"
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
