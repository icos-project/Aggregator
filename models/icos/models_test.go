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
