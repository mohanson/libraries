package main

import (
	"expvar"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/mohanson/libraries/go/expvpp"
)

var Expv = struct {
	Average *expvpp.Average
	Call    *expvar.Int
	Hits    *expvar.Int
	Percent *expvar.Func
}{
	Average: expvpp.NewAverage("Average", 64),
	Hits:    expvar.NewInt("Hits"),
	Call:    expvar.NewInt("Call"),
	Percent: expvpp.NewPercent("Percent", "Hits", "Call"),
}

func main() {
	go func() {
		for range time.NewTicker(time.Millisecond * 125).C {
			n := rand.Uint32N(256)
			Expv.Average.Append(float64(n))
			Expv.Call.Add(1)
			if n < 8 {
				Expv.Hits.Add(1)
			}
		}
	}()
	go func() {
		log.Println("main: listen and serve on 127.0.0.1:8080")
		http.ListenAndServe("127.0.0.1:8080", expvpp.ServeMux())
	}()
	select {}
}
