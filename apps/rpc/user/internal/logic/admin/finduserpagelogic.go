package adminlogic

import (
	"context"
	"encoding/json"
	"errors"
	"go-zero/models"
	"time"

	"go-zero/apps/rpc/user/internal/svc"
	"go-zero/apps/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type FindUserPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserPageLogic {
	return &FindUserPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserPageLogic) FindUserPage(in *user.AdminFindUserPageReq) (*user.AdminFindUserPageResp, error) {

	incomingContext, ok := metadata.FromIncomingContext(l.ctx)
	if !ok {
		l.Logger.Error("metadata FromIncomingContext error")
		return nil, errors.New("metadata FromIncomingContext error")
	}
	userId := incomingContext.Get("userId")[0]

	MDB := l.svcCtx.MDB

	offset := (in.Page - 1) * in.Size

	type u struct {
		UserId    string
		UserName  string
		UserEmail string
		JoinDate  time.Time
		Role      string
		Enable    int
	}

	var data []u
	var total int64

	if in.Role == "all" {
		MDB.Model(&models.Users{}).Where("role != ? and user_id != ?", "superAdmin", userId).Offset(int(offset)).Limit(int(in.Size)).Scan(&data)
		MDB.Model(&models.Users{}).Where("role != ? and user_id != ?", "superAdmin", userId).Count(&total)
	} else if in.Role == "admin" {
		MDB.Model(&models.Users{}).Where("role = ? and user_id != ?", "admin", userId).Offset(int(offset)).Limit(int(in.Size)).Scan(&data)
		MDB.Model(&models.Users{}).Where("role = ? and user_id != ?", "admin", userId).Count(&total)
	} else {
		MDB.Model(&models.Users{}).Where("role = ? and user_id != ?", "user", userId).Offset(int(offset)).Limit(int(in.Size)).Scan(&data)
		MDB.Model(&models.Users{}).Where("role = ? and user_id != ?", "user", userId).Count(&total)
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	pages := int32(total) / in.Size
	if int32(total)%in.Size != 0 {
		pages = (int32(total) / in.Size) + 1
	}

	return &user.AdminFindUserPageResp{
		Records: string(marshal),
		Total:   int32(total),
		Size:    in.Size,
		Current: in.Page,
		Pages:   pages,
	}, nil
}
