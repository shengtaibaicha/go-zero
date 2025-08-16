package admin

import (
	"context"
	"go-zero/apps/rpc/user/user"
	"go-zero/common/middleware"
	"go-zero/common/result"

	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/apps/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type ChangeUserStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeUserStatusLogic {
	return &ChangeUserStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeUserStatusLogic) ChangeUserStatus(req *types.ChangeUserStatusReq) (resp *result.Result, err error) {

	md := metadata.New(map[string]string{"role": middleware.GetRoleFromCtx(l.ctx)})
	outgoingContext := metadata.NewOutgoingContext(l.ctx, md)

	status, _ := l.svcCtx.AdminClient.ChangeUserStatus(outgoingContext, &user.ChangeUserStatusReq{
		UserId: req.UserId,
	})

	if status.Success {
		return result.Ok().SetMsg(status.Msg), nil
	}
	return result.Err().SetMsg(status.Msg), nil
}
