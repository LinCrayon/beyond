package main

import (
	"flag"
	"fmt"
	"github.com/LinCrayon/beyond/pkg/consul"
	"github.com/LinCrayon/beyond/pkg/interceptors"

	"github.com/LinCrayon/beyond/application/user/rpc/internal/config"
	"github.com/LinCrayon/beyond/application/user/rpc/internal/server"
	"github.com/LinCrayon/beyond/application/user/rpc/internal/svc"
	"github.com/LinCrayon/beyond/application/user/rpc/service"

	"github.com/zeromicro/go-zero/core/conf"
	cs "github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "application/user/rpc/etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		service.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == cs.DevMode || c.Mode == cs.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	//自定义拦截器  grpc serve 端的拦截器(把用户自定义的错误码写入到GRPC的错误中)
	s.AddUnaryInterceptors(interceptors.ServerErrorInterceptor())

	// 服务注册
	err := consul.Register(c.Consul, fmt.Sprintf("%s:%d", c.ServiceConf.Prometheus.Host, c.ServiceConf.Prometheus.Port))
	if err != nil {
		fmt.Printf("register consul error: %v\n", err)
	}

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
