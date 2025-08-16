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

	// 开启gorm事务（修复：检查事务初始化错误）
	db := l.svcCtx.MDB.Begin()
	if db.Error != nil {
		return &file.CollectFileResp{
			Success: false,
			Msg:     "事务启动失败: " + db.Error.Error(),
		}, nil
	}
	defer func() {
		if r := recover(); r != nil {
			// 恐慌时回滚，并记录日志
			if err := db.Rollback().Error; err != nil {
				l.Logger.Error("恐慌时回滚失败: ", err)
			}
		}
	}()

	incomingContext, ok := metadata.FromIncomingContext(l.ctx)
	if !ok {
		return &file.CollectFileResp{
			Success: false,
			Msg:     "操作失败！",
		}, nil
	}
	userId := incomingContext.Get("userId")[0]

	var r models.Collect
	tx := db.Model(&models.Collect{}).Where("user_id = ? and file_id = ?", userId, in.FileId).First(&r)
	if tx.RowsAffected == 0 || errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		// 当收藏表里面没有匹配的记录则写入
		create := db.Model(&models.Collect{}).Create(&models.Collect{FileId: in.FileId, UserId: userId})
		if create.Error != nil {
			db.Rollback()
		}
		update := db.Model(&models.Users{}).Where("user_id = ?", userId).Update("collect_number", gorm.Expr("collect_number + ?", 1))
		if err := update.Error; err != nil {
			db.Rollback()
		}
	} else {
		// 当收藏表里面有匹配的记录则删除
		t := db.Where("user_id = ? and file_id = ?", userId, in.FileId).Delete(&models.Collect{})
		if err := t.Error; err != nil {
			db.Rollback()
		}
		update := db.Model(&models.Users{}).Where("user_id = ?", userId).Update("collect_number", gorm.Expr("collect_number - 1"))
		if err := update.Error; err != nil {
			db.Rollback()
		}
	}

	db.Commit()

	return &file.CollectFileResp{
		Success: true,
		Msg:     "操作成功！",
	}, nil
}
