package service

import (
	"context"
	"github.com/dapr-platform/common"
	"vis-service/entity"
)

func FetchDatabase(ctx context.Context, req *entity.QueryDbReq) (data []map[string]any, err error) {
	data, err = common.CustomSql[map[string]any](ctx, common.GetDaprClient(), req.Select, req.From, req.Where)
	return
}
