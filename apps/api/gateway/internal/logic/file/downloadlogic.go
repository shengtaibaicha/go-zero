package file

import (
	"context"
	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/apps/api/gateway/internal/types"
	"go-zero/apps/rpc/file/file"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadLogic {
	return &DownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadLogic) Download(req *types.DownloadFileReq) (resp *file.DownloadFileResp, err error) {

	downloadFile, err := l.svcCtx.FileClient.DownloadFile(context.Background(), &file.DownloadFileReq{
		FileName: req.FileName,
	})
	if err != nil {
		return nil, err
	}
	return downloadFile, nil
}
