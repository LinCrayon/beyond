package svc

import (
	"github.com/LinCrayon/beyond/application/user/rpc/internal/config"
	"github.com/LinCrayon/beyond/application/user/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource) //创建数据库连接
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(conn, c.CacheRedis),
	}
}
