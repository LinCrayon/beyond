// Code generated by goctl. DO NOT EDIT.
// Source: like.proto

package like

import (
	"context"

	"github.com/LinCrayon/beyond/application/like/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	IsThumbupRequest  = pb.IsThumbupRequest
	IsThumbupResponse = pb.IsThumbupResponse
	ThumbupRequest    = pb.ThumbupRequest
	ThumbupResponse   = pb.ThumbupResponse
	UserThumbup       = pb.UserThumbup

	Like interface {
		Thumbup(ctx context.Context, in *ThumbupRequest, opts ...grpc.CallOption) (*ThumbupResponse, error)
		IsThumbup(ctx context.Context, in *IsThumbupRequest, opts ...grpc.CallOption) (*IsThumbupResponse, error)
	}

	defaultLike struct {
		cli zrpc.Client
	}
)

func NewLike(cli zrpc.Client) Like {
	return &defaultLike{
		cli: cli,
	}
}

func (m *defaultLike) Thumbup(ctx context.Context, in *ThumbupRequest, opts ...grpc.CallOption) (*ThumbupResponse, error) {
	client := pb.NewLikeClient(m.cli.Conn())
	return client.Thumbup(ctx, in, opts...)
}

func (m *defaultLike) IsThumbup(ctx context.Context, in *IsThumbupRequest, opts ...grpc.CallOption) (*IsThumbupResponse, error) {
	client := pb.NewLikeClient(m.cli.Conn())
	return client.IsThumbup(ctx, in, opts...)
}
