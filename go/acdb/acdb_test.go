package acdb

import (
	"testing"
)

func TestClient(t *testing.T) {
	for _, client := range []*Client{Mem(), Doc(t.TempDir()), Lru(4), Map(t.TempDir())} {
		client.Log(0)
		client.SetEncode("n", 1)
		n, err := client.GetInt("n")
		if err != nil || n != 1 {
			t.FailNow()
		}
		client.SetEncode("s", "Hello World!")
		s, err := client.GetString("s")
		if err != nil || s != "Hello World!" {
			t.FailNow()
		}
	}
}
