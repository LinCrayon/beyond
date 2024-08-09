package interceptors

import (
	"context"
	"github.com/LinCrayon/beyond/pkg/xcode"

	"google.golang.org/grpc"
)

func ServerErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		return resp, xcode.FromError(err).Err() //handler执行的error转化为grpc能识别的error
	}
}
