// Package ui provides colored terminal output helpers for AryFlow CLI.
package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorBold   = "\033[1m"
)

// noColor returns true if colored output should be disabled.
func noColor() bool {
	return os.Getenv("ARYFLOW_NO_COLOR") == "1" || os.Getenv("NO_COLOR") != ""
}

// autoYes returns true if prompts should be auto-accepted.
func autoYes() bool {
	return os.Getenv("ARYFLOW_YES") == "1"
}

// colorize wraps text in ANSI color codes if colors are enabled.
func colorize(color, text string) string {
	if noColor() {
		return text
	}
	return color + text + colorReset
}

// Success prints a green checkmark message.
func Success(msg string) {
	fmt.Printf("  %s %s\n", colorize(colorGreen, "\u2713"), msg)
}

// Error prints a red X message.
func Error(msg string) {
	fmt.Printf("  %s %s\n", colorize(colorRed, "\u2717"), msg)
}

// Warning prints a yellow warning message.
func Warning(msg string) {
	fmt.Printf("  %s %s\n", colorize(colorYellow, "\u26a0"), msg)
}

// Info prints a blue info message.
func Info(msg string) {
	fmt.Printf("  %s %s\n", colorize(colorBlue, "i"), msg)
}

// Prompt asks a Y/n question and returns true if the user accepts.
// If ARYFLOW_YES=1, returns true without asking.
func Prompt(msg string) bool {
	if autoYes() {
		return true
	}

	fmt.Printf("    %s [Y/n] ", msg)
	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "" || answer == "y" || answer == "yes"
}

// PromptDefaultNo asks a y/N question and returns true only if user types y.
// If ARYFLOW_YES=1, returns true without asking.
func PromptDefaultNo(msg string) bool {
	if autoYes() {
		return true
	}

	fmt.Printf("    %s [y/N] ", msg)
	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}

// Header prints a bold section header.
func Header(msg string) {
	fmt.Println()
	fmt.Println(colorize(colorBold, msg))
	fmt.Println()
}

// Suggestion prints an indented suggestion line.
func Suggestion(msg string) {
	fmt.Printf("    %s %s\n", colorize(colorBlue, "\u2192"), msg)
}

// FormatSuccess returns a formatted success string without printing.
func FormatSuccess(msg string) string {
	return fmt.Sprintf("  %s %s", colorize(colorGreen, "\u2713"), msg)
}

// FormatError returns a formatted error string without printing.
func FormatError(msg string) string {
	return fmt.Sprintf("  %s %s", colorize(colorRed, "\u2717"), msg)
}

// FormatWarning returns a formatted warning string without printing.
func FormatWarning(msg string) string {
	return fmt.Sprintf("  %s %s", colorize(colorYellow, "\u26a0"), msg)
}
