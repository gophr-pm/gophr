package main

import "fmt"

func startSpinner(message string) {
	fmt.Print(message + "...")
}

func stopSpinner(operationSuccessful bool) {
	if operationSuccessful {
		fmt.Println(" done.")
	} else {
		fmt.Println(" failed.")
	}
}
