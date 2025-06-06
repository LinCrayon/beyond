package logic

import (
	"context"
	"github.com/LinCrayon/beyond/application/like/rpc/pb"

	"github.com/LinCrayon/beyond/application/like/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type IsThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsThumbupLogic {
	return &IsThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsThumbupLogic) IsThumbup(in *pb.IsThumbupRequest) (*pb.IsThumbupResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.IsThumbupResponse{}, nil
}
