package models

type NodeData struct {
	NodeID           string `json:"nodeID"`
	Wallet           string `json:"wallet"`
	UpToDate         bool   `json:"upToDate"`
	Version          string `json:"version"`
	AllowedVersion   string `json:"allowedVersion"`
	QuicStatus       string `json:"quicStatus"`
	DiskSpace        DiskSpace
	Bandwidth        Bandwidth
	Satellites       []Satellite `json:"satellites"`
	ConfiguredPort   string      `json:"configuredPort"`
	LastPinged       string      `json:"lastPinged"`
	StartedAt        string      `json:"startedAt"`
	LastQuicPingedAt string      `json:"lastQuicPingedAt"`
}

type Satellite struct {
	ID                 string  `json:"id"`
	URL                string  `json:"url"`
	Disqualified       *string `json:"disqualified"`
	Suspended          *string `json:"suspended"`
	CurrentStorageUsed int64   `json:"currentStorageUsed"`
}

type DiskSpace struct {
	Used      int64 `json:"used"`
	Available int64 `json:"available"`
	Trash     int64 `json:"trash"`
	Overused  int64 `json:"overused"`
}

type Bandwidth struct {
	Used      int64 `json:"used"`
	Available int64 `json:"available"`
}
