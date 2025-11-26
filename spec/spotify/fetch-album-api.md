# Spotify Get Album API 仕様

## Request Example

AlbumURL https://open.spotify.com/intl-ja/album/0iiVne9c8LZC0iuhOBiTiL?si=qIxvNLRwT-iAyc2xqc9V4g

Endpoint GET https://api.spotify.com/v1/albums/{id}

id 0iiVne9c8LZC0iuhOBiTiL

## Response Example

```json
{
  "album_type": "album",
  "total_tracks": 12,
  "available_markets": ["JP"],
  "external_urls": {
    "spotify": "https://open.spotify.com/album/0iiVne9c8LZC0iuhOBiTiL"
  },
  "href": "https://api.spotify.com/v1/albums/0iiVne9c8LZC0iuhOBiTiL?locale=ja%3Bq%3D0.7",
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
  "tracks": {
    "href": "https://api.spotify.com/v1/albums/0iiVne9c8LZC0iuhOBiTiL/tracks?offset=0&limit=50&locale=ja;q%3D0.7",
    "limit": 50,
    "next": null,
    "offset": 0,
    "previous": null,
    "total": 12,
    "items": [
      {
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
          }
        ],
        "available_markets": ["JP"],
        "disc_number": 1,
        "duration_ms": 120440,
        "explicit": false,
        "external_urls": {
          "spotify": "https://open.spotify.com/track/2g5OzQBXNqDLjzm1SeOnf1"
        },
        "href": "https://api.spotify.com/v1/tracks/2g5OzQBXNqDLjzm1SeOnf1",
        "id": "2g5OzQBXNqDLjzm1SeOnf1",
        "name": "未練タラレバ",
        "preview_url": null,
        "track_number": 1,
        "type": "track",
        "uri": "spotify:track:2g5OzQBXNqDLjzm1SeOnf1",
        "is_local": false
      },
      {
        "artists": [
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
        "duration_ms": 171640,
        "explicit": false,
        "external_urls": {
          "spotify": "https://open.spotify.com/track/2tIWQRu9lEsQh6OO8NON23"
        },
        "href": "https://api.spotify.com/v1/tracks/2tIWQRu9lEsQh6OO8NON23",
        "id": "2tIWQRu9lEsQh6OO8NON23",
        "name": "出来心",
        "preview_url": null,
        "track_number": 2,
        "type": "track",
        "uri": "spotify:track:2tIWQRu9lEsQh6OO8NON23",
        "is_local": false
      },
      {
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
        "external_urls": {
          "spotify": "https://open.spotify.com/track/22ev7LxXzh9gZ274L5UG9c"
        },
        "href": "https://api.spotify.com/v1/tracks/22ev7LxXzh9gZ274L5UG9c",
        "id": "22ev7LxXzh9gZ274L5UG9c",
        "name": "Wave",
        "preview_url": null,
        "track_number": 3,
        "type": "track",
        "uri": "spotify:track:22ev7LxXzh9gZ274L5UG9c",
        "is_local": false
      },
      {
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
          }
        ],
        "available_markets": ["JP"],
        "disc_number": 1,
        "duration_ms": 151360,
        "explicit": false,
        "external_urls": {
          "spotify": "https://open.spotify.com/track/0n388D18Id0BHF5sBul3vk"
        },
        "href": "https://api.spotify.com/v1/tracks/0n388D18Id0BHF5sBul3vk",
        "id": "0n388D18Id0BHF5sBul3vk",
        "name": "イデア",
        "preview_url": null,
        "track_number": 4,
        "type": "track",
        "uri": "spotify:track:0n388D18Id0BHF5sBul3vk",
        "is_local": false
      },
      {
        "artists": [
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
        "duration_ms": 236280,
        "explicit": false,
        "external_urls": {
          "spotify": "https://open.spotify.com/track/0hWr6Pl1ddL9pAlAGoq3YP"
        },
        "href": "https://api.spotify.com/v1/tracks/0hWr6Pl1ddL9pAlAGoq3YP",
        "id": "0hWr6Pl1ddL9pAlAGoq3YP",
        "name": "花咲娘",
        "preview_url": null,
        "track_number": 5,
        "type": "track",
        "uri": "spotify:track:0hWr6Pl1ddL9pAlAGoq3YP",
        "is_local": false
      },
      {
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
          }
        ],
        "available_markets": ["JP"],
        "disc_number": 1,
        "duration_ms": 181080,
        "explicit": false,
        "external_urls": {
          "spotify": "https://open.spotify.com/track/7afyxolFUK4AmgbeeJSRAJ"
        },
        "href": "https://api.spotify.com/v1/tracks/7afyxolFUK4AmgbeeJSRAJ",
        "id": "7afyxolFUK4AmgbeeJSRAJ",
        "name": "Play the Future",
        "preview_url": null,
        "track_number": 6,
        "type": "track",
        "uri": "spotify:track:7afyxolFUK4AmgbeeJSRAJ",
        "is_local": false
      },
      {
        "artists": [
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
        "duration_ms": 201666,
        "explicit": false,
        "external_urls": {
          "spotify": "https://open.spotify.com/track/4XlUdDWSQgNHRqqsDiZsxq"
        },
        "href": "https://api.spotify.com/v1/tracks/4XlUdDWSQgNHRqqsDiZsxq",
        "id": "4XlUdDWSQgNHRqqsDiZsxq",
        "name": "Confetti",
        "preview_url": null,
        "track_number": 7,
        "type": "track",
        "uri": "spotify:track:4XlUdDWSQgNHRqqsDiZsxq",
        "is_local": false
      },
      {
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
          }
        ],
        "available_markets": ["JP"],
        "disc_number": 1,
        "duration_ms": 261453,
        "explicit": false,
        "external_urls": {
          "spotify": "https://open.spotify.com/track/0hl7Fzxccz3lIbwMh3qyEp"
        },
        "href": "https://api.spotify.com/v1/tracks/0hl7Fzxccz3lIbwMh3qyEp",
        "id": "0hl7Fzxccz3lIbwMh3qyEp",
        "name": "静寂",
        "preview_url": null,
        "track_number": 8,
        "type": "track",
        "uri": "spotify:track:0hl7Fzxccz3lIbwMh3qyEp",
        "is_local": false
      },
      {
        "artists": [
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
        "duration_ms": 136026,
        "explicit": false,
        "external_urls": {
          "spotify": "https://open.spotify.com/track/1CZGdLoQcVFGQyGDKHasYh"
        },
        "href": "https://api.spotify.com/v1/tracks/1CZGdLoQcVFGQyGDKHasYh",
        "id": "1CZGdLoQcVFGQyGDKHasYh",
        "name": "signal",
        "preview_url": null,
        "track_number": 9,
        "type": "track",
        "uri": "spotify:track:1CZGdLoQcVFGQyGDKHasYh",
        "is_local": false
      },
      {
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
          }
        ],
        "available_markets": ["JP"],
        "disc_number": 1,
        "duration_ms": 223466,
        "explicit": false,
        "external_urls": {
          "spotify": "https://open.spotify.com/track/2d21NJpFMtyb0Ix7L2eVgo"
        },
        "href": "https://api.spotify.com/v1/tracks/2d21NJpFMtyb0Ix7L2eVgo",
        "id": "2d21NJpFMtyb0Ix7L2eVgo",
        "name": "仮面ノ少女",
        "preview_url": null,
        "track_number": 10,
        "type": "track",
        "uri": "spotify:track:2d21NJpFMtyb0Ix7L2eVgo",
        "is_local": false
      },
      {
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
        "duration_ms": 143826,
        "explicit": false,
        "external_urls": {
          "spotify": "https://open.spotify.com/track/5MIZW9FYY9aIbyKir7nYSB"
        },
        "href": "https://api.spotify.com/v1/tracks/5MIZW9FYY9aIbyKir7nYSB",
        "id": "5MIZW9FYY9aIbyKir7nYSB",
        "name": "恋愛ロジック",
        "preview_url": null,
        "track_number": 11,
        "type": "track",
        "uri": "spotify:track:5MIZW9FYY9aIbyKir7nYSB",
        "is_local": false
      },
      {
        "artists": [
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
        "duration_ms": 173226,
        "explicit": false,
        "external_urls": {
          "spotify": "https://open.spotify.com/track/1fhAGRPE6Q9p08knlqAIqa"
        },
        "href": "https://api.spotify.com/v1/tracks/1fhAGRPE6Q9p08knlqAIqa",
        "id": "1fhAGRPE6Q9p08knlqAIqa",
        "name": "ムカつく",
        "preview_url": null,
        "track_number": 12,
        "type": "track",
        "uri": "spotify:track:1fhAGRPE6Q9p08knlqAIqa",
        "is_local": false
      }
    ]
  },
  "copyrights": [
    {
      "text": "© 2025 MARUMOCHI LABEL",
      "type": "C"
    },
    {
      "text": "℗ 2025 MARUMOCHI LABEL",
      "type": "P"
    }
  ],
  "external_ids": {
    "upc": "4571640817465"
  },
  "genres": [],
  "label": "MARUMOCHI LABEL",
  "popularity": 18
}
```

API 仕様は https://developer.spotify.com/documentation/web-api/reference/get-an-album を参照してください
