package cstty

import (
	"golang.org/x/sys/windows"
)

// EnableVirtualTerminalProcessing enables ANSI sequences on a given file descriptor.
// This only works on Windows 10+ but we can't do anything about older versions.
func EnableVirtualTerminalProcessing(fd uintptr) error {
	var mode uint32
	handle := windows.Handle(fd)
	if err := windows.GetConsoleMode(handle, &mode); err != nil {
		return err
	}
	mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	return windows.SetConsoleMode(handle, mode)
}
