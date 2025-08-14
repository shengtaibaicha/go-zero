package userlogic

import (
	"context"
	"go-zero/models"

	"go-zero/apps/rpc/user/internal/svc"
	"go-zero/apps/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *user.InfoReq) (*user.InfoResp, error) {

	MDB := l.svcCtx.MDB

	// 查询用户的基本信息
	var userInfo models.Users
	MDB.Model(&models.Users{}).Where("user_id = ?", in.UserId).First(&userInfo)

	return &user.InfoResp{
		UserName:       userInfo.UserName,
		UserEmail:      userInfo.UserEmail,
		UserAvatar:     userInfo.UserAvatar,
		JoinDate:       userInfo.JoinDate.String(),
		UploadNumber:   userInfo.UploadNumber,
		DownloadNumber: userInfo.DownloadNumber,
		CollectNumber:  userInfo.CollectNumber,
		Role:           userInfo.Role,
		Enable:         int32(userInfo.Enable),
	}, nil
}
