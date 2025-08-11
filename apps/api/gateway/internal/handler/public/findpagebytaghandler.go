package public

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero/apps/api/gateway/internal/logic/public"
	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/apps/api/gateway/internal/types"
)

func FindPageByTagHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FindPageByTagReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := public.NewFindPageByTagLogic(r.Context(), svcCtx)
		resp, err := l.FindPageByTag(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
