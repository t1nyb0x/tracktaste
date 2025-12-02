# Changelog

## [2.0.1](https://github.com/t1nyb0x/tracktaste/compare/v2.0.0...v2.0.1) (2025-12-02)


### Bug Fixes

* 検索結果が見つからない場合のログ出力を改善 ([4d80f9e](https://github.com/t1nyb0x/tracktaste/commit/4d80f9ef9fa7c92cd3daab2c892f4eaaa00e3213))
* 検索結果が見つからない場合のログ出力を改善 ([5eec1c0](https://github.com/t1nyb0x/tracktaste/commit/5eec1c0828709aa81edecb6d436f04fb151acef7))

## [2.0.0](https://github.com/t1nyb0x/tracktaste/compare/v1.4.0...v2.0.0) (2025-12-01)


### ⚠ BREAKING CHANGES

* レコメンドのエンドポイントをv2に変更

### Features

* Deezer, MusicBrainzから楽曲の傾向情報を取得するレコメンドエンジンを追加 ([64b3d70](https://github.com/t1nyb0x/tracktaste/commit/64b3d7025bf8144efd427149d1faf1a69778d58c))
* last.fm, Youtube musicから情報を取得する処理を追加 ([2b17df6](https://github.com/t1nyb0x/tracktaste/commit/2b17df6bf24ea15a3a54d25974e25748d8a2c603))
* レコメンドのエンドポイントをv2に変更 ([759499e](https://github.com/t1nyb0x/tracktaste/commit/759499ea5e23ba80d6529d0d20e6b2990dcb10fe))


### Bug Fixes

* Lintエラー解消 ([4a8d29f](https://github.com/t1nyb0x/tracktaste/commit/4a8d29f028c2557cae0a6262a6ad006087116a82))
* 一部情報が欠落していたのを修正 ([556af79](https://github.com/t1nyb0x/tracktaste/commit/556af7966428da2686cbe8e6652360e346398c08))

## [1.4.0](https://github.com/t1nyb0x/tracktaste/compare/v1.3.1...v1.4.0) (2025-11-30)


### Features

* レコメンドエンジン実装 ([ac0d067](https://github.com/t1nyb0x/tracktaste/commit/ac0d0674af9a750eecc0ec7b92bc3182156bd67b))

## [1.3.1](https://github.com/t1nyb0x/tracktaste/compare/v1.3.0...v1.3.1) (2025-11-30)


### Bug Fixes

* トークンが古い場合に再取得するように修正 ([2b853f4](https://github.com/t1nyb0x/tracktaste/commit/2b853f4ebd224021ddce879615efff084f73b786))
* トークンが古い場合に再取得するように修正 ([06badd0](https://github.com/t1nyb0x/tracktaste/commit/06badd0abbfa051c72bdb12700df91bbb5af78dd))
* トークンが古い場合に再取得するように修正 ([267e25a](https://github.com/t1nyb0x/tracktaste/commit/267e25a9686da76c6d979c863a20b8df6bde5c98))

## [1.3.0](https://github.com/t1nyb0x/tracktaste/compare/v1.2.1...v1.3.0) (2025-11-30)


### Features

* 類似トラック取得にexplicitとduration_msを追加 ([ea79791](https://github.com/t1nyb0x/tracktaste/commit/ea79791424ebab0424fdc151806aa1789f317b5e))
* 類似トラック取得にexplicitとduration_msを追加 ([2574308](https://github.com/t1nyb0x/tracktaste/commit/25743089f4320e8395091d7a1cf20d7a8e5fd1d2))

## [1.2.1](https://github.com/t1nyb0x/tracktaste/compare/v1.2.0...v1.2.1) (2025-11-30)


### Bug Fixes

* 検索APIのレスポンスで検索結果をresult.items[]の中に入れるように修正 ([ecfb691](https://github.com/t1nyb0x/tracktaste/commit/ecfb69139c4ae1c329e308e97883b917dd746b82))
* 検索APIのレスポンスで検索結果をresult.items[]の中に入れるように修正 ([5c4dee7](https://github.com/t1nyb0x/tracktaste/commit/5c4dee7aae7570e00038e7178275bc99d25cd039))

## [1.2.0](https://github.com/t1nyb0x/tracktaste/compare/v1.1.0...v1.2.0) (2025-11-30)


### Features

* トラック情報の返却にduration_msを追加 ([0996957](https://github.com/t1nyb0x/tracktaste/commit/09969577c85478748b4df38bbfc08fd034617c91))
* トラック情報の返却にduration_msを追加 ([dbc05db](https://github.com/t1nyb0x/tracktaste/commit/dbc05dbe44ea8364cf4aef2fbab8c333bc10d94f))

## [1.1.0](https://github.com/t1nyb0x/tracktaste/compare/v1.0.0...v1.1.0) (2025-11-30)


### Features

* トラック情報の返却にexplicit情報を追加 ([03c84cd](https://github.com/t1nyb0x/tracktaste/commit/03c84cd10a56da7f5ea1c2600e3218ca3bcf38a3))
* トラック情報の返却にexplicit情報を追加 ([b1255f7](https://github.com/t1nyb0x/tracktaste/commit/b1255f7d236dc65d35573316c599b643e54e78dd))

## 1.0.0 (2025-11-27)


### Features

* キャッシュを2層式に ([cb23082](https://github.com/t1nyb0x/tracktaste/commit/cb230824dbe04df4ab9358392c22377153539a6c))
* トラック情報、アーティスト情報、アルバム情報、類似トラック情報取得API実装 ([30f5c12](https://github.com/t1nyb0x/tracktaste/commit/30f5c12d937c108445b0bb9a197fefdfd700b42f))


### Bug Fixes

* 実装ミス修正 ([23f33aa](https://github.com/t1nyb0x/tracktaste/commit/23f33aaf48085af771155300e9275881d4e08c77))
