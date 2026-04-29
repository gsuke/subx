# subx - SUBtitles eXtractor

字幕ファイルからテキスト部分のみを抽出するCLIスクリプト。

対応形式

* ASS
* SRT

## インストール

```
go install github.com/gsuke/subx@latest
```

## 使い方

```shell
# 単一ファイル
subx anime01.ass

# 複数ファイル（出力先フォルダ指定）
subx *.srt -o ./extracted
```

## サンプルデータ & テストケース

`/samples` を参照。

## テスト

自動テスト

```shell
go test -v
```

手動テスト

```shell
go run . # エラー(ヘルプ)
go run . ./samples/sample1-in.ass # 単一ファイル
go run . ./samples/*-in.* -o extracted # 複数ファイル
```

## 開発

新しい字幕形式に対応する場合は、 `extractor.go` の `SubtitlesExtractor` インターフェースを実装し、`DetectAndExtract` にそれを追加してください。
