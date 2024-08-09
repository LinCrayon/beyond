package logic

import (
	"context"
	"fmt"
	"github.com/LinCrayon/beyond/application/like/mq/internal/svc"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

type ThumbupLogic struct {
	ctx         context.Context     //用于传递请求范围的数据、取消信号、截止时间等
	svcCtx      *svc.ServiceContext //服务的上下文，通常包含配置和依赖注入
	logx.Logger                     //嵌入的日志记录器
}

// NewThumbupLogic /*通过构造函数 NewThumbupLogic，可以确保 ThumbupLogic 结构体在创建时具备所有必要的依赖项和上下文信息*/
func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Consume ThumbupLogic 结构体的方法  （消息队列中获取的键和值）
func (l *ThumbupLogic) Consume(key, val string) error {
	fmt.Printf("get key: %s val: %s\n", key, val)
	return nil
}

// Consumers 用于消费消息队列中的消息
func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	//创建了一个新的消息队列消费者并将其添加到服务切片中
	return []service.Service{ //调用 NewThumbupLogic 函数，创建一个新的 ThumbupLogic 实例
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, NewThumbupLogic(ctx, svcCtx)),
	}
}
