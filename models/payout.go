package models

type PayoutData struct {
    EgressBandwidth         int64   `json:"egressBandwidth"`
    EgressBandwidthPayout   float64 `json:"egressBandwidthPayout"`
    EgressRepairAudit       int64   `json:"egressRepairAudit"`
    EgressRepairAuditPayout float64 `json:"egressRepairAuditPayout"`
    DiskSpace               float64 `json:"diskSpace"`
    DiskSpacePayout         float64 `json:"diskSpacePayout"`
    HeldRate                float64 `json:"heldRate"`
    Payout                  float64 `json:"payout"`
    Held                    float64 `json:"held"`
}

type PayoutResponse struct {
    CurrentMonth  PayoutData `json:"currentMonth"`
    PreviousMonth PayoutData `json:"previousMonth"`
    CurrentMonthExpectations float64 `json:"currentMonthExpectations"`
}
