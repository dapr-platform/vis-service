package api

import (
	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
	"vis-service/model"
)

func InitDashboardRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/dashboard/page", DashboardPageListHandler)
	r.Get(common.BASE_CONTEXT+"/dashboard", DashboardListHandler)
	r.Post(common.BASE_CONTEXT+"/dashboard", UpsertDashboardHandler)
	r.Delete(common.BASE_CONTEXT+"/dashboard/{id}", DeleteDashboardHandler)
	r.Post(common.BASE_CONTEXT+"/dashboard/batch-delete", batchDeleteDashboardHandler)
	r.Post(common.BASE_CONTEXT+"/dashboard/batch-upsert", batchUpsertDashboardHandler)
	r.Get(common.BASE_CONTEXT+"/dashboard/groupby", DashboardGroupbyHandler)
}

// @Summary GroupBy
// @Description GroupBy, for example,  _select=level, then return  {level_val1:sum1,level_val2:sum2}, _where can input status=0
// @Tags Dashboard
// @Param _select query string true "_select"
// @Param _where query string false "_where"
// @Produce  json
// @Success 200 {object} common.Response{data=[]map[string]any} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /dashboard/groupby [get]
func DashboardGroupbyHandler(w http.ResponseWriter, r *http.Request) {

	common.CommonGroupby(w, r, common.GetDaprClient(), "o_dashboard")
}

// @Summary batch update
// @Description batch update
// @Tags Dashboard
// @Accept  json
// @Param entities body []map[string]any true "objects array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /dashboard/batch-upsert [post]
func batchUpsertDashboardHandler(w http.ResponseWriter, r *http.Request) {

	var entities []map[string]any
	err := common.ReadRequestBody(r, &entities)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	if len(entities) == 0 {
		common.HttpResult(w, common.ErrParam.AppendMsg("len of entities is 0"))
		return
	}

	err = common.DbBatchUpsert[map[string]any](r.Context(), common.GetDaprClient(), entities, model.DashboardTableInfo.Name, model.Dashboard_FIELD_NAME_id)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}

// @Summary page query
// @Description page query, _page(from 1 begin), _page_size, _order, and others fields, status=1, name=$like.%CAM%
// @Tags Dashboard
// @Param _page query int true "current page"
// @Param _page_size query int true "page size"
// @Param _order query string false "order"
// @Param id query string false "id"
// @Param name query string false "name"
// @Param app query string false "app"
// @Param layout query string false "layout"
// @Param created_by query string false "created_by"
// @Param created_time query string false "created_time"
// @Param updated_by query string false "updated_by"
// @Param updated_time query string false "updated_time"
// @Produce  json
// @Success 200 {object} common.Response{data=common.Page{items=[]model.Dashboard}} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /dashboard/page [get]
func DashboardPageListHandler(w http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("_page")
	pageSize := r.URL.Query().Get("_page_size")
	if page == "" || pageSize == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("page or pageSize is empty"))
		return
	}
	common.CommonPageQuery[model.Dashboard](w, r, common.GetDaprClient(), "o_dashboard", "id")

}

// @Summary query objects
// @Description query objects
// @Tags Dashboard
// @Param _select query string false "_select"
// @Param _order query string false "order"
// @Param id query string false "id"
// @Param name query string false "name"
// @Param app query string false "app"
// @Param layout query string false "layout"
// @Param created_by query string false "created_by"
// @Param created_time query string false "created_time"
// @Param updated_by query string false "updated_by"
// @Param updated_time query string false "updated_time"
// @Produce  json
// @Success 200 {object} common.Response{data=[]model.Dashboard} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /dashboard [get]
func DashboardListHandler(w http.ResponseWriter, r *http.Request) {
	common.CommonQuery[model.Dashboard](w, r, common.GetDaprClient(), "o_dashboard", "id")
}

// @Summary save
// @Description save
// @Tags Dashboard
// @Accept       json
// @Param item body model.Dashboard true "object"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Dashboard} "object"
// @Failure 500 {object} common.Response ""
// @Router /dashboard [post]
func UpsertDashboardHandler(w http.ResponseWriter, r *http.Request) {
	var val model.Dashboard
	err := common.ReadRequestBody(r, &val)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	beforeHook, exists := common.GetUpsertBeforeHook("Dashboard")
	if exists {
		v, err1 := beforeHook(r, val)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
		val = v.(model.Dashboard)
	}

	err = common.DbUpsert[model.Dashboard](r.Context(), common.GetDaprClient(), val, model.DashboardTableInfo.Name, "id")
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(val))
}

// @Summary delete
// @Description delete
// @Tags Dashboard
// @Param id  path string true "实例id"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Dashboard} "object"
// @Failure 500 {object} common.Response ""
// @Router /dashboard/{id} [delete]
func DeleteDashboardHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	beforeHook, exists := common.GetDeleteBeforeHook("Dashboard")
	if exists {
		_, err1 := beforeHook(r, id)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	common.CommonDelete(w, r, common.GetDaprClient(), "o_dashboard", "id", "id")
}

// @Summary batch delete
// @Description batch delete
// @Tags Dashboard
// @Accept  json
// @Param ids body []string true "id array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /dashboard/batch-delete [post]
func batchDeleteDashboardHandler(w http.ResponseWriter, r *http.Request) {

	var ids []string
	err := common.ReadRequestBody(r, &ids)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	if len(ids) == 0 {
		common.HttpResult(w, common.ErrParam.AppendMsg("len of ids is 0"))
		return
	}
	beforeHook, exists := common.GetBatchDeleteBeforeHook("Dashboard")
	if exists {
		_, err1 := beforeHook(r, ids)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	idstr := strings.Join(ids, ",")
	err = common.DbDeleteByOps(r.Context(), common.GetDaprClient(), "o_dashboard", []string{"id"}, []string{"in"}, []any{idstr})
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}
