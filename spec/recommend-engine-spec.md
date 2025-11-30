# レコメンドエンジン仕様書

## 概要

Spotify Audio Features を活用したハイブリッドレコメンドエンジンの仕様書です。

現行の KKBOX レコメンドに加え、楽曲の音響特徴量（テンポ、エネルギー、明るさなど）に基づく類似度計算を行い、「関連性がある」かつ「雰囲気が似ている」楽曲を提案します。

## 現状の課題

| 課題             | 詳細                                                   |
| ---------------- | ------------------------------------------------------ |
| 関連性のみ       | KKBOX レコメンドはアーティスト・ジャンルの関連性が中心 |
| 雰囲気の考慮不足 | テンポや曲調の類似性が反映されにくい                   |
| 結果の偏り       | 同一アーティストや同一アルバムの曲が多くなりがち       |
| ジャンル逸脱     | ニッチなジャンル（アニソン等）から J-POP/K-POP に偏る  |

### ジャンル逸脱の具体例

```
入力: アニソン「電脳スペクタクル」（ワールドダイスター）

KKBOXレコメンド結果:
├── YOASOBI「夜に駆ける」     → J-POP
├── 米津玄師「Lemon」         → J-POP
├── HoneyWorks「可愛くてごめん」→ J-POP
├── Ado「唱」                 → J-POP
└── なにわ男子「初心LOVE」     → アイドル

→ 同じ「アニソン」ジャンルの曲が1曲も出てこない
```

KKBOX は「日本で人気の曲」を関連として出す傾向があり、アニソン・ボカロ・同人音楽などのニッチなジャンルは無視されがち。

## 解決アプローチ

### ハイブリッドレコメンド

1. **関連性ベース**: 既存の KKBOX/Spotify レコメンドで候補を取得
2. **雰囲気ベース**: Audio Features による類似度でランキング

## API エンドポイント

### 新規エンドポイント

```
GET /v1/track/recommend?url={spotify_url}&mode={mode}
```

### パラメータ

| パラメータ | 必須 | デフォルト | 説明                 |
| ---------- | ---- | ---------- | -------------------- |
| `url`      | ○    | -          | Spotify トラック URL |
| `mode`     | -    | `balanced` | レコメンドモード     |
| `limit`    | -    | `20`       | 返却件数（最大 30）  |

### レコメンドモード

| モード     | 説明       | 重み付け                         |
| ---------- | ---------- | -------------------------------- |
| `similar`  | 雰囲気重視 | Audio Features 類似度を最優先    |
| `related`  | 関連性重視 | アーティスト・ジャンル関連を優先 |
| `balanced` | バランス   | 両方を均等に考慮（デフォルト）   |

### レスポンス

```json
{
  "status": 200,
  "result": {
    "seed_track": {
      "id": "xxx",
      "name": "入力楽曲名",
      "artists": [{"name": "アーティスト名"}],
      "audio_features": {
        "tempo": 128.0,
        "energy": 0.85,
        "danceability": 0.72,
        "valence": 0.65
      }
    },
    "items": [
      {
        "id": "yyy",
        "name": "レコメンド楽曲名",
        "artists": [{"name": "アーティスト名"}],
        "album": {...},
        "url": "https://open.spotify.com/track/yyy",
        "similarity_score": 0.92,
        "audio_features": {
          "tempo": 126.0,
          "energy": 0.82,
          "danceability": 0.70,
          "valence": 0.68
        },
        "match_reasons": ["tempo", "energy", "same_genre"]
      }
    ],
    "mode": "balanced"
  }
}
```

## Spotify Audio Features

### 取得エンドポイント

```
GET https://api.spotify.com/v1/audio-features/{id}
GET https://api.spotify.com/v1/audio-features?ids={ids}  // バッチ取得（最大100件）
```

### 特徴量一覧

| 特徴量             | 型    | 範囲     | 説明                         |
| ------------------ | ----- | -------- | ---------------------------- |
| `tempo`            | float | 0-250    | BPM（テンポ）                |
| `energy`           | float | 0.0-1.0  | エネルギッシュさ、激しさ     |
| `danceability`     | float | 0.0-1.0  | 踊りやすさ                   |
| `valence`          | float | 0.0-1.0  | 明るさ、ポジティブさ         |
| `acousticness`     | float | 0.0-1.0  | アコースティック感           |
| `instrumentalness` | float | 0.0-1.0  | ボーカルの少なさ             |
| `speechiness`      | float | 0.0-1.0  | 話し言葉の多さ               |
| `liveness`         | float | 0.0-1.0  | ライブ感                     |
| `loudness`         | float | -60 to 0 | 平均音量（dB）               |
| `key`              | int   | 0-11     | 音楽のキー（0=C, 1=C#, ...） |
| `mode`             | int   | 0 or 1   | メジャー(1) / マイナー(0)    |
| `time_signature`   | int   | 3-7      | 拍子（4 = 4/4 拍子）         |

### レスポンス例

```json
{
  "id": "4uLU6hMCjMI75M1A2tKUQC",
  "tempo": 128.04,
  "energy": 0.854,
  "danceability": 0.721,
  "valence": 0.652,
  "acousticness": 0.0412,
  "instrumentalness": 0.00015,
  "speechiness": 0.0534,
  "liveness": 0.124,
  "loudness": -5.234,
  "key": 7,
  "mode": 1,
  "time_signature": 4,
  "duration_ms": 240000
}
```

## 類似度計算アルゴリズム

### 使用する特徴量

レコメンドに使用する主要特徴量:

| 特徴量         | 重要度 | 理由                   |
| -------------- | ------ | ---------------------- |
| `tempo`        | 高     | 曲の速さは雰囲気に直結 |
| `energy`       | 高     | 激しさ・落ち着きの指標 |
| `valence`      | 高     | 明るい/暗いの印象      |
| `danceability` | 中     | リズム感の類似性       |
| `acousticness` | 中     | 音色の傾向             |
| `key` + `mode` | 低     | 調性（オプション）     |

### 類似度計算式

#### 1. 重み付きユークリッド距離

```
distance = sqrt(
  w_tempo * ((tempo_a - tempo_b) / 250)^2 +
  w_energy * (energy_a - energy_b)^2 +
  w_valence * (valence_a - valence_b)^2 +
  w_danceability * (danceability_a - danceability_b)^2 +
  w_acousticness * (acousticness_a - acousticness_b)^2
)

similarity = 1 / (1 + distance)
```

#### 2. デフォルト重み設定

```go
type FeatureWeights struct {
    Tempo        float64 // デフォルト: 1.5
    Energy       float64 // デフォルト: 1.5
    Valence      float64 // デフォルト: 1.2
    Danceability float64 // デフォルト: 1.0
    Acousticness float64 // デフォルト: 0.8
}
```

#### 3. モード別重み調整

| モード     | Tempo | Energy | Valence | Danceability | Acousticness |
| ---------- | ----- | ------ | ------- | ------------ | ------------ |
| `similar`  | 2.0   | 2.0    | 1.5     | 1.2          | 1.0          |
| `related`  | 0.5   | 0.5    | 0.5     | 0.5          | 0.5          |
| `balanced` | 1.5   | 1.5    | 1.2     | 1.0          | 0.8          |

## ジャンルフィルタリング

KKBOX レコメンドのジャンル逸脱問題を解決するため、Spotify のアーティストジャンル情報を活用したフィルタリングを行う。

### Spotify Artist Genres API

```bash
GET https://api.spotify.com/v1/artists/{id}
```

```json
{
  "name": "猫足 蕾 (CV.芹澤 優)",
  "genres": ["anime", "japanese vgm"]
}
```

### ジャンルグループ定義

Spotify のジャンルタグは細分化されているため、関連ジャンルをグループ化して扱う:

```go
var genreGroups = map[string][]string{
    // オタク系
    "otaku": {
        "anime",
        "japanese vgm",        // ゲーム音楽
        "otacore",
        "anime rock",
        "japanese vocaloid",
        "vocaloid",
        "japanese electropop",
        "denpa",
        "touhou",              // 東方
    },

    // J-POP系
    "jpop": {
        "j-pop",
        "japanese pop",
        "japanese teen pop",
        "city pop",
        "shibuya-kei",
    },

    // ロック系
    "rock": {
        "j-rock",
        "japanese rock",
        "visual kei",
        "alternative rock",
        "japanese metal",
    },

    // K-POP系
    "kpop": {
        "k-pop",
        "korean pop",
        "k-pop boy group",
        "k-pop girl group",
    },

    // アイドル系
    "idol": {
        "japanese idol",
        "japanese idol pop",
        "johnnys",
    },
}
```

### ジャンルボーナス計算

```go
func CalculateGenreBonus(seedGenres, candidateGenres []string) float64 {
    seedGroup := getGenreGroup(seedGenres)
    candidateGroup := getGenreGroup(candidateGenres)

    switch {
    case hasExactMatch(seedGenres, candidateGenres):
        return 1.5  // 完全一致: anime ↔ anime

    case seedGroup == candidateGroup:
        return 1.3  // 同グループ: anime ↔ japanese vgm

    case isRelatedGroup(seedGroup, candidateGroup):
        return 1.0  // 関連グループ: otaku ↔ jpop

    default:
        return 0.5  // 無関係: otaku ↔ kpop（ペナルティ）
    }
}
```

### ジャンルボーナス一覧

| 条件         | ボーナス | 例                                       |
| ------------ | -------- | ---------------------------------------- |
| 完全一致     | 1.5      | seed=anime, candidate=anime              |
| 同グループ   | 1.3      | seed=anime, candidate=japanese vgm       |
| 関連グループ | 1.0      | seed=anime(otaku), candidate=j-pop(jpop) |
| 無関係       | 0.5      | seed=anime(otaku), candidate=k-pop(kpop) |

### 最終スコア計算

```go
// 最終スコア = Audio Features 類似度 × ジャンルボーナス
finalScore := audioSimilarity * genreBonus
```

### Before / After 比較

#### Before（ジャンルフィルタリングなし）

```
入力: 「電脳スペクタクル」(anime)

結果:
1. YOASOBI - 夜に駆ける (j-pop)      audio=0.90
2. 米津玄師 - Lemon (j-pop)          audio=0.88
3. HoneyWorks - 可愛くてごめん (j-pop) audio=0.85
4. Ado - 唱 (j-pop)                  audio=0.83
5. Official髭男dism - Pretender      audio=0.82

→ 全部 J-POP、アニソンなし
```

#### After（ジャンルフィルタリングあり）

```
入力: 「電脳スペクタクル」(anime)

結果:
1. LiSA - 紅蓮華 (anime)       audio=0.85 × genre=1.5 = 1.28
2. Aimer - 残響散歌 (anime)     audio=0.83 × genre=1.5 = 1.25
3. ClariS - コネクト (anime)    audio=0.82 × genre=1.5 = 1.23
4. fripSide - only my railgun  audio=0.80 × genre=1.5 = 1.20
5. 藍井エイル - IGNITE (anime)  audio=0.78 × genre=1.5 = 1.17
6. YOASOBI - 夜に駆ける (j-pop) audio=0.90 × genre=1.0 = 0.90

→ アニソンが上位に、J-POPも良い曲は残る
```

## レコメンドフロー

```
┌─────────────────────────────────────────────────────────────────┐
│ 1. 入力トラック                                                   │
│    └── Spotify Track ID を抽出                                   │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 2. シードトラック情報取得                                          │
│    ├── Spotify: トラック詳細 + Audio Features                     │
│    ├── Spotify: アーティストのジャンル情報 ← 【新規追加】           │
│    └── ISRC 取得                                                 │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 3. 候補トラック収集（並列実行）                                     │
│    ├── KKBOX: SearchByISRC → GetRecommendedTracks              │
│    ├── Spotify: Recommendations API (seed_tracks)               │
│    └── Spotify: アーティストのトップトラック（オプション）            │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 4. 候補の重複除去・フィルタリング                                   │
│    ├── 入力トラックと同一の曲を除外                                 │
│    ├── ISRC で重複を除去                                         │
│    └── 候補を最大50件に制限                                       │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 5. Audio Features バッチ取得                                     │
│    └── GET /v1/audio-features?ids={ids} (最大100件)             │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 6. 候補のアーティストジャンル取得 ← 【新規追加】                     │
│    └── GET /v1/artists?ids={artist_ids} (バッチ取得)            │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 7. スコア計算                                                     │
│    ├── Audio Features 類似度 (audioSimilarity)                  │
│    ├── ジャンルボーナス (genreBonus) ← 【新規追加】                │
│    ├── 最終スコア = audioSimilarity × genreBonus                 │
│    └── match_reasons の判定                                      │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 8. ランキング・返却                                               │
│    ├── finalScore 降順でソート                                    │
│    └── 上位 N 件を返却                                           │
└─────────────────────────────────────────────────────────────────┘
```

## Spotify Recommendations API

### エンドポイント

```
GET https://api.spotify.com/v1/recommendations
```

### パラメータ

| パラメータ        | 説明                               |
| ----------------- | ---------------------------------- |
| `seed_tracks`     | シードトラック ID（最大 5 件）     |
| `seed_artists`    | シードアーティスト ID（最大 5 件） |
| `seed_genres`     | シードジャンル（最大 5 件）        |
| `limit`           | 返却件数（最大 100）               |
| `target_*`        | 目標値（例: `target_energy=0.8`）  |
| `min_*` / `max_*` | 範囲指定（例: `min_tempo=120`）    |

### 使用例

```
GET /v1/recommendations?seed_tracks=4uLU6hMCjMI75M1A2tKUQC&limit=20&target_energy=0.8&target_valence=0.6
```

## ドメインモデル

### AudioFeatures

```go
// domain/audio_features.go

type AudioFeatures struct {
    TrackID          string
    Tempo            float64  // BPM
    Energy           float64  // 0.0-1.0
    Danceability     float64  // 0.0-1.0
    Valence          float64  // 0.0-1.0
    Acousticness     float64  // 0.0-1.0
    Instrumentalness float64  // 0.0-1.0
    Speechiness      float64  // 0.0-1.0
    Liveness         float64  // 0.0-1.0
    Loudness         float64  // dB
    Key              int      // 0-11
    Mode             int      // 0=minor, 1=major
    TimeSignature    int      // 拍子
}
```

### RecommendedTrack

```go
// domain/recommend.go

type RecommendedTrack struct {
    Track           Track
    SimilarityScore float64
    MatchReasons    []string
    AudioFeatures   *AudioFeatures
}

type RecommendResult struct {
    SeedTrack     Track
    SeedFeatures  AudioFeatures
    Items         []RecommendedTrack
    Mode          RecommendMode
}

type RecommendMode string

const (
    RecommendModeSimilar  RecommendMode = "similar"
    RecommendModeRelated  RecommendMode = "related"
    RecommendModeBalanced RecommendMode = "balanced"
)
```

## インターフェース

### SpotifyAPI 拡張

```go
// port/external/spotify.go

type SpotifyAPI interface {
    // 既存
    GetTrackByID(ctx context.Context, id string) (*domain.Track, error)
    SearchByISRC(ctx context.Context, isrc string) (*domain.Track, error)
    SearchTracks(ctx context.Context, query string) ([]domain.Track, error)
    GetArtistByID(ctx context.Context, id string) (*domain.Artist, error)
    GetAlbumByID(ctx context.Context, id string) (*domain.Album, error)

    // 新規追加
    GetAudioFeatures(ctx context.Context, trackID string) (*domain.AudioFeatures, error)
    GetAudioFeaturesBatch(ctx context.Context, trackIDs []string) ([]domain.AudioFeatures, error)
    GetRecommendations(ctx context.Context, params RecommendationParams) ([]domain.Track, error)
    GetArtistGenres(ctx context.Context, artistID string) ([]string, error)
    GetArtistGenresBatch(ctx context.Context, artistIDs []string) (map[string][]string, error)
}

type RecommendationParams struct {
    SeedTracks   []string
    SeedArtists  []string
    SeedGenres   []string
    Limit        int
    TargetTempo  *float64
    TargetEnergy *float64
    // ... 他のターゲットパラメータ
}
```

## UseCase

```go
// usecase/recommend.go

type RecommendUseCase struct {
    spotifyAPI external.SpotifyAPI
    kkboxAPI   external.KKBOXAPI
    calculator *SimilarityCalculator
}

func (uc *RecommendUseCase) GetRecommendations(
    ctx context.Context,
    trackID string,
    mode domain.RecommendMode,
    limit int,
) (*domain.RecommendResult, error) {
    // 1. シードトラック情報取得
    // 2. 候補収集（並列）
    // 3. 重複除去
    // 4. Audio Features バッチ取得
    // 5. 類似度計算
    // 6. ランキング・返却
}
```

## 実装フェーズ

### Phase 1: 基盤実装（2 日）

- [ ] AudioFeatures ドメインモデル追加
- [ ] Spotify Audio Features API 実装
- [ ] Spotify Recommendations API 実装
- [ ] 単体テスト

### Phase 2: 類似度計算・ジャンルフィルタリング（2 日）

- [ ] SimilarityCalculator 実装
- [ ] 重み付け設定
- [ ] GenreMatcher 実装
- [ ] ジャンルグループ定義
- [ ] ジャンルボーナス計算ロジック
- [ ] match_reasons 判定ロジック
- [ ] 単体テスト

### Phase 3: UseCase・Handler（1.5 日）

- [ ] RecommendUseCase 実装
- [ ] 候補収集の並列化
- [ ] RecommendHandler 実装
- [ ] ルーティング追加

### Phase 4: テスト・チューニング（2 日）

- [ ] 結合テスト
- [ ] 実データでのチューニング
- [ ] 重み調整
- [ ] ドキュメント更新

### 合計: 約 7.5 日

## 今後の拡張案

1. **ユーザー嗜好学習**: フィードバックを元に重みを調整
2. ~~**ジャンル考慮**: 同一ジャンル内での類似度計算~~ → ジャンルフィルタリングとして実装予定
3. **キャッシュ**: Audio Features のキャッシュ（Redis）
4. **プレイリスト対応**: 複数曲を入力としたレコメンド
5. **フィルタリングオプション**: 年代、国、Explicit 除外など
6. **ジャンルグループの拡張**: ユーザーからのフィードバックでグループ定義を改善

## 参考

- [Spotify Audio Features API](https://developer.spotify.com/documentation/web-api/reference/get-audio-features)
- [Spotify Recommendations API](https://developer.spotify.com/documentation/web-api/reference/get-recommendations)
- [KKBOX Recommended Tracks API](https://developer.kkbox.com/docs)
