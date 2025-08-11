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

type FindPageByNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindPageByNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindPageByNameLogic {
	return &FindPageByNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindPageByNameLogic) FindPageByName(req *types.FindPageByNameReq) (resp *result.Result, err error) {

	nameData, _ := l.svcCtx.FileClient.FindPageByName(
		l.ctx,
		&file.FindPageByNameReq{
			Page: req.Page,
			Size: req.Size,
			Name: req.Name,
		},
	)

	var data []models.Files

	json.Unmarshal([]byte(nameData.Records), &data)
	r := map[string]any{}
	r["records"] = data
	r["total"] = nameData.Total
	r["current"] = nameData.Current
	r["pages"] = nameData.Pages
	r["size"] = nameData.Size

	return result.Ok().SetData(r), nil
}
