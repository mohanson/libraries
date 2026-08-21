package getpass

import (
	"fmt"
	"math"
	"os"
	"syscall"
	"unsafe"

	"github.com/mohanson/libraries/go/lognoln"
)

func GetPass(prompt string) string {
	lognoln.Print("getpass: " + prompt)
	termios := syscall.Termios{}
	fd := os.Stdout.Fd()
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCGETA), uintptr(unsafe.Pointer(&termios)))
	termios.Lflag &= syscall.ECHO ^ math.MaxUint32
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(&termios)))
	secrets := ""
	fmt.Scanln(&secrets)
	fmt.Println()
	termios.Lflag |= syscall.ECHO
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(&termios)))
	return secrets
}
