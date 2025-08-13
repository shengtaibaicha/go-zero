package file

import (
	"context"
	"go-zero/apps/rpc/file/file"
	"go-zero/common/result"

	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/apps/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.DeleteReq) (resp *result.Result, err error) {

	// 拿到上下文的userId并通过metadata传递到rpc
	md := metadata.New(map[string]string{
		"userId": l.ctx.Value("userId").(string),
	})
	l.ctx = metadata.NewOutgoingContext(l.ctx, md)
	// 调用rpc服务
	deleteFile, _ := l.svcCtx.FileClient.DeleteFile(l.ctx, &file.DeleteFileReq{
		FileId: req.FileId,
	})
	// 判断rpc返回的结果，根据结果返回
	if deleteFile.Success != true {
		return result.Err().SetMsg(deleteFile.GetMsg()), nil
	}

	return result.Ok().SetMsg(deleteFile.GetMsg()), nil
}
