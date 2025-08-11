package filelogic

import (
	"context"
	"go-zero/models"
	"io"
	"path/filepath"

	"go-zero/apps/rpc/file/file"
	"go-zero/apps/rpc/file/internal/svc"

	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/logx"
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
	l.svcCtx.MDB.Where("file_name = ?", in.FileName).Find(&files)
	l.svcCtx.MDB.Model(&models.Files{}).Where("file_name = ?", in.FileName).Update("number", files.Number+1)

	return &file.DownloadFileResp{
		Content:     content,
		FileName:    filepath.Base(in.FileName),
		ContentType: stat.ContentType,
	}, nil
}
