package expvpp

import (
	"encoding/json"
	"expvar"
	"math"
	"net/http"
	"strconv"
	"sync/atomic"
)

// Average is a structure to maintain a running average using expvar.Float.
type Average struct {
	f atomic.Uint64
	l float64
}

// String returns the string representation of the average.
func (v *Average) String() string {
	return strconv.FormatFloat(math.Float64frombits(v.f.Load()), 'g', -1, 64)
}

// Adds a new value to the running average.
func (v *Average) Append(delta float64) {
	for {
		cur := v.f.Load()
		curVal := math.Float64frombits(cur)
		nxtVal := curVal + (delta-curVal)/v.l
		nxt := math.Float64bits(nxtVal)
		if v.f.CompareAndSwap(cur, nxt) {
			return
		}
	}
}

// NewAverage creates and initializes a new Average instance.
func NewAverage(name string, length int) *Average {
	average := &Average{
		f: atomic.Uint64{},
		l: float64(length),
	}
	expvar.Publish(name, average)
	return average
}

// NewPercent creates a new expvar.Func that calculates the ratio of two expvar.Int or expvar.Float metrics.
func NewPercent(name string, n string, d string) *expvar.Func {
	f := expvar.Func(func() any {
		v, _ := strconv.ParseFloat(expvar.Get(n).String(), 64)
		w, _ := strconv.ParseFloat(expvar.Get(d).String(), 64)
		return float64(v) / float64(max(1, w))
	})
	expvar.Publish(name, f)
	return &f
}

// ServeMux returns a new http.ServeMux that removes cmdline and memstats from expvar default exports.
// See: https://github.com/golang/go/issues/29105
func ServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/vars", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		vars := new(expvar.Map).Init()
		expvar.Do(func(kv expvar.KeyValue) {
			vars.Set(kv.Key, kv.Value)
		})
		vars.Delete("cmdline")
		vars.Delete("memstats")
		msg := map[string]any{}
		err := json.Unmarshal([]byte(vars.String()), &msg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		enc := json.NewEncoder(w)
		enc.SetIndent("", "    ")
		enc.Encode(msg)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.DefaultServeMux.ServeHTTP(w, r)
	})
	return mux
}
