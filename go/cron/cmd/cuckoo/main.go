package main

import (
	"log"
	"time"

	"github.com/mohanson/libraries/go/cron"
)

func main() {
	for range cron.Cron(time.Second*10, time.Second*5) {
		log.Println("main: cuckoo")
	}
}
