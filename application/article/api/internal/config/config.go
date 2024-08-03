package config

import (
	"github.com/LinCrayon/beyond/pkg/consul"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret  string
		AccessExpire  int64
		RefreshSecret string
		RefreshExpire int64
		RefreshAfter  int64
	}
	Oss struct {
		Endpoint         string
		AccessKeyId      string
		AccessKeySecret  string
		BucketName       string
		ConnectTimeout   int64 `json:",optional"`
		ReadWriteTimeout int64 `json:",optional"`
	}
	ArticleRPC zrpc.RpcClientConf
	UserRPC    zrpc.RpcClientConf
	Consul     consul.Conf
}
