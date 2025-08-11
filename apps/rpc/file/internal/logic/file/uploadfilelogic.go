package filelogic

import (
	"bytes"
	"context"
	"go-zero/apps/rpc/file/file"
	"go-zero/apps/rpc/file/internal/svc"
	"go-zero/apps/rpc/file/tools/image"
	"go-zero/models"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 单请求上传（适合中小文件）
func (l *UploadFileLogic) UploadFile(in *file.UploadReq) (*file.UploadResponse, error) {
	MDB := l.svcCtx.MDB

	fileList := strings.Split(in.Filename, ".")

	// metadata接收数据
	//incomingContext, ok := metadata.FromIncomingContext(l.ctx)
	//if !ok {
	//	log.Printf("metadata.FromIncomingContext() fail")
	//}
	//log.Printf("metadata取到的数据%s\n", incomingContext)

	// 存储原文件到minio
	filename := strconv.FormatInt(time.Now().UnixMilli(), 10) + "_" + in.Filename
	response, err := up(l, filename, in.File, in.Size, in.MimeType)
	if err != nil {
		return response, err
	}

	mfile := &models.Files{}

	// 保存压缩后的图片
	if fileList[1] != "png" {
		// 生成图片缩略图
		// 不知道为什么png的图片用这个库压缩后反而更大了，所以png图片就不压缩了
		compressImage, _ := image.CompressImage(in.File, fileList[1], 20)
		FileName := filename + "SE." + fileList[1]
		fileSE, _ := up(l, FileName, compressImage, int64(len(compressImage)), in.MimeType)

		// 将压缩后的记录插入数据库
		mfile = &models.Files{
			FileId:     uuid.New().String(),
			FileUrl:    response.FileUrl,
			UploadTime: time.Now().Format("2006-01-02 15:04:05.000"),
			Status:     "未审核",
			UserId:     in.UserId,
			FileName:   filename,
			FileTitle:  fileList[0],
			FileUrlse:  fileSE.FileUrl,
			Number:     0,
		}
	}

	mfile = &models.Files{
		FileId:     uuid.New().String(),
		FileUrl:    response.FileUrl,
		UploadTime: time.Now().Format("2006-01-02 15:04:05.000"),
		Status:     "未审核",
		UserId:     in.UserId,
		FileName:   filename,
		FileTitle:  fileList[0],
		FileUrlse:  response.FileUrl,
		Number:     0,
		Deleted:    0,
	}

	// 开启gorm事务
	db := MDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			db.Rollback() // 回滚事务
		}
	}()

	err = db.Create(mfile).Error
	if err != nil {
		// 失败回滚
		db.Rollback()
		return &file.UploadResponse{
			Success: false,
			Message: "存储上传信息失败！",
		}, nil
	}

	err = db.Create(&models.TagAndFile{
		Id:     uuid.New().String(),
		TagId:  in.TagId,
		FileId: mfile.FileId,
	}).Error
	if err != nil {
		// 失败回滚
		db.Rollback()
		return &file.UploadResponse{
			Success: false,
			Message: "存储上传信息失败！",
		}, nil
	}

	// 提交事务
	db.Commit()

	return response, err
}

// 生成文件访问URL
func (l *UploadFileLogic) generateFileUrl(objectName string) string {
	// 如果MinIO配置了公开访问，可以直接使用基础URL
	return "http://" + l.svcCtx.Config.Minio.Endpoint + "/" + l.svcCtx.Config.Minio.Bucket + "/" + objectName
}

func up(l *UploadFileLogic, Filename string, File []byte, Size int64, MimeType string) (*file.UploadResponse, error) {
	minioClient := l.svcCtx.MinioClient
	_, err := minioClient.PutObject(
		l.ctx,
		l.svcCtx.Config.Minio.Bucket, // 存储桶名称
		Filename,                     // 存储在MinIO中的对象名称
		bytes.NewReader(File),        // 文件数据流
		Size,                         // 文件大小
		minio.PutObjectOptions{
			ContentType: MimeType, // 设置文件MIME类型
		})
	if err != nil {
		return &file.UploadResponse{
			Success: false,
			Message: "上传图片失败！",
		}, nil
	}

	url := l.generateFileUrl(Filename)

	return &file.UploadResponse{
		Success: true,
		FileUrl: url,
		FileId:  "0",
		Message: "文件上传成功！",
	}, nil
}
