package admin

import (
	"context"
	"go-zero/apps/rpc/file/file"
	"go-zero/common/middleware"
	"go-zero/common/result"

	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/apps/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type AdminAuditFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminAuditFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminAuditFileLogic {
	return &AdminAuditFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminAuditFileLogic) AdminAuditFile(req *types.AuditReq) (resp *result.Result, err error) {

	md := metadata.New(map[string]string{"role": middleware.GetRoleFromCtx(l.ctx)})
	outgoingContext := metadata.NewOutgoingContext(l.ctx, md)

	auditFile, _ := l.svcCtx.FileClient.AuditFile(outgoingContext, &file.AuditFileReq{
		FileId:  req.FileId,
		Audited: req.Audited,
	})
	if auditFile.GetSuccess() {
		return result.Ok().SetMsg(auditFile.GetMsg()), nil
	} else {
		l.Logger.Error(auditFile.Msg)
	}

	return result.Err().SetMsg(auditFile.GetMsg()), nil
}
