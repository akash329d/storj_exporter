package collectors

import (
	"log"

	"github.com/akash329d/storj_exporter/api"
	"github.com/akash329d/storj_exporter/models"

	"github.com/prometheus/client_golang/prometheus"
)

type SatelliteCollector struct {
	clients []*api.ApiClient
	metrics map[string]*prometheus.Desc
}

func NewSatelliteCollector(clients []*api.ApiClient) *SatelliteCollector {
	return &SatelliteCollector{
		clients: clients,
		metrics: map[string]*prometheus.Desc{
			"satelliteInfo": prometheus.NewDesc(
				"storj_satellite_info",
				"Storj satellite information",
				[]string{"node_id", "satellite_id", "satellite_url", "satellite_name"},
				nil,
			),
			"satelliteStorageSummary": prometheus.NewDesc(
				"storj_satellite_storage_summary_bytes",
				"Total amount of storage used by the node as reported by the satellite",
				[]string{"node_id", "satellite_id"},
				nil,
			),
			"satelliteAverageUsage": prometheus.NewDesc(
				"storj_satellite_average_usage_bytes",
				"Average storage usage in bytes as reported by the satellite",
				[]string{"node_id", "satellite_id"},
				nil,
			),
			"satelliteBandwidthSummary": prometheus.NewDesc(
				"storj_satellite_bandwidth_summary_bytes",
				"Total bandwidth used by the node as reported by the satellite",
				[]string{"node_id", "satellite_id"},
				nil,
			),
			"satelliteEgressSummary": prometheus.NewDesc(
				"storj_satellite_egress_summary_bytes",
				"Total egress bandwidth used by the node as reported by the satellite",
				[]string{"node_id", "satellite_id"},
				nil,
			),
			"satelliteIngressSummary": prometheus.NewDesc(
				"storj_satellite_ingress_summary_bytes",
				"Total ingress bandwidth used by the node as reported by the satellite",
				[]string{"node_id", "satellite_id"},
				nil,
			),
			"satelliteCurrentStorageUsed": prometheus.NewDesc(
				"storj_satellite_current_storage_used_bytes",
				"Current storage used by the node as reported by the satellite",
				[]string{"node_id", "satellite_id"},
				nil,
			),
			"satelliteAuditScore": prometheus.NewDesc(
				"storj_satellite_audit_score",
				"Audit score of the node as reported by the satellite",
				[]string{"node_id", "satellite_id"},
				nil,
			),
			"satelliteSuspensionScore": prometheus.NewDesc(
				"storj_satellite_suspension_score",
				"Suspension score of the node as reported by the satellite",
				[]string{"node_id", "satellite_id"},
				nil,
			),
			"satelliteOnlineScore": prometheus.NewDesc(
				"storj_satellite_online_score",
				"Online score of the node as reported by the satellite",
				[]string{"node_id", "satellite_id"},
				nil,
			),
			"satelliteDailyStorage": prometheus.NewDesc(
				"storj_satellite_storage",
				"Storage used by the node as reported by the satellite",
				[]string{"node_id", "satellite_id", "category"},
				nil,
			),
			"satelliteDailyBandwidth": prometheus.NewDesc(
				"storj_satellite_bandwidth_bytes",
				"Bandwidth used by the node as reported by the satellite",
				[]string{"node_id", "satellite_id", "type", "category"},
				nil,
			),
			"satelliteAuditStatus": prometheus.NewDesc(
				"storj_satellite_audit_status",
				"Current total and online audits as reported by the satellite",
				[]string{"node_id", "satellite_id", "audit_type"},
				nil,
			),
			"satellitePriceModel": prometheus.NewDesc(
				"storj_satellite_price_model",
				"Price model for the satellite",
				[]string{"node_id", "satellite_id", "type"},
				nil,
			),
			"satelliteNodeJoinedAt": prometheus.NewDesc(
				"storj_satellite_node_joined_timestamp",
				"Timestamp when the node joined the satellite",
				[]string{"node_id", "satellite_id"},
				nil,
			),
		},
	}
}

func (c *SatelliteCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric
	}
}

func (c *SatelliteCollector) Collect(ch chan<- prometheus.Metric) {
	for _, client := range c.clients {
		node, err := client.Node()
		if err != nil {
			log.Printf("Error collecting node data: %v", err)
			continue
		}

		for _, satellite := range node.Satellites {
			satelliteData, err := client.Satellite(satellite.ID)
			if err != nil {
				log.Printf("Error collecting satellite data: %v", err)
				continue
			}

			c.collectSatelliteMetrics(ch, client.NodeID, satellite.URL, &satelliteData)
		}
	}
}

func (c *SatelliteCollector) collectSatelliteMetrics(ch chan<- prometheus.Metric, nodeID string, satelliteURL string, data *models.SatelliteResponse) {
	ch <- prometheus.MustNewConstMetric(
		c.metrics["satelliteInfo"],
		prometheus.GaugeValue,
		1,
		nodeID,
		data.ID,
		satelliteURL,
		data.Audits.SatelliteName,
	)

	ch <- prometheus.MustNewConstMetric(c.metrics["satelliteStorageSummary"], prometheus.GaugeValue, data.StorageSummary, nodeID, data.ID)
	ch <- prometheus.MustNewConstMetric(c.metrics["satelliteAverageUsage"], prometheus.GaugeValue, data.AverageUsageBytes, nodeID, data.ID)
	ch <- prometheus.MustNewConstMetric(c.metrics["satelliteBandwidthSummary"], prometheus.GaugeValue, float64(data.BandwidthSummary), nodeID, data.ID)
	ch <- prometheus.MustNewConstMetric(c.metrics["satelliteEgressSummary"], prometheus.GaugeValue, float64(data.EgressSummary), nodeID, data.ID)
	ch <- prometheus.MustNewConstMetric(c.metrics["satelliteIngressSummary"], prometheus.GaugeValue, float64(data.IngressSummary), nodeID, data.ID)
	ch <- prometheus.MustNewConstMetric(c.metrics["satelliteCurrentStorageUsed"], prometheus.GaugeValue, float64(data.CurrentStorageUsed), nodeID, data.ID)
	ch <- prometheus.MustNewConstMetric(c.metrics["satelliteAuditScore"], prometheus.GaugeValue, data.Audits.AuditScore, nodeID, data.ID)
	ch <- prometheus.MustNewConstMetric(c.metrics["satelliteSuspensionScore"], prometheus.GaugeValue, data.Audits.SuspensionScore, nodeID, data.ID)
	ch <- prometheus.MustNewConstMetric(c.metrics["satelliteOnlineScore"], prometheus.GaugeValue, data.Audits.OnlineScore, nodeID, data.ID)

	// Storj API reports storageDaily as an array of values with an intervalStart and intervalInHours.
	// We allow Prometheus to handle the time series aspect of this for us and only return the latest valid value.
	var newestStorageDaily *models.StorageDaily
	for _, storageDaily := range data.StorageDaily {
		if storageDaily.IntervalInHours > 0 {
			if newestStorageDaily == nil || storageDaily.IntervalStart.After(newestStorageDaily.IntervalStart) {
				newestStorageDaily = &storageDaily
			}
		}
	}

	ch <- prometheus.MustNewConstMetric(
		c.metrics["satelliteDailyStorage"],
		prometheus.GaugeValue,
		newestStorageDaily.AtRestTotalBytes,
		nodeID,
		data.ID,
		"at_rest_total_bytes",
	)

	ch <- prometheus.MustNewConstMetric(
		c.metrics["satelliteDailyStorage"],
		prometheus.GaugeValue,
		newestStorageDaily.AtRestTotal,
		nodeID,
		data.ID,
		"at_rest_total",
	)

	// Storj API reports storageDaily as an array of values with an intervalStart and intervalInHours.
	// We allow Prometheus to handle the time series aspect of this for us and only return the latest valid value.
	var newestBandwidthDaily *models.BandwidthDaily
	for _, bandwidthDaily := range data.BandwidthDaily {
		if newestBandwidthDaily == nil || bandwidthDaily.IntervalStart.After(newestBandwidthDaily.IntervalStart) {
			newestBandwidthDaily = &bandwidthDaily
		}
	}
    c.collectDailyBandwidth(ch, nodeID, data.ID, newestBandwidthDaily)

    // Storj API reports audits as an array of values with windowStart.
	// We allow Prometheus to handle the time series aspect of this for us and only return the latest valid value.
	var newestAuditHistory *models.Window
	for _, window := range data.AuditHistory.Windows {
		if newestAuditHistory == nil || window.WindowStart.After(newestAuditHistory.WindowStart) {
			newestAuditHistory = &window
		}
	}

    ch <- prometheus.MustNewConstMetric(
        c.metrics["satelliteAuditStatus"],
        prometheus.GaugeValue,
        float64(newestAuditHistory.OnlineCount),
        nodeID,
        data.ID,
        "online",
    )

    ch <- prometheus.MustNewConstMetric(
        c.metrics["satelliteAuditStatus"],
        prometheus.GaugeValue,
        float64(newestAuditHistory.TotalCount),
        nodeID,
        data.ID,
        "total",
    )

	c.collectPriceModel(ch, nodeID, data.ID, &data.PriceModel)

	ch <- prometheus.MustNewConstMetric(
		c.metrics["satelliteNodeJoinedAt"],
		prometheus.GaugeValue,
		float64(data.NodeJoinedAt.Unix()),
		nodeID,
		data.ID,
	)
}

func (c *SatelliteCollector) collectDailyBandwidth(ch chan<- prometheus.Metric, nodeID, satelliteID string, daily *models.BandwidthDaily) {
	egressMetrics := map[string]int64{
		"repair": daily.Egress.Repair,
		"audit":  daily.Egress.Audit,
		"usage":  daily.Egress.Usage,
	}

	ingressMetrics := map[string]int64{
		"repair": daily.Ingress.Repair,
		"usage":  daily.Ingress.Usage,
	}

	for category, value := range egressMetrics {
		ch <- prometheus.MustNewConstMetric(
			c.metrics["satelliteDailyBandwidth"],
			prometheus.GaugeValue,
			float64(value),
			nodeID,
			satelliteID,
			"egress",
			category,
		)
	}

	for category, value := range ingressMetrics {
		ch <- prometheus.MustNewConstMetric(
			c.metrics["satelliteDailyBandwidth"],
			prometheus.GaugeValue,
			float64(value),
			nodeID,
			satelliteID,
			"ingress",
			category,
		)
	}

	ch <- prometheus.MustNewConstMetric(
		c.metrics["satelliteDailyBandwidth"],
		prometheus.GaugeValue,
		float64(daily.Delete),
		nodeID,
		satelliteID,
		"delete",
		"total",
	)
}

func (c *SatelliteCollector) collectPriceModel(ch chan<- prometheus.Metric, nodeID, satelliteID string, priceModel *models.PriceModel) {
	ch <- prometheus.MustNewConstMetric(c.metrics["satellitePriceModel"], prometheus.GaugeValue, float64(priceModel.EgressBandwidth), nodeID, satelliteID, "egress_bandwidth")
	ch <- prometheus.MustNewConstMetric(c.metrics["satellitePriceModel"], prometheus.GaugeValue, float64(priceModel.RepairBandwidth), nodeID, satelliteID, "repair_bandwidth")
	ch <- prometheus.MustNewConstMetric(c.metrics["satellitePriceModel"], prometheus.GaugeValue, float64(priceModel.AuditBandwidth), nodeID, satelliteID, "audit_bandwidth")
	ch <- prometheus.MustNewConstMetric(c.metrics["satellitePriceModel"], prometheus.GaugeValue, float64(priceModel.DiskSpace), nodeID, satelliteID, "disk_space")
}
