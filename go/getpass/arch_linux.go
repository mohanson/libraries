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
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&termios)))
	termios.Lflag &= 0xfffffff7
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&termios)))
}

func sttyEchoOn() {
	termios := syscall.Termios{}
	fd := os.Stdout.Fd()
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&termios)))
	termios.Lflag |= 0x00000008
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&termios)))
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
