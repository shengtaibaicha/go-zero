package filelogic

import (
	"context"
	"errors"
	"go-zero/models"

	"go-zero/apps/rpc/file/file"
	"go-zero/apps/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type CollectFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCollectFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CollectFileLogic {
	return &CollectFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CollectFileLogic) CollectFile(in *file.CollectFileReq) (*file.CollectFileResp, error) {

	MDB := l.svcCtx.MDB

	incomingContext, ok := metadata.FromIncomingContext(l.ctx)
	if !ok {
		return &file.CollectFileResp{
			Success: false,
			Msg:     "操作失败！",
		}, nil
	}
	userId := incomingContext.Get("userId")

	var r models.Collect
	tx := MDB.Model(&models.Collect{}).Where("user_id = ? and file_id = ?", userId, in.FileId).First(&r)
	if tx.RowsAffected == 0 || errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		// 当收藏表里面没有匹配的记录则写入
		MDB.Model(&models.Collect{}).Create(&models.Collect{FileId: in.FileId, UserId: userId[0]})
		MDB.Model(&models.Users{}).Where("user_id = ?", userId[0]).Update("collect_number", gorm.Expr("collect_number + 1"))
	} else {
		// 当收藏表里面有匹配的记录则删除
		MDB.Where("user_id = ? and file_id = ?", userId, in.FileId).Delete(&models.Collect{})
		MDB.Model(&models.Users{}).Where("user_id = ?", userId[0]).Update("collect_number", gorm.Expr("collect_number - 1"))
	}

	return &file.CollectFileResp{
		Success: true,
		Msg:     "操作成功！",
	}, nil
}
