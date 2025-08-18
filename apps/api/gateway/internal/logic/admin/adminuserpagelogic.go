package admin

import (
	"context"
	"encoding/json"
	"go-zero/apps/rpc/user/user"
	"go-zero/common/middleware"
	"go-zero/common/result"
	"time"

	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/apps/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type AdminUserPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUserPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUserPageLogic {
	return &AdminUserPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUserPageLogic) AdminUserPage(req *types.AdminUserPageReq) (resp *result.Result, err error) {

	md := metadata.New(map[string]string{"userId": middleware.GetUserIdFromCtx(l.ctx)})
	outgoingContext := metadata.NewOutgoingContext(l.ctx, md)

	type u struct {
		UserId    string    `json:"userId"`
		UserName  string    `json:"userName"`
		UserEmail string    `json:"userEmail"`
		JoinDate  time.Time `json:"joinDate"`
		Role      string    `json:"role"`
		Enable    int       `json:"enable"`
	}

	var data []u

	page, err := l.svcCtx.AdminClient.FindUserPage(outgoingContext, &user.AdminFindUserPageReq{
		Page: req.Page,
		Size: req.Size,
		Role: req.Role,
	})
	if err != nil {
		l.Logger.Error(err)
		return result.Err().SetMsg("获取用户列表错误！"), nil
	}

	err = json.Unmarshal([]byte(page.GetRecords()), &data)
	if err != nil {
		l.Logger.Error(err)
		return result.Err().SetMsg("解析json数据错误！"), nil
	}

	r := map[string]any{}
	if data == nil {
		data = []u{}
	}
	r["records"] = data
	r["total"] = page.GetTotal()
	r["current"] = page.GetCurrent()
	r["pages"] = page.GetPages()
	r["size"] = page.GetSize()

	return result.Ok().SetData(r), nil
}
