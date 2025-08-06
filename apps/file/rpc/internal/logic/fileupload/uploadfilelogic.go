package fileuploadlogic

import (
	"context"

	"go-zero/apps/file/rpc/file"
	"go-zero/apps/file/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 单请求上传（适合中小文件）
func (l *UploadFileLogic) UploadFile(in *file.UploadReq) (*file.UploadResponse, error) {
	// todo: add your logic here and delete this line

	return &file.UploadResponse{}, nil
}
