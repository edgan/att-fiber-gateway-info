package color

// checkColorTerminal determines if the terminal supports color
func CheckColorTerminal() bool {
	colorTerminal := isColorTerminal()
	return colorTerminal
}
