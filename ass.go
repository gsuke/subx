package main

import (
	"regexp"
	"strings"
)

// ASS形式の字幕を抽出する
type ASSExtractor struct{}

// ASS形式かどうかを判定する
func (e *ASSExtractor) CanExtract(content string) bool {
	return strings.Contains(content, "[Events]") &&
		strings.Contains(content, "Dialogue:")
}

// ASS形式の字幕からテキストを抽出する
func (e *ASSExtractor) Extract(content string) ([]string, error) {
	lines := strings.Split(content, "\n")

	var textParts []string
	for _, line := range lines {
		if strings.HasPrefix(line, "Dialogue: ") {
			text := extractTextFromDialogue(line)
			if text != "" {
				textParts = append(textParts, text)
			}
		}
	}

	return textParts, nil
}

// Dialogue行からテキスト部分(Text列)を抽出する
func extractTextFromDialogue(line string) string {
	// 行の書式: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
	parts := strings.SplitN(line, ",", 10)
	if len(parts) < 10 {
		return ""
	}

	text := parts[9]

	// 抽出したText列から、不要なものを除去
	text = removeASSMetadata(text)
	text = replaceNewlineCode(text)
	text = strings.TrimSpace(text)

	return text
}

// ASS形式のメタデータ（{\pos(...)}など）を除去する
// ※ SRT形式ファイルにも混入していることがある
func removeASSMetadata(text string) string {
	// {\...} 形式のタグを除去
	re := regexp.MustCompile(`\{[^}]*\}`)
	return re.ReplaceAllString(text, "")
}

// \N を改行に置換する
func replaceNewlineCode(text string) string {
	return strings.ReplaceAll(text, `\N`, "\n")
}
