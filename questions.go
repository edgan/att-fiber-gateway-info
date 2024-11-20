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
func askYesNo(colorMode bool, flags *Flags, question string, warning string) bool {
	if colorMode {
		red := color.New(color.FgRed)
		warning = red.Sprint(warning)
	}

	if !*flags.AnswerYes {
		reader := bufio.NewReader(os.Stdin)
		for {
			fullQuestion := question + " " + warning + " (yes/no): "
			if !*flags.AnswerNo {
				fmt.Print(fullQuestion)
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
			} else {
				fmt.Println(fullQuestion)
				return false
			}
		}
	}

	return true
}
