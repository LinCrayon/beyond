package svc

import (
	"github.com/LinCrayon/beyond/application/follow/rpc/internal/config"
	"github.com/LinCrayon/beyond/application/follow/rpc/internal/model"
	"github.com/LinCrayon/beyond/pkg/orm"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config           config.Config
	DB               *orm.DB
	FollowModel      *model.FollowModel
	FollowCountModel *model.FollowCountModel
	BizRedis         *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := orm.MustNewMysql(&orm.Config{
		DSN:          c.DB.DataSource,
		MaxOpenConns: c.DB.MaxOpenConns,
		MaxIdleConns: c.DB.MaxIdleConns,
		MaxLifetime:  c.DB.MaxLifetime,
	})
	//用Must 强依赖（错误直接panic()）
	rds := redis.MustNewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	})

	return &ServiceContext{
		Config:           c,
		DB:               db,
		FollowModel:      model.NewFollowModel(db.DB),
		FollowCountModel: model.NewFollowCountModel(db.DB),
		BizRedis:         rds,
	}
}
