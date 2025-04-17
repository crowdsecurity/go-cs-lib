package cstty

import (
	isatty "github.com/mattn/go-isatty"
)

// IsTTY returns true if the given file is an interactive terminal.
func IsTTY(fd uintptr) bool {
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}
