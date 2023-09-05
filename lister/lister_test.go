package lister

import (
	"os"
	"reflect"
	"testing"
)

func TestListMetrics(t *testing.T) { // TODO: Check if {"up","promhttp_metric_handler_requests_total"} is contained in the labels list
	os.Setenv("PROMETHEUS_ADDRESS", "http://localhost:10902") // thanos-query

	want := []string{"up"}

	labels := ListMetrics()

	for _, label := range labels {

		t.Log(label)
	}

	if !reflect.DeepEqual(want, labels) {
		t.Errorf("lister.ListMetrics() error")
	}
}
