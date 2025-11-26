# Spotify Get artist API 仕様

## Request Example

ArtistURL https://open.spotify.com/intl-ja/artist/76pJHTMyTukyWJNR6yRrZS?si=wIyvE4OmT6ejCkLw0ZbRAg

Endpoint GET https://api.spotify.com/v1/artists/{id}

id 76pJHTMyTukyWJNR6yRrZS

## Response Example

```json
{
  "external_urls": {
    "spotify": "https://open.spotify.com/artist/76pJHTMyTukyWJNR6yRrZS"
  },
  "followers": {
    "href": null,
    "total": 47909
  },
  "genres": ["j-pop"],
  "href": "https://api.spotify.com/v1/artists/76pJHTMyTukyWJNR6yRrZS?locale=ja%3Bq%3D0.7",
  "id": "76pJHTMyTukyWJNR6yRrZS",
  "images": [
    {
      "url": "https://i.scdn.co/image/ab67616d0000b273034e9d8fd0c27383982e3b99",
      "height": 640,
      "width": 640
    },
    {
      "url": "https://i.scdn.co/image/ab67616d00001e02034e9d8fd0c27383982e3b99",
      "height": 300,
      "width": 300
    },
    {
      "url": "https://i.scdn.co/image/ab67616d00004851034e9d8fd0c27383982e3b99",
      "height": 64,
      "width": 64
    }
  ],
  "name": "Kotoha",
  "popularity": 45,
  "type": "artist",
  "uri": "spotify:artist:76pJHTMyTukyWJNR6yRrZS"
}
```

API 仕様は https://developer.spotify.com/documentation/web-api/reference/get-an-artist を参照してください
