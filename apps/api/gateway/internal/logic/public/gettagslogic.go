package public

import (
	"context"
	"encoding/json"
	"go-zero/apps/rpc/tag/tag"
	"go-zero/common/result"
	"go-zero/models"

	"go-zero/apps/api/gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTagsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTagsLogic {
	return &GetTagsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTagsLogic) GetTags() (resp *result.Result, err error) {

	tags, _ := l.svcCtx.TagClient.GetTags(l.ctx, &tag.GetTagsReq{})
	var data []models.Tags
	json.Unmarshal([]byte(tags.Tags), &data)

	return result.Ok().SetData(data), nil
}
