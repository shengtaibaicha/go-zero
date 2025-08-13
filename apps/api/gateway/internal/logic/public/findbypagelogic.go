package public

import (
	"context"
	"encoding/json"
	"go-zero/apps/rpc/file/file"
	"go-zero/common/result"
	"go-zero/models"

	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/apps/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindByPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByPageLogic {
	return &FindByPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindByPageLogic) FindByPage(req *types.FindByPageReq, auth string) (resp *result.Result, err error) {
	page, err := l.svcCtx.FileClient.FindByPage(l.ctx, &file.FindByPageReq{
		Page: req.Page,
		Size: req.Size,
	})

	if err != nil {
		return nil, err
	}

	c := page.Collect

	type res struct {
		models.Files
		Collect bool `json:"collect"`
	}

	// 将json字符串序列化为对象
	var data []res
	json.Unmarshal([]byte(page.Records), &data)

	if auth != "" {
		// 使用索引遍历，直接访问原切片中的元素
		for i := range data {
			// 通过 data[i] 访问原元素（不是副本）
			e := &data[i] // 获取元素的指针，避免再次拷贝
			v, ex := c[e.FileId]
			if ex {
				e.Collect = v // 直接修改原元素的字段
			}
		}
	}

	r := map[string]any{}
	r["records"] = data
	r["total"] = page.Total
	r["current"] = page.Current
	r["pages"] = page.Pages
	r["size"] = page.Size

	return result.Ok().SetData(r), nil
}
