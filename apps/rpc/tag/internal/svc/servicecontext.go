package svc

import (
	"go-zero/apps/rpc/tag/internal/config"
	"go-zero/common/db"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	MDB    *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		MDB:    db.Init(c.Mysql.DataSource),
	}
}
