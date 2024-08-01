/*
Copyright © 2022-2024 EVIDEN

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package querier

import (
	"os"
	"testing"
)

func TestQuery(t *testing.T) {

	os.Setenv("PROMETHEUS_ADDRESS", "http://query.192.168.137.175.nip.io/") // thanos-query

	q := PromQLQuery{
		Metric: "vector(100)",
		Params: map[string]string{}}

	want := []string{"{}", "100"}
	queryResult := Query(q.String())
	got := []string{queryResult[0].Metric.String(), queryResult[0].Value.String()}

	if got[0] != want[0] || got[1] != want[1] {
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
