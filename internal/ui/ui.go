package ui

import (
	"fmt"

	"github.com/fatih/color"
)

type ColorWriter func(a ...interface{}) string

type ColorWriters map[string]ColorWriter

var DefaultColorWriters = ColorWriters{
	"black":     color.New(color.FgBlack).SprintFunc(),
	"red":       color.New(color.FgRed).SprintFunc(),
	"green":     color.New(color.FgGreen).SprintFunc(),
	"yellow":    color.New(color.FgYellow).SprintFunc(),
	"blue":      color.New(color.FgBlue).SprintFunc(),
	"magenta":   color.New(color.FgMagenta).SprintFunc(),
	"cyan":      color.New(color.FgCyan).SprintFunc(),
	"white":     color.New(color.FgWhite).SprintFunc(),
	"hiBlack":   color.New(color.FgHiBlack).SprintFunc(),
	"hiRed":     color.New(color.FgHiRed).SprintFunc(),
	"hiGreen":   color.New(color.FgHiGreen).SprintFunc(),
	"hiYellow":  color.New(color.FgHiYellow).SprintFunc(),
	"hiBlue":    color.New(color.FgHiBlue).SprintFunc(),
	"hiMagenta": color.New(color.FgHiMagenta).SprintFunc(),
	"hiCyan":    color.New(color.FgHiCyan).SprintFunc(),
	"hiWhite":   color.New(color.FgHiWhite).SprintFunc(),
}

func PrintDashes() {
	fmt.Printf("—————————————————————————————————————————————————————————\n")
}

func PrintMsgWithStatus(status string, color string, format string, a ...any) {
	writer, ok := DefaultColorWriters[color]
	if !ok {
		fmt.Printf("%s — %s", status, fmt.Sprintf(format, a...))
		return
	}
	fmt.Printf("%s — %s", writer(status), fmt.Sprintf(format, a...))
}
