package getpass

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func sttyEchoNo() {
	termios := syscall.Termios{}
	fd := os.Stdout.Fd()
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCGETA), uintptr(unsafe.Pointer(&termios)))
	termios.Lflag &^= syscall.ECHO
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(&termios)))
}

func sttyEchoOn() {
	termios := syscall.Termios{}
	fd := os.Stdout.Fd()
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCGETA), uintptr(unsafe.Pointer(&termios)))
	termios.Lflag |= syscall.ECHO
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(&termios)))
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
