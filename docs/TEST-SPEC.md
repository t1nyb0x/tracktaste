# テスト仕様書

## 概要

本ドキュメントは trackTaste プロジェクトのテスト仕様を定義します。

### テスト方針

- テーブル駆動テスト（Table-Driven Tests）を採用
- httptest を使用した HTTP モック
- インターフェースベースのモックによる依存性の分離
- 正常系・異常系の網羅的なテスト

### カバレッジ目標

- 全体: 80%以上
- ビジネスロジック (usecase/v2): 75%以上 ✅ (現状: 76%)
- ビジネスロジック (usecase/v1): 80%以上
- ハンドラー (handler): 85%以上

---

## テストファイル一覧

| パッケージ                  | テストファイル         | 対象                         |
| --------------------------- | ---------------------- | ---------------------------- |
| adapter/handler             | track_test.go          | トラック API                 |
| adapter/handler             | artist_test.go         | アーティスト API             |
| adapter/handler             | album_test.go          | アルバム API                 |
| adapter/handler             | extract_test.go        | URL 抽出ユーティリティ       |
| adapter/handler             | recommend_test.go      | レコメンド API               |
| adapter/gateway/spotify     | gateway_test.go        | Spotify API クライアント     |
| adapter/gateway/kkbox       | gateway_test.go        | KKBOX API クライアント       |
| adapter/gateway/deezer      | gateway_test.go        | Deezer API クライアント      |
| adapter/gateway/musicbrainz | gateway_test.go        | MusicBrainz API クライアント |
| adapter/gateway/cache       | repository_test.go     | キャッシュリポジトリ         |
| usecase/v1                  | track_test.go          | トラックユースケース         |
| usecase/v1                  | artist_test.go         | アーティストユースケース     |
| usecase/v1                  | album_test.go          | アルバムユースケース         |
| usecase/v1                  | similar_tracks_test.go | 類似トラックユースケース     |
| usecase/v2                  | recommend_test.go      | レコメンドユースケース V2    |
| usecase/v2                  | similarity_test.go     | 類似度計算 V2                |
| usecase                     | genre_matcher_test.go  | ジャンルマッチャー           |
| domain                      | errors_test.go         | ドメインエラー               |
| config                      | config_test.go         | 設定                         |
| util/logger                 | logger_test.go         | ロガー                       |
| testutil                    | mock_test.go           | テストユーティリティ         |

---

## 1. Handler テスト

### 1.1 TrackHandler (track_test.go)

#### TestTrackHandler_FetchByURL

トラック取得エンドポイント `/v1/track/fetch` のテスト

| No  | テストケース               | 入力                                            | 期待結果                | ステータス              |
| --- | -------------------------- | ----------------------------------------------- | ----------------------- | ----------------------- |
| 1   | 正常系: 有効な URL         | `https://open.spotify.com/track/abc123`         | トラック情報取得成功    | 200 OK                  |
| 2   | 正常系: intl-ja 付き URL   | `https://open.spotify.com/intl-ja/track/abc123` | トラック情報取得成功    | 200 OK                  |
| 3   | 異常系: 空の URL           | `""`                                            | EMPTY_PARAM             | 400 Bad Request         |
| 4   | 異常系: Spotify 以外の URL | `https://music.apple.com/track/abc123`          | NOT_SPOTIFY_URL         | 400 Bad Request         |
| 5   | 異常系: artist の URL      | `https://open.spotify.com/artist/abc123`        | DIFFERENT_SPOTIFY_URL   | 400 Bad Request         |
| 6   | 異常系: API エラー         | 有効な URL + API エラー                         | SOMETHING_SPOTIFY_ERROR | 503 Service Unavailable |

#### TestTrackHandler_Search

トラック検索エンドポイント `/v1/track/search` のテスト

| No  | テストケース         | 入力                  | 期待結果                | ステータス              |
| --- | -------------------- | --------------------- | ----------------------- | ----------------------- |
| 1   | 正常系: 有効なクエリ | `q=test`              | 検索結果リスト          | 200 OK                  |
| 2   | 正常系: 結果 0 件    | `q=nonexistent`       | 空のリスト              | 200 OK                  |
| 3   | 異常系: クエリなし   | なし                  | EMPTY_PARAM             | 400 Bad Request         |
| 4   | 異常系: API エラー   | `q=test` + API エラー | SOMETHING_SPOTIFY_ERROR | 503 Service Unavailable |

#### TestTrackHandler_GetSimilarTracks

類似トラック取得エンドポイント `/v1/track/similar` のテスト

| No  | テストケース               | 入力            | 期待結果           | ステータス      |
| --- | -------------------------- | --------------- | ------------------ | --------------- |
| 1   | 正常系: 有効な URL         | Spotify URL     | 類似トラックリスト | 200 OK          |
| 2   | 異常系: 空の URL           | なし            | EMPTY_PARAM        | 400 Bad Request |
| 3   | 異常系: Spotify 以外の URL | Apple Music URL | NOT_SPOTIFY_URL    | 400 Bad Request |

### 1.2 ArtistHandler (artist_test.go)

#### TestArtistHandler_FetchByURL

アーティスト取得エンドポイント `/v1/artist/fetch` のテスト

| No  | テストケース               | 入力                                             | 期待結果                | ステータス              |
| --- | -------------------------- | ------------------------------------------------ | ----------------------- | ----------------------- |
| 1   | 正常系: 有効な URL         | `https://open.spotify.com/artist/abc123`         | アーティスト情報        | 200 OK                  |
| 2   | 正常系: intl-ja 付き URL   | `https://open.spotify.com/intl-ja/artist/abc123` | アーティスト情報        | 200 OK                  |
| 3   | 異常系: 空の URL           | `""`                                             | EMPTY_PARAM             | 400 Bad Request         |
| 4   | 異常系: Spotify 以外の URL | Apple Music URL                                  | NOT_SPOTIFY_URL         | 400 Bad Request         |
| 5   | 異常系: track の URL       | `https://open.spotify.com/track/abc123`          | DIFFERENT_SPOTIFY_URL   | 400 Bad Request         |
| 6   | 異常系: API エラー         | 有効な URL + API エラー                          | SOMETHING_SPOTIFY_ERROR | 503 Service Unavailable |

### 1.3 AlbumHandler (album_test.go)

#### TestAlbumHandler_FetchByURL

アルバム取得エンドポイント `/v1/album/fetch` のテスト

| No  | テストケース               | 入力                                            | 期待結果                | ステータス              |
| --- | -------------------------- | ----------------------------------------------- | ----------------------- | ----------------------- |
| 1   | 正常系: 有効な URL         | `https://open.spotify.com/album/abc123`         | アルバム情報            | 200 OK                  |
| 2   | 正常系: intl-ja 付き URL   | `https://open.spotify.com/intl-ja/album/abc123` | アルバム情報            | 200 OK                  |
| 3   | 異常系: 空の URL           | `""`                                            | EMPTY_PARAM             | 400 Bad Request         |
| 4   | 異常系: Spotify 以外の URL | Apple Music URL                                 | NOT_SPOTIFY_URL         | 400 Bad Request         |
| 5   | 異常系: track の URL       | `https://open.spotify.com/track/abc123`         | DIFFERENT_SPOTIFY_URL   | 400 Bad Request         |
| 6   | 異常系: API エラー         | 有効な URL + API エラー                         | SOMETHING_SPOTIFY_ERROR | 503 Service Unavailable |

### 1.4 Extract (extract_test.go)

#### TestExtractSpotifyID

Spotify URL から ID を抽出するユーティリティのテスト

| No  | テストケース             | 入力                                            | 期待結果                |
| --- | ------------------------ | ----------------------------------------------- | ----------------------- |
| 1   | 正常系: track URL        | `https://open.spotify.com/track/abc123`         | `abc123`, `track`, nil  |
| 2   | 正常系: intl-ja 付き URL | `https://open.spotify.com/intl-ja/track/abc123` | `abc123`, `track`, nil  |
| 3   | 正常系: artist URL       | `https://open.spotify.com/artist/xyz789`        | `xyz789`, `artist`, nil |
| 4   | 正常系: album URL        | `https://open.spotify.com/album/def456`         | `def456`, `album`, nil  |
| 5   | 異常系: 無効な URL       | `not-a-url`                                     | エラー                  |
| 6   | 異常系: Spotify 以外     | `https://music.apple.com/track/abc123`          | エラー                  |

---

## 2. Usecase テスト

### 2.1 TrackUseCase (track_test.go)

#### TestTrackUseCase_FetchByID

| No  | テストケース       | 入力                 | 期待結果         |
| --- | ------------------ | -------------------- | ---------------- |
| 1   | 正常系: 有効な ID  | `track123`           | トラック情報     |
| 2   | 異常系: 空の ID    | `""`                 | ErrTrackNotFound |
| 3   | 異常系: API エラー | 任意 ID + API エラー | エラー           |

#### TestTrackUseCase_Search

| No  | テストケース         | 入力                    | 期待結果        |
| --- | -------------------- | ----------------------- | --------------- |
| 1   | 正常系: 有効なクエリ | `test query`            | トラックリスト  |
| 2   | 正常系: 結果 0 件    | `nonexistent`           | 空リスト        |
| 3   | 異常系: 空のクエリ   | `""`                    | ErrInvalidQuery |
| 4   | 異常系: API エラー   | 任意クエリ + API エラー | エラー          |

### 2.2 ArtistUseCase (artist_test.go)

#### TestArtistUseCase_FetchByID

| No  | テストケース       | 入力                 | 期待結果          |
| --- | ------------------ | -------------------- | ----------------- |
| 1   | 正常系: 有効な ID  | `artist123`          | アーティスト情報  |
| 2   | 異常系: 空の ID    | `""`                 | ErrArtistNotFound |
| 3   | 異常系: API エラー | 任意 ID + API エラー | エラー            |

### 2.3 AlbumUseCase (album_test.go)

#### TestAlbumUseCase_FetchByID

| No  | テストケース       | 入力                 | 期待結果         |
| --- | ------------------ | -------------------- | ---------------- |
| 1   | 正常系: 有効な ID  | `album123`           | アルバム情報     |
| 2   | 異常系: 空の ID    | `""`                 | ErrAlbumNotFound |
| 3   | 異常系: API エラー | 任意 ID + API エラー | エラー           |

### 2.4 SimilarTracksUseCase (similar_tracks_test.go)

#### TestSimilarTracksUseCase_GetSimilarTracks

| No  | テストケース             | 入力                     | 期待結果           |
| --- | ------------------------ | ------------------------ | ------------------ |
| 1   | 正常系: ISRC あり        | trackID (ISRC 付き)      | 類似トラックリスト |
| 2   | 正常系: KKBOX 結果 0 件  | trackID                  | 空リスト           |
| 3   | 異常系: Spotify 取得失敗 | trackID + Spotify エラー | エラー             |
| 4   | 異常系: KKBOX 検索失敗   | trackID + KKBOX エラー   | エラー             |

### 2.5 RecommendUseCase V2 (usecase/v2/recommend_test.go)

#### TestRecommendUseCase_GetRecommendations

| No  | テストケース           | 入力                     | 期待結果                 |
| --- | ---------------------- | ------------------------ | ------------------------ |
| 1   | 正常系: 候補あり       | 有効な trackID           | レコメンドリスト         |
| 2   | 正常系: ISRC なし      | ISRC なしトラック        | 空リスト                 |
| 3   | 正常系: KKBOX 結果 nil | trackID + KKBOX nil 応答 | 空リスト（パニックなし） |

#### TestNewRecommendUseCaseWithLastFM / TestNewRecommendUseCaseFull

| No  | テストケース              | 期待結果               |
| --- | ------------------------- | ---------------------- |
| 1   | LastFM 付きコンストラクタ | lastfmAPI が設定される |
| 2   | Full コンストラクタ       | 全 API が設定される    |

#### TestRecommendUseCase_CollectFromLastFM

| No  | テストケース     | 入力                   | 期待結果                 |
| --- | ---------------- | ---------------------- | ------------------------ |
| 1   | 正常系           | 有効なトラック         | Last.fm 候補が収集される |
| 2   | アーティストなし | Artists が空のトラック | スキップ（エラーなし）   |
| 3   | API エラー       | Last.fm API エラー     | スキップ（エラーなし）   |

#### TestRecommendUseCase_CollectFromYouTubeMusic

| No  | テストケース     | 入力                    | 期待結果                 |
| --- | ---------------- | ----------------------- | ------------------------ |
| 1   | 正常系           | 有効なトラック          | YouTube Music 候補が収集 |
| 2   | 検索エラー       | SearchTracks エラー     | スキップ（エラーなし）   |
| 3   | 検索結果なし     | 空の検索結果            | スキップ（エラーなし）   |
| 4   | 類似曲取得エラー | GetSimilarTracks エラー | スキップ（エラーなし）   |

#### TestSearchSpotifyWithFallback

| No  | テストケース             | 入力                    | 期待結果     |
| --- | ------------------------ | ----------------------- | ------------ |
| 1   | 厳密検索成功             | 通常の曲名/アーティスト | トラック情報 |
| 2   | 検索結果なし             | 存在しない曲            | nil          |
| 3   | fuzzy アーティストマッチ | 部分一致アーティスト名  | トラック情報 |

#### TestSanitizeSearchQuery

| No  | テストケース     | 入力                       | 期待結果                |
| --- | ---------------- | -------------------------- | ----------------------- |
| 1   | 日本語括弧除去   | `聖戦と死神 第1部「銀色」` | `聖戦と死神 第1部 銀色` |
| 2   | 括弧除去         | `Track (Remix)`            | `Track Remix`           |
| 3   | 特殊記号除去     | `Track～Version`           | `Track Version`         |
| 4   | 複数スペース圧縮 | `Track    Name`            | `Track Name`            |

#### TestSimplifyTrackName

| No  | テストケース | 入力                         | 期待結果     |
| --- | ------------ | ---------------------------- | ------------ |
| 1   | feat. 除去   | `Track Name (feat. Artist)`  | `Track Name` |
| 2   | Remix 除去   | `Track Name - Remix Version` | `Track Name` |
| 3   | 短い名前保持 | `AB`                         | `AB`         |

#### TestFuzzyMatchArtist

| No  | テストケース     | 入力                               | 期待結果 |
| --- | ---------------- | ---------------------------------- | -------- |
| 1   | 完全一致         | `Artist Name`, `Artist Name`       | true     |
| 2   | 大文字小文字無視 | `ARTIST NAME`, `artist name`       | true     |
| 3   | 部分一致         | `The Artist`, `Artist`             | true     |
| 4   | & と and 正規化  | `Artist & Band`, `Artist and Band` | true     |
| 5   | 不一致           | `Artist A`, `Artist B`             | false    |

### 2.6 SimilarityCalculator V2 (usecase/v2/similarity_test.go)

#### TestSimilarityCalculator_Calculate

| No  | テストケース     | 入力                        | 期待結果 |
| --- | ---------------- | --------------------------- | -------- |
| 1   | nil seed         | seed=nil                    | 0.5      |
| 2   | nil candidate    | candidate=nil               | 0.5      |
| 3   | 同一特徴量       | 同じ BPM/Duration/Gain/Tags | ~1.0     |
| 4   | 類似 BPM         | BPM 差 ±5                   | 0.9-1.0  |
| 5   | 非常に異なる特徴 | 大きな BPM/Duration/Gain 差 | <0.7     |

#### TestSimilarityCalculator_hasVoiceActorRelation

| No  | テストケース     | 入力                  | 期待結果 |
| --- | ---------------- | --------------------- | -------- |
| 1   | nil アーティスト | nil, nil              | false    |
| 2   | 共通声優（MBID） | 同じ voice actor MBID | true     |
| 3   | 共通声優（名前） | 同じ voice actor 名   | true     |
| 4   | 声優関係なし     | 異なる voice actor    | false    |

#### TestSimilarityCalculator_hasCollaborationRelation

| No  | テストケース            | 入力                   | 期待結果 |
| --- | ----------------------- | ---------------------- | -------- |
| 1   | nil アーティスト        | nil, nil               | false    |
| 2   | A が B とコラボ（MBID） | collaboration relation | true     |
| 3   | B が A とコラボ（名前） | collaborator relation  | true     |
| 4   | コラボ関係なし          | 無関係                 | false    |

---

## 3. Gateway テスト

### 3.1 Spotify Gateway (gateway_test.go)

#### TestNewGateway

Gateway の初期化テスト

| No  | テストケース | 期待結果                                        |
| --- | ------------ | ----------------------------------------------- |
| 1   | 正常系       | clientID, secret, tokenRepo, httpc が正しく設定 |

#### TestRawTrack_ToDomain

Spotify API 応答からドメインモデルへの変換テスト

| No  | テストケース         | 入力                 | 期待結果            |
| --- | -------------------- | -------------------- | ------------------- |
| 1   | 正常系: 全フィールド | 完全な rawTrack      | 正しい domain.Track |
| 2   | 正常系: ISRC なし    | ISRC なしの rawTrack | ISRC が nil         |

#### TestGateway_GetTrackByID

| No  | テストケース       | モック応答                      | 期待結果         |
| --- | ------------------ | ------------------------------- | ---------------- |
| 1   | 正常系: 有効な ID  | 200 OK + トラック JSON          | トラック情報     |
| 2   | 異常系: 404        | 404 Not Found                   | ErrTrackNotFound |
| 3   | 異常系: 認証エラー | 401 Unauthorized + リトライ成功 | トラック情報     |

#### TestGateway_SearchTracks

| No  | テストケース     | モック応答             | 期待結果       |
| --- | ---------------- | ---------------------- | -------------- |
| 1   | 正常系: 結果あり | 200 OK + 検索結果 JSON | トラックリスト |
| 2   | 正常系: 結果なし | 200 OK + 空結果        | 空リスト       |

#### TestGateway_GetArtistByID

| No  | テストケース | モック応答                 | 期待結果          |
| --- | ------------ | -------------------------- | ----------------- |
| 1   | 正常系       | 200 OK + アーティスト JSON | アーティスト情報  |
| 2   | 異常系: 404  | 404 Not Found              | ErrArtistNotFound |

#### TestGateway_GetAlbumByID

| No  | テストケース | モック応答             | 期待結果         |
| --- | ------------ | ---------------------- | ---------------- |
| 1   | 正常系       | 200 OK + アルバム JSON | アルバム情報     |
| 2   | 異常系: 404  | 404 Not Found          | ErrAlbumNotFound |

### 3.2 KKBOX Gateway (gateway_test.go)

#### TestKKBOXGateway_SearchByISRC

| No  | テストケース              | モック応答                | 期待結果       |
| --- | ------------------------- | ------------------------- | -------------- |
| 1   | 正常系: ISRC 発見         | 200 OK + トラック JSON    | KKBOXTrackInfo |
| 2   | 正常系: ISRC 見つからない | 200 OK + 空結果           | nil            |
| 3   | 異常系: API エラー        | 500 Internal Server Error | エラー         |

#### TestKKBOXGateway_GetRecommendedTracks

| No  | テストケース     | モック応答               | 期待結果              |
| --- | ---------------- | ------------------------ | --------------------- |
| 1   | 正常系           | 200 OK + レコメンド JSON | KKBOXTrackInfo リスト |
| 2   | 正常系: 結果なし | 200 OK + 空結果          | 空リスト              |

#### TestKKBOXGateway_GetTrackDetail

| No  | テストケース | モック応答             | 期待結果       |
| --- | ------------ | ---------------------- | -------------- |
| 1   | 正常系       | 200 OK + トラック JSON | KKBOXTrackInfo |
| 2   | 異常系: 404  | 404 Not Found          | エラー         |

### 3.3 Cache Repository (repository_test.go)

#### TestCacheRepository_Get/Set

| No  | テストケース                | 操作                   | 期待結果             |
| --- | --------------------------- | ---------------------- | -------------------- |
| 1   | 正常系: L1 キャッシュヒット | Get (L1 存在)          | データ取得           |
| 2   | 正常系: L2 キャッシュヒット | Get (L1 なし, L2 存在) | データ取得 + L1 昇格 |
| 3   | 正常系: キャッシュミス      | Get (両方なし)         | nil                  |
| 4   | 正常系: 保存                | Set                    | L1・L2 両方に保存    |

---

## 4. Domain テスト

### 4.1 Errors (errors_test.go)

#### TestDomainErrors

| No  | テストケース      | 期待結果                    |
| --- | ----------------- | --------------------------- |
| 1   | ErrTrackNotFound  | エラーメッセージ確認        |
| 2   | ErrArtistNotFound | エラーメッセージ確認        |
| 3   | ErrAlbumNotFound  | エラーメッセージ確認        |
| 4   | ErrInvalidQuery   | エラーメッセージ確認        |
| 5   | エラー比較        | errors.Is()で正しく比較可能 |

---

## 5. Config テスト

### 5.1 Config (config_test.go)

#### TestLoadConfig

| No  | テストケース           | 環境変数               | 期待結果          |
| --- | ---------------------- | ---------------------- | ----------------- |
| 1   | 正常系: 全環境変数設定 | 全て設定               | Config 構造体取得 |
| 2   | 正常系: デフォルト値   | PORT なし              | 8080 がデフォルト |
| 3   | 異常系: 必須項目なし   | SPOTIFY_CLIENT_ID なし | エラー            |

---

## 6. ユーティリティテスト

### 6.1 Logger (logger_test.go)

#### TestLogger

| No  | テストケース | 操作       | 期待結果   |
| --- | ------------ | ---------- | ---------- |
| 1   | InitLogger   | 初期化     | エラーなし |
| 2   | Info         | Info 出力  | 正常出力   |
| 3   | Error        | Error 出力 | 正常出力   |
| 4   | Debug        | Debug 出力 | 正常出力   |

### 6.2 Mock (mock_test.go)

#### TestMockSpotifyAPI

| No  | テストケース                | 期待結果         |
| --- | --------------------------- | ---------------- |
| 1   | GetTrackByID - 関数設定あり | 設定した値を返す |
| 2   | GetTrackByID - 関数設定なし | nil, nil         |
| 3   | GetArtistByID               | 設定した値を返す |
| 4   | GetAlbumByID                | 設定した値を返す |
| 5   | SearchTracks                | 設定した値を返す |
| 6   | SearchByISRC                | 設定した値を返す |

#### TestMockKKBOXAPI

| No  | テストケース         | 期待結果         |
| --- | -------------------- | ---------------- |
| 1   | SearchByISRC         | 設定した値を返す |
| 2   | GetRecommendedTracks | 設定した値を返す |
| 3   | GetTrackDetail       | 設定した値を返す |

#### TestMockTokenRepository

| No  | テストケース | 期待結果       |
| --- | ------------ | -------------- |
| 1   | SaveToken    | 正常保存       |
| 2   | GetToken     | 保存した値取得 |
| 3   | IsTokenValid | 存在確認       |

#### TestHelper 関数

| No  | テストケース     | 期待結果             |
| --- | ---------------- | -------------------- |
| 1   | StringPtr        | 文字列ポインタ生成   |
| 2   | IntPtr           | 整数ポインタ生成     |
| 3   | CreateTestTrack  | テスト用 Track 生成  |
| 4   | CreateTestArtist | テスト用 Artist 生成 |
| 5   | CreateTestAlbum  | テスト用 Album 生成  |

---

## テスト実行方法

### 全テスト実行

```bash
go test ./...
```

### カバレッジ付き実行

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### 特定パッケージのテスト

```bash
go test ./internal/usecase/...
go test ./internal/adapter/handler/...
```

### 詳細出力

```bash
go test ./... -v
```

---

## 備考

### テスト除外対象

以下はインフラ層のため自動テストから除外:

- `cmd/server/main.go` - エントリーポイント
- `internal/adapter/server/server.go` - サーバー起動
- `internal/adapter/gateway/redis/repository.go` - Redis 接続

### モック戦略

- インターフェースベースのモック (`testutil.MockSpotifyAPI`, `testutil.MockKKBOXAPI`)
- httptest.Server による外部 API モック
- 依存性注入によるテスタビリティ確保
