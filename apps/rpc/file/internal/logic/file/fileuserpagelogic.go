package filelogic

import (
	"context"
	"encoding/json"
	"go-zero/models"

	"go-zero/apps/rpc/file/file"
	"go-zero/apps/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type FileUserPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFileUserPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUserPageLogic {
	return &FileUserPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FileUserPageLogic) FileUserPage(in *file.FileUserPageReq) (*file.FileUserPageResp, error) {

	MDB := l.svcCtx.MDB
	// 计算需要跳过的记录
	offset := (in.Page - 1) * in.Size

	incomingContext, ok := metadata.FromIncomingContext(l.ctx)
	if !ok {
		return nil, nil
	}

	userId := incomingContext.Get("userId")

	// 查询用户上传的图片
	var r []models.Files
	var total int64
	MDB.Model(&models.Files{}).Where("user_id = ?", userId).Offset(int(offset)).Limit(int(in.Size)).Find(&r)
	MDB.Model(&models.Files{}).Where("user_id = ?", userId).Count(&total)

	marshal, _ := json.Marshal(r)

	pages := int32(total) / in.Size
	if int32(total)%in.Size != 0 {
		pages = (int32(total) / in.Size) + 1
	}
	return &file.FileUserPageResp{
		Records: string(marshal),
		Total:   int32(total),
		Size:    in.Size,
		Current: in.Page,
		Pages:   pages,
	}, nil
}
