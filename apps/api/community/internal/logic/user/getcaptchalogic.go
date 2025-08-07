package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero/apps/api/community/internal/svc"
	"go-zero/apps/api/community/internal/types"
	"go-zero/apps/api/community/tools/captcha"
	"go-zero/common/result"
)

type GetCaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCaptchaLogic {
	return &GetCaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCaptchaLogic) GetCaptcha(req *types.GetCaptchaReq) (resp *result.Result, redisKey string, err error) {
	// todo: add your logic here and delete this line
	captcha, code := myCaptcha.GetCaptcha()
	// 生成redisKey
	id := uuid.New()
	redisKey = "captcha:" + id.String()
	// 将验证码存入redis
	l.svcCtx.RedisClient.Setex(redisKey, code, 60*5)
	resp = result.Ok().SetData(captcha)
	return
}
