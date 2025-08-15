package filelogic

import (
	"bytes"
	"context"
	"errors"
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
	"gorm.io/gorm"
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

	fileList := strings.Split(in.Filename, ".")

	MDB := l.svcCtx.MDB

	// 判断图片是否已经被当前用户上传过
	var f models.Files
	tx := MDB.Model(&models.Files{}).Where("user_id = ? and file_title = ?", in.UserId, fileList[0]).First(&f)
	if tx.RowsAffected != 0 || !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return &file.UploadResponse{
			Success: false,
			Msg:     "您已经上传过该图片了！",
		}, nil
	}

	// 存储原文件到minio
	filename := strconv.FormatInt(time.Now().UnixMilli(), 10) + "_" + fileList[0]
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
			UploadTime: time.Now(),
			Status:     "未审核",
			UserId:     in.UserId,
			FileName:   filename,
			FileTitle:  fileList[0],
			FileUrlse:  fileSE.FileUrl,
			Number:     0,
			Deleted:    0,
			FileSize:   in.Size,
		}

		// 开启gorm事务（修复：检查事务初始化错误）
		db := MDB.Begin()
		if db.Error != nil {
			return &file.UploadResponse{
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

		// 创建Files记录
		err = db.Create(mfile).Error
		if err != nil {
			if rollbackErr := db.Rollback().Error; rollbackErr != nil {
				l.Logger.Error("创建Files回滚失败: ", rollbackErr)
			}
			return &file.UploadResponse{
				Success: false,
				Msg:     "存储上传信息失败（Files）: " + err.Error(),
			}, nil
		}

		// 创建TagAndFile记录
		err = db.Create(&models.TagAndFile{
			Id:     uuid.New().String(),
			TagId:  in.TagId,
			FileId: mfile.FileId,
		}).Error
		if err != nil {
			if rollbackErr := db.Rollback().Error; rollbackErr != nil {
				l.Logger.Error("创建TagAndFile回滚失败: ", rollbackErr)
			}

			return &file.UploadResponse{
				Success: false,
				Msg:     "存储上传信息失败（TagAndFile）: " + err.Error(),
			}, nil
		}

		// 更新用户上传数量（修复：检查更新错误）
		result := db.Model(&models.Users{}).Where("user_id = ?", in.UserId).Update("upload_number", gorm.Expr("upload_number + ?", 1))
		if result.Error != nil {
			if rollbackErr := db.Rollback().Error; rollbackErr != nil {
				l.Logger.Error("更新用户数量回滚失败: ", rollbackErr)
			}
			return &file.UploadResponse{
				Success: false,
				Msg:     "更新用户上传数量失败: " + result.Error.Error(),
			}, nil
		}

		// 提交事务（检查提交是否成功）
		if err := db.Commit().Error; err != nil {
			return &file.UploadResponse{
				Success: false,
				Msg:     "事务提交失败: " + err.Error(),
			}, nil
		}

		return response, err
	}

	mfile = &models.Files{
		FileId:     uuid.New().String(),
		FileUrl:    response.FileUrl,
		UploadTime: time.Now(),
		Status:     "未审核",
		UserId:     in.UserId,
		FileName:   filename,
		FileTitle:  fileList[0],
		FileUrlse:  response.FileUrl,
		Number:     0,
		Deleted:    0,
		FileSize:   in.Size,
	}

	// 开启gorm事务（修复：检查事务初始化错误）
	db := MDB.Begin()
	if db.Error != nil {
		return &file.UploadResponse{
			Success: false,
			Msg:     "事务启动失败: " + db.Error.Error(),
		}, nil
	}
	defer func() {
		if r := recover(); r != nil {
			if err := db.Rollback().Error; err != nil {
				l.Logger.Error("恐慌时回滚失败: ", err)
			}
		}
	}()

	// 创建Files记录
	err = db.Create(mfile).Error
	if err != nil {
		if rollbackErr := db.Rollback().Error; rollbackErr != nil {
			l.Logger.Error("创建Files回滚失败: ", rollbackErr)
		}
		return &file.UploadResponse{
			Success: false,
			Msg:     "存储上传信息失败（Files）: " + err.Error(),
		}, nil
	}

	// 创建TagAndFile记录
	err = db.Create(&models.TagAndFile{
		Id:     uuid.New().String(),
		TagId:  in.TagId,
		FileId: mfile.FileId,
	}).Error
	if err != nil {
		if rollbackErr := db.Rollback().Error; rollbackErr != nil {
			l.Logger.Error("创建TagAndFile回滚失败: ", rollbackErr)
		}
		return &file.UploadResponse{
			Success: false,
			Msg:     "存储上传信息失败（TagAndFile）: " + err.Error(),
		}, nil
	}

	// 更新用户上传数量（修复：检查更新错误）
	result := db.Model(&models.Users{}).Where("user_id = ?", in.UserId).Update("upload_number", gorm.Expr("upload_number + ?", 1))
	if result.Error != nil {
		if rollbackErr := db.Rollback().Error; rollbackErr != nil {
			l.Logger.Error("更新用户数量回滚失败: ", rollbackErr)
		}
		return &file.UploadResponse{
			Success: false,
			Msg:     "更新用户上传数量失败: " + result.Error.Error(),
		}, nil
	}

	// 提交事务（检查提交是否成功）
	if err := db.Commit().Error; err != nil {
		return &file.UploadResponse{
			Success: false,
			Msg:     "事务提交失败: " + err.Error(),
		}, nil
	}

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
			Msg:     "上传图片失败！",
		}, nil
	}

	url := l.generateFileUrl(Filename)

	return &file.UploadResponse{
		Success: true,
		FileUrl: url,
		FileId:  "0",
		Msg:     "文件上传成功！",
	}, nil
}
