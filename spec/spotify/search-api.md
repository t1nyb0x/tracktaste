# Spotify Search for item API 仕様

## Request Example

ENDPOINT GET https://api.spotify.com/v1/search?q={query}&type=track&limit={limit}&offset={offset}

query isrc:TCJPC2483393 (ucp の場合は upc:{upc})

limit 100

offset 0

## Response Example

```json
{
  "tracks": {
    "href": "https://api.spotify.com/v1/search?offset=0&limit=20&query=isrc%3ATCJPC2483393&type=track&locale=ja;q%3D0.7",
    "limit": 20,
    "next": null,
    "offset": 0,
    "previous": null,
    "total": 1,
    "items": [
      {
        "album": {
          "album_type": "album",
          "total_tracks": 12,
          "available_markets": ["JP"],
          "external_urls": {
            "spotify": "https://open.spotify.com/album/0iiVne9c8LZC0iuhOBiTiL"
          },
          "href": "https://api.spotify.com/v1/albums/0iiVne9c8LZC0iuhOBiTiL",
          "id": "0iiVne9c8LZC0iuhOBiTiL",
          "images": [
            {
              "url": "https://i.scdn.co/image/ab67616d0000b2733c24633d162fc1fad1a9ce4e",
              "height": 640,
              "width": 640
            },
            {
              "url": "https://i.scdn.co/image/ab67616d00001e023c24633d162fc1fad1a9ce4e",
              "height": 300,
              "width": 300
            },
            {
              "url": "https://i.scdn.co/image/ab67616d000048513c24633d162fc1fad1a9ce4e",
              "height": 64,
              "width": 64
            }
          ],
          "name": "箱庭共鳴-ハコニワレゾナンス-Hanon×Kotoha 歌唱版",
          "release_date": "2025-01-29",
          "release_date_precision": "day",
          "type": "album",
          "uri": "spotify:album:0iiVne9c8LZC0iuhOBiTiL",
          "artists": [
            {
              "external_urls": {
                "spotify": "https://open.spotify.com/artist/3mpsjaIGwvF17DDOof3njV"
              },
              "href": "https://api.spotify.com/v1/artists/3mpsjaIGwvF17DDOof3njV",
              "id": "3mpsjaIGwvF17DDOof3njV",
              "name": "Hanon",
              "type": "artist",
              "uri": "spotify:artist:3mpsjaIGwvF17DDOof3njV"
            },
            {
              "external_urls": {
                "spotify": "https://open.spotify.com/artist/76pJHTMyTukyWJNR6yRrZS"
              },
              "href": "https://api.spotify.com/v1/artists/76pJHTMyTukyWJNR6yRrZS",
              "id": "76pJHTMyTukyWJNR6yRrZS",
              "name": "Kotoha",
              "type": "artist",
              "uri": "spotify:artist:76pJHTMyTukyWJNR6yRrZS"
            }
          ],
          "is_playable": true
        },
        "artists": [
          {
            "external_urls": {
              "spotify": "https://open.spotify.com/artist/3mpsjaIGwvF17DDOof3njV"
            },
            "href": "https://api.spotify.com/v1/artists/3mpsjaIGwvF17DDOof3njV",
            "id": "3mpsjaIGwvF17DDOof3njV",
            "name": "Hanon",
            "type": "artist",
            "uri": "spotify:artist:3mpsjaIGwvF17DDOof3njV"
          },
          {
            "external_urls": {
              "spotify": "https://open.spotify.com/artist/76pJHTMyTukyWJNR6yRrZS"
            },
            "href": "https://api.spotify.com/v1/artists/76pJHTMyTukyWJNR6yRrZS",
            "id": "76pJHTMyTukyWJNR6yRrZS",
            "name": "Kotoha",
            "type": "artist",
            "uri": "spotify:artist:76pJHTMyTukyWJNR6yRrZS"
          }
        ],
        "available_markets": ["JP"],
        "disc_number": 1,
        "duration_ms": 221200,
        "explicit": false,
        "external_ids": {
          "isrc": "TCJPC2483393"
        },
        "external_urls": {
          "spotify": "https://open.spotify.com/track/22ev7LxXzh9gZ274L5UG9c"
        },
        "href": "https://api.spotify.com/v1/tracks/22ev7LxXzh9gZ274L5UG9c",
        "id": "22ev7LxXzh9gZ274L5UG9c",
        "is_playable": true,
        "name": "Wave",
        "popularity": 9,
        "preview_url": null,
        "track_number": 3,
        "type": "track",
        "uri": "spotify:track:22ev7LxXzh9gZ274L5UG9c",
        "is_local": false
      }
    ]
  }
}
```

API 仕様は https://developer.spotify.com/documentation/web-api/reference/search を参照してください
