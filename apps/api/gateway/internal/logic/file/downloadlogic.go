package file

import (
	"context"
	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/apps/api/gateway/internal/types"
	"go-zero/apps/rpc/file/file"
	"go-zero/common/middleware"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
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

	md := metadata.New(map[string]string{"userId": middleware.GetUserIdFromCtx(l.ctx)})
	outgoingContext := metadata.NewOutgoingContext(l.ctx, md)

	downloadFile, err := l.svcCtx.FileClient.DownloadFile(outgoingContext, &file.DownloadFileReq{
		FileName: req.FileName,
	})
	if err != nil {
		return nil, err
	}
	return downloadFile, nil
}
