package util

import (
	"fmt"
	"runtime"
)

const (
	TextBlack = iota + 30
	TextRed
	TextGreen
	TextYellow
	TextBlue
	TextMagenta
	TextCyan
	TextWhite
)

func Black(str interface{}) string {
	return textColor(TextBlack, fmt.Sprintf("%v", str))
}

func Red(str interface{}) string {
	return textColor(TextRed, fmt.Sprintf("%v", str))
}

func Green(str interface{}) string {
	return textColor(TextGreen, fmt.Sprintf("%v", str))
}

func Yellow(str interface{}) string {
	return textColor(TextYellow, fmt.Sprintf("%v", str))
}

func Blue(str interface{}) string {
	return textColor(TextBlue, fmt.Sprintf("%v", str))
}

func Magenta(str interface{}) string {
	return textColor(TextMagenta, fmt.Sprintf("%v", str))
}

func Cyan(str interface{}) string {
	return textColor(TextCyan, fmt.Sprintf("%v", str))
}

func White(str interface{}) string {
	return textColor(TextWhite, fmt.Sprintf("%v", str))
}

func textColor(color int, str string) string {
	if IsWindows() {
		return str
	}

	switch color {
	case TextBlack:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextBlack, str)
	case TextRed:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextRed, str)
	case TextGreen:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextGreen, str)
	case TextYellow:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextYellow, str)
	case TextBlue:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextBlue, str)
	case TextMagenta:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextMagenta, str)
	case TextCyan:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextCyan, str)
	case TextWhite:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextWhite, str)
	default:
		return str
	}
}

func IsWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	} else {
		return false
	}
}
