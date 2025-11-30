# 類似トラック情報取得 API 仕様書

## 処理フロー

以下の順で情報を取得します

1. SpotifyAPI からトラック情報を取得する
2. トラック情報に含まれている isrc と upc を抽出する（存在しない場合もある）
3. KKBOXAPI から Search various objects を使って情報を取得する
4. 3 で取得した情報から id を抽出する
5. KKBOXAPI から Get recommended tracks for a given track を使って情報を取得する
6. TrackID を抽出して KKBOX の Get a Track API を使って取得する
7. 6 で取得した情報から isrc と upc を抽出する（存在しない場合もある）
8. 7 で抽出した値を利用して SpotifyAPI で Search for item を使って検索を行う
   - ISRC の場合: `q=isrc:{value}&type=track`
   - UPC の場合: `q=upc:{value}&type=track`
9. 検索結果が 0 件の場合はスキップ
10. 最大 30 件まで取得する。重複分は取り除く
11. 7 で検索していった結果をレスポンスの形に加工して配列にまとめる
12. まとめた配列を返す

## 仕様

- 検索は並列処理で行う
  - 最大同時接続数は 5
  - 各リクエストのタイムアウトは 5 秒
  - 全体のタイムアウトは 30 秒
- 同じ曲が複数回出た場合は重複分を取り除く（isrc で判定する）
- ソート順は popularity 降順

## KKBOXAPI 仕様

### Search various objects

[KKBOXAPI Search various objects API 仕様](./kkbox/kkbox-search-api.md) を参照してください

### Get recommend tracks for a given track

[# KKBOX Get recommend tracks for a given track API 仕様](./kkbox/kkbox-get-recommend-tracks-api.md)を参照してください

## SpotifyAPI 仕様

### Search for item

[Spotify Search for item API 仕様](./spotify/search-api.md)を参照してください

## リクエスト

ENDPOINT: /track/similar?url={url}

METHOD: GET

Param example: https://open.spotify.com/intl-ja/track/22ev7LxXzh9gZ274L5UG9c?si=160b83922bf7467c

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
                        "id": string,
                    }
                ]
            },
            "isrc": string|null,
            "upc": string|null,
            "url": string,
            "id": string,
            "name": string,
            "popularity": int|null,
            "track_number": int,
            "duration_ms": int,
            "explicit": bool
        }
    ]
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

### KKBOXAPI に何らかの異常が発生しているとき

```json
{
  "status": 503,
  "message": "KKBOX APIで問題が発生しているようです",
  "code": "SOMETHING_KKBOX_ERROR"
}
```

### タイムアウト時

```json
{
  "status": 504,
  "message": "処理がタイムアウトしました",
  "code": "REQUEST_TIMEOUT"
}
```
