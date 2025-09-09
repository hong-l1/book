package timer

import (
	"testing"
	"time"
)

func TestTiker(t *testing.T) {
	temp := time.NewTicker(1 * time.Second)
	defer temp.Stop()
	for temp := range temp.C {
		t.Log(temp)
	}
}
