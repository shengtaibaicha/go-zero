package adminlogic

import (
	"context"
	"encoding/json"
	"go-zero/models"
	"time"

	"go-zero/apps/rpc/user/internal/svc"
	"go-zero/apps/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
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
		MDB.Model(&models.Users{}).Where("role != ?", "superAdmin").Offset(int(offset)).Limit(int(in.Size)).Scan(&data)
		MDB.Model(&models.Users{}).Where("role != ?", "superAdmin").Count(&total)
	} else if in.Role == "admin" {
		MDB.Model(&models.Users{}).Where("role = ?", "admin").Offset(int(offset)).Limit(int(in.Size)).Scan(&data)
		MDB.Model(&models.Users{}).Where("role = ?", "admin").Count(&total)
	} else {
		MDB.Model(&models.Users{}).Where("role = ?", "user").Offset(int(offset)).Limit(int(in.Size)).Scan(&data)
		MDB.Model(&models.Users{}).Where("role = ?", "user").Count(&total)
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
