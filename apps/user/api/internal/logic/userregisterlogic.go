package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero/apps/user/api/internal/svc"
	"go-zero/apps/user/api/internal/types"
	"go-zero/apps/user/pkg/result"
	"go-zero/apps/user/rpc/user"
)

type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterReq) (resp *result.Result, err error) {
	// todo: add your logic here and delete this line
	register, err := l.svcCtx.Regiter.UserRegister(l.ctx, &user.RegisterReq{
		UserName:     req.UserName,
		UserPassword: req.Password,
		UserEmail:    req.UserEmail,
	})
	if err != nil {
		return result.Err().SetMsg(err.Error()), nil
	}
	return result.Ok().SetData(register), nil
}
