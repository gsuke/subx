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
subx [オプション] <字幕ファイル...>

# 単一ファイル
subx . anime01.ass

# 複数ファイル（出力先フォルダ指定）
subx . *.srt -o ./extracted
```

## サンプルデータ

`/samples` を参照。

## テスト

```shell
# 自動テスト
go test -v

# 手動テスト
go run . # エラー(ヘルプ)
go run . ./samples/sample1-in.ass
go run . ./samples/sample1-in.ass -o extracted # 単一ファイルのフォルダ出力
go run . ./samples/sample1-in.ass ./samples/sample2-in.srt -o extracted # 複数ファイルのフォルダ出力
```

## 開発

新しい字幕形式に対応する場合は、 `extractor.go` の `SubtitlesExtractor` インターフェースを実装し、`DetectAndExtract` にそれを追加してください。
