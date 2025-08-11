package public

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero/apps/api/gateway/internal/logic/public"
	"go-zero/apps/api/gateway/internal/svc"
)

func GetTagsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewGetTagsLogic(r.Context(), svcCtx)
		resp, err := l.GetTags()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
