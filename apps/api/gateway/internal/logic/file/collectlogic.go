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

type CollectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCollectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CollectLogic {
	return &CollectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CollectLogic) Collect(req *types.CollectReq) (resp *result.Result, err error) {

	md := metadata.New(map[string]string{"userId": l.ctx.Value("userId").(string)})
	outgoingContext := metadata.NewOutgoingContext(l.ctx, md)

	collectFile, _ := l.svcCtx.FileClient.CollectFile(outgoingContext, &file.CollectFileReq{FileId: req.FileId})
	if collectFile.Success != true {
		return result.Err().SetMsg(collectFile.Msg), nil
	}

	return result.Ok().SetMsg(collectFile.Msg), nil
}
