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

type DeleteFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteFileLogic) DeleteFile(in *file.DeleteFileReq) (*file.DeleteFileResp, error) {

	db := l.svcCtx.MDB.Begin()
	if db.Error != nil {
		return &file.DeleteFileResp{
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
		return &file.DeleteFileResp{Success: false, Msg: "解析metadata数据失败！"}, nil
	}
	userId := incomingContext.Get("userId")[0]
	role := incomingContext.Get("role")[0]

	// 判断操作的用户是不是superAdmin或者admin
	if role != "superAdmin" && role != "admin" {

		// 先查询数据库是否有用户id和文件id匹配的记录
		var r models.Files
		t := db.Where("user_id = ? and file_id = ?", userId, in.FileId).First(&r)
		if errors.Is(t.Error, gorm.ErrRecordNotFound) {
			var a models.Files
			b := db.Where("file_id = ?", in.FileId).First(&a)
			// 当没查到匹配的记录则单独查询fileId符合的记录
			if (b.RowsAffected == 0) || (errors.Is(b.Error, gorm.ErrRecordNotFound)) {
				return &file.DeleteFileResp{Success: false, Msg: "此文件不存在，fileId：" + in.FileId}, nil
			}
			return &file.DeleteFileResp{Success: false, Msg: "此文件不属于当前登陆用户："}, nil
		} else if t.Error != nil {
			return &file.DeleteFileResp{Success: false, Msg: "查询匹配记录出错：" + t.Error.Error()}, nil
		}

		// 对用户id和文件id匹配的记录进行逻辑删除
		tx := db.Where("user_id = ? and file_id = ?", userId, in.FileId).Delete(&models.Files{})
		if tx.Error != nil {
			db.Rollback()
			return &file.DeleteFileResp{Success: false, Msg: "操作数据库删除失败：" + tx.Error.Error()}, nil
		}

		// 对用户的上传数量减去
		update := db.Model(&models.Users{}).Where("user_id = ?", userId).Update("upload_number", gorm.Expr("upload_number - ?", 1))
		if update.Error != nil {
			db.Rollback()
			return &file.DeleteFileResp{Success: false, Msg: "操作数据库删除失败：" + tx.Error.Error()}, nil
		}
	}

	// 当身份未管理员时，直接对文件id匹配的记录进行逻辑删除
	var r models.Files
	tx := db.Where("file_id = ?", in.FileId).First(&r)
	if tx.Error != nil {
		db.Rollback()
		return &file.DeleteFileResp{Success: false, Msg: "操作数据库删除失败：" + tx.Error.Error()}, nil
	}
	t := db.Where("file_id = ?", in.FileId).Delete(&models.Files{})
	if t.Error != nil {
		db.Rollback()
		return &file.DeleteFileResp{Success: false, Msg: "操作数据库删除失败：" + tx.Error.Error()}, nil
	}
	update := db.Model(&models.Users{}).Where("user_id = ?", r.UserId).Update("upload_number", gorm.Expr("upload_number - ?", 1))
	if update.Error != nil {
		db.Rollback()
		return &file.DeleteFileResp{Success: false, Msg: "操作数据库删除失败：" + tx.Error.Error()}, nil
	}

	db.Commit()

	return &file.DeleteFileResp{Success: true, Msg: "删除文件成功！"}, nil
}
