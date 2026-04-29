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
		{
			name:    "フィールド不足_Dialogue行",
			content: "Dialogue: 0,0:00:00.00,0:00:02.00,Default,,0,0,0", // 10フィールド未満
			want:    "",
		},
		{
			name:    "メタデータのみ空文字化",
			content: "Dialogue: 0,0:00:00.00,0:00:02.00,Default,,0,0,0,,{\\fad(100,200)}",
			want:    "",
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

func TestExtractTextFromDialogue(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{
			name: "通常10フィールド",
			line: "Dialogue: 0,0:00:00.00,0:00:02.00,Default,,0,0,0,,テキスト部分",
			want: "テキスト部分",
		},
		{
			name: "9フィールドのみ",
			line: "Dialogue: 0,0:00:00.00,0:00:02.00,Default,,0,0,0", // Text列なし
			want: "",
		},
		{
			name: "Text部分にカンマ含む",
			line: "Dialogue: 0,0:00:00.00,0:00:02.00,Default,,0,0,0,,テキスト,含有,カンマ",
			want: "テキスト,含有,カンマ",
		},
		{
			name: "Text部分空",
			line: "Dialogue: 0,0:00:00.00,0:00:02.00,Default,,0,0,0,,",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTextFromDialogue(tt.line)
			if got != tt.want {
				t.Errorf("extractTextFromDialogue() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRemoveASSMetadata(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "posタグ除去",
			input: "{\\pos(100,200)}テキスト",
			want:  "テキスト",
		},
		{
			name:  "fadeタグ除去",
			input: "テキスト{\\fad(100,200)}続き",
			want:  "テキスト続き",
		},
		{
			name:  "タグなし",
			input: "タグなしテキスト",
			want:  "タグなしテキスト",
		},
		{
			name:  "空タグ除去",
			input: "テキスト{}続き",
			want:  "テキスト続き", // {}は{}全体がマッチし除去される
		},
		{
			name:  "複数タグ",
			input: "{\\pos(100,200)}先頭{\\fad(100,200)}末尾",
			want:  "先頭末尾",
		},
		{
			name:  "括弧開きのみ",
			input: "テキスト{途中",
			want:  "テキスト{途中", // }がないため除去されない
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeASSMetadata(tt.input)
			if got != tt.want {
				t.Errorf("removeASSMetadata() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReplaceNewlineCode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "単一置換",
			input: "一行\\N二行",
			want:  "一行\n二行",
		},
		{
			name:  "複数置換",
			input: "1行\\N2行\\N3行",
			want:  "1行\n2行\n3行",
		},
		{
			name:  "タグなし置換なし",
			input: "タグなし",
			want:  "タグなし",
		},
		{
			name:  "空文字列",
			input: "",
			want:  "",
		},
		{
			name:  "\\n小文字は無視",
			input: "改行\\n確認", // \n (小文字) は置換されない
			want:  "改行\\n確認",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := replaceNewlineCode(tt.input)
			if got != tt.want {
				t.Errorf("replaceNewlineCode() = %q, want %q", got, tt.want)
			}
		})
	}
}
