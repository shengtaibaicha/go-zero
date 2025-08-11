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

type FindPageByTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindPageByTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindPageByTagLogic {
	return &FindPageByTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindPageByTagLogic) FindPageByTag(req *types.FindPageByTagReq) (resp *result.Result, err error) {

	tagData, _ := l.svcCtx.FileClient.FindPageByTag(l.ctx, &file.FindPageByTagReq{
		Page:  req.Page,
		Size:  req.Size,
		TagId: req.TagId,
	})

	var data []models.Files

	json.Unmarshal([]byte(tagData.Records), &data)
	r := map[string]any{}
	r["records"] = data
	r["total"] = tagData.Total
	r["current"] = tagData.Current
	r["pages"] = tagData.Pages
	r["size"] = tagData.Size

	return result.Ok().SetData(r), nil
}
