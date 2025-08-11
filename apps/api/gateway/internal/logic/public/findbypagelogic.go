package public

import (
	"context"
	"encoding/json"
	"go-zero/apps/rpc/file/file"
	"go-zero/common/result"
	"go-zero/models"

	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/apps/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindByPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByPageLogic {
	return &FindByPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindByPageLogic) FindByPage(req *types.FindByPageReq) (resp *result.Result, err error) {
	page, err := l.svcCtx.FileClient.FindByPage(l.ctx, &file.FindByPageReq{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	// 将json字符串序列化为对象
	var data []models.Files
	json.Unmarshal([]byte(page.Records), &data)
	r := map[string]any{}
	r["records"] = data
	r["total"] = page.Total
	r["current"] = page.Current
	r["pages"] = page.Pages
	r["size"] = page.Size
	return result.Ok().SetData(r), nil
}
