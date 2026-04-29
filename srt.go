package main

import (
	"regexp"
	"strings"
)

// SRT形式の字幕を抽出する
type SRTExtractor struct{}

// SRT形式かどうかを判定する
func (e *SRTExtractor) CanExtract(content string) bool {
	// SRT形式の特徴: タイムスタンプ行（00:00:00,000 --> 00:00:00,000）が存在する
	timestampPattern := regexp.MustCompile(`\d+:\d{2}:\d{2},\d{3}\s*-->\s*\d+:\d{2}:\d{2},\d{3}`)
	return timestampPattern.MatchString(content)
}

// SRT形式の字幕からテキストを抽出する
//
// SRT形式はブロック構造を持つ:
//
//	1                   <- シーケンス番号
//	00:00:01,000 --> 00:00:02,000  <- タイムスタンプ
//	テキスト1行目        <- テキスト（複数行可）
//	テキスト2行目
//
//	2                   <- 次のブロック開始（空行で区切られる）
//	00:00:02,000 --> 00:00:03,000
//	次のテキスト
//
// この関数は各ブロックのテキスト行を1つのスライス要素にまとめ、[]stringで返す。
// ブロック内の複数行はLF(\n)で結合される。
func (e *SRTExtractor) Extract(content string) ([]string, error) {
	lines := strings.Split(content, "\n")

	var textParts []string   // 抽出したスライス（1ブロック = 1要素）
	var currentText []string // 現在のブロックで収集中のテキスト行

	// flushText: 現在のブロックのテキストをスライスに追加し、currentTextをクリアする
	flushText := func() {
		if len(currentText) == 0 {
			return
		}
		// ブロック内の複数行をLF(\n)で結合して1要素にする
		textParts = append(textParts, strings.Join(currentText, "\n"))
		currentText = nil
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 空行はブロック境界
		if line == "" {
			flushText()
			continue
		}

		// シーケンス番号行（数字のみ）: 新しいブロック開始
		if matched, _ := regexp.MatchString(`^\d+$`, line); matched {
			flushText()
			continue
		}

		// タイムスタンプ行: 新しいブロック開始
		if matched, _ := regexp.MatchString(`^\d+:\d{2}:\d{2},\d{3}\s*-->\s*\d+:\d{2}:\d{2},\d{3}$`, line); matched {
			flushText()
			continue
		}

		// ASS形式のメタデータが混入していることがあるので除去
		line = removeASSMetadata(line)

		// テキスト行: 現在のブロックに追加
		currentText = append(currentText, line)
	}

	// 最後のブロックをflush
	flushText()

	return textParts, nil
}
