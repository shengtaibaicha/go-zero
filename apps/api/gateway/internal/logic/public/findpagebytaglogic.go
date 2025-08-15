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

type FindPageByTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindPageByTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindPageByTagLogic {
	return &FindPageByTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindPageByTagLogic) FindPageByTag(req *types.FindPageByTagReq, auth string) (resp *result.Result, err error) {

	tagData, _ := l.svcCtx.FileClient.FindPageByTag(l.ctx, &file.FindPageByTagReq{
		Page:  req.Page,
		Size:  req.Size,
		TagId: req.TagId,
	})

	c := tagData.Collect

	type res struct {
		models.Files
		Collect bool `json:"collect"`
	}

	var data []res

	jsonErr := json.Unmarshal([]byte(tagData.GetRecords()), &data)
	if jsonErr != nil {
		l.Logger.Error("json序列化失败", jsonErr.Error())
		return result.Err().SetMsg("json序列化失败！"), err
	}

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
	r["total"] = tagData.GetTotal()
	r["current"] = tagData.GetCurrent()
	r["pages"] = tagData.GetPages()
	r["size"] = tagData.GetSize()

	return result.Ok().SetData(r), nil
}
