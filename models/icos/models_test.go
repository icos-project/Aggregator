package models_icos

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestControler(t *testing.T) {

	want := []byte(`{"type":"MetaOrchestrator","name":"ICOS1","location":{"name":"BCN"},"serviceLevelAgreement":{},"API":{}}`)

	c := Controller{Type: "MetaOrchestrator",
		Name:     "ICOS1",
		Location: Location{Name: "BCN"},
	}

	got, err := json.Marshal(c)

	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(got))
		t.Log(string(want))
	}

	if bytes.Compare(want, got) != 0 {
		t.Errorf("Controller model error")
	}

}
