package getpass

import (
	"bytes"
	"log"
)

func lognoln(s string) {
	w := new(bytes.Buffer)
	log.New(w, "", log.Flags()).Println(s)
	log.Default().Writer().Write(w.Bytes()[:len(w.Bytes())-1])
}
