package public

import (
	"context"
	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/apps/api/gateway/internal/types"
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

func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterReq, redisKey string) (resp *result.Result, err error) {
	// todo: add your logic here and delete this line
	// 先验证接收到的验证码与redis中的验证码是否相等
	get, _ := l.svcCtx.RedisClient.Get(redisKey)
	if get == "" || get != req.CaptchaCode {
		return result.Err().SetMsg("验证码不正确"), nil
	}
	// 验证成功后删除redis中的验证码
	l.svcCtx.RedisClient.DelCtx(l.ctx, redisKey)
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
