package registerlogic

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go-zero/apps/rpc/user/internal/svc"
	"go-zero/apps/rpc/user/user"
	"go-zero/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserRegisterLogic) UserRegister(in *user.RegisterReq) (*user.RegisterResp, error) {
	// todo: add your logic here and delete this line
	MDB := l.svcCtx.MDB
	var re models.Users
	tx := MDB.Where("user_name = ?", in.UserName).First(&re)
	// 如果数据库查到了这个用户名的数据
	if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("该用户名已注册！")
	}
	// 将密码加密后创建user实例并且写入数据库
	password, _ := bcrypt.GenerateFromPassword([]byte(in.UserPassword), 12)
	userId := uuid.New().String()
	create := MDB.Create(&models.Users{
		UserId:       userId,
		UserName:     in.UserName,
		UserEmail:    in.UserEmail,
		UserAvatar:   "",
		UserPassword: string(password),
	})
	if create.Error != nil {
		return nil, errors.New("用户注册失败！")
	}
	return &user.RegisterResp{
		UserId: userId,
	}, nil
}
