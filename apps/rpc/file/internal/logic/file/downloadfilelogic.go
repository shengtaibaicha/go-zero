package filelogic

import (
	"context"
	"errors"
	"go-zero/models"
	"io"
	"path/filepath"

	"go-zero/apps/rpc/file/file"
	"go-zero/apps/rpc/file/internal/svc"

	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type DownloadFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDownloadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadFileLogic {
	return &DownloadFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DownloadFileLogic) DownloadFile(in *file.DownloadFileReq) (*file.DownloadFileResp, error) {

	incomingContext, Ok := metadata.FromIncomingContext(l.ctx)
	if !Ok {
		logx.Error("在下载文件任务中从metadata获取userId失败，请排查原因！")
		return nil, errors.New("在下载文件任务中从metadata获取userId失败，请排查原因！")
	}
	userId := incomingContext.Get("userId")[0]

	db := l.svcCtx.MDB.Begin()
	if db.Error != nil {
		return nil, db.Error
	}
	defer func() {
		if r := recover(); r != nil {
			// 恐慌时回滚，并记录日志
			if err := db.Rollback().Error; err != nil {
				l.Logger.Error("恐慌时回滚失败: ", err)
			}
		}
	}()

	// 判断用户状态是否被禁用enable = 0
	var enable models.Users
	isEnable := db.Model(&models.Users{}).Where("enable = ? and user_id = ?", 0, userId).First(&enable)
	if isEnable.RowsAffected != 0 {
		return nil, errors.New("当前用户已被禁用，请联系管理员！")
	}

	// 从 MinIO 获取文件对象
	object, err := l.svcCtx.MinioClient.GetObject(
		l.ctx,
		l.svcCtx.Config.Minio.Bucket,
		in.FileName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		l.Errorf("Get object failed: %v", err)
		return nil, err
	}
	defer object.Close()

	// 获取文件信息
	stat, err := object.Stat()
	if err != nil {
		l.Errorf("Get object stat failed: %v", err)
		return nil, err
	}

	// 读取文件内容到字节数组（适合小文件，大文件会占用大量内存）
	content, err := io.ReadAll(object)
	if err != nil {
		l.Errorf("Read file content failed: %v", err)
		return nil, err
	}
	// 将下载数量增加
	var files models.Files
	db.Where("file_name = ?", in.FileName).Find(&files)
	tx := db.Model(&models.Files{}).Where("file_name = ?", in.FileName).Update("number", files.Number+1)
	if tx.Error != nil {
		db.Rollback()
	}

	// 更新用户下载数量
	update := db.Model(&models.Users{}).Where("user_id = ?", userId).Update("download_number", gorm.Expr("download_number + ?", 1))
	if update.Error != nil {
		db.Rollback()
	}

	db.Commit()

	return &file.DownloadFileResp{
		Content:     content,
		FileName:    filepath.Base(in.FileName),
		ContentType: stat.ContentType,
	}, nil
}
