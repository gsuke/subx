package main

import (
	"fmt"
	"strings"
)

// 字幕抽出のインターフェース
type SubtitlesExtractor interface {
	Extract(content string) ([]string, error)
	CanExtract(content string) bool
}

// 字幕形式を自動判別して抽出する
func DetectAndExtract(content string) (string, error) {
	// BOMを除去
	content = strings.TrimPrefix(content, "\xef\xbb\xbf")

	extractors := []SubtitlesExtractor{
		&ASSExtractor{},
		&SRTExtractor{},
	}

	for _, extractor := range extractors {
		if !extractor.CanExtract(content) {
			continue
		}

		parts, err := extractor.Extract(content)
		if err != nil {
			return "", err
		}
		return deduplicateLines(parts), nil
	}

	return "", fmt.Errorf("未対応の字幕形式です")
}

// 連続重複行をひとつにまとめる
func deduplicateLines(lines []string) string {
	if len(lines) == 0 {
		return ""
	}

	var result []string
	result = append(result, lines[0])

	for i := 1; i < len(lines); i++ {
		if lines[i] != lines[i-1] {
			result = append(result, lines[i])
		}
	}

	return strings.Join(result, "\n")
}
