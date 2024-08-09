package svc

import (
	"github.com/LinCrayon/beyond/application/article/api/internal/config"
	"github.com/LinCrayon/beyond/application/article/rpc/article"
	"github.com/LinCrayon/beyond/application/user/rpc/user"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/zeromicro/go-zero/zrpc"
)

const (
	defaultOssConnectTimeout   = 1
	defaultOssReadWriteTimeout = 3
)

type ServiceContext struct {
	Config     config.Config
	OssClient  *oss.Client
	ArticleRPC article.Article
	UserRPC    user.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	if c.Oss.ConnectTimeout == 0 {
		c.Oss.ConnectTimeout = defaultOssConnectTimeout
	}
	if c.Oss.ReadWriteTimeout == 0 {
		c.Oss.ReadWriteTimeout = defaultOssReadWriteTimeout
	}
	//创建 Oss 客户端
	//创建 Oss 客户端
	oc, err := oss.New(c.Oss.Endpoint, c.Oss.AccessKeyId, c.Oss.AccessKeySecret,
		oss.Timeout(c.Oss.ConnectTimeout, c.Oss.ReadWriteTimeout)) //超时设置，包括连接超时和读写超时
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:     c,
		OssClient:  oc,
		ArticleRPC: article.NewArticle(zrpc.MustNewClient(c.ArticleRPC)), //实例化 Article 服务，传入RPC客户端
		UserRPC:    user.NewUser(zrpc.MustNewClient(c.UserRPC)),
	}
}
