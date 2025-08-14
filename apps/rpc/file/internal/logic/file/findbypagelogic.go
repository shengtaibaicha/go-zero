package filelogic

import (
	"context"
	"encoding/json"
	"fmt"
	"go-zero/models"

	"go-zero/apps/rpc/file/file"
	"go-zero/apps/rpc/file/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByPageLogic {
	return &FindByPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByPageLogic) FindByPage(in *file.FindByPageReq) (*file.FindByPageResp, error) {

	// 从redis中获取数据
	key := fmt.Sprint(l.svcCtx.Config.Redis.Key, "FindByPage:page:", in.Page, "size:", in.Size)
	value, err := l.svcCtx.RedisClient.GetCtx(context.Background(), key)
	if err != nil {
		logx.Error(err)
	}
	if value == "" {
		MDB := l.svcCtx.MDB
		var total int64

		// 计算需要跳过的记录
		offset := (in.Page - 1) * in.Size
		var fileData []models.Files
		MDB.Where("status = ?", "已审核").Offset(int(offset)).Limit(int(in.Size)).Find(&fileData)

		// 获取记录总条数
		MDB.Where("status = ?", "已审核").Model(&models.Files{}).Count(&total)

		// 获取当前用户的收藏信息
		var num []string
		MDB.Model(&models.Collect{}).Select("file_id").Scan(&num)
		r := make(map[string]bool)
		for _, v := range num {
			r[v] = true
		}

		// 将查询到的数据解析为json格式
		marshaled, _ := json.Marshal(fileData)

		pages := int32(total) / in.Size
		if int32(total)%in.Size != 0 {
			pages = (int32(total) / in.Size) + 1
		}

		result := &file.FindByPageResp{
			Records: string(marshaled),
			Size:    in.Size,
			Total:   int32(total),
			Current: in.Page,
			Pages:   pages,
			Collect: r,
		}

		// 将数据序列化后存入redis
		redisData, _ := json.Marshal(result)
		l.svcCtx.RedisClient.SetexCtx(context.Background(), key, string(redisData), l.svcCtx.Config.RedisExpires)

		return result, nil
	}
	var result file.FindByPageResp
	json.Unmarshal([]byte(value), &result)
	return &result, nil
}
