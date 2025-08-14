package public

import (
	"context"
	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/apps/api/gateway/internal/types"
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

func (l *UserLoginLogic) UserLogin(req *types.UserLoginReq, redisKey string) (resp *result.Result, err error) {
	// 先验证接收到的验证码与redis中的验证码是否相等
	get, _ := l.svcCtx.RedisClient.Get(redisKey)
	if get == "" || get != req.CaptchaCode {
		return result.Err().SetMsg("验证码不正确"), nil
	}
	// 验证成功后删除redis中的验证码
	l.svcCtx.RedisClient.DelCtx(l.ctx, redisKey)
	token, err := l.svcCtx.UserClient.UserLogin(l.ctx, &user.LoginReq{
		UserName:     req.UserName,
		UserPassword: req.UserPassword,
	})
	if err != nil {
		return result.Err().SetMsg(err.Error()), nil
	}
	return result.Ok().SetData(token), nil
}
