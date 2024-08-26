package entity

import "github.com/prometheus/common/model"

type QueryMetrics struct {
	Id         string      `json:"id"`
	Query      string      `json:"query"`
	ResultDefs []ResultDef `json:"result_defs"`
}
type ResultDef struct {
	Name          model.LabelName `json:"name"`
	Type          string          `json:"type"` //metric||value
	DecimalPlaces int             `json:"decimal_places"`
}
type QueryRangeMetrics struct {
	Query string `json:"query"`
	// The boundaries of the time range.
	Start string `json:"start"` //格式 2006-01-02 15:04:05
	End   string `json:"end"`   //格式 2006-01-02 15:04:05
	// The maximum time between two slices within the boundaries.
	Step int `json:"step"` //单位秒
}
type QueryMetricsCombin struct {
	Reqs []QueryMetrics `json:"reqs"`
}

type QueryMetricsRangeCombin struct {
	Reqs  []QueryMetrics `json:"reqs"`
	Start string         `json:"start"` //格式 2006-01-02 15:04:05
	End   string         `json:"end"`   //格式 2006-01-02 15:04:05
	// The maximum time between two slices within the boundaries.
	Step int `json:"step"` //单位秒

}

type DashboardDataResp struct {
	Data map[string]any `json:"data"`
}
