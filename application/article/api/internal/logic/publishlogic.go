package logic

import (
	"context"
	"encoding/json"
	"github.com/LinCrayon/beyond/application/article/api/internal/code"
	"github.com/LinCrayon/beyond/application/article/rpc/pb"
	"github.com/LinCrayon/beyond/pkg/xcode"

	"github.com/LinCrayon/beyond/application/article/api/internal/svc"
	"github.com/LinCrayon/beyond/application/article/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const minContentLen = 80

type PublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishLogic) Publish(req *types.PublishRequest) (resp *types.PublishResponse, err error) {
	if len(req.Title) == 0 {
		return nil, code.ArticleTitleEmpty
	}
	if len(req.Content) < minContentLen {
		return nil, code.ArticleContentTooFewWords
	}
	if len(req.Cover) == 0 {
		return nil, code.ArticleCoverEmpty
	}
	//判断登录状态 ,从 context 中取 userId
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("l.ctx.Value error: %v", err)
		return nil, xcode.NoLogin
	}
	//调用RPC服务的Publish,并传出参数
	pret, err := l.svcCtx.ArticleRPC.Publish(l.ctx, &pb.PublishRequest{
		UserId:      userId,
		Title:       req.Title,
		Content:     req.Content,
		Description: req.Description,
		Cover:       req.Cover,
	})
	if err != nil {
		logx.Errorf("l.svcCtx.ArticleRPC.Publish req: %v userId: %d error: %v", req, userId, err)
		return nil, err
	}

	return &types.PublishResponse{ArticleId: pret.ArticleId}, nil
}
