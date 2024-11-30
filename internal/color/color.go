// Package color is for text color management
package color

// CheckColorTerminal determines if the terminal supports color
func CheckColorTerminal() bool {
	colorTerminal := isColorTerminal()
	return colorTerminal
}
