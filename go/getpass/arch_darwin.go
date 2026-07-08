package getpass

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"
)

func GetPass(prompt string) string {
	log.Println("getpass: " + prompt)
	termios := syscall.Termios{}
	fd := os.Stdout.Fd()
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCGETA), uintptr(unsafe.Pointer(&termios)))
	oldflag := termios.Lflag
	termios.Lflag &^= syscall.ECHO
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(&termios)))
	secrets := ""
	fmt.Scanln(&secrets)
	termios.Lflag = oldflag
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(&termios)))
	fmt.Println()
	return secrets
}
