//go:build unix

package cstty

func EnableVirtualTerminalProcessing(fd uintptr) error {
	return nil
}
