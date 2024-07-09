package collectors

import (
	"log"

	"github.com/akash329d/storj_exporter/api"

	"github.com/prometheus/client_golang/prometheus"
)

type SatelliteCollector struct {
    StorjCollector
    storageSummary     *prometheus.Desc
    averageUsageBytes  *prometheus.Desc
    bandwidthSummary   *prometheus.Desc
    egressSummary      *prometheus.Desc
    ingressSummary     *prometheus.Desc
    currentStorageUsed *prometheus.Desc
    auditScore         *prometheus.Desc
    suspensionScore    *prometheus.Desc
    onlineScore        *prometheus.Desc
    monthlyEgress      *prometheus.Desc
    monthlyIngress     *prometheus.Desc
    dayStorage           *prometheus.Desc
}

func NewSatelliteCollector(client *api.ApiClient) *SatelliteCollector {
    return &SatelliteCollector{
        StorjCollector: StorjCollector{client: client},
        storageSummary: prometheus.NewDesc("storj_storage_summary", "Total amount of storage used by the node reported by satellite", []string{"satellite_id", "node_id", "satellite_url"}, nil),
        averageUsageBytes: prometheus.NewDesc("storj_average_usage_bytes", "Average storage usage in bytes reported by satellite", []string{"satellite_id", "node_id", "satellite_url"}, nil),
        bandwidthSummary: prometheus.NewDesc("storj_bandwidth_summary", "Total bandwidth used by the node reported by satellite", []string{"satellite_id", "node_id", "satellite_url"}, nil),
        egressSummary: prometheus.NewDesc("storj_egress_summary", "Total egress bandwidth used by the node reported by satellite", []string{"satellite_id", "node_id", "satellite_url"}, nil),
        ingressSummary: prometheus.NewDesc("storj_ingress_summary", "Total ingress bandwidth used by the node reported by satellite", []string{"satellite_id", "node_id", "satellite_url"}, nil),
        currentStorageUsed: prometheus.NewDesc("storj_current_storage_used", "Current storage used by the node reported by satellite", []string{"satellite_id", "node_id", "satellite_url"}, nil),
        auditScore: prometheus.NewDesc("storj_audit_score", "Audit score of the node reported by satellite", []string{"satellite_id", "node_id", "satellite_url"}, nil),
        suspensionScore: prometheus.NewDesc("storj_suspension_score", "Suspension score of the node reported by satellite", []string{"satellite_id", "node_id", "satellite_url"}, nil),
        onlineScore: prometheus.NewDesc("storj_online_score", "Online score of the node reported by satellite", []string{"satellite_id", "node_id", "satellite_url"}, nil),
        monthlyEgress: prometheus.NewDesc("storj_monthly_egress", "Total egress bandwidth used by the node for the current month reported by satellite", []string{"satellite_id", "node_id", "satellite_url"}, nil),
        monthlyIngress: prometheus.NewDesc("storj_monthly_ingress", "Total ingress bandwidth used by the node for the current month reported by satellite", []string{"satellite_id", "node_id", "satellite_url"}, nil),
        dayStorage: prometheus.NewDesc("storj_day_storage", "Total data stored on disk since current day reported by satellite", []string{"satellite_id", "node_id", "satellite_url"}, nil),
    }
}

func (c *SatelliteCollector) Describe(ch chan<- *prometheus.Desc) {
    ch <- c.storageSummary
    ch <- c.averageUsageBytes
    ch <- c.bandwidthSummary
    ch <- c.egressSummary
    ch <- c.ingressSummary
    ch <- c.currentStorageUsed
    ch <- c.auditScore
    ch <- c.suspensionScore
    ch <- c.onlineScore
    ch <- c.monthlyEgress
    ch <- c.monthlyIngress
    ch <- c.dayStorage
}

func (c *SatelliteCollector) Collect(ch chan<- prometheus.Metric) {
    for _, satellite := range c.client.Satellites {
        satelliteData, err := c.client.Satellite(satellite.ID)
        if err != nil {
            log.Printf("Error collecting satellite data: %v", err)
            continue
        }

		var totalMonthlyEgress int64
		var totalMonthlyIngress int64
		for _, bandwidthDaily := range satelliteData.BandwidthDaily {
			totalMonthlyEgress += bandwidthDaily.Egress.Repair
			totalMonthlyEgress += bandwidthDaily.Egress.Audit
			totalMonthlyEgress += bandwidthDaily.Egress.Usage
			totalMonthlyIngress += bandwidthDaily.Ingress.Repair
			totalMonthlyIngress += bandwidthDaily.Ingress.Usage
		}

        ch <- prometheus.MustNewConstMetric(c.storageSummary, prometheus.GaugeValue, float64(satelliteData.StorageSummary), satellite.ID, c.client.NodeID, satellite.URL)
        ch <- prometheus.MustNewConstMetric(c.averageUsageBytes, prometheus.GaugeValue, float64(satelliteData.AverageUsageBytes), satellite.ID, c.client.NodeID, satellite.URL)
        ch <- prometheus.MustNewConstMetric(c.bandwidthSummary, prometheus.GaugeValue, float64(satelliteData.BandwidthSummary), satellite.ID, c.client.NodeID, satellite.URL)
        ch <- prometheus.MustNewConstMetric(c.egressSummary, prometheus.GaugeValue, float64(satelliteData.EgressSummary), satellite.ID, c.client.NodeID, satellite.URL)
        ch <- prometheus.MustNewConstMetric(c.ingressSummary, prometheus.GaugeValue, float64(satelliteData.IngressSummary), satellite.ID, c.client.NodeID, satellite.URL)
        ch <- prometheus.MustNewConstMetric(c.currentStorageUsed, prometheus.GaugeValue, float64(satelliteData.CurrentStorageUsed), satellite.ID, c.client.NodeID, satellite.URL)
        ch <- prometheus.MustNewConstMetric(c.auditScore, prometheus.GaugeValue, float64(satelliteData.Audits.AuditScore), satellite.ID, c.client.NodeID, satellite.URL)
        ch <- prometheus.MustNewConstMetric(c.suspensionScore, prometheus.GaugeValue, float64(satelliteData.Audits.SuspensionScore), satellite.ID, c.client.NodeID, satellite.URL)
        ch <- prometheus.MustNewConstMetric(c.onlineScore, prometheus.GaugeValue, float64(satelliteData.Audits.OnlineScore), satellite.ID, c.client.NodeID, satellite.URL)
        ch <- prometheus.MustNewConstMetric(c.monthlyEgress, prometheus.GaugeValue, float64(totalMonthlyEgress), satellite.ID, c.client.NodeID, satellite.URL)
        ch <- prometheus.MustNewConstMetric(c.monthlyIngress, prometheus.GaugeValue, float64(totalMonthlyIngress), satellite.ID, c.client.NodeID, satellite.URL)
        ch <- prometheus.MustNewConstMetric(c.dayStorage, prometheus.GaugeValue, float64(satelliteData.StorageDaily[len(satelliteData.StorageDaily)-1].AtRestTotal), satellite.ID, c.client.NodeID, satellite.URL)
    }
}