package filelogic

import (
	"context"
	"go-zero/models"

	"go-zero/apps/rpc/file/file"
	"go-zero/apps/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
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

	incomingContext, ok := metadata.FromIncomingContext(l.ctx)
	if !ok {
		l.Logger.Error("metadata.FromIncomingContext")
	}

	role := incomingContext.Get("role")[0]

	db := l.svcCtx.MDB.Begin()

	var status string
	if in.Audited == "未审核" {
		status = "已审核"
	} else {
		status = "未审核"
	}

	if role == "superAdmin" || role == "admin" {
		update := db.Model(&models.Files{}).Where("file_id = ?", in.FileId).Update("status", status)
		if update.Error != nil {
			db.Rollback()
			return &file.AuditFileResp{
				Success: false,
				Msg:     "审核失败！",
			}, nil
		}
	} else {
		return &file.AuditFileResp{
			Success: false,
			Msg:     "权限不足！",
		}, nil
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
