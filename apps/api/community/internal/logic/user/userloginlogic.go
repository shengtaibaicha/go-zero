package user

import (
	"context"
	"go-zero/apps/api/community/internal/svc"
	"go-zero/apps/api/community/internal/types"
	"go-zero/apps/rpc/user/user"
	"go-zero/common/result"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.UserLoginReq) (resp *result.Result, err error) {
	// todo: add your logic here and delete this line
	token, err := l.svcCtx.LoginClient.UserLogin(l.ctx, &user.LoginReq{
		UserName:     req.UserName,
		UserPassword: req.UserPassword,
	})
	if err != nil {
		return result.Err().SetMsg(err.Error()), nil
	}
	return result.Ok().SetData(token), nil
}
