package monitor_client

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dapr-platform/common"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var VictoriaMetricsUrl = "http://n9e-victoriametrics:8428"
var client = &http.Client{}

func Query(ctx context.Context, query string, queryTime time.Time) (result *QueryResult, err error) {
	if queryTime.IsZero() {
		queryTime = time.Now()
	}
	values := url.Values{}
	values.Add("query", query)
	values.Add("time", formatTime(queryTime))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, VictoriaMetricsUrl+"/api/v1/query", nil)
	if err != nil {
		common.Logger.Error("Error creating HTTP request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Add the URL values to the HTTP request.
	req.URL.RawQuery = values.Encode()
	//common.Logger.Debug(req.URL)
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
func formatTime(t time.Time) string {
	return strconv.FormatFloat(float64(t.Unix()), 'f', -1, 64)
}

func QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (result *QueryResult, err error) {
	u, err := url.Parse(VictoriaMetricsUrl + "/api/v1/query_range")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("query", query)
	q.Set("start", formatTime(start))
	q.Set("end", formatTime(end))
	q.Set("step", strconv.FormatFloat(float64(step), 'f', -1, 64))
	//values.Add("time", queryTime.Format("2006-01-02T15:04:05.000Z"))
	common.Logger.Debug("q=" + q.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), strings.NewReader(q.Encode()))
	if err != nil {
		common.Logger.Error("Error creating HTTP request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// Add the URL values to the HTTP request.
	//req.URL.RawQuery = values.Encode()
	resp, err := client.Do(req)
	if err != nil {
		common.Logger.Error("Error sending HTTP request:", err)
		return
	}
	//common.Logger.Debug("status=" + resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		common.Logger.Error("read body error ", err)
		return nil, err
	}
	defer resp.Body.Close()
	//common.Logger.Debug("body=" + string(body))
	var metricsResp QueryResultResp
	if err = json.Unmarshal(body, &metricsResp); err != nil {
		common.Logger.Error("Error reading HTTP response body:", err)
		return
	}
	common.Logger.Debug("metricsResp=", metricsResp)
	if metricsResp.Status != "success" {
		common.Logger.Error("Error reading response status:", metricsResp.Status)
		return nil, errors.New(metricsResp.Status)
	}
	return &metricsResp.Data, nil
}
