package utils

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

func PrintFancyText(s string) {
	figure.NewColorFigure(s, "basic", "red", true).Print()
}

var Red = color.New(color.FgRed).SprintFunc()
var Green = color.New(color.FgGreen).SprintFunc()
var Yellow = color.New(color.FgYellow).SprintFunc()
