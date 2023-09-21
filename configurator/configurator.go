package configurator

import (
	"icos/server/querier"
	"os"

	"gopkg.in/yaml.v3"
)

type test struct {
	text string
}

func NewConfigQuery() *querier.PromQLQuery {

	yamlData, err := os.ReadFile("../../config.yaml")
	if err != nil {
		panic(err)
	}

	q := &querier.PromQLQuery{}

	if err := yaml.Unmarshal([]byte(yamlData), q); err != nil {
		panic(err)
	}

	return q
}
