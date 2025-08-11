package file

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"go-zero/apps/api/gateway/internal/logic/file"
	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/apps/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DownloadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DownloadFileReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := file.NewDownloadLogic(r.Context(), svcCtx)
		data, err := l.Download(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}

		if data == nil {
			httpx.Ok(w)
		}

		//w.Header().Set("Content-Disposition", "attachment; filename=\""+data.FileName+"\"") // 强制下载并指定文件名
		// 正确设置下载响应头
		w.Header().Set("Content-Disposition", func() string {
			// 对文件名进行 URL 编码，处理中文等特殊字符
			encodedFileName := url.QueryEscape(data.FileName)
			// 将编码后的 + 替换为 %20（空格的标准编码）
			encodedFileName = strings.ReplaceAll(encodedFileName, "+", "%20")
			// 拼接 Content-Disposition 头
			return "attachment; filename=\"" + encodedFileName + "\""
		}())
		w.Header().Set("Content-Type", data.ContentType)                  // 设置MIME类型
		w.Header().Set("Content-Length", strconv.Itoa(len(data.Content))) // 设置文件大小
		w.WriteHeader(http.StatusOK)                                      // 200状态码

		// 4. 将二进制流写入响应体
		if _, err := w.Write(data.Content); err != nil {
			l.Errorf("写入文件流失败: %v", err)
		}

	}
}
