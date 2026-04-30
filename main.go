package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag" // 引数とフラグの順番を自由にさせるため導入 (SetInterspersed)
)

func main() {
	// フラグの定義
	flag.SetInterspersed(true)
	outputDir := flag.StringP("outputdir", "o", "", "出力先フォルダ（入力ファイルが複数の場合は必須）")
	flag.Usage = printHelp
	flag.Parse()

	// 入力ファイルの取得
	inputFiles := flag.Args()

	// 入力ファイル数0のとき、終了
	if len(inputFiles) == 0 {
		printHelp()
		os.Exit(1)
	}

	// 入力ファイルが複数の場合は -o オプションが必須
	if len(inputFiles) > 1 && *outputDir == "" {
		fmt.Fprintln(os.Stderr, "エラー: 複数ファイル入力時は -o オプションで出力先フォルダを指定してください")
		os.Exit(1)
	}

	// 単一ファイルかつ -o 未指定なら標準出力して終了
	if len(inputFiles) == 1 && *outputDir == "" {
		result, err := processFile(inputFiles[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(result)
		return
	}

	// 複数ファイル処理: 出力先フォルダを作成
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: 出力先フォルダを作成できません: %v\n", err)
		os.Exit(1)
	}

	// 各ファイルを処理して出力フォルダに書き込み
	for _, inputFile := range inputFiles {
		result, err := processFile(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー [%s]: %v\n", inputFile, err)
			continue
		}

		outputFile := getOutputPath(inputFile, *outputDir)
		if err := os.WriteFile(outputFile, []byte(result+"\n"), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "エラー [%s]: 出力ファイルの書き込みに失敗: %v\n", inputFile, err)
			continue
		}
		fmt.Printf("%s -> %s\n", inputFile, outputFile)
	}
}

// ファイルを処理して抽出結果を返す
func processFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("ファイルを開けません: %v", err)
	}

	result, err := DetectAndExtract(string(content))
	if err != nil {
		return "", fmt.Errorf("処理に失敗しました: %v", err)
	}

	return result, nil
}

// 出力ファイルパスを生成する（拡張子を.txtに変更）
func getOutputPath(inputFile, outputDir string) string {
	base := filepath.Base(inputFile)
	ext := filepath.Ext(base)
	nameWithoutExt := strings.TrimSuffix(base, ext)
	return filepath.Join(outputDir, nameWithoutExt+".txt")
}

func printHelp() {
	fmt.Print(`subx - SUBtitles eXtractor

Usage:
  subx <file>
  subx <files...> -o <dir>

Examples:
  subx anime01.ass
  subx *.srt -o ./extracted

Supported formats:
  * .ass (Advanced SubStation Alpha)
  * .srt (SubRip Text)
  * .vtt (WebVTT)

Description:
  字幕ファイルからテキスト部分のみを抽出します。
  メタデータやASSタグを除去し、純粋なテキストを出力します。
`)
}
