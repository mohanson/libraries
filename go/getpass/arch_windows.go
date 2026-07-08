package getpass

import (
	"fmt"
	"log"
	"math"
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
	log.Println("getpass: " + prompt + string([]byte{0x1b, 0x37}))
	log.Writer().Write([]byte{0x1b, 0x38})
	fd := os.Stdin.Fd()
	mode := uint32(0)
	procGetConsoleMode.Call(fd, uintptr(unsafe.Pointer(&mode)))
	mode &= 0x00000004 ^ math.MaxUint32
	procSetConsoleMode.Call(fd, uintptr(mode))
	secrets := ""
	fmt.Scanln(&secrets)
	fmt.Println()
	mode |= 0x00000004
	procSetConsoleMode.Call(fd, uintptr(mode))
	return secrets
}
