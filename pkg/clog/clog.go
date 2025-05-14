package clog

import (
	"fmt"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

// Banner prints a banner
func Banner(str string) {
	myFigure := figure.NewColorFigure(str, "", "green", true)
	myFigure.Print()
	fmt.Println()
}

// Line prints a line with a new line
func Line(str string) {
	color.Blue("=====\t%v\t=====\n", str)
}

// Info prints a string with a new line with time
func Info(format string, args ...any) {
	color.Green(fmt.Sprintf(format, args...))
}

// Warn prints a string with a new line with time
func Warn(format string, args ...any) {
	color.Yellow(fmt.Sprintf(format, args...))
}

// Error prints a string with a new line with time
func Error(format string, args ...any) {
	color.Red(fmt.Sprintf(format, args...))
}

func Panic(format string, args ...any) {
	color.Red(fmt.Sprintf(format, args...))
	os.Exit(1)
}

func Item(s string) {
	color.Blue(fmt.Sprintf("  - %s", s))
}

func Console(s string) {
	fmt.Print(s)
}
