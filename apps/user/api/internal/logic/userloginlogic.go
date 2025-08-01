package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero/apps/user/api/internal/svc"
	"go-zero/apps/user/api/internal/types"
	"go-zero/apps/user/pkg/result"
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

	return
}
