package file

import (
	"go-zero/apps/api/community/internal/logic/file"
	"go-zero/apps/api/community/internal/svc"
	"go-zero/apps/api/community/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := file.NewUploadLogic(r.Context(), svcCtx)
		resp, err := l.Upload(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
