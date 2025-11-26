# KKBOXAPI Get a Track API 仕様

## Request

ENDPOINT GET https://api.kkbox.com/v1.1/tracks/{id}?territory={territory}

id HamlOuN1E8K3Y617Ab

territory JP

## Response

```json
{
  "id": "HamlOuN1E8K3Y617Ab",
  "name": "感電",
  "duration": 264542,
  "isrc": "JPU902001053",
  "url": "https://www.kkbox.com/jp/ja/song/HamlOuN1E8K3Y617Ab",
  "track_number": 3,
  "explicitness": false,
  "available_territories": ["TW", "HK", "SG", "MY", "JP"],
  "album": {
    "id": "4kUU08ZXF4KmWQhM86",
    "name": "STRAY SHEEP",
    "url": "https://www.kkbox.com/jp/ja/album/4kUU08ZXF4KmWQhM86",
    "explicitness": false,
    "available_territories": ["TW", "HK", "SG", "MY", "JP"],
    "release_date": "2020-08-05",
    "images": [
      {
        "height": 160,
        "width": 160,
        "url": "https://i.kfs.io/album/global/83019209,2v1/fit/160x160.jpg"
      },
      {
        "height": 500,
        "width": 500,
        "url": "https://i.kfs.io/album/global/83019209,2v1/fit/500x500.jpg"
      },
      {
        "height": 1000,
        "width": 1000,
        "url": "https://i.kfs.io/album/global/83019209,2v1/fit/1000x1000.jpg"
      }
    ],
    "artist": {
      "id": "8oRs6ttmEEYzg2nyTU",
      "name": "米津玄師",
      "url": "https://www.kkbox.com/jp/ja/artist/8oRs6ttmEEYzg2nyTU",
      "images": [
        {
          "height": 160,
          "width": 160,
          "url": "https://i.kfs.io/artist/global/6653779,0v27/fit/160x160.jpg"
        },
        {
          "height": 300,
          "width": 300,
          "url": "https://i.kfs.io/artist/global/6653779,0v27/fit/300x300.jpg"
        }
      ]
    }
  }
}
```
