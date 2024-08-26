package monitor_client

import (
	"context"
	"errors"
	"github.com/dapr-platform/common"
	"github.com/spf13/cast"
	"net/http"
	"net/url"
	"time"
)

var LokiUrl = "http://loki:3100"
var lokiClient = &http.Client{}

func LokiQuery(ctx context.Context, query string, limit int) (result *QueryResult, err error) {

	values := url.Values{}
	values.Add("query", query)
	values.Add("limit", cast.ToString(limit))
	//values.Add("time", queryTime.Format("2006-01-02T15:04:05.000Z"))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, LokiUrl+"/loki/api/v1/query", nil)
	if err != nil {
		common.Logger.Error("Error creating HTTP request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Add the URL values to the HTTP request.
	req.URL.RawQuery = values.Encode()
	common.Logger.Debug(req.URL)
	// Send the HTTP request.
	resp, err := client.Do(req)
	if err != nil {
		common.Logger.Error("Error sending HTTP request:", err)
		return
	}

	var metricsResp QueryResultResp
	if err = common.ReadResponseBody(resp, &metricsResp); err != nil {
		common.Logger.Error("Error reading HTTP response body:", err)
		return
	}
	if metricsResp.Status != "success" {
		common.Logger.Error("Error reading HTTP response body:", metricsResp.Status)
		return nil, errors.New(metricsResp.Status)
	}
	return &metricsResp.Data, nil
}

func LokiStreamQuery(ctx context.Context, query string, limit int, preHours int) (result *LokiQueryResult, err error) {

	values := url.Values{}
	values.Add("query", query)
	values.Add("limit", cast.ToString(limit))
	end := time.Now().UnixNano()
	start := time.Now().Add(time.Duration(-1*preHours) * time.Hour).UnixNano()
	values.Add("start", cast.ToString(start))
	values.Add("end", cast.ToString(end))
	//values.Add("time", queryTime.Format("2006-01-02T15:04:05.000Z"))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, LokiUrl+"/loki/api/v1/query_range", nil)
	if err != nil {
		common.Logger.Error("Error creating HTTP request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Add the URL values to the HTTP request.
	req.URL.RawQuery = values.Encode()
	common.Logger.Debug(req.URL)
	// Send the HTTP request.
	resp, err := client.Do(req)
	if err != nil {
		common.Logger.Error("Error sending HTTP request:", err)
		return
	}

	var lokiResp LokiQueryResultResp
	if err = common.ReadResponseBody(resp, &lokiResp); err != nil {
		common.Logger.Error("Error reading HTTP response body:", err)
		return
	}
	if lokiResp.Status != "success" {
		common.Logger.Error("Error reading HTTP response body:", lokiResp.Status)
		return nil, errors.New(lokiResp.Status)
	}
	return &lokiResp.Data, nil
}

func LokiLabelValues(ctx context.Context, label string) (result []string, err error) {

	urlSuffix := "/loki/api/v1/label/" + label + "/values"
	//values.Add("time", queryTime.Format("2006-01-02T15:04:05.000Z"))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, LokiUrl+urlSuffix, nil)
	if err != nil {
		common.Logger.Error("Error creating HTTP request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the HTTP request.
	resp, err := client.Do(req)
	if err != nil {
		common.Logger.Error("Error sending HTTP request:", err)
		return
	}

	var lokiResp LokiLabelValueResp
	if err = common.ReadResponseBody(resp, &lokiResp); err != nil {
		common.Logger.Error("Error reading HTTP response body:", err)
		return
	}
	if lokiResp.Status != "success" {
		common.Logger.Error("Error reading HTTP response body:", lokiResp.Status)
		return nil, errors.New(lokiResp.Status)
	}
	return lokiResp.Data, nil
}
