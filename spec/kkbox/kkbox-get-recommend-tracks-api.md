# KKBOX Get recommend tracks for a given track API 仕様

## Request

ENDPOINT GET https://api.kkbox.com/v1.1/tracks/{id}/recommended-tracks?territory={territory}&limit=50

id 8tPxi6LzsnPsEkgirY
territory JP

##### Response

```json
{
  "seed_track": {
    "id": "8tPxi6LzsnPsEkgirY",
    "name": "Wave",
    "duration": 221000,
    "isrc": "TCJPC2483393",
    "url": "https://www.kkbox.com/jp/ja/song/8tPxi6LzsnPsEkgirY",
    "track_number": 3,
    "explicitness": false,
    "available_territories": ["TW", "HK", "SG", "MY", "JP"]
  },
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
      },
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
      },
      {
        "id": "4rZeUztQSSJLr5ICaK",
        "name": "MIRROR",
        "duration": 178024,
        "isrc": "JPPO02402466",
        "url": "https://www.kkbox.com/jp/ja/song/4rZeUztQSSJLr5ICaK",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "Ha8_4luP_dvJJj3Fzg",
          "name": "MIRROR",
          "url": "https://www.kkbox.com/jp/ja/album/Ha8_4luP_dvJJj3Fzg",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2024-05-31",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/270253222,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/270253222,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/270253222,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "8pAtcVODL2RadSu8CT",
            "name": "Ado",
            "url": "https://www.kkbox.com/jp/ja/artist/8pAtcVODL2RadSu8CT",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/24169439,0v7/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/24169439,0v7/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "KtsZxJBFIXgP1fMzFa",
        "name": "浅草キッド",
        "duration": 244453,
        "isrc": "JPN001700243",
        "url": "https://www.kkbox.com/jp/ja/song/KtsZxJBFIXgP1fMzFa",
        "track_number": 10,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "5XU3zt6RTkffLtvSRt",
          "name": "PLAY(Special Edition)",
          "url": "https://www.kkbox.com/jp/ja/album/5XU3zt6RTkffLtvSRt",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2018-03-21",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/33038018,3v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/33038018,3v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/33038018,3v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "Kq5oXVJYQwYgk1wGZ4",
            "name": "Masaki Suda",
            "url": "https://www.kkbox.com/jp/ja/artist/Kq5oXVJYQwYgk1wGZ4",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/8466036,0v12/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/8466036,0v12/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "1X-0YeRF9YZDyfaoji",
        "name": "Stand by me, Stand by you.",
        "duration": 193018,
        "isrc": "JPB602002630",
        "url": "https://www.kkbox.com/jp/ja/song/1X-0YeRF9YZDyfaoji",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["JP"],
        "album": {
          "id": "8oFnH3LA3_AF4f78yX",
          "name": "Stand by me, Stand by you.",
          "url": "https://www.kkbox.com/jp/ja/album/8oFnH3LA3_AF4f78yX",
          "explicitness": false,
          "available_territories": ["JP"],
          "release_date": "2020-09-09",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/86634096,1v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/86634096,1v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/86634096,1v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "8mZW8InFFn4evpNRR8",
            "name": "平井 大",
            "url": "https://www.kkbox.com/jp/ja/artist/8mZW8InFFn4evpNRR8",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/631761,0v34/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/631761,0v34/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "8mhMzPGwJRsalPE-QJ",
        "name": "魔法にかけられて",
        "duration": 254955,
        "isrc": "JPJ222201954",
        "url": "https://www.kkbox.com/jp/ja/song/8mhMzPGwJRsalPE-QJ",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["MY"],
        "album": {
          "id": "GrDGZDyJSVBLbTkXHz",
          "name": "魔法にかけられて",
          "url": "https://www.kkbox.com/jp/ja/album/GrDGZDyJSVBLbTkXHz",
          "explicitness": false,
          "available_territories": ["MY"],
          "release_date": "2022-03-25",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/161478023,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/161478023,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/161478023,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "TYOh8oJpHjUVzwIkE5",
            "name": "Saucy Dog",
            "url": "https://www.kkbox.com/jp/ja/artist/TYOh8oJpHjUVzwIkE5",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/8319846,0v9/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/8319846,0v9/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "X_8dsPx-xGbE4t_WX8",
        "name": "RAIN",
        "duration": 306886,
        "isrc": "JPTF01706801",
        "url": "https://www.kkbox.com/jp/ja/song/X_8dsPx-xGbE4t_WX8",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["JP"],
        "album": {
          "id": "9-ri5ZhymwqbJObxu8",
          "name": "RAIN",
          "url": "https://www.kkbox.com/jp/ja/album/9-ri5ZhymwqbJObxu8",
          "explicitness": false,
          "available_territories": ["JP"],
          "release_date": "2017-07-05",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/26704367,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/26704367,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/26704367,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "Cqak7AAUuRkVZSrh-H",
            "name": "SEKAI NO OWARI",
            "url": "https://www.kkbox.com/jp/ja/artist/Cqak7AAUuRkVZSrh-H",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/186181,0v14/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/186181,0v14/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "On869Wm7o9axNB5WAI",
        "name": "鬼ノ宴",
        "duration": 176013,
        "isrc": "JPB602400319",
        "url": "https://www.kkbox.com/jp/ja/song/On869Wm7o9axNB5WAI",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["JP"],
        "album": {
          "id": "8rKc--TbLZoqMAPLMs",
          "name": "鬼ノ宴",
          "url": "https://www.kkbox.com/jp/ja/album/8rKc--TbLZoqMAPLMs",
          "explicitness": false,
          "available_territories": ["JP"],
          "release_date": "2024-01-10",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/265537066,1v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/265537066,1v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/265537066,1v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "SkR_uwHgT10-mHUgtJ",
            "name": "友成空",
            "url": "https://www.kkbox.com/jp/ja/artist/SkR_uwHgT10-mHUgtJ",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/27718904,0v5/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/27718904,0v5/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "8tBJWhMDsnPsHd4O3S",
        "name": "HEARTRIS",
        "duration": 180897,
        "isrc": "US5TA2300166",
        "url": "https://www.kkbox.com/jp/ja/song/8tBJWhMDsnPsHd4O3S",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["JP"],
        "album": {
          "id": "GlxEHsXgs-18BeHunx",
          "name": "Press Play",
          "url": "https://www.kkbox.com/jp/ja/album/GlxEHsXgs-18BeHunx",
          "explicitness": false,
          "available_territories": ["JP"],
          "release_date": "2023-10-30",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/262886051,1v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/262886051,1v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/262886051,1v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "GsoCZLHOdyFlInKI6c",
            "name": "NiziU",
            "url": "https://www.kkbox.com/jp/ja/artist/GsoCZLHOdyFlInKI6c",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/21757555,0v22/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/21757555,0v22/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "9aJvErBt1jmuK1VZbA",
        "name": "一目惚れ",
        "duration": 201952,
        "isrc": "GXD7G2432027",
        "url": "https://www.kkbox.com/jp/ja/song/9aJvErBt1jmuK1VZbA",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "CqWdwNERmpvC-9kRxj",
          "name": "一目惚れ",
          "url": "https://www.kkbox.com/jp/ja/album/CqWdwNERmpvC-9kRxj",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2024-05-14",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/269213346,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/269213346,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/269213346,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "_Y5M6UP5N1J1scqxjb",
            "name": "舟津真翔",
            "url": "https://www.kkbox.com/jp/ja/artist/_Y5M6UP5N1J1scqxjb",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/12940422,0v1/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/12940422,0v1/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "Gr6GbsEyk9ZkzFusj0",
        "name": "ツキヨミ",
        "duration": 215144,
        "isrc": "JPPO02204079",
        "url": "https://www.kkbox.com/jp/ja/song/Gr6GbsEyk9ZkzFusj0",
        "track_number": 13,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "_-_sifYhQDQcWx71d5",
          "name": "Mr.5 - Special Edition",
          "url": "https://www.kkbox.com/jp/ja/album/_-_sifYhQDQcWx71d5",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2023-04-19",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/269941075,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/269941075,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/269941075,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "Gl8SEY85xUb9_F63HQ",
            "name": "King & Prince",
            "url": "https://www.kkbox.com/jp/ja/artist/Gl8SEY85xUb9_F63HQ",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/55724271,0v8/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/55724271,0v8/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "LXxMqR2hU4Yl5ltTN7",
        "name": "Speaking",
        "duration": 229227,
        "isrc": "JPPO01508774",
        "url": "https://www.kkbox.com/jp/ja/song/LXxMqR2hU4Yl5ltTN7",
        "track_number": 2,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "OsYkDOA4qcVQbX489E",
          "name": "TWELVE",
          "url": "https://www.kkbox.com/jp/ja/album/OsYkDOA4qcVQbX489E",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2016-01-13",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/12801559,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/12801559,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/12801559,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "4naOIR-K2KLL59lD-Q",
            "name": "Mrs. GREEN APPLE",
            "url": "https://www.kkbox.com/jp/ja/artist/4naOIR-K2KLL59lD-Q",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/3222118,0v38/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/3222118,0v38/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "CsBPw65N3ZiY3pWDKz",
        "name": "ヨワネハキ",
        "duration": 166974,
        "isrc": "JPU902101067",
        "url": "https://www.kkbox.com/jp/ja/song/CsBPw65N3ZiY3pWDKz",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "Gpu-JRyDWop22TXh1f",
          "name": "ヨワネハキ",
          "url": "https://www.kkbox.com/jp/ja/album/Gpu-JRyDWop22TXh1f",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2021-05-19",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/116698105,1v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/116698105,1v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/116698105,1v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "Kp9mrlhVg8xEvmL3nF",
            "name": "MAISONdes",
            "url": "https://www.kkbox.com/jp/ja/artist/Kp9mrlhVg8xEvmL3nF",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/26895224,0v2/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/26895224,0v2/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "Cmv1wX91tBWGpjuVql",
        "name": "第六感",
        "duration": 191634,
        "isrc": "JPVI02001689",
        "url": "https://www.kkbox.com/jp/ja/song/Cmv1wX91tBWGpjuVql",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["JP"],
        "album": {
          "id": "DXBpBpNfYpMkhtaUCs",
          "name": "第六感",
          "url": "https://www.kkbox.com/jp/ja/album/DXBpBpNfYpMkhtaUCs",
          "explicitness": false,
          "available_territories": ["JP"],
          "release_date": "2020-07-27",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/81568015,1v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/81568015,1v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/81568015,1v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "4tyBsbjJ8jt7HTtChY",
            "name": "Reol",
            "url": "https://www.kkbox.com/jp/ja/artist/4tyBsbjJ8jt7HTtChY",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/3508304,0v26/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/3508304,0v26/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "0tBXU31gKUfrDlg5B2",
        "name": "ランデヴー",
        "duration": 237795,
        "isrc": "TCJPV2377398",
        "url": "https://www.kkbox.com/jp/ja/song/0tBXU31gKUfrDlg5B2",
        "track_number": 2,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "Kp1mHsqvaopkBepnJJ",
          "name": "誘拐 / ランデヴー",
          "url": "https://www.kkbox.com/jp/ja/album/Kp1mHsqvaopkBepnJJ",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2023-04-25",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/240922126,1v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/240922126,1v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/240922126,1v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "H_0WxuqX05ZsYX9Tbh",
            "name": "シャイトープ",
            "url": "https://www.kkbox.com/jp/ja/artist/H_0WxuqX05ZsYX9Tbh",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/40826062,0v4/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/40826062,0v4/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "4n50C4no3TYuq1fPX1",
        "name": "いのちの名前",
        "duration": 269871,
        "isrc": "JPP302300593",
        "url": "https://www.kkbox.com/jp/ja/song/4n50C4no3TYuq1fPX1",
        "track_number": 3,
        "explicitness": false,
        "available_territories": ["JP"],
        "album": {
          "id": "DZykuDrP1aDjG3FQ_8",
          "name": "スタジオジブリ トリビュートアルバム「ジブリをうたう」",
          "url": "https://www.kkbox.com/jp/ja/album/DZykuDrP1aDjG3FQ_8",
          "explicitness": false,
          "available_territories": ["JP"],
          "release_date": "2023-11-01",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/262829833,2v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/262829833,2v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/262829833,2v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "LY5scHimqDphGqy41I",
            "name": "Various Artists",
            "url": "https://www.kkbox.com/jp/ja/artist/LY5scHimqDphGqy41I",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/21023,0v2755/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/21023,0v2755/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "CqXGhfOyUh-79-E7CJ",
        "name": "椿",
        "duration": 327471,
        "isrc": "JPB602070849",
        "url": "https://www.kkbox.com/jp/ja/song/CqXGhfOyUh-79-E7CJ",
        "track_number": 1,
        "explicitness": true,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "Kk31XEVgLfA_PvRaHY",
          "name": "椿",
          "url": "https://www.kkbox.com/jp/ja/album/Kk31XEVgLfA_PvRaHY",
          "explicitness": true,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2020-07-10",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/126235790,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/126235790,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/126235790,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "Cpy2-7tnOiWI6FTuHT",
            "name": "Torauma",
            "url": "https://www.kkbox.com/jp/ja/artist/Cpy2-7tnOiWI6FTuHT",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/571878,0v1/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/571878,0v1/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "H_Bhre0o6JmyxVoq_p",
        "name": "おもかげ -self cover-",
        "duration": 209946,
        "isrc": "JPP302101989",
        "url": "https://www.kkbox.com/jp/ja/song/H_Bhre0o6JmyxVoq_p",
        "track_number": 4,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "CpzqbJUkaw7YBHa9nP",
          "name": "裸の勇者",
          "url": "https://www.kkbox.com/jp/ja/album/CpzqbJUkaw7YBHa9nP",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2022-02-21",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/148725316,4v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/148725316,4v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/148725316,4v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "WnJzCp5yTvH9-d9veg",
            "name": "Vaundy",
            "url": "https://www.kkbox.com/jp/ja/artist/WnJzCp5yTvH9-d9veg",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/18646961,0v7/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/18646961,0v7/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "Ck9ZKh8KnNz69isyoz",
        "name": "幾億光年",
        "duration": 276967,
        "isrc": "JPU902305217",
        "url": "https://www.kkbox.com/jp/ja/song/Ck9ZKh8KnNz69isyoz",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "8tWx_lOMZgqWfseZQ1",
          "name": "幾億光年",
          "url": "https://www.kkbox.com/jp/ja/album/8tWx_lOMZgqWfseZQ1",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2024-01-24",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/265855174,2v2/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/265855174,2v2/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/265855174,2v2/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "GqA6ZnigyaDAnsugN1",
            "name": "Omoinotake",
            "url": "https://www.kkbox.com/jp/ja/artist/GqA6ZnigyaDAnsugN1",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/7229128,0v16/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/7229128,0v16/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "GtjULKHO-2a5-PtwVq",
        "name": "水平線",
        "duration": 285413,
        "isrc": "JPPO02100907",
        "url": "https://www.kkbox.com/jp/ja/song/GtjULKHO-2a5-PtwVq",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "LZd3jche2Z7b1wkUzM",
          "name": "水平線",
          "url": "https://www.kkbox.com/jp/ja/album/LZd3jche2Z7b1wkUzM",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2021-08-13",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/131129721,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/131129721,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/131129721,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "GtX1PP-Fw6HSAe1lb2",
            "name": "back number",
            "url": "https://www.kkbox.com/jp/ja/artist/GtX1PP-Fw6HSAe1lb2",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/752759,0v17/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/752759,0v17/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "0p8LqPM1MinpW5BBAt",
        "name": "別の人の彼女になったよ",
        "duration": 303020,
        "isrc": "JPU901802225",
        "url": "https://www.kkbox.com/jp/ja/song/0p8LqPM1MinpW5BBAt",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "4ptNwIvc3PHEZ_lPm9",
          "name": "別の人の彼女になったよ",
          "url": "https://www.kkbox.com/jp/ja/album/4ptNwIvc3PHEZ_lPm9",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2018-08-22",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/38160125,3v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/38160125,3v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/38160125,3v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "CnAmawi5Ft4y5GbB3n",
            "name": "wacci",
            "url": "https://www.kkbox.com/jp/ja/artist/CnAmawi5Ft4y5GbB3n",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/638112,0v10/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/638112,0v10/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "HXOfPBLoaXJ8aRrSpm",
        "name": "おかえり",
        "duration": 225854,
        "isrc": "TCJPT2254648",
        "url": "https://www.kkbox.com/jp/ja/song/HXOfPBLoaXJ8aRrSpm",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "9aAuBF0zUhQp-MxKat",
          "name": "おかえり",
          "url": "https://www.kkbox.com/jp/ja/album/9aAuBF0zUhQp-MxKat",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2022-11-09",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/201257819,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/201257819,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/201257819,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "0qujmjV7Pu5uv4TSJo",
            "name": "Tani Yuuki",
            "url": "https://www.kkbox.com/jp/ja/artist/0qujmjV7Pu5uv4TSJo",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/21989644,0v4/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/21989644,0v4/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "8p8WBNLr2c2emMG_9O",
        "name": "今はいいんだよ。 (feat. 可不)",
        "duration": 147095,
        "isrc": "TCJPU2217956",
        "url": "https://www.kkbox.com/jp/ja/song/8p8WBNLr2c2emMG_9O",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "DZcWiVr_1aDjFDJ2Kp",
          "name": "今はいいんだよ。 (feat. 可不)",
          "url": "https://www.kkbox.com/jp/ja/album/DZcWiVr_1aDjFDJ2Kp",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2022-12-23",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/213669042,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/213669042,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/213669042,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "-r5ZvqzjktQaXdHx-m",
            "name": "MIMI",
            "url": "https://www.kkbox.com/jp/ja/artist/-r5ZvqzjktQaXdHx-m",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/18578816,0v1/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/18578816,0v1/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "Kop7eJSBP2VFvMH7id",
        "name": "おもかげ (produced by Vaundy)",
        "duration": 188081,
        "isrc": "JPU902104741",
        "url": "https://www.kkbox.com/jp/ja/song/Kop7eJSBP2VFvMH7id",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "JP"],
        "album": {
          "id": "0mCdmm57qw4Imt-r0M",
          "name": "おもかげ (produced by Vaundy)",
          "url": "https://www.kkbox.com/jp/ja/album/0mCdmm57qw4Imt-r0M",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "JP"],
          "release_date": "2021-12-17",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/146853209,2v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/146853209,2v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/146853209,2v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "Ck9wH3DdLnwAfyE1zv",
            "name": "milet, Aimer, 幾田りら",
            "url": "https://www.kkbox.com/jp/ja/artist/Ck9wH3DdLnwAfyE1zv",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/34664978,0v2/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/34664978,0v2/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "DYNoXOsbri5YTW07YC",
        "name": "あなたに",
        "duration": 208796,
        "isrc": "JPPO02303151",
        "url": "https://www.kkbox.com/jp/ja/song/DYNoXOsbri5YTW07YC",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "WtSNzcDfYOSQlUdwUj",
          "name": "あなたに",
          "url": "https://www.kkbox.com/jp/ja/album/WtSNzcDfYOSQlUdwUj",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2023-10-25",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/262801224,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/262801224,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/262801224,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "4naOIR-K2KLL59lD-Q",
            "name": "Mrs. GREEN APPLE",
            "url": "https://www.kkbox.com/jp/ja/artist/4naOIR-K2KLL59lD-Q",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/3222118,0v38/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/3222118,0v38/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "PY3G1goiNHxJu8r1l3",
        "name": "us",
        "duration": 254862,
        "isrc": "JPU901901783",
        "url": "https://www.kkbox.com/jp/ja/song/PY3G1goiNHxJu8r1l3",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "4m4vGfF02kVG50L2RK",
          "name": "us",
          "url": "https://www.kkbox.com/jp/ja/album/4m4vGfF02kVG50L2RK",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2019-08-19",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/59723864,13v2/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/59723864,13v2/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/59723864,13v2/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "GlHwvl8JxUb98BqzNI",
            "name": "milet",
            "url": "https://www.kkbox.com/jp/ja/artist/GlHwvl8JxUb98BqzNI",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/11493889,0v24/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/11493889,0v24/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "WrBk3jgnHVANKP5BcT",
        "name": "唱",
        "duration": 189779,
        "isrc": "JPPO02302806",
        "url": "https://www.kkbox.com/jp/ja/song/WrBk3jgnHVANKP5BcT",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "-rtg_SYp1XBjLi5YP1",
          "name": "唱",
          "url": "https://www.kkbox.com/jp/ja/album/-rtg_SYp1XBjLi5YP1",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2023-09-06",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/261298616,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/261298616,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/261298616,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "8pAtcVODL2RadSu8CT",
            "name": "Ado",
            "url": "https://www.kkbox.com/jp/ja/artist/8pAtcVODL2RadSu8CT",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/24169439,0v7/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/24169439,0v7/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "KonqrqfRP2VFsQnZlX",
        "name": "パブリック",
        "duration": 222818,
        "isrc": "JPPO01510085",
        "url": "https://www.kkbox.com/jp/ja/song/KonqrqfRP2VFsQnZlX",
        "track_number": 3,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "OsYkDOA4qcVQbX489E",
          "name": "TWELVE",
          "url": "https://www.kkbox.com/jp/ja/album/OsYkDOA4qcVQbX489E",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2016-01-13",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/12801559,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/12801559,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/12801559,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "4naOIR-K2KLL59lD-Q",
            "name": "Mrs. GREEN APPLE",
            "url": "https://www.kkbox.com/jp/ja/artist/4naOIR-K2KLL59lD-Q",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/3222118,0v38/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/3222118,0v38/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "SskoPNJfDXYxbPznI1",
        "name": "首のない天使",
        "duration": 218174,
        "isrc": "JPU902401289",
        "url": "https://www.kkbox.com/jp/ja/song/SskoPNJfDXYxbPznI1",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "JP"],
        "album": {
          "id": "9ZWFu6kn_BzuNxGjWE",
          "name": "首のない天使",
          "url": "https://www.kkbox.com/jp/ja/album/9ZWFu6kn_BzuNxGjWE",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "JP"],
          "release_date": "2024-05-01",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/268878047,3v2/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/268878047,3v2/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/268878047,3v2/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "0nnL1geN5hLUKmQaKa",
            "name": "Queen Bee",
            "url": "https://www.kkbox.com/jp/ja/artist/0nnL1geN5hLUKmQaKa",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/621582,0v13/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/621582,0v13/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "1aTeey-X0t4rx2EO5x",
        "name": "エジソン",
        "duration": 193515,
        "isrc": "JPWP02200075",
        "url": "https://www.kkbox.com/jp/ja/song/1aTeey-X0t4rx2EO5x",
        "track_number": 3,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "CoTMsw4-1JHry5bs8T",
          "name": "ネオン",
          "url": "https://www.kkbox.com/jp/ja/album/CoTMsw4-1JHry5bs8T",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2022-05-25",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/172038581,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/172038581,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/172038581,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "__JSYECtv1wl2cCn7s",
            "name": "水曜日のカンパネラ",
            "url": "https://www.kkbox.com/jp/ja/artist/__JSYECtv1wl2cCn7s",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/5045691,0v15/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/5045691,0v15/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "CoN2AofNHOVZgXAF1T",
        "name": "君に贈る歌 ～Song For You",
        "duration": 341611,
        "isrc": "JPPO01500206",
        "url": "https://www.kkbox.com/jp/ja/song/CoN2AofNHOVZgXAF1T",
        "track_number": 3,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "0sMkayr_ljsme-KdeB",
          "name": "シェネル・ワールド",
          "url": "https://www.kkbox.com/jp/ja/album/0sMkayr_ljsme-KdeB",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2015-02-11",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/5771061,0v2/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/5771061,0v2/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/5771061,0v2/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "CrG5e6x2Y7UDDziYLl",
            "name": "シェネル",
            "url": "https://www.kkbox.com/jp/ja/artist/CrG5e6x2Y7UDDziYLl",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/28134,0v9/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/28134,0v9/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "0lA5MA53s1XUeM8ryL",
        "name": "向日葵",
        "duration": 259645,
        "isrc": "JPPO02302087",
        "url": "https://www.kkbox.com/jp/ja/song/0lA5MA53s1XUeM8ryL",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "8qHLrislOgKcsgaWz4",
          "name": "向日葵",
          "url": "https://www.kkbox.com/jp/ja/album/8qHLrislOgKcsgaWz4",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2023-07-11",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/255966138,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/255966138,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/255966138,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "8pAtcVODL2RadSu8CT",
            "name": "Ado",
            "url": "https://www.kkbox.com/jp/ja/artist/8pAtcVODL2RadSu8CT",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/24169439,0v7/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/24169439,0v7/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "4tZnZNdYLX0lFKh38-",
        "name": "いつか",
        "duration": 274886,
        "isrc": "JPZ921704027",
        "url": "https://www.kkbox.com/jp/ja/song/4tZnZNdYLX0lFKh38-",
        "track_number": 3,
        "explicitness": false,
        "available_territories": ["MY"],
        "album": {
          "id": "HarrW_qf_dvJJNGYUN",
          "name": "カントリーロード",
          "url": "https://www.kkbox.com/jp/ja/album/HarrW_qf_dvJJNGYUN",
          "explicitness": false,
          "available_territories": ["MY"],
          "release_date": "2017-05-24",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/25294843,1v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/25294843,1v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/25294843,1v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "TYOh8oJpHjUVzwIkE5",
            "name": "Saucy Dog",
            "url": "https://www.kkbox.com/jp/ja/artist/TYOh8oJpHjUVzwIkE5",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/8319846,0v9/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/8319846,0v9/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "Go7SziQgV7_OdVQbwn",
        "name": "BANG BANG BANG",
        "duration": 222145,
        "isrc": "JPB601502279",
        "url": "https://www.kkbox.com/jp/ja/song/Go7SziQgV7_OdVQbwn",
        "track_number": 2,
        "explicitness": false,
        "available_territories": ["JP"],
        "album": {
          "id": "4sfetQsyxHs9Dr1LAv",
          "name": "MADE SERIES",
          "url": "https://www.kkbox.com/jp/ja/album/4sfetQsyxHs9Dr1LAv",
          "explicitness": false,
          "available_territories": ["JP"],
          "release_date": "2016-02-03",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/13241904,1v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/13241904,1v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/13241904,1v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "_Z047IW0OU_ZyemxdA",
            "name": "BIGBANG",
            "url": "https://www.kkbox.com/jp/ja/artist/_Z047IW0OU_ZyemxdA",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/43934,0v7/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/43934,0v7/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "SsyjxwJPDXYxYX6WyF",
        "name": "Moshi Moshi (feat. 百足)",
        "duration": 176953,
        "isrc": "TCJPY2493758",
        "url": "https://www.kkbox.com/jp/ja/song/SsyjxwJPDXYxYX6WyF",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "_Z_TFV1VKVWPvt_5l0",
          "name": "Moshi Moshi (feat. 百足)",
          "url": "https://www.kkbox.com/jp/ja/album/_Z_TFV1VKVWPvt_5l0",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2024-01-24",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/266208261,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/266208261,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/266208261,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "L-fwbNZyGWw2pNVtFD",
            "name": "Nozomi Kitay, GAL D",
            "url": "https://www.kkbox.com/jp/ja/artist/L-fwbNZyGWw2pNVtFD",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/55118171,0v1/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/55118171,0v1/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "OpkCqscQzUmSYUJ8OI",
        "name": "SWEET NONFICTION",
        "duration": 199706,
        "isrc": "JPU902400070",
        "url": "https://www.kkbox.com/jp/ja/song/OpkCqscQzUmSYUJ8OI",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "JP"],
        "album": {
          "id": "Wkp8kMCTIvVnkqy-iZ",
          "name": "SWEET NONFICTION",
          "url": "https://www.kkbox.com/jp/ja/album/Wkp8kMCTIvVnkqy-iZ",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "JP"],
          "release_date": "2024-03-14",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/267119372,4v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/267119372,4v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/267119372,4v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "GsoCZLHOdyFlInKI6c",
            "name": "NiziU",
            "url": "https://www.kkbox.com/jp/ja/artist/GsoCZLHOdyFlInKI6c",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/21757555,0v22/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/21757555,0v22/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "CsNJX8793ZiY045gvg",
        "name": "前前前世 - movie ver.",
        "duration": 285930,
        "isrc": "JPPO01617097",
        "url": "https://www.kkbox.com/jp/ja/song/CsNJX8793ZiY045gvg",
        "track_number": 8,
        "explicitness": false,
        "available_territories": ["JP"],
        "album": {
          "id": "9ZG3TJhn_BzuPCGpGk",
          "name": "君の名は。",
          "url": "https://www.kkbox.com/jp/ja/album/9ZG3TJhn_BzuPCGpGk",
          "explicitness": false,
          "available_territories": ["JP"],
          "release_date": "2016-08-24",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/74837916,1v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/74837916,1v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/74837916,1v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "Csgubao952obU7R3NE",
            "name": "RADWIMPS",
            "url": "https://www.kkbox.com/jp/ja/artist/Csgubao952obU7R3NE",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/22036,0v15/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/22036,0v15/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "Co3ZWnXdHOVZjJ00IQ",
        "name": "まちがいさがし",
        "duration": 267467,
        "isrc": "JPU902001657",
        "url": "https://www.kkbox.com/jp/ja/song/Co3ZWnXdHOVZjJ00IQ",
        "track_number": 9,
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
      },
      {
        "id": "Os0YCxmE8UtsykNejP",
        "name": "Tot Musica",
        "duration": 195239,
        "isrc": "JPPO02202654",
        "url": "https://www.kkbox.com/jp/ja/song/Os0YCxmE8UtsykNejP",
        "track_number": 5,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "5Yhcqfr3sLpPwzHP2L",
          "name": "ウタの歌 ONE PIECE FILM RED",
          "url": "https://www.kkbox.com/jp/ja/album/5Yhcqfr3sLpPwzHP2L",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2022-08-10",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/186412914,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/186412914,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/186412914,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "8pAtcVODL2RadSu8CT",
            "name": "Ado",
            "url": "https://www.kkbox.com/jp/ja/artist/8pAtcVODL2RadSu8CT",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/24169439,0v7/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/24169439,0v7/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "8rHUjGKE59bDSdoLBt",
        "name": "Masterplan",
        "duration": 212506,
        "isrc": "JPB602401000",
        "url": "https://www.kkbox.com/jp/ja/song/8rHUjGKE59bDSdoLBt",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["JP"],
        "album": {
          "id": "5XPIMW-BTkffIiX4AR",
          "name": "Masterplan",
          "url": "https://www.kkbox.com/jp/ja/album/5XPIMW-BTkffIiX4AR",
          "explicitness": false,
          "available_territories": ["JP"],
          "release_date": "2024-04-24",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/268689837,3v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/268689837,3v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/268689837,3v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "XXa9_Ho6vXP53Ojmw3",
            "name": "BE:FIRST",
            "url": "https://www.kkbox.com/jp/ja/artist/XXa9_Ho6vXP53Ojmw3",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/31762760,0v24/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/31762760,0v24/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "LZUJG4GsMFHagKgxYI",
        "name": "花詩",
        "duration": 209652,
        "isrc": "JPU902401232",
        "url": "https://www.kkbox.com/jp/ja/song/LZUJG4GsMFHagKgxYI",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "GsMGEU_zG1L3H3MvSh",
          "name": "HOMETOWN",
          "url": "https://www.kkbox.com/jp/ja/album/GsMGEU_zG1L3H3MvSh",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2024-05-08",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/268649348,2v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/268649348,2v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/268649348,2v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "0qujmjV7Pu5uv4TSJo",
            "name": "Tani Yuuki",
            "url": "https://www.kkbox.com/jp/ja/artist/0qujmjV7Pu5uv4TSJo",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/21989644,0v4/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/21989644,0v4/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "8pZTK5Tr2c2em0f6qn",
        "name": "Yessir (feat. Eric.B.Jr.)",
        "duration": 172408,
        "isrc": "TCJPM2024795",
        "url": "https://www.kkbox.com/jp/ja/song/8pZTK5Tr2c2em0f6qn",
        "track_number": 4,
        "explicitness": true,
        "available_territories": [],
        "album": {
          "id": "DZukMtpv1aDjFQdDlI",
          "name": "Jungle",
          "url": "https://www.kkbox.com/jp/ja/album/DZukMtpv1aDjFQdDlI",
          "explicitness": true,
          "available_territories": [],
          "release_date": "2020-08-21",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/84180857,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/84180857,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/84180857,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "0nHmmbeN5hLULQdwgc",
            "name": "¥ellow Bucks",
            "url": "https://www.kkbox.com/jp/ja/artist/0nHmmbeN5hLULQdwgc",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/15902860,0v3/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/15902860,0v3/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "D-IFdxDmldMilf3dCK",
        "name": "トゲめくスピカ",
        "duration": 279536,
        "isrc": "JPPO01904517",
        "url": "https://www.kkbox.com/jp/ja/song/D-IFdxDmldMilf3dCK",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "9XX5dhlXC0_tn1vj08",
          "name": "トゲめくスピカ",
          "url": "https://www.kkbox.com/jp/ja/album/9XX5dhlXC0_tn1vj08",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2019-12-20",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/67249971,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/67249971,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/67249971,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "Wq_UNa6bMxpS3m0j5E",
            "name": "ポルカドットスティングレイ",
            "url": "https://www.kkbox.com/jp/ja/artist/Wq_UNa6bMxpS3m0j5E",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/6885503,0v18/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/6885503,0v18/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "WnaSCGfzwvj1xaaXaF",
        "name": "You",
        "duration": 222354,
        "isrc": "TCJPZ2434874",
        "url": "https://www.kkbox.com/jp/ja/song/WnaSCGfzwvj1xaaXaF",
        "track_number": 6,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "TY5biZz_OR7QjDTEWR",
          "name": "MAGIC HOUR",
          "url": "https://www.kkbox.com/jp/ja/album/TY5biZz_OR7QjDTEWR",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2024-03-20",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/267481213,1v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/267481213,1v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/267481213,1v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "0suzt2UBC5Tq4NxPdU",
            "name": "808",
            "url": "https://www.kkbox.com/jp/ja/artist/0suzt2UBC5Tq4NxPdU",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/55238276,0v1/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/55238276,0v1/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "5_2GqErzMf04hJacCq",
        "name": "Supernatural",
        "duration": 191007,
        "isrc": "USA2P2417518",
        "url": "https://www.kkbox.com/jp/ja/song/5_2GqErzMf04hJacCq",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["JP"],
        "album": {
          "id": "SmTT8-qygNqlRceLXU",
          "name": "Supernatural",
          "url": "https://www.kkbox.com/jp/ja/album/SmTT8-qygNqlRceLXU",
          "explicitness": false,
          "available_territories": ["JP"],
          "release_date": "2024-06-21",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/271105636,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/271105636,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/271105636,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "GrJOTzYkMQFlP-jK7k",
            "name": "NewJeans",
            "url": "https://www.kkbox.com/jp/ja/artist/GrJOTzYkMQFlP-jK7k",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/41595877,0v7/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/41595877,0v7/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "5-CijyTei4L5pzzHAl",
        "name": "ショコラカタブラ",
        "duration": 183925,
        "isrc": "JPPO02304564",
        "url": "https://www.kkbox.com/jp/ja/song/5-CijyTei4L5pzzHAl",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "GrMN8DzpSVBLaztKXL",
          "name": "ショコラカタブラ",
          "url": "https://www.kkbox.com/jp/ja/album/GrMN8DzpSVBLaztKXL",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2024-01-31",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/266292411,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/266292411,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/266292411,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "8pAtcVODL2RadSu8CT",
            "name": "Ado",
            "url": "https://www.kkbox.com/jp/ja/artist/8pAtcVODL2RadSu8CT",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/24169439,0v7/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/24169439,0v7/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "Gs8L2JC0Z-QUaaHPqj",
        "name": "napori",
        "duration": 203807,
        "isrc": "JPR652000114",
        "url": "https://www.kkbox.com/jp/ja/song/Gs8L2JC0Z-QUaaHPqj",
        "track_number": 9,
        "explicitness": false,
        "available_territories": ["JP"],
        "album": {
          "id": "TaoPuxOD9CYeFPjeOk",
          "name": "strobo",
          "url": "https://www.kkbox.com/jp/ja/album/TaoPuxOD9CYeFPjeOk",
          "explicitness": false,
          "available_territories": ["JP"],
          "release_date": "2020-05-27",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/262047503,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/262047503,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/262047503,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "WnJzCp5yTvH9-d9veg",
            "name": "Vaundy",
            "url": "https://www.kkbox.com/jp/ja/artist/WnJzCp5yTvH9-d9veg",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/18646961,0v7/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/18646961,0v7/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "Onw8NV7Lo9axMffLRO",
        "name": "足りない",
        "duration": 180297,
        "isrc": "JPB602104577",
        "url": "https://www.kkbox.com/jp/ja/song/Onw8NV7Lo9axMffLRO",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["JP"],
        "album": {
          "id": "CnXtWRT14d6Z3JjPVy",
          "name": "足りない",
          "url": "https://www.kkbox.com/jp/ja/album/CnXtWRT14d6Z3JjPVy",
          "explicitness": false,
          "available_territories": ["JP"],
          "release_date": "2022-06-01",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/175087103,1v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/175087103,1v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/175087103,1v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "5_sGir5eN1tMAftjO-",
            "name": "DUSTCELL",
            "url": "https://www.kkbox.com/jp/ja/artist/5_sGir5eN1tMAftjO-",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/20293835,0v3/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/20293835,0v3/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "T--hlQiu-xXCViMF2q",
        "name": "バニラ",
        "duration": 251100,
        "isrc": "TCJPR2220366",
        "url": "https://www.kkbox.com/jp/ja/song/T--hlQiu-xXCViMF2q",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "PZEpxjPo_n53VYbbOG",
          "name": "バニラ",
          "url": "https://www.kkbox.com/jp/ja/album/PZEpxjPo_n53VYbbOG",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2022-03-09",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/159231515,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/159231515,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/159231515,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "CmSX5cnD1F5sMnZLvK",
            "name": "きゃない",
            "url": "https://www.kkbox.com/jp/ja/artist/CmSX5cnD1F5sMnZLvK",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/21939929,0v2/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/21939929,0v2/fit/300x300.jpg"
              }
            ]
          }
        }
      },
      {
        "id": "H-kHM8raCJ6qyCMo17",
        "name": "Talking Box (Dirty Pop Remix)",
        "duration": 128992,
        "isrc": "JPB452205940",
        "url": "https://www.kkbox.com/jp/ja/song/H-kHM8raCJ6qyCMo17",
        "track_number": 1,
        "explicitness": false,
        "available_territories": ["TW", "HK", "SG", "MY", "JP"],
        "album": {
          "id": "KqyQIWxaOnJP4LZoE3",
          "name": "Talking Box (Dirty Pop Remix) / SPACESHIP",
          "url": "https://www.kkbox.com/jp/ja/album/KqyQIWxaOnJP4LZoE3",
          "explicitness": false,
          "available_territories": ["TW", "HK", "SG", "MY", "JP"],
          "release_date": "2022-04-01",
          "images": [
            {
              "height": 160,
              "width": 160,
              "url": "https://i.kfs.io/album/global/160266204,0v1/fit/160x160.jpg"
            },
            {
              "height": 500,
              "width": 500,
              "url": "https://i.kfs.io/album/global/160266204,0v1/fit/500x500.jpg"
            },
            {
              "height": 1000,
              "width": 1000,
              "url": "https://i.kfs.io/album/global/160266204,0v1/fit/1000x1000.jpg"
            }
          ],
          "artist": {
            "id": "GookN1Y8PJM_uxgccX",
            "name": "WurtS",
            "url": "https://www.kkbox.com/jp/ja/artist/GookN1Y8PJM_uxgccX",
            "images": [
              {
                "height": 160,
                "width": 160,
                "url": "https://i.kfs.io/artist/global/25059618,0v6/fit/160x160.jpg"
              },
              {
                "height": 300,
                "width": 300,
                "url": "https://i.kfs.io/artist/global/25059618,0v6/fit/300x300.jpg"
              }
            ]
          }
        }
      }
    ],
    "paging": {
      "offset": 0,
      "limit": 50,
      "previous": null,
      "next": null
    },
    "summary": {
      "total": 50
    }
  }
}
```

API 仕様は　https://docs-zhtw.kkbox.codes/#get-/me/recommended-seed-tracks/{track_id}　を参照してください
