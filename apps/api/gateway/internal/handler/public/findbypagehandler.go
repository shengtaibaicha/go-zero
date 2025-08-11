package public

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero/apps/api/gateway/internal/logic/public"
	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/apps/api/gateway/internal/types"
)

func FindByPageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FindByPageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := public.NewFindByPageLogic(r.Context(), svcCtx)
		resp, err := l.FindByPage(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
