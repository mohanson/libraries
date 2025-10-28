package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand/v2"
	"runtime"

	"github.com/mohanson/libraries/go/gool"
)

var aimpub map[string]uint8 = map[string]uint8{
	"7e8c088760bfde1dddcf32c17f209b8242ee52aaf131facd88d0ea2c6d0b06f2": 1,
	"649999ced9a400a9a42aabd361def70587997f519731abdee3e17d70ebaceafa": 1,
	"48c01b5059005455d9dcb0c6bcecdcb4fb5b2eabc1a9a82b57392baaa40f04e6": 1,
	"6ee038ae2053ad43b1fc5af7193f964a90fed6f6d3c5485fa9eed62e588e2ca9": 1,
	"a60f3e403a0de93850b22307d100dd7e06b977734337032abfb27322f3938993": 1,
	"e0a523797137b5ea61145a023d1ae936131fdb54a13706102c89ce6ccb816cfe": 1,
	"6d533c613bc953533082c39a295d656741ff6d1fb5e266295cd72236bdc93b16": 1,
	"08dc30ce2b661931daea75affcf984c016e1cb3d4f6ecd488b6db83f877abdc7": 1,
	"5ff0e80b514452c2bd3a45179290b6e5cad00704ba2705a0cb9556e483f67190": 1,
	"c185afd92c0fc7109f0c46efd6573883480626f20e15cdb2c66bb0379270a7e2": 1,
}

func once() {
	readerSeed := [32]byte{}
	for i := range readerSeed {
		readerSeed[i] = uint8(rand.Uint64())
	}
	reader := rand.NewChaCha8(readerSeed)
	for range 1024 * 32 {
		prikeySrc := make([]byte, 32)
		reader.Read(prikeySrc)
		prikey := ed25519.NewKeyFromSeed(prikeySrc)
		pubkey := prikey.Public().(ed25519.PublicKey)
		pubkeyHex := hex.EncodeToString(pubkey)
		log.Println(pubkeyHex)
		if _, win := aimpub[pubkeyHex]; win {
			log.Panicln("main: done", fmt.Sprintf("prikey=%s", hex.EncodeToString(prikey)))
		}
	}
}

func main() {
	log.Println("main:", runtime.NumCPU(), "logical cpus usable by the current process")
	grun := gool.Cpu()
	for {
		log.Println("main: once")
		grun.Call(once)
	}
}
