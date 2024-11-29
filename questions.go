package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func formatQuestion(task string, resource string) string {
	return fmt.Sprintf("Do you want to %s the %s?", task, resource)
}

func readYesNoInput(reader *bufio.Reader) (bool, error) {
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input. Please try again.")
		return false, err
	}

	// Trim whitespace and convert to lowercase
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "yes" || input == "y" {
		return true, nil
	} else if input == "no" || input == "n" {
		return false, nil
	}

	fmt.Println("Invalid input. Please type 'yes' or 'no'.")
	return false, nil
}

func askYesNo(colorMode bool, flags *flags, question string, warning string) bool {
	if colorMode {
		red := color.New(color.FgRed)
		warning = red.Sprint(warning)
	}

	fullQuestion := question + " " + warning + " (yes/no): "

	if *flags.AnswerNo {
		fmt.Println(fullQuestion)
		return false
	}

	if *flags.AnswerYes {
		return true
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(fullQuestion)

		result, err := readYesNoInput(reader)

		if err != nil {
			continue
		}

		return result
	}
}
