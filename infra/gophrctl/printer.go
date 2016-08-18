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
	green.Println(append([]interface{}{"✓"}, args...)...)
}

func printInfo(args ...interface{}) {
	blue.Println(append([]interface{}{"ℹ"}, args...)...)
}

func printError(args ...interface{}) {
	red.Println(append([]interface{}{"✗"}, args...)...)
}
