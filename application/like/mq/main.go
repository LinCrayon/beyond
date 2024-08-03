package main

import (
	"context"
	"flag"
	"github.com/LinCrayon/beyond/application/like/mq/internal/config"
	"github.com/LinCrayon/beyond/application/like/mq/internal/logic"
	"github.com/LinCrayon/beyond/application/like/mq/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "application/like/mq/etc/like.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	svcCtx := svc.NewServiceContext(c)
	ctx := context.Background() //在没有更具体的上下文可用时作为根上下文
	/*服务组（serviceGroup）通常用于管理一组服务的生命周期，提供启动、停止和管理多个服务的功能*/
	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	//如何遍历从 Consumers 函数获取的消息队列消费者，并将它们添加到服务组中,确保所有消费者都被正确地启动和管理
	for _, mq := range logic.Consumers(ctx, svcCtx) {
		serviceGroup.Add(mq) //使用 serviceGroup 管理多个服务，确保它们可以统一启动和停止
	}

	serviceGroup.Start()
}
