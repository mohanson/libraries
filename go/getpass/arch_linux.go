package getpass

import (
	"fmt"
	"log"
	"math"
	"os"
	"syscall"
	"unsafe"
)

func GetPass(prompt string) string {
	log.Println("getpass: " + prompt + string([]byte{0x1b, 0x37}))
	log.Writer().Write([]byte{0x1b, 0x38})
	termios := syscall.Termios{}
	fd := os.Stdout.Fd()
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&termios)))
	termios.Lflag &= syscall.ECHO ^ math.MaxUint32
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&termios)))
	secrets := ""
	fmt.Scanln(&secrets)
	fmt.Println()
	termios.Lflag |= syscall.ECHO
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&termios)))
	return secrets
}
