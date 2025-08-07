package user

import (
	"context"
	"go-zero/apps/api/community/internal/svc"
	"go-zero/apps/api/community/internal/types"
	"go-zero/apps/rpc/user/user"
	"go-zero/common/result"

	"github.com/zeromicro/go-zero/core/logx"
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
	// 调用rpc层进行验证和数据库操作
	register, err := l.svcCtx.RegisterClient.UserRegister(l.ctx, &user.RegisterReq{
		UserName:     req.UserName,
		UserPassword: req.Password,
		UserEmail:    req.UserEmail,
	})
	if err != nil {
		return result.Err().SetMsg(err.Error()), nil
	}
	return result.Ok().SetData(register), nil
}
