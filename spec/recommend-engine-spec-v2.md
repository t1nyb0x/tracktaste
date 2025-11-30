# レコメンドエンジン仕様書 v2

## 概要

入力されたトラックと音楽的に類似した曲をレコメンドするエンジン。**Deezer API（BPM/Duration/Gain）** と **MusicBrainz API（タグ/アーティスト関連性）** を活用して、楽曲の特性に基づいた精度の高いレコメンドを実現する。

### 変更履歴

| 日付       | 変更内容                                                          |
| ---------- | ----------------------------------------------------------------- |
| 2025-01-XX | v2: Spotify Audio Features → Deezer + MusicBrainz に移行          |
| 2025-01-XX | v1: 初版（Spotify Audio Features ベース）※ API 廃止により使用不可 |

### 重要事項

> ⚠️ **Spotify Audio Features API は 2024 年 11 月に廃止されました。**
>
> - GET /v1/audio-features/{id} → 403 Forbidden
> - GET /v1/recommendations → 404 Not Found
>
> 本仕様書では代替として **Deezer + MusicBrainz** を使用します。

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

## 解決アプローチ

### ハイブリッドレコメンド

1. **関連性ベース**: 既存の KKBOX/Spotify レコメンドで候補を取得
2. **雰囲気ベース**: Deezer BPM/Gain + MusicBrainz タグによる類似度でランキング
3. **ジャンル考慮**: MusicBrainz タグ + Spotify アーティストジャンルでフィルタリング

### アーキテクチャ

```
┌─────────────────────────────────────────────────────────────────┐
│                    TrackTaste Recommend Engine v2               │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────┐    ┌───────────┐    ┌─────────────────────────┐  │
│  │  KKBOX   │    │  Spotify  │    │      Deezer API         │  │
│  │ Recommend│    │  (Track)  │    │ ┌─────────────────────┐ │  │
│  └────┬─────┘    └─────┬─────┘    │ │ BPM / Duration /    │ │  │
│       │                │          │ │ Gain                │ │  │
│       └────────┬───────┘          │ └─────────────────────┘ │  │
│                ↓                  └────────────┬────────────┘  │
│       ┌────────────────┐                       │               │
│       │ 候補トラック   │                       │               │
│       │ 収集・重複除去 │                       │               │
│       └───────┬────────┘                       │               │
│               │                                │               │
│               ↓                                ↓               │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │            Similarity Calculator                        │   │
│  │  ┌────────────────┐  ┌────────────────────────────────┐│   │
│  │  │ Deezer 特徴量   │  │ MusicBrainz タグ              ││   │
│  │  │ - BPM類似度     │  │ - Jaccard類似度               ││   │
│  │  │ - Duration類似度│  │ - アーティスト関連ボーナス    ││   │
│  │  │ - Gain類似度    │  │                               ││   │
│  │  └────────────────┘  └────────────────────────────────┘│   │
│  │                                                         │   │
│  │  最終スコア = Σ(重み × 各類似度) × ジャンルボーナス     │   │
│  └─────────────────────────────────────────────────────────┘   │
│                               │                                │
│                               ↓                                │
│                    ┌───────────────────┐                       │
│                    │ スコア順でソート  │                       │
│                    │ 上位N件を返却     │                       │
│                    └───────────────────┘                       │
└─────────────────────────────────────────────────────────────────┘
```

## API エンドポイント

### リクエスト

```
GET /v1/track/recommend?url={spotify_url}&mode={mode}&limit={limit}
```

### パラメータ

| パラメータ | 型     | 必須 | デフォルト | 説明                                               |
| ---------- | ------ | ---- | ---------- | -------------------------------------------------- |
| `url`      | string | ✓    | -          | Spotify トラック URL                               |
| `mode`     | string | -    | `balanced` | レコメンドモード: `similar`, `related`, `balanced` |
| `limit`    | int    | -    | `10`       | 返却件数（1〜50）                                  |

### モード詳細

| モード     | 説明                          | 用途                             |
| ---------- | ----------------------------- | -------------------------------- |
| `similar`  | Deezer 特徴量（BPM/Gain）重視 | テンポや音圧が似た曲を探したい   |
| `related`  | MusicBrainz タグ・関連性重視  | 同じジャンル・スタイルの曲を探す |
| `balanced` | バランス型（デフォルト）      | 総合的に似た曲を探す             |

### レスポンス

```json
{
  "success": true,
  "data": {
    "seed_track": {
      "id": "spotify:track:xxx",
      "name": "電脳スペクタクル",
      "artist": {
        "name": "鳳ここな(CV.長谷川育美)"
      },
      "album": {
        "name": "ワールドダイスター"
      }
    },
    "seed_features": {
      "bpm": 175.0,
      "duration_seconds": 245,
      "gain": -7.2,
      "tags": ["anime", "jpop", "female vocalist"]
    },
    "items": [
      {
        "track": {
          "id": "spotify:track:yyy",
          "name": "Fly Me to the Star",
          "artist": { "name": "新妻八恵(CV.石見舞菜香)" }
        },
        "similarity_score": 0.92,
        "match_reasons": ["similar_bpm", "same_tag:anime", "artist_relation"],
        "features": {
          "bpm": 172.0,
          "duration_seconds": 238,
          "gain": -6.8,
          "tags": ["anime", "jpop", "female vocalist"]
        }
      }
    ],
    "mode": "balanced"
  }
}
```

---

## データソース

### 1. Deezer API

ISRC を使って Deezer からトラック情報を取得し、BPM・Duration・Gain を取得する。

#### エンドポイント

```bash
# ISRC検索
GET https://api.deezer.com/track/isrc:{isrc}

# トラック検索（ISRCがない場合）
GET https://api.deezer.com/search/track?q=track:"{title}" artist:"{artist}"
```

#### 取得できる特徴量

| フィールド        | 型    | 説明                           |
| ----------------- | ----- | ------------------------------ |
| `bpm`             | float | テンポ (BPM)                   |
| `duration`        | int   | 曲の長さ（秒）                 |
| `gain`            | float | ReplayGain 値 (dB)、音圧の指標 |
| `explicit_lyrics` | bool  | 歌詞に explicit な内容があるか |

#### レスポンス例

```json
{
  "id": 123456789,
  "title": "電脳スペクタクル",
  "isrc": "JPAB12345678",
  "duration": 245,
  "bpm": 175.0,
  "gain": -7.2,
  "explicit_lyrics": false,
  "artist": {
    "id": 98765,
    "name": "鳳ここな(CV.長谷川育美)"
  }
}
```

#### 注意事項

- **レート制限**: 50 requests / 5 seconds
- **認証不要**: API キー不要で利用可能
- **カバレッジ**: 日本のアニソンは収録率が低い場合がある

---

### 2. MusicBrainz API

Recording（楽曲）と Artist（アーティスト）のタグ・関連情報を取得。

#### エンドポイント

```bash
# ISRC検索
GET https://musicbrainz.org/ws/2/isrc/{isrc}?inc=recordings+artists+tags&fmt=json

# Recording詳細（タグ含む）
GET https://musicbrainz.org/ws/2/recording/{mbid}?inc=tags+artist-rels&fmt=json

# Artist詳細（タグ・関連アーティスト含む）
GET https://musicbrainz.org/ws/2/artist/{mbid}?inc=tags+artist-rels&fmt=json
```

#### 取得できる情報

| 情報             | 説明                                                   |
| ---------------- | ------------------------------------------------------ |
| `tags`           | ユーザー投票によるタグ（ジャンル、スタイル、特徴など） |
| `artist-rels`    | アーティスト間の関連（メンバー、別名義、コラボなど）   |
| `recording-rels` | 楽曲間の関連（カバー、リミックスなど）                 |

#### タグの例

```json
{
  "tags": [
    { "name": "anime", "count": 15 },
    { "name": "female vocalist", "count": 8 },
    { "name": "japanese", "count": 12 },
    { "name": "jpop", "count": 5 }
  ]
}
```

#### アーティスト関連の例

```json
{
  "relations": [
    {
      "type": "member of band",
      "artist": { "name": "ワールドダイスター" }
    },
    {
      "type": "voice actor",
      "artist": { "name": "長谷川育美" }
    }
  ]
}
```

#### 注意事項

- **レート制限**: 1 request / second（User-Agent 必須）
- **認証不要**: API キー不要
- **カバレッジ**: コミュニティベースのため、新しい曲・マイナーな曲は情報が少ない場合がある

---

## 類似度計算アルゴリズム

### 特徴量一覧

| 特徴量             | ソース      | 範囲/型       | 説明                  |
| ------------------ | ----------- | ------------- | --------------------- |
| `bpm`              | Deezer      | 0-250 (float) | テンポ                |
| `duration_seconds` | Deezer      | 0-∞ (int)     | 曲の長さ（秒）        |
| `gain`             | Deezer      | -20 to 0 (dB) | 音圧（ReplayGain）    |
| `tags`             | MusicBrainz | []string      | ジャンル/スタイルタグ |
| `artist_relations` | MusicBrainz | []Relation    | アーティスト関連情報  |

### 類似度計算式

#### 1. BPM 類似度

```go
func bpmSimilarity(bpmA, bpmB float64) float64 {
    // BPMの差分を正規化（0-250の範囲を想定）
    diff := math.Abs(bpmA - bpmB) / 250.0
    return 1.0 - diff
}
```

#### 2. Duration 類似度

```go
func durationSimilarity(durA, durB int) float64 {
    // 曲の長さの差分を正規化（最大10分=600秒を想定）
    diff := math.Abs(float64(durA - durB)) / 600.0
    if diff > 1.0 {
        diff = 1.0
    }
    return 1.0 - diff
}
```

#### 3. Gain 類似度

```go
func gainSimilarity(gainA, gainB float64) float64 {
    // Gainの差分を正規化（-20 to 0 dBの範囲）
    diff := math.Abs(gainA - gainB) / 20.0
    return 1.0 - diff
}
```

#### 4. タグ類似度（Jaccard 係数）

```go
func tagSimilarity(tagsA, tagsB []string) float64 {
    setA := toSet(tagsA)
    setB := toSet(tagsB)

    intersection := setIntersection(setA, setB)
    union := setUnion(setA, setB)

    if len(union) == 0 {
        return 0.5  // タグがない場合はニュートラル
    }
    return float64(len(intersection)) / float64(len(union))
}
```

#### 5. アーティスト関連ボーナス

```go
func artistRelationBonus(seedArtist, candidateArtist ArtistInfo) float64 {
    switch {
    case isSameArtist(seedArtist, candidateArtist):
        return 1.5  // 同一アーティスト

    case hasSameGroup(seedArtist, candidateArtist):
        return 1.3  // 同じグループ/ユニットのメンバー

    case hasCollaboration(seedArtist, candidateArtist):
        return 1.2  // コラボ経験あり

    case hasSameVoiceActor(seedArtist, candidateArtist):
        return 1.2  // 同じ声優（アニソン用）

    case hasSameLabelOrProducer(seedArtist, candidateArtist):
        return 1.1  // 同じレーベル/プロデューサー

    default:
        return 1.0  // 関連なし
    }
}
```

### 重み設定

#### デフォルト重み

```go
type FeatureWeights struct {
    BPM            float64  // デフォルト: 1.5
    Duration       float64  // デフォルト: 0.5
    Gain           float64  // デフォルト: 1.2
    TagSimilarity  float64  // デフォルト: 2.0
}
```

#### モード別重み

| モード     | BPM | Duration | Gain | TagSimilarity |
| ---------- | --- | -------- | ---- | ------------- |
| `similar`  | 2.0 | 0.8      | 1.5  | 1.0           |
| `related`  | 0.5 | 0.3      | 0.5  | 3.0           |
| `balanced` | 1.5 | 0.5      | 1.2  | 2.0           |

### 最終スコア計算

```go
func calculateFinalScore(seedFeatures, candidateFeatures TrackFeatures, mode RecommendMode, artistInfo ArtistInfo) float64 {
    weights := weightsForMode(mode)

    // 各特徴量の類似度
    bpmSim := bpmSimilarity(seedFeatures.BPM, candidateFeatures.BPM)
    durSim := durationSimilarity(seedFeatures.Duration, candidateFeatures.Duration)
    gainSim := gainSimilarity(seedFeatures.Gain, candidateFeatures.Gain)
    tagSim := tagSimilarity(seedFeatures.Tags, candidateFeatures.Tags)

    // 重み付き平均
    totalWeight := weights.BPM + weights.Duration + weights.Gain + weights.TagSimilarity
    baseSimilarity := (
        weights.BPM * bpmSim +
        weights.Duration * durSim +
        weights.Gain * gainSim +
        weights.TagSimilarity * tagSim
    ) / totalWeight

    // アーティスト関連ボーナス
    artistBonus := artistRelationBonus(seedFeatures.Artist, candidateFeatures.Artist)

    // ジャンルボーナス（既存のGenreMatcher使用）
    genreBonus := genreMatcher.CalculateBonus(seedFeatures.Tags, candidateFeatures.Tags)

    return baseSimilarity * artistBonus * genreBonus
}
```

---

## ジャンル/タグマッチング

### タググループ定義

MusicBrainz のタグを Spotify のジャンルグループと同様にグループ化:

```go
var tagGroups = map[string][]string{
    // アニメ/ゲーム系
    "anime": {
        "anime",
        "anime song",
        "anison",
        "game music",
        "video game music",
        "vgm",
        "visual novel",
        "denpa",
        "touhou",
        "vocaloid",
        "otaku",
    },

    // J-POP系
    "jpop": {
        "jpop",
        "j-pop",
        "japanese pop",
        "city pop",
        "shibuya-kei",
    },

    // ロック系
    "rock": {
        "japanese rock",
        "j-rock",
        "visual kei",
        "alternative rock",
        "japanese metal",
    },

    // K-POP系
    "kpop": {
        "kpop",
        "k-pop",
        "korean pop",
    },

    // アイドル系
    "idol": {
        "japanese idol",
        "idol",
        "akb48",
        "johnny's",
    },
}
```

### タググループ関連性

```go
var tagGroupRelations = map[string][]string{
    "anime": {"jpop", "idol"},  // アニソンはJ-POP、アイドルと親和性あり
    "jpop":  {"anime", "idol", "rock"},
    "rock":  {"jpop"},
    "kpop":  {"idol"},
    "idol":  {"jpop", "anime", "kpop"},
}
```

### ジャンルボーナス計算

```go
func CalculateGenreBonus(seedTags, candidateTags []string) float64 {
    seedGroup := getTagGroup(seedTags)
    candidateGroup := getTagGroup(candidateTags)

    switch {
    case hasExactTagMatch(seedTags, candidateTags):
        return 1.5  // 完全一致タグあり: anime ↔ anime

    case seedGroup == candidateGroup:
        return 1.3  // 同グループ: anime ↔ vocaloid

    case isRelatedGroup(seedGroup, candidateGroup):
        return 1.0  // 関連グループ: anime ↔ jpop

    default:
        return 0.5  // 無関係: anime ↔ kpop（ペナルティ）
    }
}
```

---

## レコメンドフロー

```
┌─────────────────────────────────────────────────────────────────┐
│ 1. 入力トラック                                                   │
│    └── Spotify Track ID を抽出                                   │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 2. シードトラック情報取得（並列）                                  │
│    ├── Spotify: トラック詳細 + ISRC                              │
│    ├── Deezer: ISRC検索 → BPM / Duration / Gain                 │
│    ├── MusicBrainz: ISRC検索 → タグ / アーティスト関連           │
│    └── Spotify: アーティストのジャンル情報                        │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 3. 候補トラック収集（並列実行）                                     │
│    ├── KKBOX: SearchByISRC → GetRecommendedTracks              │
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
│ 5. 候補の特徴量取得（並列実行）                                     │
│    ├── Deezer: ISRC → BPM / Duration / Gain                     │
│    └── MusicBrainz: ISRC → タグ / アーティスト関連               │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 6. スコア計算                                                     │
│    ├── Deezer 特徴量類似度（BPM, Duration, Gain）                │
│    ├── タグ類似度（Jaccard係数）                                  │
│    ├── アーティスト関連ボーナス                                    │
│    ├── ジャンルボーナス                                           │
│    ├── 最終スコア = baseSimilarity × artistBonus × genreBonus    │
│    └── match_reasons の判定                                      │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│ 7. ランキング・返却                                               │
│    ├── finalScore 降順でソート                                    │
│    └── 上位 N 件を返却                                           │
└─────────────────────────────────────────────────────────────────┘
```

---

## ドメインモデル

### TrackFeatures（新規）

```go
// domain/track_features.go

// TrackFeatures は Deezer + MusicBrainz から取得した楽曲特徴量
type TrackFeatures struct {
    TrackID         string
    ISRC            string

    // Deezer から取得
    BPM             float64  // テンポ (0-250)
    DurationSeconds int      // 曲の長さ（秒）
    Gain            float64  // ReplayGain (dB)

    // MusicBrainz から取得
    Tags            []Tag    // ジャンル/スタイルタグ
    ArtistMBID      string   // MusicBrainz Artist ID
}

type Tag struct {
    Name  string
    Count int  // 投票数（重み付けに使用可能）
}
```

### ArtistRelation（新規）

```go
// domain/artist_relation.go

type ArtistRelation struct {
    ArtistMBID      string
    ArtistName      string
    Relations       []Relation
}

type Relation struct {
    Type       string  // "member of band", "voice actor", "collaboration", etc.
    TargetMBID string
    TargetName string
}
```

### RecommendedTrack（更新）

```go
// domain/recommend.go

type RecommendedTrack struct {
    Track           Track
    SimilarityScore float64
    MatchReasons    []string
    Features        *TrackFeatures  // AudioFeatures から TrackFeatures に変更
}

type RecommendResult struct {
    SeedTrack    Track
    SeedFeatures TrackFeatures  // AudioFeatures から TrackFeatures に変更
    Items        []RecommendedTrack
    Mode         RecommendMode
}
```

---

## インターフェース

### DeezerAPI（新規）

```go
// port/external/deezer.go

type DeezerAPI interface {
    // ISRC でトラック検索
    GetTrackByISRC(ctx context.Context, isrc string) (*DeezerTrack, error)

    // タイトル+アーティストで検索（ISRCがない場合のフォールバック）
    SearchTrack(ctx context.Context, title, artist string) (*DeezerTrack, error)

    // バッチ取得（複数ISRC）
    GetTracksByISRCBatch(ctx context.Context, isrcs []string) (map[string]*DeezerTrack, error)
}

type DeezerTrack struct {
    ID              int64
    Title           string
    ISRC            string
    BPM             float64
    DurationSeconds int
    Gain            float64
}
```

### MusicBrainzAPI（新規）

```go
// port/external/musicbrainz.go

type MusicBrainzAPI interface {
    // ISRC で Recording 検索
    GetRecordingByISRC(ctx context.Context, isrc string) (*MBRecording, error)

    // Recording の詳細（タグ含む）
    GetRecordingWithTags(ctx context.Context, mbid string) (*MBRecording, error)

    // Artist の詳細（タグ・関連含む）
    GetArtistWithRelations(ctx context.Context, mbid string) (*MBArtist, error)

    // バッチ取得
    GetRecordingsByISRCBatch(ctx context.Context, isrcs []string) (map[string]*MBRecording, error)
}

type MBRecording struct {
    MBID       string
    Title      string
    ISRC       string
    Tags       []Tag
    ArtistMBID string
}

type MBArtist struct {
    MBID      string
    Name      string
    Tags      []Tag
    Relations []Relation
}
```

### SpotifyAPI（既存を維持）

```go
// port/external/spotify.go

type SpotifyAPI interface {
    // 既存メソッド（維持）
    GetTrackByID(ctx context.Context, id string) (*domain.Track, error)
    SearchByISRC(ctx context.Context, isrc string) (*domain.Track, error)
    SearchTracks(ctx context.Context, query string) ([]domain.Track, error)
    GetArtistByID(ctx context.Context, id string) (*domain.Artist, error)
    GetAlbumByID(ctx context.Context, id string) (*domain.Album, error)

    // ジャンル取得（維持）
    GetArtistGenres(ctx context.Context, artistID string) ([]string, error)
    GetArtistGenresBatch(ctx context.Context, artistIDs []string) (map[string][]string, error)

    // 以下は削除（API廃止のため）
    // GetAudioFeatures → 削除
    // GetAudioFeaturesBatch → 削除
    // GetRecommendations → 削除
}
```

---

## UseCase

```go
// usecase/recommend.go

type RecommendUseCase struct {
    spotifyAPI     external.SpotifyAPI
    kkboxAPI       external.KKBOXAPI
    deezerAPI      external.DeezerAPI       // 新規
    musicBrainzAPI external.MusicBrainzAPI  // 新規
    calculator     *SimilarityCalculator
    genreMatcher   *GenreMatcher
}

func (uc *RecommendUseCase) GetRecommendations(
    ctx context.Context,
    trackID string,
    mode domain.RecommendMode,
    limit int,
) (*domain.RecommendResult, error) {
    // 1. シードトラック情報取得（Spotify）
    // 2. シードの特徴量取得（Deezer + MusicBrainz、並列）
    // 3. 候補収集（KKBOX）
    // 4. 重複除去
    // 5. 候補の特徴量バッチ取得（Deezer + MusicBrainz、並列）
    // 6. 類似度計算
    // 7. ランキング・返却
}
```

---

## キャッシュ戦略

### キャッシュ対象

| データ               | TTL  | 理由                                   |
| -------------------- | ---- | -------------------------------------- |
| Deezer トラック情報  | 7 日 | BPM/Duration/Gain は変わらない         |
| MusicBrainz タグ     | 1 日 | コミュニティ投票で徐々に変化           |
| MusicBrainz 関連情報 | 7 日 | アーティスト関連は頻繁に変わらない     |
| Spotify ジャンル     | 7 日 | アーティストジャンルは頻繁に変わらない |

### キャッシュキー設計

```
deezer:track:{isrc}
musicbrainz:recording:{isrc}
musicbrainz:artist:{mbid}
spotify:artist:genres:{artist_id}
```

---

## エラーハンドリング

### フォールバック戦略

| 状況                               | フォールバック                          |
| ---------------------------------- | --------------------------------------- |
| Deezer で ISRC が見つからない      | タイトル+アーティスト検索を試行         |
| Deezer で BPM が取得できない       | BPM 類似度の重みを 0 にして計算         |
| MusicBrainz でタグがない           | Spotify アーティストジャンルで代替      |
| MusicBrainz で ISRC が見つからない | タグ類似度をニュートラル(0.5)として計算 |
| 両 API からデータ取得不可          | 類似度 0.5（ニュートラル）として処理    |

---

## レート制限対策

| API         | 制限           | 対策                               |
| ----------- | -------------- | ---------------------------------- |
| Deezer      | 50 req / 5 sec | キャッシュ活用、並列リクエスト制限 |
| MusicBrainz | 1 req / sec    | User-Agent 設定、直列リクエスト    |
| Spotify     | 制限緩め       | 既存の実装を維持                   |
| KKBOX       | 制限緩め       | 既存の実装を維持                   |

### MusicBrainz 特別対応

```go
// MusicBrainz はレート制限が厳しいため、専用の rate limiter を使用
type MusicBrainzClient struct {
    httpClient *http.Client
    limiter    *rate.Limiter  // 1 req/sec
    userAgent  string         // 必須: "TrackTaste/1.0 (contact@example.com)"
}
```

---

## 実装フェーズ

### Phase 1: Deezer Gateway 実装（1 日）

- [ ] DeezerAPI インターフェース定義
- [ ] Deezer Gateway 実装（ISRC 検索、タイトル検索）
- [ ] キャッシュ統合
- [ ] 単体テスト

### Phase 2: MusicBrainz Gateway 実装（1.5 日）

- [ ] MusicBrainzAPI インターフェース定義
- [ ] MusicBrainz Gateway 実装
- [ ] Rate Limiter 実装
- [ ] キャッシュ統合
- [ ] 単体テスト

### Phase 3: TrackFeatures ドメイン実装（0.5 日）

- [ ] TrackFeatures ドメインモデル
- [ ] ArtistRelation ドメインモデル
- [ ] 既存 AudioFeatures からの移行

### Phase 4: SimilarityCalculator 更新（1 日）

- [ ] 新しい類似度計算式の実装
- [ ] Jaccard 係数実装
- [ ] アーティスト関連ボーナス実装
- [ ] モード別重み調整
- [ ] 単体テスト

### Phase 5: RecommendUseCase 更新（1 日）

- [ ] Deezer/MusicBrainz API 統合
- [ ] 並列処理の最適化
- [ ] エラーハンドリング・フォールバック
- [ ] 結合テスト

### Phase 6: テスト・チューニング（1.5 日）

- [ ] 実データでのテスト
- [ ] 重み調整
- [ ] パフォーマンス最適化
- [ ] ドキュメント更新

### 合計: 約 6.5 日

---

## Before / After 比較

### Before（v1: Spotify Audio Features）

```
問題: Spotify Audio Features API が 403 を返す
結果: すべての類似度が 0.5（デフォルト値）になる

入力: 「電脳スペクタクル」

結果:
1. Track A - similarity: 0.50
2. Track B - similarity: 0.50
3. Track C - similarity: 0.50
→ すべて同じスコア、ランダムな順序
```

### After（v2: Deezer + MusicBrainz）

```
入力: 「電脳スペクタクル」(anime, BPM=175, Gain=-7.2dB)

結果:
1. Fly Me to the Star (anime, BPM=172)
   - BPM類似度: 0.99
   - タグ類似度: 0.85 (anime, jpop共通)
   - アーティスト関連: 1.3 (同ユニット)
   - 最終スコア: 0.92

2. 紅蓮華 (anime, BPM=180)
   - BPM類似度: 0.98
   - タグ類似度: 0.75 (anime共通)
   - アーティスト関連: 1.0
   - 最終スコア: 0.85

3. 夜に駆ける (jpop, BPM=165)
   - BPM類似度: 0.96
   - タグ類似度: 0.30 (共通タグなし)
   - ジャンルボーナス: 1.0 (anime↔jpop関連)
   - 最終スコア: 0.72

→ アニソンが上位、J-POPは関連グループとして下位に
```

---

## 参考リンク

- [Deezer API Documentation](https://developers.deezer.com/api)
- [MusicBrainz API Documentation](https://musicbrainz.org/doc/MusicBrainz_API)
- [MusicBrainz Rate Limiting](https://musicbrainz.org/doc/MusicBrainz_API/Rate_Limiting)
- [Spotify Web API](https://developer.spotify.com/documentation/web-api) - トラック/アーティスト情報のみ使用
- [KKBOX Open API](https://developer.kkbox.com/docs)
