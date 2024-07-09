package models

import "time"

type SatelliteResponse struct {
	ID                 string         `json:"id"`
	StorageDaily       []StorageDaily `json:"storageDaily"`
	BandwidthDaily     []BandwidthDaily `json:"bandwidthDaily"`
	StorageSummary     float64        `json:"storageSummary"`
	AverageUsageBytes  float64        `json:"averageUsageBytes"`
	BandwidthSummary   int64          `json:"bandwidthSummary"`
	EgressSummary      int64          `json:"egressSummary"`
	IngressSummary     int64          `json:"ingressSummary"`
	CurrentStorageUsed int64          `json:"currentStorageUsed"`
	Audits             Audits         `json:"audits"`
	AuditHistory       AuditHistory   `json:"auditHistory"`
	PriceModel         PriceModel     `json:"priceModel"`
	NodeJoinedAt       time.Time      `json:"nodeJoinedAt"`
}

type StorageDaily struct {
	AtRestTotal       float64   `json:"atRestTotal"`
	AtRestTotalBytes  float64   `json:"atRestTotalBytes"`
	IntervalInHours   int       `json:"intervalInHours"`
	IntervalStart     time.Time `json:"intervalStart"`
}

type BandwidthDaily struct {
	Egress        Egress    `json:"egress"`
	Ingress       Ingress   `json:"ingress"`
	Delete        int64     `json:"delete"`
	IntervalStart time.Time `json:"intervalStart"`
}

type Egress struct {
	Repair int64 `json:"repair"`
	Audit  int64 `json:"audit"`
	Usage  int64 `json:"usage"`
}

type Ingress struct {
	Repair int64 `json:"repair"`
	Usage  int64 `json:"usage"`
}

type Audits struct {
	AuditScore      float64 `json:"auditScore"`
	SuspensionScore float64 `json:"suspensionScore"`
	OnlineScore     float64 `json:"onlineScore"`
	SatelliteName   string  `json:"satelliteName"`
}

type AuditHistory struct {
	Score   float64  `json:"score"`
	Windows []Window `json:"windows"`
}

type Window struct {
	WindowStart time.Time `json:"windowStart"`
	TotalCount  int       `json:"totalCount"`
	OnlineCount int       `json:"onlineCount"`
}

type PriceModel struct {
	EgressBandwidth  int `json:"EgressBandwidth"`
	RepairBandwidth  int `json:"RepairBandwidth"`
	AuditBandwidth   int `json:"AuditBandwidth"`
	DiskSpace        int `json:"DiskSpace"`
}