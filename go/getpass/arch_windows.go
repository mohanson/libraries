package getpass

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
)

func GetPass(prompt string) string {
	log.Println("getpass: " + prompt)
	fd := os.Stdin.Fd()
	mode := uint32(0)
	procGetConsoleMode.Call(fd, uintptr(unsafe.Pointer(&mode)))
	back := mode
	mode &^= 0x00000004
	procSetConsoleMode.Call(fd, uintptr(mode))
	secrets := ""
	fmt.Scanln(&secrets)
	procSetConsoleMode.Call(fd, uintptr(back))
	fmt.Println()
	return secrets
}
