/*
Copyright 2023 Bull SAS

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
package server_http

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestHTTPRequest(t *testing.T) {

	w := []byte(`
	{"10.42.0.63:8080":{
		"name":"10.42.0.63:8080",
		"location":{},
		"serviceLevelAgreement":{},
		"API":{},
		"node":{
			"ocm-worker1.bull1.ari-imet.eu":{
				"name":"ocm-worker1.bull1.ari-imet.eu",
				"staticMetrics":{
					"cpuCores":4
				},
				"dynamicMetrics":{}
			}
		}
	},"
	10.42.1.7:8080":{
		"name":"10.42.1.7:8080",
		"location":{},
		"serviceLevelAgreement":{},
		"API":{},
		"node":{
			"k3s-node1":{
				"name":"k3s-node1",
				"staticMetrics":{
					"cpuCores":4
				},
				"dynamicMetrics":{}
			},"
			k3s-node2":{
				"name":"k3s-node2",
				"staticMetrics":{
					"cpuCores":4
				},
				"dynamicMetrics":{}
			}
		}
	}
	}`)
	w1 := []byte(strings.ReplaceAll(string(w), "\t", ""))
	w2 := []byte(strings.ReplaceAll(string(w1), "\n", ""))
	want := []byte(strings.ReplaceAll(string(w2), " ", ""))

	var wg sync.WaitGroup
	wg.Add(1)
	go CreateServer(&wg, "icos", "8080")
	time.Sleep(time.Second) // Wait for server to start up

	c := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := c.Get("http://localhost:8080/")

	if err != nil {
		t.Errorf("GET error: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status error: %v", resp.StatusCode)
	}

	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Read body: %v", err)
	}

	t.Log(string(got))
	t.Log(string(want))

	if bytes.Compare(want, got) != 0 {
		t.Errorf("Server error")
	}
}
