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

type FindPageByNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindPageByNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindPageByNameLogic {
	return &FindPageByNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindPageByNameLogic) FindPageByName(req *types.FindPageByNameReq, auth string) (resp *result.Result, err error) {

	nameData, _ := l.svcCtx.FileClient.FindPageByName(
		l.ctx,
		&file.FindPageByNameReq{
			Page: req.Page,
			Size: req.Size,
			Name: req.Name,
		},
	)

	c := nameData.Collect

	type res struct {
		models.Files
		Collect bool `json:"collect"`
	}

	var data []res

	jsonErr := json.Unmarshal([]byte(nameData.GetRecords()), &data)
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
	r["total"] = nameData.GetTotal()
	r["current"] = nameData.GetCurrent()
	r["pages"] = nameData.GetPages()
	r["size"] = nameData.GetSize()

	return result.Ok().SetData(r), nil
}
