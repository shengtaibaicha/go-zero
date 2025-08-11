package filelogic

import (
	"context"
	"encoding/json"
	"go-zero/models"

	"go-zero/apps/rpc/file/file"
	"go-zero/apps/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindPageByTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindPageByTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindPageByTagLogic {
	return &FindPageByTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindPageByTagLogic) FindPageByTag(in *file.FindPageByTagReq) (*file.FindPageByTagResp, error) {
	MDB := l.svcCtx.MDB
	var total int32
	// 计算需要跳过的记录
	offset := (in.Page - 1) * in.Size
	var files []models.Files
	sql := "SELECT * FROM files WHERE file_id IN (SELECT file_id FROM tagAndFile WHERE tag_id = ?) AND status = '已审核' LIMIT ? OFFSET ?"
	MDB.Raw(sql, in.TagId, in.Size, offset).Scan(&files)
	// 获取记录总条数
	MDB.Raw("SELECT count(*) FROM files WHERE file_id IN (SELECT file_id FROM tagAndFile WHERE tag_id = ?) AND status = '已审核'", in.TagId).Scan(&total)
	// 将查询到的数据解析为json格式
	marshaled, _ := json.Marshal(files)

	pages := total / in.Size
	if total%in.Size != 0 {
		pages = (total / in.Size) + 1
	}
	return &file.FindPageByTagResp{
		Records: string(marshaled),
		Size:    in.Size,
		Total:   total,
		Current: in.Page,
		Pages:   pages,
	}, nil
}
