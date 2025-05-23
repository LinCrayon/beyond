// Code generated by goctl. DO NOT EDIT.
// Source: follow.proto

package server

import (
	"context"

	"github.com/LinCrayon/beyond/application/follow/rpc/internal/logic"
	"github.com/LinCrayon/beyond/application/follow/rpc/internal/svc"
	"github.com/LinCrayon/beyond/application/follow/rpc/pb"
)

type FollowServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedFollowServer
}

func NewFollowServer(svcCtx *svc.ServiceContext) *FollowServer {
	return &FollowServer{
		svcCtx: svcCtx,
	}
}

// 关注
func (s *FollowServer) Follow(ctx context.Context, in *pb.FollowRequest) (*pb.FollowResponse, error) {
	l := logic.NewFollowLogic(ctx, s.svcCtx)
	return l.Follow(in)
}

// 取消关注
func (s *FollowServer) UnFollow(ctx context.Context, in *pb.UnFollowRequest) (*pb.UnFollowResponse, error) {
	l := logic.NewUnFollowLogic(ctx, s.svcCtx)
	return l.UnFollow(in)
}

// 关注列表
func (s *FollowServer) FollowList(ctx context.Context, in *pb.FollowListRequest) (*pb.FollowListResponse, error) {
	l := logic.NewFollowListLogic(ctx, s.svcCtx)
	return l.FollowList(in)
}

// 粉丝列表
func (s *FollowServer) FansList(ctx context.Context, in *pb.FansListRequest) (*pb.FansListResponse, error) {
	l := logic.NewFansListLogic(ctx, s.svcCtx)
	return l.FansList(in)
}
