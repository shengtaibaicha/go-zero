package taglogic

import (
	"context"
	"encoding/json"
	"go-zero/models"

	"go-zero/apps/rpc/tag/internal/svc"
	"go-zero/apps/rpc/tag/tag"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTagsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTagsLogic {
	return &GetTagsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTagsLogic) GetTags(in *tag.GetTagsReq) (*tag.GetTagsResp, error) {

	MDB := l.svcCtx.MDB
	var tags []models.Tags
	MDB.Find(&tags)
	marshal, _ := json.Marshal(tags)

	return &tag.GetTagsResp{
		Tags: string(marshal),
	}, nil
}
