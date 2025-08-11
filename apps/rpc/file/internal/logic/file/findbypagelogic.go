package filelogic

import (
	"context"
	"encoding/json"
	"go-zero/models"

	"go-zero/apps/rpc/file/file"
	"go-zero/apps/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByPageLogic {
	return &FindByPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByPageLogic) FindByPage(in *file.FindByPageReq) (*file.FindByPageResp, error) {
	MDB := l.svcCtx.MDB
	var total int64
	// 计算需要跳过的记录
	offset := (in.Page - 1) * in.Size
	var files []models.Files
	MDB.Offset(int(offset)).Limit(int(in.Size)).Find(&files)
	// 获取记录总条数
	MDB.Model(&files).Count(&total)
	// 将查询到的数据解析为json格式
	marshaled, _ := json.Marshal(files)

	pages := int32(total) / in.Size
	if int32(total)%in.Size != 0 {
		pages = (int32(total) / in.Size) + 1
	}
	return &file.FindByPageResp{
		Records: string(marshaled),
		Size:    in.Size,
		Total:   int32(total),
		Current: in.Page,
		Pages:   pages,
	}, nil
}
