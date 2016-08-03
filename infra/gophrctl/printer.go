package main

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	red   = color.New(color.FgRed)
	blue  = color.New(color.FgBlue)
	green = color.New(color.FgGreen)
)

func print(args ...interface{}) {
	fmt.Println(args...)
}

func printSuccess(args ...interface{}) {
	green.Print("✓ ")
	green.Println(args...)
}

func printInfo(args ...interface{}) {
	blue.Print("ℹ ")
	blue.Println(args...)
}

func printError(args ...interface{}) {
	red.Print("✗ ")
	red.Println(args...)
}
