package public

import (
	"context"
	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/apps/api/gateway/internal/types"
	"go-zero/apps/api/gateway/tools/captcha"
	"go-zero/common/result"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
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
	captcha, code := myCaptcha.GetCaptcha()
	// 生成redisKey
	id := uuid.New()
	redisKey = "captcha:" + id.String()
	// 将验证码存入redis
	l.svcCtx.RedisClient.Setex(redisKey, code, 60*5)
	resp = result.Ok().SetData(captcha)
	return
}
