package file

import (
	"context"
	"go-zero/apps/api/gateway/internal/svc"
	rpcfile "go-zero/apps/rpc/file/file"
	"go-zero/common/result"
	"io"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
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
	// 限制文件大小
	err = r.ParseMultipartForm(30 << 20) // 10*1024*1024
	if err != nil {
		return result.Err().SetMsg("文件大小超过限制!"), nil
	}
	// 1. 解析前端上传的文件（表单字段名为"file"）
	file, header, err := r.FormFile("file")
	// 判断文件格式
	contentType := header.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		// 非允许的格式，返回错误
		return result.Err().SetMsg("文件格式不符合!"), nil
	}

	tagId := r.FormValue("tagId")
	if err != nil {
		return result.Err().SetMsg("解析文件失败！"), err
	}
	defer file.Close() // 确保文件流关闭

	// 2. 读取文件内容为字节数组（适合中小文件）
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return result.Err().SetMsg("转换文件格式失败！"), err
	}

	atoi, _ := strconv.Atoi(tagId)
	f := &rpcfile.UploadReq{
		File:     fileContent,
		Filename: header.Filename,
		Size:     header.Size,
		MimeType: header.Header.Get("Content-Type"),
		UserId:   l.ctx.Value("userId").(string),
		TagId:    int32(atoi),
	}

	// metadata传递数据
	//md := metadata.New(map[string]string{"aaa": "dada"})
	//l.ctx = metadata.NewOutgoingContext(l.ctx, md)

	uploadFile, err := l.svcCtx.FileClient.UploadFile(l.ctx, f)

	if uploadFile.Success != true {
		if uploadFile != nil {
			return result.Err().SetMsg(uploadFile.Message), err
		}
	}

	return result.Ok().SetData(map[string]string{
		"url": uploadFile.FileUrl,
	}), nil
}
