package getpass

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
)

func sttyEchoNo() {
	handle := os.Stdin.Fd()
	mode := uint32(0)
	procGetConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode)))
	mode &= 0xfffffffb
	procSetConsoleMode.Call(handle, uintptr(mode))
}

func sttyEchoOn() {
	handle := os.Stdin.Fd()
	mode := uint32(0)
	procGetConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode)))
	mode |= 0x00000004
	procSetConsoleMode.Call(handle, uintptr(mode))
}

func GetPass(prompt string) string {
	pass := ""
	fmt.Print(prompt)
	sttyEchoNo()
	fmt.Scanln(&pass)
	sttyEchoOn()
	fmt.Println()
	return pass
}
