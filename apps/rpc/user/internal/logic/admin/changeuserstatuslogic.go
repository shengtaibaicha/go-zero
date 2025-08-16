package adminlogic

import (
	"context"
	"go-zero/models"

	"go-zero/apps/rpc/user/internal/svc"
	"go-zero/apps/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type ChangeUserStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeUserStatusLogic {
	return &ChangeUserStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangeUserStatusLogic) ChangeUserStatus(in *user.ChangeUserStatusReq) (*user.CommonResp, error) {

	db := l.svcCtx.MDB.Begin()

	incomingContext, ok := metadata.FromIncomingContext(l.ctx)
	if !ok {
		l.Logger.Error("metadata FromIncomingContext error")
		return &user.CommonResp{
			Success: false,
			Msg:     "metadata FromIncomingContext error",
		}, nil
	}

	role := incomingContext.Get("role")[0]

	var u models.Users
	first := db.Model(&models.Users{}).Where("user_id = ?", in.UserId).First(&u)
	if first.Error != nil {
		db.Rollback()
		return &user.CommonResp{
			Success: false,
			Msg:     "操作的用户不存在！",
		}, nil
	}

	if (role == "superAdmin" && (u.Role == "user" || u.Role == "admin")) || role == "admin" && u.Role == "user" {
		if u.Enable == 0 {
			u.Enable = 1
		} else {
			u.Enable = 0
		}
	} else {
		return &user.CommonResp{
			Success: false,
			Msg:     "用户权限和执行的操作不匹配！",
		}, nil
	}

	tx := db.Model(&models.Users{}).Where("user_id = ?", in.UserId).Update("enable", u.Enable)
	if tx.Error != nil {
		db.Rollback()
		return &user.CommonResp{
			Success: false,
			Msg:     tx.Error.Error(),
		}, nil
	}

	db.Commit()

	return &user.CommonResp{
		Success: true,
		Msg:     "修改用户状态成功！",
	}, nil
}
