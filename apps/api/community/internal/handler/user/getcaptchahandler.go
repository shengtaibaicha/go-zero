package user

import (
	"go-zero/apps/api/community/internal/logic/user"
	"go-zero/apps/api/community/internal/svc"
	"go-zero/apps/api/community/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetCaptchaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCaptchaReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewGetCaptchaLogic(r.Context(), svcCtx)
		resp, redisKey, err := l.GetCaptcha(&req)
		w.Header().Set("redisKey", redisKey)

		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
