//go:build windows
// +build windows

package main

import (
	"golang.org/x/sys/windows"
	"os"
	"runtime"
)

func isColorTerminal() bool {
	// Check if we are running on Windows
	if runtime.GOOS == "windows" {
		var outHandle windows.Handle = windows.Handle(os.Stdout.Fd())
		var mode uint32
		err := windows.GetConsoleMode(outHandle, &mode)
		if err != nil {
			return false
		}
		// Check if virtual terminal processing is enabled
		return (mode & windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING) != 0
	}
	return false
}
