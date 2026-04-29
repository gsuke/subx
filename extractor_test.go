package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// samplesフォルダ内の全サンプルについて変換結果を検証する
func TestDetectAndExtract_Samples(t *testing.T) {
	samplesDir := "samples"

	// 入力ファイルを検索（*-in.* パターン）
	entries, err := os.ReadDir(samplesDir)
	if err != nil {
		t.Fatalf("samplesフォルダの読み込みに失敗: %v", err)
	}

	for _, entry := range entries {
		name := entry.Name()

		// 入力ファイル（*-in.*）のみを対象
		if !strings.Contains(name, "-in.") {
			continue
		}

		t.Run(name, func(t *testing.T) {
			// 入力ファイルのパス
			inputPath := filepath.Join(samplesDir, name)

			// 期待出力ファイルのパスを生成（sample1-in.ass → sample1-out.txt）
			baseName := strings.Split(name, "-in.")[0]
			expectedPath := filepath.Join(samplesDir, baseName+"-out.txt")

			// 入力ファイルを読み込む
			inputContent, err := os.ReadFile(inputPath)
			if err != nil {
				t.Fatalf("入力ファイルの読み込みに失敗: %v", err)
			}

			// 期待出力ファイルを読み込む
			expectedContent, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Fatalf("期待出力ファイルの読み込みに失敗: %v", err)
			}

			// 変換を実行
			result, err := DetectAndExtract(string(inputContent))
			if err != nil {
				t.Fatalf("変換に失敗: %v", err)
			}

			// 結果を比較（末尾の空白・改行を正規化して比較）
			expected := strings.TrimSpace(string(expectedContent))
			actual := strings.TrimSpace(result)

			if actual != expected {
				t.Errorf("変換結果が期待値と一致しません\n期待:\n%s\n\n実際:\n%s", expected, actual)
			}
		})
	}
}

func TestDetectAndExtract(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
		wantErr bool
	}{
		{
			name:    "BOM除去_ASS",
			content: "\xef\xbb\xbf[Script Info]\n[Events]\nDialogue: 0,0:00:00.00,0:00:02.00,Default,,0,0,0,,テスト",
			want:    "テスト",
			wantErr: false,
		},
		{
			name:    "BOM除去_SRT",
			content: "\xef\xbb\xbf1\n00:00:00,000 --> 00:00:02,000\nテスト\n",
			want:    "テスト",
			wantErr: false,
		},
		{
			name:    "未対応形式",
			content: "これは何かのテキスト\n改行含む\nだが字幕形式ではない",
			want:    "",
			wantErr: true,
		},
		{
			name:    "空文字列",
			content: "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "連続重複_SRT",
			content: "1\n00:00:01,000 --> 00:00:02,000\nテスト\n\n2\n00:00:02,000 --> 00:00:03,000\nテスト\n",
			want:    "テスト",
			wantErr: false,
		},
		{
			name:    "連続重複_ASS",
			content: "[Script Info]\n[Events]\nDialogue: 0,0:00:00.00,0:00:02.00,Default,,0,0,0,,テスト\nDialogue: 0,0:00:02.00,0:00:04.00,Default,,0,0,0,,テスト\nDialogue: 0,0:00:04.00,0:00:06.00,Default,,0,0,0,,mage",
			want:    "テスト\nmage",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DetectAndExtract(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectAndExtract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DetectAndExtract() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDeduplicateLines(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  string
	}{
		{
			name:  "重複なし",
			lines: []string{"a", "b", "c"},
			want:  "a\nb\nc",
		},
		{
			name:  "連続重複",
			lines: []string{"a", "a", "b", "b", "b", "c"},
			want:  "a\nb\nc",
		},
		{
			name:  "すべて同一",
			lines: []string{"x", "x", "x"},
			want:  "x",
		},
		{
			name:  "空スライス",
			lines: []string{},
			want:  "",
		},
		{
			name:  "単一要素",
			lines: []string{"only"},
			want:  "only",
		},
		{
			name:  "非連続重複は除去しない",
			lines: []string{"a", "b", "a"},
			want:  "a\nb\na",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deduplicateLines(tt.lines)
			if got != tt.want {
				t.Errorf("deduplicateLines() = %q, want %q", got, tt.want)
			}
		})
	}
}

func equalStringSlice(a, b []string) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
