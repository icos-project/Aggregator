package querier

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type PromQLQuery struct {
	Metric string
	Params map[string]string
}

func (q PromQLQuery) String() string {

	params := make([]string, 0, len(q.Params))
	for key, value := range q.Params {
		params = append(params, key+`="`+value+`"`)
	}
	if len(params) > 0 {
		return q.Metric + "{" + strings.Join(params, ", ") + "}"
	} else {
		return q.Metric
	}

}

func Query(query string) model.Vector {

	// create prometheus API client
	client, err := api.NewClient(api.Config{
		Address: os.Getenv("PROMETHEUS_ADDRESS"),
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
	}

	// create prometheus API object
	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := v1api.Query(ctx, query, time.Now(), v1.WithTimeout(5*time.Second))
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}

	// match the response to vector and print the response values
	switch r := result.(type) {
	case model.Vector:

		if r.Len() == 0 {
			fmt.Printf("PromQL Query Result length is zero: %s\n", query)
		}

		return r

	default:
		panic(errors.New("not implemented"))
	}
}

func ServeQuery(w http.ResponseWriter, r *http.Request) {
	Query(`up{container="prometheus"}`) // TODO: pass query string from API call
}
