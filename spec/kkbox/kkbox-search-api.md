# KKBOXAPI Search various objects API 仕様

## Request

ENDPOINT GET https://api.kkbox.com/v1.1/search?q={query}&type=track%2Calbum%2Cartist%2Cplaylist&territory=JP&limit=15

query isrc:TCJPC2483393 (upc の場合は upc:{upc})

## Response

```json
{
  "tracks": {
    "data": [
      {
        "id": "8tPxi6LzsnPsEkgirY",
        "name": "Wave",
        "duration": 221000,
        "isrc": "TCJPC2483393",
        "url": "https://www.kkbox.com/jp/ja/song/8tPxi6LzsnPsEkgirY",
        "track_number": 3,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "-lRRV7YGp_79zz3WjE",
          "name": "箱庭共鳴-ハコニワレゾナンス-Hanon×Kotoha 歌唱版",
          "url": "https://www.kkbox.com/jp/ja/album/-lRRV7YGp_79zz3WjE",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2025-01-29",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/278308256,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/278308256,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/278308256,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "PYMxNnImScLQmroQi3",
            "name": "Hanon, Kotoha",
            "url": "https://www.kkbox.com/jp/ja/artist/PYMxNnImScLQmroQi3",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/30463871,0v1/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/30463871,0v1/fit/300x300.jpg"
              }
            ]
          }
        }
      }
    ],
    "paging": {
      "offset": 0,
      "limit": 15,
      "previous": null,
      "next": null
    },
    "summary": {
      "total": 1
    }
  },
  "summary": {
    "total": 1
  },
  "paging": {
    "offset": 0,
    "limit": 15,
    "previous": null,
    "next": null
  }
}
```

API 仕様は https://docs-zhtw.kkbox.codes/#get-/search を参照
