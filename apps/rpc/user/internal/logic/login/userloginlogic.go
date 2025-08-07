package loginlogic

import (
	"context"
	"errors"
	"go-zero/apps/rpc/user/internal/svc"
	"go-zero/apps/rpc/user/tools"
	"go-zero/apps/rpc/user/user"
	"go-zero/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserLoginLogic) UserLogin(in *user.LoginReq) (*user.LoginResp, error) {
	// todo: add your logic here and delete this line
	MDB := l.svcCtx.MDB
	var users models.Users
	tx := MDB.Where("user_name = ?", in.UserName).First(&users)
	// 数据库没有匹配的数据，代表登陆失败，返回token为0
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("用户名不存在！")
	}
	// 验证数据的密码和传入的密码是否相等
	err2 := bcrypt.CompareHashAndPassword([]byte(users.UserPassword), []byte(in.UserPassword))
	if err2 != nil {
		return nil, errors.New("用户名或密码不正确！")
	}
	// 查找到对应的数据，我们生成token返回
	token, err := tools.GenerateToken(l.svcCtx.Jwt.SecretKey, users.UserId, users.UserName, "0", 24)
	if err != nil {
		return nil, errors.New("create token failed")
	}
	return &user.LoginResp{
		Token: token,
	}, nil
}
