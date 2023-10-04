package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestHTTPRequest(t *testing.T) {

	want := []byte(`{
				"Cluster_type": "1",
				"Cluster_name": "B",
				"Location_zone": "Madrid",
				"ServiceLevelAgreement": "C",
				"API": "D",
				"Node": [{"Node_name": "alpa", "Node_type": "1"}, 
						{"Node_name": "beta", "Node_type": "2"}]
				}`)

	go CreateServer()
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
