package user

import (
	"context"
	"go-zero/apps/rpc/user/client/other"
	"go-zero/common/result"

	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/apps/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoReq) (resp *result.Result, err error) {

	info, _ := l.svcCtx.OtherClient.UserInfo(l.ctx, &other.InfoReq{
		UserId: l.ctx.Value("userId").(string),
	})

	if info == nil {
		return result.Err().SetMsg("查询用户信息失败！"), nil
	}

	return result.Ok().SetMsg("查询用户信息成功！").SetData(info), nil
}
