# Subtitles Extractor

字幕ファイルからテキスト部分のみを抽出するCLIスクリプト。

対応形式

* ASS
* SRT

## 使い方

```shell
go run . [オプション] <字幕ファイル...>

# 単一ファイル（標準出力 + クリップボード出力）
go run . anime01.ass

# 複数ファイル（出力先フォルダ指定）
go run . *.srt -o ./extracted
```

単一ファイルを指定した場合は、クリップボードに、LLMに渡すためのプロンプト込みで出力されます。
そのままLLMに貼り付けて与えることで、字幕の要約を得ることができます。

## サンプルデータ

`/samples` を参照。

## テスト

```shell
# 自動テスト
go test -v

# 手動テスト
go run . # エラー(ヘルプ)
go run . ./samples/sample1-in.ass # 標準出力 + クリップボード出力
go run . ./samples/sample1-in.ass -o extracted # 単一ファイルのフォルダ出力
go run . ./samples/sample1-in.ass ./samples/sample2-in.srt -o extracted # 複数ファイルのフォルダ出力
```

## 開発

新しい字幕形式に対応する場合は、 `extractor.go` の `SubtitlesExtractor` インターフェースを実装し、`DetectAndExtract` にそれを追加してください。

## Claude Code によるストーリー説明

Claude Code の Skills 機能を使って、字幕ファイルのストーリーを説明させることができます。

```
/explaining-story
```

* 入力フォルダ: `/extracted`
* 出力フォルダ: `/explained`

## サブコマンド

### concat

複数のテキストファイルを連結して、Markdown形式で出力できます。

```
go run ./cmd/concat ./cmd/concat/samples/*
```
