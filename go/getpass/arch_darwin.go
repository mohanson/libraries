package getpass

import (
	"fmt"
)

func GetPass(prompt string) string {
	pass := ""
	fmt.Print(prompt)
	fmt.Scanln(&pass)
	return pass
}
