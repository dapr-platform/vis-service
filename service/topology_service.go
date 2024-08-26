package service

import (
	"github.com/dapr-platform/common"
	"net/http"
	"time"
	"vis-service/model"
)

func init() {
	common.RegisterUpsertBeforeHook("Topology", ProcessUpsertTopology)
}

func ProcessUpsertTopology(r *http.Request, in any) (out any, err error) {
	sub, _ := common.ExtractUserSub(r)

	topology := in.(model.Topology)
	if topology.ID == "" {
		topology.ID = common.NanoId()
	}
	if topology.CreatedBy == "" {
		topology.CreatedBy = sub
	}
	zeroTime, _ := time.Parse("2006-01-02 15:04:05", "0001-01-01 00:00:00")
	if topology.CreatedTime == common.LocalTime(zeroTime) {
		topology.CreatedTime = common.LocalTime(time.Now())
	}
	topology.UpdatedBy = sub
	topology.UpdatedTime = common.LocalTime(time.Now())

	out = topology
	return
}
