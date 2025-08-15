package adminlogic

import (
	"context"
	"encoding/json"
	"go-zero/models"

	"go-zero/apps/rpc/user/internal/svc"
	"go-zero/apps/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindPageLogic {
	return &FindPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindPageLogic) FindPage(in *user.AdminFindPageReq) (*user.AdminFindPageResp, error) {

	MDB := l.svcCtx.MDB

	var datas []models.Files
	var total int64
	offset := (in.Page - 1) * in.Size
	if in.Filter == "all" {
		MDB.Model(&models.Files{}).Offset(int(offset)).Limit(int(in.Size)).Find(&datas)
		MDB.Model(&models.Files{}).Count(&total)
	} else if in.Filter == "audited" {
		MDB.Model(&models.Files{}).Where("status = ?", "已审核").Offset(int(offset)).Limit(int(in.Size)).Find(&datas)
		MDB.Model(&models.Files{}).Where("status = ?", "已审核").Count(&total)
	} else {
		MDB.Model(&models.Files{}).Where("status = ?", "未审核").Offset(int(offset)).Limit(int(in.Size)).Find(&datas)
		MDB.Model(&models.Files{}).Where("status = ?", "未审核").Count(&total)
	}

	marshal, _ := json.Marshal(datas)

	pages := int32(total) / in.Size
	if int32(total)%in.Size != 0 {
		pages = (int32(total) / in.Size) + 1
	}

	return &user.AdminFindPageResp{
		Records: string(marshal),
		Total:   int32(total),
		Size:    in.Size,
		Current: in.Page,
		Pages:   pages,
	}, nil
}
