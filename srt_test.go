package main

import (
	"testing"
)

func TestSRTExtractor_CanExtract(t *testing.T) {
	e := &SRTExtractor{}

	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "有効なSRT形式_2桁時間",
			content: "1\n00:00:01,000 --> 00:00:02,000\nテスト\n",
			want:    true,
		},
		{
			name:    "有効なSRT形式_1桁時間",
			content: "1\n0:00:01,000 --> 0:00:02,000\nテスト\n",
			want:    true,
		},
		{
			name:    "1桁時間のみ_9時間台",
			content: "1\n9:59:59,000 --> 9:59:59,999\nテスト\n",
			want:    true,
		},
		{
			name:    "タイムスタンプなし",
			content: "1\nテスト字幕\n",
			want:    false,
		},
		{
			name:    "空文字列",
			content: "",
			want:    false,
		},
		{
			name:    "時刻だけだが矢印なし",
			content: "00:00:01,000\nテスト\n",
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

func TestSRTExtractor_Extract(t *testing.T) {
	e := &SRTExtractor{}

	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "通常のSRT",
			content: "1\n00:00:01,000 --> 00:00:02,000\nテスト字幕\n",
			want:    []string{"テスト字幕"},
		},
		{
			name:    "空の字幕ブロック",
			content: "1\n00:00:01,000 --> 00:00:02,000\n\n",
			want:    nil, // 空行はスキップされる
		},
		{
			name:    "シーケンス番号のみ",
			content: "1\n",
			want:    nil,
		},
		{
			name:    "ASSタグ混入",
			content: "1\n00:00:01,000 --> 00:00:02,000\n{\\pos(100,200)}テスト\\N二行目\n",
			want:    []string{"テスト\\N二行目"},
		},
		{
			name:    "複数ブロック",
			content: "1\n00:00:01,000 --> 00:00:02,000\n一行目\n\n2\n00:00:02,000 --> 00:00:03,000\n二行目\n",
			want:    []string{"一行目", "二行目"},
		},
		{
			name:    "テキスト中の数字のみ行は無視_SRT形式として不正",
			content: "1\n00:00:01,000 --> 00:00:02,000\n12345\n",
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := e.Extract(tt.content)
			if err != nil {
				t.Fatalf("Extract() error = %v", err)
			}
			if !equalStringSlice(got, tt.want) {
				t.Errorf("Extract() = %v, want %v", got, tt.want)
			}
		})
	}
}