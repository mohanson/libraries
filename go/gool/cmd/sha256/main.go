package main

import (
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math/rand/v2"
	"runtime"
	"time"

	"github.com/mohanson/libraries/go/gool"
)

func once() {
	buffer := make([]byte, 256)
	hasher := sha256.New()
	for range 1024 {
		for i := range 32 {
			binary.LittleEndian.PutUint64(buffer[i*8:i*8+8], rand.Uint64())
		}
		hasher.Write(buffer)
		hasher.Sum(nil)
		hasher.Reset()
	}
}

// Execution on different CPUs.
//
// Intel(R) Core(TM) m3-7Y30 --- 590976
// Intel(R) Xeon(R) Gold 6133 -- 706816
// Intel(R) Core(TM) i7-9700 -- 1020672
// AMD EPYC 7K62 -------------- 1861888
func mainLoop() int {
	done := 0
	time.AfterFunc(time.Second, func() {
		done += 1
	})
	cnts := 0
	for done != 1 {
		once()
		cnts += 1
	}
	rate := cnts * 1024
	return rate
}

func mainGool() int {
	done := 0
	time.AfterFunc(time.Second, func() {
		done += 1
	})
	cnts := 0
	grun := gool.Cpu()
	for done != 1 {
		grun.Call(func() {
			once()
			grun.Lock(func() {
				cnts += 1
			})
		})
	}
	grun.Wait()
	rate := cnts * 1024
	return rate
}

func main() {
	log.Println("main:", runtime.NumCPU(), "logical cpus usable by the current process")
	log.Println("main: sha256 by loop")
	log.Println("main: sha256 by loop rate", mainLoop())
	log.Println("main: sha256 by gool")
	log.Println("main: sha256 by gool rate", mainGool())
}
