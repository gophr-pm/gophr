package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// promptChoiceArgs is the arguments struct for promptChoice.
type promptChoiceArgs struct {
	prompt             string
	choice             string
	options            []string
	defaultOptionIndex int
}

// promptChoice prompts the user to choose from a list of options. After, the
// index of the selected option is returned.
func promptChoice(args promptChoiceArgs) int {
	for {
		fmt.Println(args.prompt)
		for i, option := range args.options {
			yellow.Printf("(%d) ", i+1)
			fmt.Println(option)
		}
		fmt.Println()
		fmt.Printf("%s [%d]: ", args.choice, args.defaultOptionIndex+1)

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()

		if len(text) == 0 {
			return args.defaultOptionIndex
		} else if answer, err := strconv.Atoi(text); err != nil || answer < 1 || answer > len(args.options) {
			printError(fmt.Sprintf("\"%s\" is not a valid choice. Please try again.\n\n", text))
		} else {
			return answer - 1
		}
	}
}
