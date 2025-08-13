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

	incomingContext, ok := metadata.FromIncomingContext(l.ctx)
	if !ok {
		return &file.DeleteFileResp{Success: false, Msg: "解析metadata数据失败！"}, nil
	}
	userId := incomingContext.Get("userId")

	// 先查询数据库是否有用户id和文件id匹配的记录
	var r models.Files
	t := l.svcCtx.MDB.Where("user_id = ? and file_id = ?", userId, in.FileId).First(&r)
	if errors.Is(t.Error, gorm.ErrRecordNotFound) {
		var a models.Files
		b := l.svcCtx.MDB.Where("file_id = ?", in.FileId).First(&a)
		// 当没查到匹配的记录则单独查询fileId符合的记录
		if (b.RowsAffected == 0) || (errors.Is(b.Error, gorm.ErrRecordNotFound)) {
			return &file.DeleteFileResp{Success: false, Msg: "此文件不存在，fileId：" + in.FileId}, nil
		}
		return &file.DeleteFileResp{Success: false, Msg: "此文件不属于当前登陆用户："}, nil
	} else if t.Error != nil {
		return &file.DeleteFileResp{Success: false, Msg: "查询匹配记录出错：" + t.Error.Error()}, nil
	}

	// 对用户id和文件id匹配的记录进行逻辑删除
	tx := l.svcCtx.MDB.Where("user_id = ? and file_id = ?", userId, in.FileId).Delete(&models.Files{})
	if tx.Error != nil {
		return &file.DeleteFileResp{Success: false, Msg: "操作数据库删除失败：" + tx.Error.Error()}, nil
	}
	return &file.DeleteFileResp{Success: false, Msg: "删除文件成功！"}, nil
}
