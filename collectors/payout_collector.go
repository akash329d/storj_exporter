package collectors

import (
	"log"

	"github.com/akash329d/storj_exporter/api"

	"github.com/prometheus/client_golang/prometheus"
)

type PayoutCollector struct {
    StorjCollector
    egress_bandwidth *prometheus.Desc
    egress_bandwidth_payout *prometheus.Desc
    egress_repair_audit *prometheus.Desc
    egress_repair_audit_payout *prometheus.Desc
    disk_space *prometheus.Desc
    disk_space_payout *prometheus.Desc
    held_rate *prometheus.Desc
    payout *prometheus.Desc
    held *prometheus.Desc
    current_month_expectations *prometheus.Desc
}

func NewPayoutCollector(client *api.ApiClient) *PayoutCollector {
    return &PayoutCollector{
        StorjCollector: StorjCollector{client: client},
        egress_bandwidth: prometheus.NewDesc("storj_egress_bandwidth", "Egress bandwidth used by the node", []string{"node_id"}, nil),
        egress_bandwidth_payout: prometheus.NewDesc("storj_egress_bandwidth_payout", "Payout for the egress bandwidth used", []string{"node_id"}, nil),
        egress_repair_audit: prometheus.NewDesc("storj_egress_repair_audit", "Egress bandwidth used for repairs and audits", []string{"node_id"}, nil),
        egress_repair_audit_payout: prometheus.NewDesc("storj_egress_repair_audit_payout", "Payout for the egress bandwidth used for repairs and audits", []string{"node_id"}, nil),
        disk_space: prometheus.NewDesc("storj_disk_space", "Disk space used by the node", []string{"node_id"}, nil),
        disk_space_payout: prometheus.NewDesc("storj_disk_space_payout", "Payout for the disk space used", []string{"node_id"}, nil),
        held_rate: prometheus.NewDesc("storj_held_rate", "Percentage of payout held back by the network", []string{"node_id"}, nil),
        payout: prometheus.NewDesc("storj_payout", "Total payout for the node", []string{"node_id"}, nil),
        held: prometheus.NewDesc("storj_held", "Total amount held back by the network", []string{"node_id"}, nil),
        current_month_expectations: prometheus.NewDesc("storj_current_month_expectations", "Expected payout for the current month", []string{"node_id"}, nil),
    }
}

func (c *PayoutCollector) Describe(ch chan<- *prometheus.Desc) {
    ch <- c.egress_bandwidth
    ch <- c.egress_bandwidth_payout
    ch <- c.egress_repair_audit
    ch <- c.egress_repair_audit_payout
    ch <- c.disk_space
    ch <- c.disk_space_payout
    ch <- c.held_rate
    ch <- c.payout
    ch <- c.held
    ch <- c.current_month_expectations
}

func (c *PayoutCollector) Collect(ch chan<- prometheus.Metric) {
    payoutData, err := c.client.Payout()
    if err != nil {
        log.Printf("Error collecting node payout data: %v", err)
        return
    }

    ch <- prometheus.MustNewConstMetric(c.egress_bandwidth, prometheus.GaugeValue, float64(payoutData.CurrentMonth.EgressBandwidth), c.client.NodeID)
    ch <- prometheus.MustNewConstMetric(c.egress_bandwidth_payout, prometheus.GaugeValue, float64(payoutData.CurrentMonth.EgressBandwidthPayout), c.client.NodeID)
    ch <- prometheus.MustNewConstMetric(c.egress_repair_audit, prometheus.GaugeValue, float64(payoutData.CurrentMonth.EgressRepairAudit), c.client.NodeID)
    ch <- prometheus.MustNewConstMetric(c.egress_repair_audit_payout, prometheus.GaugeValue, float64(payoutData.CurrentMonth.EgressRepairAuditPayout), c.client.NodeID)
    ch <- prometheus.MustNewConstMetric(c.disk_space, prometheus.GaugeValue, float64(payoutData.CurrentMonth.DiskSpace), c.client.NodeID)
    ch <- prometheus.MustNewConstMetric(c.disk_space_payout, prometheus.GaugeValue, float64(payoutData.CurrentMonth.DiskSpacePayout), c.client.NodeID)
    ch <- prometheus.MustNewConstMetric(c.held_rate, prometheus.GaugeValue, float64(payoutData.CurrentMonth.HeldRate), c.client.NodeID)
    ch <- prometheus.MustNewConstMetric(c.payout, prometheus.GaugeValue, float64(payoutData.CurrentMonth.Payout), c.client.NodeID)
    ch <- prometheus.MustNewConstMetric(c.held, prometheus.GaugeValue, float64(payoutData.CurrentMonth.Held), c.client.NodeID)
    ch <- prometheus.MustNewConstMetric(c.current_month_expectations, prometheus.GaugeValue, float64(payoutData.CurrentMonthExpectations), c.client.NodeID)
}