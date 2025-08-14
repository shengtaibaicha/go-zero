package admin

import (
	"context"
	"encoding/json"
	"go-zero/apps/rpc/user/user"
	"go-zero/common/result"
	"go-zero/models"

	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/apps/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminFindPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminFindPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminFindPageLogic {
	return &AdminFindPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminFindPageLogic) AdminFindPage(req *types.FindByPageReq) (resp *result.Result, err error) {

	page, _ := l.svcCtx.AdminClient.FindPage(l.ctx, &user.AdminFindPageReq{
		Page: req.Page,
		Size: req.Size,
	})
	if page == nil {
		return result.Err().SetMsg("查询失败！"), nil
	}

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
