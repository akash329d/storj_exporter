package collectors

import (
	"log"
	"strconv"

	"github.com/akash329d/storj_exporter/api"
	"github.com/akash329d/storj_exporter/models"

	"github.com/prometheus/client_golang/prometheus"
)

type NodeCollector struct {
	StorjCollector
	nodeInfo       *prometheus.Desc
	totalDiskSpace *prometheus.Desc
	totalBandwidth *prometheus.Desc
}

func NewNodeCollector(clients []*api.ApiClient) *NodeCollector {
	return &NodeCollector{
		StorjCollector: StorjCollector{clients: clients},
		nodeInfo: prometheus.NewDesc(
			"storj_node_info",
			"Storj node info",
			[]string{"node_id", "wallet", "up_to_date", "version", "allowed_version", "quic_status"},
			nil,
		),
		totalDiskSpace: prometheus.NewDesc(
			"storj_total_diskspace_bytes",
			"Storj total diskspace metrics",
			[]string{"node_id", "type"},
			nil,
		),
		totalBandwidth: prometheus.NewDesc(
			"storj_total_bandwidth_bytes",
			"Storj total bandwidth metrics",
			[]string{"node_id", "type"},
			nil,
		),
	}
}

func (c *NodeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.nodeInfo
	ch <- c.totalDiskSpace
	ch <- c.totalBandwidth
}

func (c *NodeCollector) Collect(ch chan<- prometheus.Metric) {
	for _, client := range c.clients {
		node, err := client.Node()
		if err != nil {
			log.Printf("Error collecting node metrics: %v", err)
			return
		}

		// Collect node info
		ch <- prometheus.MustNewConstMetric(
			c.nodeInfo,
			prometheus.GaugeValue,
			1,
			client.NodeID,
			node.Wallet,
			strconv.FormatBool(node.UpToDate),
			node.Version,
			node.AllowedVersion,
			node.QuicStatus,
		)

		// Collect disk space metrics
		c.collectDiskSpaceMetrics(ch, &node, client.NodeID)
		// Collect bandwidth metrics
		c.collectBandwidthMetrics(ch, &node, client.NodeID)
	}
}

func (c *NodeCollector) collectDiskSpaceMetrics(ch chan<- prometheus.Metric, node *models.NodeData, nodeId string) {
	diskSpace := node.DiskSpace
	ch <- prometheus.MustNewConstMetric(c.totalDiskSpace, prometheus.GaugeValue, float64(diskSpace.Used), nodeId, "used")
	ch <- prometheus.MustNewConstMetric(c.totalDiskSpace, prometheus.GaugeValue, float64(diskSpace.Available), nodeId, "available")
	ch <- prometheus.MustNewConstMetric(c.totalDiskSpace, prometheus.GaugeValue, float64(diskSpace.Trash), nodeId, "trash")
}

func (c *NodeCollector) collectBandwidthMetrics(ch chan<- prometheus.Metric, node *models.NodeData, nodeId string) {
	bandwidth := node.Bandwidth
	ch <- prometheus.MustNewConstMetric(c.totalBandwidth, prometheus.GaugeValue, float64(bandwidth.Used), nodeId, "used")
	ch <- prometheus.MustNewConstMetric(c.totalBandwidth, prometheus.GaugeValue, float64(bandwidth.Available), nodeId, "available")
}
