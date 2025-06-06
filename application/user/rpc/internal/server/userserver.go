// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package server

import (
	"context"

	"github.com/LinCrayon/beyond/application/user/rpc/internal/logic"
	"github.com/LinCrayon/beyond/application/user/rpc/internal/svc"
	"github.com/LinCrayon/beyond/application/user/rpc/service"
)

type UserServer struct {
	svcCtx *svc.ServiceContext
	service.UnimplementedUserServer
}

func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
	return &UserServer{
		svcCtx: svcCtx,
	}
}

func (s *UserServer) Register(ctx context.Context, in *service.RegisterRequest) (*service.RegisterResponse, error) {
	l := logic.NewRegisterLogic(ctx, s.svcCtx)
	return l.Register(in)
}

func (s *UserServer) FindById(ctx context.Context, in *service.FindByIdRequest) (*service.FindByIdResponse, error) {
	l := logic.NewFindByIdLogic(ctx, s.svcCtx)
	return l.FindById(in)
}

func (s *UserServer) FindByMobile(ctx context.Context, in *service.FindByMobileRequest) (*service.FindByMobileResponse, error) {
	l := logic.NewFindByMobileLogic(ctx, s.svcCtx)
	return l.FindByMobile(in)
}

func (s *UserServer) SendSms(ctx context.Context, in *service.SendSmsRequest) (*service.SendSmsResponse, error) {
	l := logic.NewSendSmsLogic(ctx, s.svcCtx)
	return l.SendSms(in)
}
