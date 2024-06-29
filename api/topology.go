package api

import (
	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
	"vis-service/model"
)

func InitTopologyRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/topology/page", TopologyPageListHandler)
	r.Get(common.BASE_CONTEXT+"/topology", TopologyListHandler)
	r.Post(common.BASE_CONTEXT+"/topology", UpsertTopologyHandler)
	r.Delete(common.BASE_CONTEXT+"/topology/{id}", DeleteTopologyHandler)
	r.Post(common.BASE_CONTEXT+"/topology/batch-delete", batchDeleteTopologyHandler)
	r.Post(common.BASE_CONTEXT+"/topology/batch-upsert", batchUpsertTopologyHandler)
	r.Get(common.BASE_CONTEXT+"/topology/groupby", TopologyGroupbyHandler)
}

// @Summary GroupBy
// @Description GroupBy, for example,  _select=level, then return  {level_val1:sum1,level_val2:sum2}, _where can input status=0
// @Tags Topology
// @Param _select query string true "_select"
// @Param _where query string false "_where"
// @Produce  json
// @Success 200 {object} common.Response{data=[]map[string]any} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /topology/groupby [get]
func TopologyGroupbyHandler(w http.ResponseWriter, r *http.Request) {

	common.CommonGroupby(w, r, common.GetDaprClient(), "o_topology")
}

// @Summary batch update
// @Description batch update
// @Tags Topology
// @Accept  json
// @Param entities body []map[string]any true "objects array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /topology/batch-upsert [post]
func batchUpsertTopologyHandler(w http.ResponseWriter, r *http.Request) {

	var entities []map[string]any
	err := common.ReadRequestBody(r, &entities)
	if err != nil {
		common.HttpResult(w, common.ErrParam)
		return
	}
	if len(entities) == 0 {
		common.HttpResult(w, common.ErrParam)
		return
	}

	err = common.DbBatchUpsert[map[string]any](r.Context(), common.GetDaprClient(), entities, model.TopologyTableInfo.Name, model.Topology_FIELD_NAME_id)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}

// @Summary page query
// @Description page query, _page(from 1 begin), _page_size, _order, and others fields, status=1, name=$like.%CAM%
// @Tags Topology
// @Param _page query int true "current page"
// @Param _page_size query int true "page size"
// @Param _order query string false "order"
// @Param id query string false "id"
// @Param created_by query string false "created_by"
// @Param created_time query string false "created_time"
// @Param updated_by query string false "updated_by"
// @Param updated_time query string false "updated_time"
// @Param group_id query string false "group_id"
// @Param file_data query string false "file_data"
// @Param key_name query string false "key_name"
// @Param name query string false "name"
// @Param parent_id query string false "parent_id"
// @Param remark query string false "remark"
// @Param type query string false "type"
// @Produce  json
// @Success 200 {object} common.Response{data=common.Page{items=[]model.Topology}} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /topology/page [get]
func TopologyPageListHandler(w http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("_page")
	pageSize := r.URL.Query().Get("_page_size")
	if page == "" || pageSize == "" {
		common.HttpResult(w, common.ErrParam)
		return
	}
	common.CommonPageQuery[model.Topology](w, r, common.GetDaprClient(), "o_topology", "id")

}

// @Summary query objects
// @Description query objects
// @Tags Topology
// @Param _select query string false "_select"
// @Param _order query string false "order"
// @Param id query string false "id"
// @Param created_by query string false "created_by"
// @Param created_time query string false "created_time"
// @Param updated_by query string false "updated_by"
// @Param updated_time query string false "updated_time"
// @Param group_id query string false "group_id"
// @Param file_data query string false "file_data"
// @Param key_name query string false "key_name"
// @Param name query string false "name"
// @Param parent_id query string false "parent_id"
// @Param remark query string false "remark"
// @Param type query string false "type"
// @Produce  json
// @Success 200 {object} common.Response{data=[]model.Topology} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /topology [get]
func TopologyListHandler(w http.ResponseWriter, r *http.Request) {
	common.CommonQuery[model.Topology](w, r, common.GetDaprClient(), "o_topology", "id")
}

// @Summary save
// @Description save
// @Tags Topology
// @Accept       json
// @Param item body model.Topology true "object"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Topology} "object"
// @Failure 500 {object} common.Response ""
// @Router /topology [post]
func UpsertTopologyHandler(w http.ResponseWriter, r *http.Request) {
	var val model.Topology
	err := common.ReadRequestBody(r, &val)
	if err != nil {
		common.HttpResult(w, common.ErrParam)
		return
	}
	beforeHook, exists := common.GetUpsertBeforeHook("Topology")
	if exists {
		v, err1 := beforeHook(r, val)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
		val = v.(model.Topology)
	}

	err = common.DbUpsert[model.Topology](r.Context(), common.GetDaprClient(), val, model.TopologyTableInfo.Name, "id")
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(val))
}

// @Summary delete
// @Description delete
// @Tags Topology
// @Param id  path string true "实例id"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Topology} "object"
// @Failure 500 {object} common.Response ""
// @Router /topology/{id} [delete]
func DeleteTopologyHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	beforeHook, exists := common.GetDeleteBeforeHook("Topology")
	if exists {
		_, err1 := beforeHook(r, id)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	common.CommonDelete(w, r, common.GetDaprClient(), "o_topology", "id", "id")
}

// @Summary batch delete
// @Description batch delete
// @Tags Topology
// @Accept  json
// @Param ids body []string true "id array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /topology/batch-delete [post]
func batchDeleteTopologyHandler(w http.ResponseWriter, r *http.Request) {

	var ids []string
	err := common.ReadRequestBody(r, &ids)
	if err != nil {
		common.HttpResult(w, common.ErrParam)
		return
	}
	if len(ids) == 0 {
		common.HttpResult(w, common.ErrParam)
		return
	}
	beforeHook, exists := common.GetBatchDeleteBeforeHook("Topology")
	if exists {
		_, err1 := beforeHook(r, ids)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	idstr := strings.Join(ids, ",")
	err = common.DbDeleteByOps(r.Context(), common.GetDaprClient(), "o_topology", []string{"id"}, []string{"in"}, []any{idstr})
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}
