package filelogic

import (
	"context"
	"encoding/json"
	"fmt"
	"go-zero/apps/rpc/file/file"
	"go-zero/apps/rpc/file/internal/svc"
	"go-zero/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindPageByNameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindPageByNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindPageByNameLogic {
	return &FindPageByNameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindPageByNameLogic) FindPageByName(in *file.FindPageByNameReq) (*file.FindPageByNameResp, error) {

	// 从redis中获取数据
	key := fmt.Sprint(l.svcCtx.Config.Redis.Key, "FindPageByName:page:", in.Page, "size:", in.Size, "name:", in.Name)
	value, err := l.svcCtx.RedisClient.GetCtx(context.Background(), key)
	if err != nil {
		logx.Error(err)
	}

	if value == "" {
		MDB := l.svcCtx.MDB

		var total int32
		// 计算需要跳过的记录
		offset := (in.Page - 1) * in.Size
		var files []models.Files
		sql := "SELECT * FROM files WHERE file_title like ? AND status = '已审核' AND deleted != 1 LIMIT ? OFFSET ?"
		MDB.Raw(sql, "%"+in.Name+"%", in.Size, offset).Scan(&files)
		// 获取记录总条数
		MDB.Raw("SELECT count(*) FROM files WHERE file_title like ?  AND status = '已审核' AND deleted != 1", "%"+in.Name+"%").Scan(&total)
		// 将查询到的数据解析为json格式
		marshaled, _ := json.Marshal(files)

		pages := total / in.Size
		if total%in.Size != 0 {
			pages = (total / in.Size) + 1
		}

		// 获取当前用户的收藏信息
		var num []string
		MDB.Model(&models.Collect{}).Select("file_id").Scan(&num)
		r := make(map[string]bool)
		for _, v := range num {
			r[v] = true
		}

		result := &file.FindPageByNameResp{
			Records: string(marshaled),
			Size:    in.Size,
			Total:   total,
			Current: in.Page,
			Pages:   pages,
			Collect: r,
		}

		// 将数据序列化后存入redis
		redisData, _ := json.Marshal(result)
		l.svcCtx.RedisClient.SetexCtx(context.Background(), key, string(redisData), l.svcCtx.Config.RedisExpires*60)

		return result, nil
	}
	var result file.FindPageByNameResp
	json.Unmarshal([]byte(value), &result)
	return &result, nil
}
