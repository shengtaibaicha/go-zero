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

func (l *AdminFindPageLogic) AdminFindPage(req *types.AdminFindPageReq) (resp *result.Result, err error) {

	page, _ := l.svcCtx.AdminClient.FindPage(l.ctx, &user.AdminFindPageReq{
		Page:   req.Page,
		Size:   req.Size,
		Filter: req.Filter,
	})
	if page == nil {
		return result.Err().SetMsg("查询失败！"), nil
	}

	var data []models.Files
	jsonErr := json.Unmarshal([]byte(page.Records), &data)
	if jsonErr != nil {
		l.Logger.Error("json序列化失败：", jsonErr.Error())
	}

	r := map[string]any{}
	r["records"] = data
	r["total"] = page.GetTotal()
	r["current"] = page.GetCurrent()
	r["pages"] = page.GetPages()
	r["size"] = page.GetSize()

	return result.Ok().SetData(r), nil
}
