package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth" // go get github.com/mattn/go-runewidth
	"golang.org/x/term"             // go get golang.org/x/term
)

// ANSIエスケープコードを除去するための正規表現
var ansiStripper = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// getTerminalWidth は現在のターミナル幅（列数）を返します。
// 取得できない場合はデフォルトの80を返します。
func getTerminalWidth() int {
	fd := int(os.Stdout.Fd()) // 標準出力のファイルディスクリプタ
	if !term.IsTerminal(fd) {
		// ターミナルでない場合はデフォルト幅を返す（例: パイプ経由の出力時など）
		return 80
	}

	width, _, err := term.GetSize(fd)
	if err != nil {
		// エラーが発生した場合のフォールバック
		fmt.Fprintf(os.Stderr, "Warning: Failed to get terminal size: %v. Using default width 80.\n", err)
		return 80
	}
	return width
}

// getDisplayWidth はANSIエスケープコードを除外した文字列の、ターミナル上での表示幅を返します。
// 全角文字を2文字分として正確にカウントします。
func getDisplayWidth(s string) int {
	// ANSIエスケープコードを除去
	plainString := ansiStripper.ReplaceAllString(s, "")
	// go-runewidth を使って表示幅を計算
	return runewidth.StringWidth(plainString)
}

// centerLine は与えられた文字列をターミナル幅に合わせて中央寄せします。
// 色情報（ANSIエスケープコード）は保持されます。
func centerLine(line string, termWidth int) string {
	displayWidth := getDisplayWidth(line) // 表示上の文字幅を取得
	if displayWidth >= termWidth {
		// 行がターミナル幅以上の場合、そのまま返す
		return line
	}

	// 左右に追加するパディング（空白）の数を計算
	padding := (termWidth - displayWidth) / 2
	// 計算したパディングの数だけ半角スペースを追加し、元の行と結合
	return strings.Repeat(" ", padding) + line
}