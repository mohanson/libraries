package getpass

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func GetPass(prompt string) string {
	fmt.Print(prompt)
	termios := syscall.Termios{}
	fd := os.Stdout.Fd()
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&termios)))
	oldflag := termios.Lflag
	termios.Lflag &^= syscall.ECHO
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&termios)))
	secrets := ""
	fmt.Scanln(&secrets)
	termios.Lflag = oldflag
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&termios)))
	fmt.Println()
	return secrets
}
