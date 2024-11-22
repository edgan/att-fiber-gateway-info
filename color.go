package main

// checkColorTerminal determines if the terminal supports color
func checkColorTerminal() bool {
	colorTerminal := isColorTerminal()
	return colorTerminal
}
