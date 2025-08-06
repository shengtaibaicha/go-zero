package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero/apps/file/api/internal/svc"
	rpcfile "go-zero/apps/file/rpc/file"
	"go-zero/common/result"
	"io"
	"net/http"
)

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadLogic) Upload(r *http.Request) (resp *result.Result, err error) {
	// todo: add your logic here and delete this line
	// 1. 解析前端上传的文件（表单字段名为"file"）
	file, header, err := r.FormFile("file")
	if err != nil {
		return result.Err().SetMsg("解析文件失败！"), err
	}
	defer file.Close() // 确保文件流关闭

	// 2. 读取文件内容为字节数组（适合中小文件）
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return result.Err().SetMsg("转换文件格式失败！"), err
	}

	f := &rpcfile.UploadReq{
		File:     fileContent,
		Filename: header.Filename,
		Size:     header.Size,
		MimeType: header.Header.Get("Content-Type"),
	}

	uploadFile, err := l.svcCtx.FileClient.UploadFile(l.ctx, f)
	if err != nil {
		return result.Err().SetMsg("文件上传失败！"), err
	}

	fmt.Println(uploadFile)

	return
}
