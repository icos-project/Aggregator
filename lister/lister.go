package lister

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

func ListMetrics() model.LabelValues {

	// create prometheus API client
	client, err := api.NewClient(api.Config{
		Address: os.Getenv("PROMETHEUS_ADDRESS"),
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	labels, _, err := v1api.LabelValues(ctx, "__name__", []string{}, time.Now().Add(-5*time.Minute), time.Now())

	return labels

}
