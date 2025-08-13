package file

import (
	"context"
	"encoding/json"
	"go-zero/apps/rpc/file/file"
	"go-zero/common/result"
	"go-zero/models"

	"go-zero/apps/api/gateway/internal/svc"
	"go-zero/apps/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type FileUserPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUserPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUserPageLogic {
	return &FileUserPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUserPageLogic) FileUserPage(req *types.FileUserPageReq) (resp *result.Result, err error) {

	md := metadata.New(map[string]string{"userId": l.ctx.Value("userId").(string)})
	outgoingContext := metadata.NewOutgoingContext(l.ctx, md)

	page, _ := l.svcCtx.FileClient.FileUserPage(outgoingContext, &file.FileUserPageReq{
		Page: req.Page,
		Size: req.Size,
	})

	if page == nil {
		return result.Err().SetMsg("查询用户上传的图片失败！"), nil
	}
	var data []models.Files
	json.Unmarshal([]byte(page.Records), &data)

	r := map[string]any{}
	r["records"] = data
	r["total"] = page.Total
	r["current"] = page.Current
	r["pages"] = page.Pages
	r["size"] = page.Size

	return result.Ok().SetMsg("查询用户上传的图片成功！").SetData(r), nil
}
