package admin

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero/apps/api/gateway/internal/logic/admin"
	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/apps/api/gateway/internal/types"
)

func ChangeUserStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChangeUserStatusReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := admin.NewChangeUserStatusLogic(r.Context(), svcCtx)
		resp, err := l.ChangeUserStatus(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
