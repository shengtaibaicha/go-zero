package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero/apps/user/api/internal/svc"
	"go-zero/apps/user/api/internal/types"
	"go-zero/apps/user/rpc/user"
	"go-zero/common/result"
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
	// 调用rpc层进行验证和数据库操作
	register, err := l.svcCtx.Register.UserRegister(l.ctx, &user.RegisterReq{
		UserName:     req.UserName,
		UserPassword: req.Password,
		UserEmail:    req.UserEmail,
	})
	if err != nil {
		return result.Err().SetMsg(err.Error()), nil
	}
	return result.Ok().SetData(register), nil
}
