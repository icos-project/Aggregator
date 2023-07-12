package aggregator

import (
	"os"
	"testing"
)

func TestQuery(t *testing.T) {
	os.Setenv("PROMETHEUS_ADDRESS", "http://localhost:10902")
	want := float64(1)
	if got := Query(); got != want {
		t.Errorf("Query() = %v, want %v", got, want)
	}
}
