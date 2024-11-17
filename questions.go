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

// Prompt the user for yes/no input
func askYesNo(answerYes bool, colorMode bool, question string, warning string) bool {
	if colorMode {
		red := color.New(color.FgRed)
		warning = red.Sprint(warning)
	}

	if !answerYes {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print(question + " " + warning + " (yes/no): ")
			input, err := reader.ReadString('\n')

			if err != nil {
				fmt.Println("Error reading input. Please try again.")
				continue
			}

			// Trim whitespace and convert to lowercase
			input = strings.TrimSpace(strings.ToLower(input))

			if input == "yes" || input == "y" {
				return true
			} else if input == "no" || input == "n" {
				return false
			} else {
				fmt.Println("Invalid input. Please type 'yes' or 'no'.")
			}
		}
	}

	return true
}
