package main

import (
	"testing"
)

func TestASSExtractor_CanExtract(t *testing.T) {
	e := &ASSExtractor{}

	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "有効なASS形式",
			content: "[Script Info]\n[Events]\nDialogue: 0,0:00:00.00,0:00:02.00,Default,,0,0,0,,テスト",
			want:    true,
		},
		{
			name:    "Dialogueがない",
			content: "[Script Info]\n[Events]",
			want:    false,
		},
		{
			name:    "[Events]がない",
			content: "[Script Info]\nDialogue: 0,0:00:00.00,0:00:02.00,Default,,0,0,0,,テスト",
			want:    false,
		},
		{
			name:    "空文字列",
			content: "",
			want:    false,
		},
		{
			name:    "小文字events",
			content: "[script info]\n[events]\nDialogue: 0,0:00:00.00,0:00:02.00,Default,,0,0,0,,テスト",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := e.CanExtract(tt.content)
			if got != tt.want {
				t.Errorf("CanExtract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASSExtractor_Extract(t *testing.T) {
	e := &ASSExtractor{}

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "通常のDialogue行",
			content: "Dialogue: 0,0:00:00.00,0:00:02.00,Default,,0,0,0,,テスト字幕",
			want:    "テスト字幕",
		},
		{
			name:    "空のDialogue行",
			content: "Dialogue: 0,0:00:00.00,0:00:02.00,Default,,0,0,0,,",
			want:    "",
		},
		{
			name:    "ASSタグを含む",
			content: "Dialogue: 0,0:00:00.00,0:00:02.00,Default,,0,0,0,,{\\pos(100,200)}テスト",
			want:    "テスト",
		},
		{
			name:    "文中に改行コードを含む\\N",
			content: "Dialogue: 0,0:00:00.00,0:00:02.00,Default,,0,0,0,,一行目\\N二行目",
			want:    "一行目\n二行目",
		},
		{
			name:    "末尾に改行コードを含む\\N",
			content: "Dialogue: 0,0:00:00.00,0:00:02.00,Default,,0,0,0,,テスト\\N",
			want:    "テスト",
		},
		{
			name:    "Dialogue行がない",
			content: "[Script Info]\n[Events]",
			want:    "",
		},
		{
			name:    "複数Dialogue行",
			content: "Dialogue: 0,0:00:00.00,0:00:02.00,Default,,0,0,0,,最初の行\nDialogue: 0,0:00:02.00,0:00:04.00,Default,,0,0,0,,次の行",
			want:    "最初の行\n次の行",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := e.Extract(tt.content)
			if err != nil {
				t.Fatalf("Extract() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Extract() = %q, want %q", got, tt.want)
			}
		})
	}
}
