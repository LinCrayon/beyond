// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package user

import (
	"context"

	"github.com/LinCrayon/beyond/application/user/rpc/service"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	FindByIdRequest      = service.FindByIdRequest
	FindByIdResponse     = service.FindByIdResponse
	FindByMobileRequest  = service.FindByMobileRequest
	FindByMobileResponse = service.FindByMobileResponse
	RegisterRequest      = service.RegisterRequest
	RegisterResponse     = service.RegisterResponse
	SendSmsRequest       = service.SendSmsRequest
	SendSmsResponse      = service.SendSmsResponse

	User interface {
		Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
		FindById(ctx context.Context, in *FindByIdRequest, opts ...grpc.CallOption) (*FindByIdResponse, error)
		FindByMobile(ctx context.Context, in *FindByMobileRequest, opts ...grpc.CallOption) (*FindByMobileResponse, error)
		SendSms(ctx context.Context, in *SendSmsRequest, opts ...grpc.CallOption) (*SendSmsResponse, error)
	}

	defaultUser struct {
		cli zrpc.Client
	}
)

func NewUser(cli zrpc.Client) User {
	return &defaultUser{
		cli: cli,
	}
}

func (m *defaultUser) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	client := service.NewUserClient(m.cli.Conn())
	return client.Register(ctx, in, opts...)
}

func (m *defaultUser) FindById(ctx context.Context, in *FindByIdRequest, opts ...grpc.CallOption) (*FindByIdResponse, error) {
	client := service.NewUserClient(m.cli.Conn())
	return client.FindById(ctx, in, opts...)
}

func (m *defaultUser) FindByMobile(ctx context.Context, in *FindByMobileRequest, opts ...grpc.CallOption) (*FindByMobileResponse, error) {
	client := service.NewUserClient(m.cli.Conn())
	return client.FindByMobile(ctx, in, opts...)
}

func (m *defaultUser) SendSms(ctx context.Context, in *SendSmsRequest, opts ...grpc.CallOption) (*SendSmsResponse, error) {
	client := service.NewUserClient(m.cli.Conn())
	return client.SendSms(ctx, in, opts...)
}
