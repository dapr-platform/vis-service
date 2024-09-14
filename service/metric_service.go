package service

import (
	"context"
	"fmt"
	"github.com/dapr-platform/common"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"github.com/spf13/cast"
	"strconv"
	"time"
	"vis-service/entity"
	"vis-service/monitor_client"
)

func FetchMetricData(ctx context.Context, reqs *entity.QueryMetricsCombin) (data map[string][]map[model.LabelName]any, err error) {
	data = make(map[string][]map[model.LabelName]any, 0)
	for _, req := range reqs.Reqs {
		common.Logger.Debug("req id=", req.Id)
		vector, err1 := fetchOneMetricVal(ctx, req.Query, time.Now())
		if err1 != nil {
			common.Logger.Error("fetch metric data error ", req, err1)
			err = errors.Wrap(err1, "fetch metric data error "+req.Id)
			return
		}
		oneData := make([]map[model.LabelName]any, 0)
		common.Logger.Debug("result len=", len(*vector))
		for _, v := range *vector {
			itemData := make(map[model.LabelName]any, 0)
			for _, resultDef := range req.ResultDefs {
				switch resultDef.Type {
				case "":
				case "metric":
					itemData[resultDef.Name] = cast.ToString(string(v.Metric[resultDef.Name]))
				case "value":
					bits := cast.ToInt(resultDef.DecimalPlaces)
					if bits == 0 {
						bits = 2
					}
					itemData[resultDef.Name] = toFloat64WithBits(float64(v.Value), bits)
				case "timestamp":
					itemData[resultDef.Name] = v.Timestamp.Time().In(time.Local).Format("2006-01-02 15:04:05")
				}

			}
			oneData = append(oneData, itemData)
		}
		data[req.Id] = oneData
	}
	return
}

func FetchMetricRangeData(ctx context.Context, reqs *entity.QueryMetricsRangeCombin) (data map[string][]map[model.LabelName]any, err error) {
	data = make(map[string][]map[model.LabelName]any, 0)
	ss := time.Now().Add(-time.Hour)
	se := time.Now()
	if reqs.Start != "" {
		ss, err = time.ParseInLocation("2006-01-02 15:04:05", reqs.Start, time.Local)
		if err != nil {
			err = errors.Wrap(err, "query.Start "+reqs.Start+" "+err.Error())
			return
		}
	}
	if reqs.End != "" {
		se, err = time.ParseInLocation("2006-01-02 15:04:05", reqs.End, time.Local)
		if err != nil {
			err = errors.Wrap(err, "query.End "+reqs.End+" "+err.Error())
			return
		}
	}
	if ss.After(se) {
		err = errors.Wrap(err, "start after end: "+reqs.Start+" "+reqs.End)
		return
	}
	step := time.Duration(15)
	if reqs.Step != 0 {
		step = time.Duration(reqs.Step)
	}
	for _, req := range reqs.Reqs {

		metrix, err1 := fetchOneMetricRangeVal(ctx, req.Query, ss, se, step)
		if err1 != nil {
			common.Logger.Error("fetch metric data error ", req, err1)
			err = errors.Wrap(err1, "fetch metric data error "+req.Id)
			return
		}
		oneData := make([]map[model.LabelName]any, 0)

		for _, v := range *metrix {
			itemData := make(map[model.LabelName]any, 0)
			var valueName model.LabelName
			valueBit := 2

			for _, resultDef := range req.ResultDefs {
				switch resultDef.Type {
				case "":
				case "metric":
					itemData[resultDef.Name] = cast.ToString(string(v.Metric[resultDef.Name]))
				case "value":
					valueBit = cast.ToInt(resultDef.DecimalPlaces)
					if valueBit == 0 {
						valueBit = 2
					}
					valueName = resultDef.Name

				}

			}
			newValues := make([]entity.ValuePair, 0)
			for _, vv := range v.Values {
				newValues = append(newValues, entity.ValuePair{
					Timestamp: vv.Timestamp.Time().In(time.Local).Format("2006-01-02 15:04:05"),
					Value:     toFloat64WithBits(float64(vv.Value), valueBit),
				})
			}
			itemData[valueName] = newValues

			oneData = append(oneData, itemData)
		}
		data[req.Id] = oneData
	}
	return
}

func toFloat64WithBits(val float64, bits int) (v float64) {
	// 使用格式化字符串创建指定小数位数的字符串表示
	format := fmt.Sprintf("%%.%df", bits)
	// 将浮点数格式化为字符串
	str := fmt.Sprintf(format, val)
	// 将字符串解析回浮点数
	v, _ = strconv.ParseFloat(str, 64)
	return
}

func fetchOneMetricRangeVal(ctx context.Context, query string, ss time.Time, se time.Time, step time.Duration) (result *model.Matrix, err error) {
	data, err := monitor_client.QueryRange(ctx, query, ss, se, step)
	if err != nil {
		common.Logger.Error("refreshHostMetrics err", err)
		return
	}

	switch data.Result.Type().String() {
	case "matrix":
		metrix := data.Result.(model.Matrix)
		result = &metrix
	}

	if result == nil {
		return nil, fmt.Errorf("unexpected result type: %s", data.Result.Type().String())
	}
	return
}

func fetchOneMetricVal(ctx context.Context, query string, queryTime time.Time) (result *model.Vector, err error) {
	data, err := monitor_client.Query(ctx, query, queryTime)
	if err != nil {
		common.Logger.Error("refreshHostMetrics err", err)
		return
	}

	switch data.Result.Type().String() {
	case "vector":
		vector := data.Result.(model.Vector)
		result = &vector
	}
	if result == nil {
		return nil, fmt.Errorf("unexpected result type: %s", data.Result.Type().String())
	}
	return
}
