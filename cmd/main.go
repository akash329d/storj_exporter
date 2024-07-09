package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/akash329d/storj_exporter/api"
	"github.com/akash329d/storj_exporter/collectors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run() {
	nodeURLs := getNodeURLs()
	if len(nodeURLs) == 0 {
		log.Fatal("No Storj node URLs found in environment variables.")
	}

	clients := make([]*api.ApiClient, len(nodeURLs))
	for i, url := range nodeURLs {
		clients[i] = api.NewApiClient(url)
	}
	
	prometheus.MustRegister(collectors.NewNodeCollector(clients))
	prometheus.MustRegister(collectors.NewSatelliteCollector(clients))
	prometheus.MustRegister(collectors.NewPayoutCollector(clients))

	port := 8000 // Default port
    if value, exists := os.LookupEnv("EXPORTER_PORT"); exists {
		if intValue, err := strconv.Atoi(value); err != nil {
			log.Fatalf("Invalid port number in EXPORTER_PORT: %v\n", err)
		} else {
			port = intValue
		}
    }

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting Storj Node Exporter on :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func getNodeURLs() []string {
	var urls []string
	for i := 1; ; i++ {
		NodeURL := os.Getenv(fmt.Sprintf("STORJ_NODE_%d_URL", i))
		if NodeURL == "" {
			break
		}
		url, err := url.Parse(NodeURL)
		if err != nil {
			log.Printf("Error parsing URL for node %d, %s: %v", i, NodeURL, err)
			continue
		}
		urls = append(urls, url.String())
	}
	return urls
}