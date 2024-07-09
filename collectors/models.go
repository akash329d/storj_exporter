package collectors

import "github.com/akash329d/storj_exporter/api"

type StorjCollector struct {
	clients[]   *api.ApiClient
}