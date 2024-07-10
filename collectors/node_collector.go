package collectors

import (
	"log"
	"time"

	"github.com/akash329d/storj_exporter/api"
	"github.com/akash329d/storj_exporter/models"

	"github.com/prometheus/client_golang/prometheus"
)

type NodeCollector struct {
	clients []*api.ApiClient
	metrics map[string]*prometheus.Desc
}

func NewNodeCollector(clients []*api.ApiClient) *NodeCollector {
	return &NodeCollector{
		clients: clients,
		metrics: map[string]*prometheus.Desc{
			"nodeInfo": prometheus.NewDesc(
				"storj_node_info",
				"Storj node info",
				[]string{"node_id", "wallet", "version", "configured_port"},
				nil,
			),
			"satelliteStorageUsed": prometheus.NewDesc(
				"storj_satellite_storage_used_bytes",
				"Storage used per satellite",
				[]string{"node_id", "satellite_id", "satellite_url"},
				nil,
			),
			"diskSpace": prometheus.NewDesc(
				"storj_disk_space_bytes",
				"Storj disk space metrics",
				[]string{"node_id", "type"},
				nil,
			),
			"bandwidth": prometheus.NewDesc(
				"storj_bandwidth_bytes",
				"Storj bandwidth metrics",
				[]string{"node_id", "type"},
				nil,
			),
			"lastPinged": prometheus.NewDesc(
				"storj_last_pinged_timestamp",
				"Timestamp of last ping",
				[]string{"node_id"},
				nil,
			),
			"nodeStarted": prometheus.NewDesc(
				"storj_node_started_timestamp",
				"Timestamp when the node was started",
				[]string{"node_id"},
				nil,
			),
			"lastQuicPinged": prometheus.NewDesc(
				"storj_last_quic_pinged_timestamp",
				"Timestamp of last QUIC ping",
				[]string{"node_id"},
				nil,
			),
			"nodeUpToDate": prometheus.NewDesc(
				"storj_node_up_to_date",
				"Indicates if the node is up to date",
				[]string{"node_id"},
				nil,
			),
			"quicStatus": prometheus.NewDesc(
				"storj_quic_status",
				"QUIC status of the node",
				[]string{"node_id", "status"},
				nil,
			),
			"satelliteStatus": prometheus.NewDesc(
				"storj_satellite_status",
				"Status of the satellite",
				[]string{"node_id", "satellite_id", "satellite_url", "status"},
				nil,
			),
		},
	}
}

func (c *NodeCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric
	}
}

func (c *NodeCollector) Collect(ch chan<- prometheus.Metric) {
	for _, client := range c.clients {
		node, err := client.Node()
		if err != nil {
			log.Printf("Error collecting node metrics: %v", err)
			continue
		}

		c.collectNodeInfo(ch, client.NodeID, &node)
		c.collectSatelliteMetrics(ch, client.NodeID, &node)
		c.collectDiskSpaceMetrics(ch, client.NodeID, &node)
		c.collectBandwidthMetrics(ch, client.NodeID, &node)
		c.collectTimeMetrics(ch, client.NodeID, &node)
		c.collectStatusMetrics(ch, client.NodeID, &node)
	}
}

func (c *NodeCollector) collectNodeInfo(ch chan<- prometheus.Metric, nodeID string, node *models.NodeData) {
	ch <- prometheus.MustNewConstMetric(
		c.metrics["nodeInfo"],
		prometheus.GaugeValue,
		1,
		nodeID,
		node.Wallet,
		node.Version,
		node.ConfiguredPort,
	)
}

func (c *NodeCollector) collectSatelliteMetrics(ch chan<- prometheus.Metric, nodeID string, node *models.NodeData) {
	for _, satellite := range node.Satellites {
		ch <- prometheus.MustNewConstMetric(
			c.metrics["satelliteStorageUsed"],
			prometheus.GaugeValue,
			float64(satellite.CurrentStorageUsed),
			nodeID,
			satellite.ID,
			satellite.URL,
		)

		c.collectSatelliteStatus(ch, nodeID, &satellite)
	}
}

func (c *NodeCollector) collectSatelliteStatus(ch chan<- prometheus.Metric, nodeID string, satellite *models.Satellite) {
	statuses := map[string]float64{"active": 1, "disqualified": 0, "suspended": 0}

	if satellite.Disqualified != nil {
		statuses["disqualified"] = 1
		statuses["active"] = 0
	} else if satellite.Suspended != nil {
		statuses["suspended"] = 1
		statuses["active"] = 0
	}

	for status, value := range statuses {
		ch <- prometheus.MustNewConstMetric(
			c.metrics["satelliteStatus"],
			prometheus.GaugeValue,
			value,
			nodeID,
			satellite.ID,
			satellite.URL,
			status,
		)
	}
}

func (c *NodeCollector) collectDiskSpaceMetrics(ch chan<- prometheus.Metric, nodeID string, node *models.NodeData) {
	diskSpace := node.DiskSpace
	ch <- prometheus.MustNewConstMetric(c.metrics["diskSpace"], prometheus.GaugeValue, float64(diskSpace.Used), nodeID, "used")
	ch <- prometheus.MustNewConstMetric(c.metrics["diskSpace"], prometheus.GaugeValue, float64(diskSpace.Available), nodeID, "available")
	ch <- prometheus.MustNewConstMetric(c.metrics["diskSpace"], prometheus.GaugeValue, float64(diskSpace.Trash), nodeID, "trash")
	ch <- prometheus.MustNewConstMetric(c.metrics["diskSpace"], prometheus.GaugeValue, float64(diskSpace.Overused), nodeID, "overused")
}

func (c *NodeCollector) collectBandwidthMetrics(ch chan<- prometheus.Metric, nodeID string, node *models.NodeData) {
	bandwidth := node.Bandwidth
	ch <- prometheus.MustNewConstMetric(c.metrics["bandwidth"], prometheus.GaugeValue, float64(bandwidth.Used), nodeID, "used")
	ch <- prometheus.MustNewConstMetric(c.metrics["bandwidth"], prometheus.GaugeValue, float64(bandwidth.Available), nodeID, "available")
}

func (c *NodeCollector) collectTimeMetrics(ch chan<- prometheus.Metric, nodeID string, node *models.NodeData) {
	lastPinged, _ := time.Parse(time.RFC3339Nano, node.LastPinged)
	ch <- prometheus.MustNewConstMetric(c.metrics["lastPinged"], prometheus.GaugeValue, float64(lastPinged.Unix()), nodeID)

	startedAt, _ := time.Parse(time.RFC3339Nano, node.StartedAt)
	ch <- prometheus.MustNewConstMetric(c.metrics["nodeStarted"], prometheus.GaugeValue, float64(startedAt.Unix()), nodeID)

	lastQuicPinged, _ := time.Parse(time.RFC3339Nano, node.LastQuicPingedAt)
	ch <- prometheus.MustNewConstMetric(c.metrics["lastQuicPinged"], prometheus.GaugeValue, float64(lastQuicPinged.Unix()), nodeID)
}

func (c *NodeCollector) collectStatusMetrics(ch chan<- prometheus.Metric, nodeID string, node *models.NodeData) {
	ch <- prometheus.MustNewConstMetric(
		c.metrics["nodeUpToDate"],
		prometheus.GaugeValue,
		boolToFloat64(node.UpToDate),
		nodeID,
	)

	ch <- prometheus.MustNewConstMetric(
		c.metrics["quicStatus"],
		prometheus.GaugeValue,
		1,
		nodeID,
		node.QuicStatus,
	)
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
