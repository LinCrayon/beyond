package main

import (
	"flag"
	"fmt"
	"github.com/LinCrayon/beyond/pkg/consul"
	"github.com/LinCrayon/beyond/pkg/xcode"
	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/LinCrayon/beyond/application/article/api/internal/config"
	"github.com/LinCrayon/beyond/application/article/api/internal/handler"
	"github.com/LinCrayon/beyond/application/article/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "application/article/api/etc/article-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	//加载配置文件解析到 c 结构体
	conf.MustLoad(*configFile, &c)

	//创建 REST 服务器导入配置
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	//创建服务上下文
	ctx := svc.NewServiceContext(c)
	//注册处理程序 ,服务上下文 ctx 被创建并传递给处理程序以注册处理程序
	handler.RegisterHandlers(server, ctx) //servicecontext传到handle中

	// 自定义错误处理方法
	httpx.SetErrorHandler(xcode.ErrHandler)

	// 服务注册
	err := consul.Register(c.Consul, fmt.Sprintf("%s:%d", c.ServiceConf.Prometheus.Host, c.ServiceConf.Prometheus.Port))
	if err != nil {
		fmt.Printf("register consul error: %v\n", err)
	}

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
