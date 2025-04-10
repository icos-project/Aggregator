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
package server_http

import (
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestHealthzEndpoint(t *testing.T) {
	// Start the server on a different port to avoid conflicts with other tests
	var wg sync.WaitGroup
	wg.Add(1)
	go CreateServer(&wg, "icos", "8081")
	time.Sleep(time.Second) // Wait for server to start up

	// Create HTTP client
	c := &http.Client{
		Timeout: 15 * time.Second,
	}

	// Make request to healthz endpoint
	resp, err := c.Get("http://localhost:8081/healthz")
	if err != nil {
		t.Fatalf("GET error: %v", err)
	}
	defer resp.Body.Close()

	// Verify status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK; got %v", resp.StatusCode)
	}

	// Read and verify response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	// Check that response contains the expected health check message
	// The exact format depends on responses.JSON implementation, but should contain this message
	expectedMessage := "Aggregator working properly!"
	if !strings.Contains(string(body), expectedMessage) {
		t.Errorf("Response body does not contain expected message.\nGot: %s\nExpected to contain: %s",
			string(body), expectedMessage)
	}
}
