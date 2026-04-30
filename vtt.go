package main

import (
	"regexp"
	"strings"
)

// WebVTT形式の字幕を抽出する
type VTTExtractor struct{}

// VTT形式かどうかを判定する
func (e *VTTExtractor) CanExtract(content string) bool {
	return strings.Contains(content, "WEBVTT")
}

// VTT形式の字幕からテキストを抽出する
// SRTに変換してからSRTExtractorを使用
func (e *VTTExtractor) Extract(content string) ([]string, error) {
	// VTT -> SRT変換: .→,に変換、空行そのまま
	srtContent := convertVTTToSRT(content)

	// SRTとして処理
	srtExtractor := &SRTExtractor{}
	return srtExtractor.Extract(srtContent)
}

// VTTをSRTに変換する
// * WEBVTTヘッダー、X-TIMESTAMP-MAPを読み捨てる
// * タイムスタンプ: .→, に変換
func convertVTTToSRT(vttContent string) string {
	lines := strings.Split(vttContent, "\n")

	var out []string
	started := false

	for _, line := range lines {
		// シーケンス「1」が来るまで読み捨てる
		if !started {
			if matched, _ := regexp.MatchString(`^1$`, line); matched {
				started = true
			}
			continue
		}

		// タイムスタンプ行: .→, に変換
		if matched, _ := regexp.MatchString(`^\d+:\d+:\d{2}\.\d{3}\s*-->\s*\d+:\d+:\d{2}\.\d{3}$`, line); matched {
			line = strings.ReplaceAll(line, ".", ",")
		}

		out = append(out, line)
	}

	return strings.Join(out, "\n")
}
