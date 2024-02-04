package helpers

import (
	"os"
	"strings"

	"golang.org/x/sys/windows"
)

const (
	RESET     = "\033[0m"
	BOLD      = "\033[1m"
	UNDERLINE = "\033[4m"
	STRIKE    = "\033[9m"
	ITALIC    = "\033[3m"
)

const (
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
	BLUE   = "\033[34m"
	PURPLE = "\033[35m"
	CYAN   = "\033[36m"
	GRAY   = "\033[37m"
	WHITE  = "\033[37m"
)

func init() {
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32

	windows.GetConsoleMode(stdout, &originalMode)
	windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}

func Format(text string, styles ...string) string {
	return strings.Join(styles, "") + text + RESET
}

func ErrorF(err string) string {
	if err != "" {
		return Format("Failed: ", RED) + Format(err, ITALIC)
	}
	return Format("FAILED", RED)
}

func SuccessF(msg string) string {
	return Format(msg, GREEN, BOLD)
}
