package aggregator

import (
	"os"
	"testing"
)

func TestQuery(t *testing.T) {
	os.Setenv("PROMETHEUS_ADDRESS", "http://localhost:10902") // thanos-query

	q := PromQLQuery{
		Metric: "up",
		Params: map[string]string{"container": "prometheus"}}

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
