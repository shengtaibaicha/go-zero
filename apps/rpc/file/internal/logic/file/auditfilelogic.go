package filelogic

import (
	"context"
	"go-zero/models"

	"go-zero/apps/rpc/file/file"
	"go-zero/apps/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuditFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAuditFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuditFileLogic {
	return &AuditFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AuditFileLogic) AuditFile(in *file.AuditFileReq) (*file.AuditFileResp, error) {

	db := l.svcCtx.MDB.Begin()

	tx := db.First(&models.Files{})
	if tx.Error != nil {
		return &file.AuditFileResp{
			Success: false,
			Msg:     "审核失败：" + tx.Error.Error(),
		}, nil
	}

	var status string
	if in.Audited == "未审核" {
		status = "已审核"
	} else {
		status = "未审核"
	}

	update := db.Model(&models.Files{}).Where("file_id = ?", in.FileId).Update("status", status)
	if update.Error != nil {
		db.Rollback()
	}

	db.Commit()

	if status == "未审核" {
		return &file.AuditFileResp{
			Success: true,
			Msg:     "审核成功！",
		}, nil
	} else {
		return &file.AuditFileResp{
			Success: true,
			Msg:     "取消审核成功！",
		}, nil
	}
}
