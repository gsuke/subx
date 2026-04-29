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
func (e *SRTExtractor) Extract(content string) ([]string, error) {
	lines := strings.Split(content, "\n")
	timestampPattern := regexp.MustCompile(`^\d+:\d{2}:\d{2},\d{3}\s*-->\s*\d+:\d{2}:\d{2},\d{3}$`)
	sequencePattern := regexp.MustCompile(`^\d+$`)

	var textParts []string
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 空行、シーケンス番号、タイムスタンプ行はスキップ
		if line == "" || sequencePattern.MatchString(line) || timestampPattern.MatchString(line) {
			continue
		}

		// ASS形式のメタデータが混入していることがあるので、それを除去
		line = removeASSMetadata(line)

		// それ以外はテキスト行
		textParts = append(textParts, line)
	}

	return textParts, nil
}
