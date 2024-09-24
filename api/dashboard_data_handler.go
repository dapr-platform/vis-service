package api

import (
	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
	"net/http"
	"vis-service/entity"
	"vis-service/service"
)

func InitDashboardDataRoute(r chi.Router) {
	r.Post(common.BASE_CONTEXT+"/dashboard-data/metric", dashboardMetricHandler)
	r.Post(common.BASE_CONTEXT+"/dashboard-data/metric-range", dashboardMetricRangeHandler)
	r.Post(common.BASE_CONTEXT+"/dashboard-data/db-query", dashboardDbHandler)

}

// @Summary DB
// @Description 从DB中获取数据
// @Tags DashboardData
// @Produce  json
// @Param reqs body entity.QueryDbReq true "请求参数"
// @Success 200 {object} common.Response{data=[]map[string]any} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /dashboard-data/db-query [post]
func dashboardDbHandler(w http.ResponseWriter, r *http.Request) {
	req := &entity.QueryDbReq{}
	err := common.ReadRequestBody(r, req)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	data, err := service.FetchDatabase(r.Context(), req)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(data))
}

// @Summary Metric
// @Description 从metric中获取数据
// @Tags DashboardData
// @Produce  json
// @Param reqs body entity.QueryMetricsCombin true "请求参数"
// @Success 200 {object} common.Response{data=[]map[string]any} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /dashboard-data/metric [post]
func dashboardMetricHandler(w http.ResponseWriter, r *http.Request) {
	req := &entity.QueryMetricsCombin{}
	err := common.ReadRequestBody(r, req)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	data, err := service.FetchMetricData(r.Context(), req)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(data))
}

// @Summary Metric Range
// @Description 从metric中获取批量数据
// @Tags DashboardData
// @Produce  json
// @Param reqs body entity.QueryMetricsRangeCombin true "请求参数"
// @Success 200 {object} common.Response{data=[]map[string]any} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /dashboard-data/metric-range [post]
func dashboardMetricRangeHandler(w http.ResponseWriter, r *http.Request) {
	req := &entity.QueryMetricsRangeCombin{}
	err := common.ReadRequestBody(r, req)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	data, err := service.FetchMetricRangeData(r.Context(), req)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(data))
}
