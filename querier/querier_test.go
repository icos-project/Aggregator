package querier

import (
	"os"
	"testing"
)

func TestQuery(t *testing.T) {
	os.Setenv("PROMETHEUS_ADDRESS", "http://query.192.168.137.175.nip.io/") // thanos-query

	q := PromQLQuery{
		Metric: "up{container='prometheus'}",
		Params: map[string]string{}}

	want := float64(1)

	if got := Query(q.String()); got != want {
		t.Errorf("Query() = %v, want %v", got, want)
	}
}

func TestPromQL(t *testing.T) {

	query := PromQLQuery{
		Metric: "up",
		Params: map[string]string{"container": "prometheus"},
	}

	want := `up{container="prometheus"}`

	if got := query.String(); got != want {
		t.Errorf("PromQL() = %v, want %v", got, want)
	}
}
