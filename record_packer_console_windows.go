//go:build windows

package rainbowlog

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	consoleHandle      = syscall.Handle(os.Stdout.Fd())
)

const (
	enableVirtualTerminalProcessing = 0x0004
)

func init() {
	// 尝试启用虚拟终端处理（ANSI 转义序列支持）
	enableVirtualTerminal()
}

func enableVirtualTerminal() {
	var mode uint32
	procGetConsoleMode.Call(uintptr(consoleHandle), uintptr(unsafe.Pointer(&mode)))
	procSetConsoleMode.Call(uintptr(consoleHandle), uintptr(mode|enableVirtualTerminalProcessing))
}
