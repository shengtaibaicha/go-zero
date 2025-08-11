package public

import (
	"go-zero/apps/api/gateway/internal/logic/public"
	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/apps/api/gateway/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := public.NewUserLoginLogic(r.Context(), svcCtx)
		resp, err := l.UserLogin(&req, r.Header.Get("redisKey"))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
