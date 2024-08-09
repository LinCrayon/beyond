package interceptors

import (
	"context"
	"github.com/LinCrayon/beyond/pkg/xcode"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// ClientErrorInterceptor 自定义拦截器用于捕捉 gRPC 调用中的错误，将 gRPC 错误状态转换为自定义的错误代码（xcode）格式
func ClientErrorInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			grpcStatus, _ := status.FromError(err)
			xc := xcode.GrpcStatusToXCode(grpcStatus) //将 grpcStatus 转化成自定义的xcode
			err = errors.WithMessage(xc, grpcStatus.Message())
		}

		return err
	}
}
