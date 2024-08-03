package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/LinCrayon/beyond/application/article/rpc/internal/svc"
	"github.com/LinCrayon/beyond/application/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDetailLogic {
	return &ArticleDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticleDetailLogic) ArticleDetail(in *pb.ArticleDetailRequest) (*pb.ArticleDetailResponse, error) {
	article, err := l.svcCtx.ArticleModel.FindOne(l.ctx, in.ArticleId)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) { //errors.Is 检查错误是否是sqlx.ErrNotFound错误
			return &pb.ArticleDetailResponse{}, nil
		}
		return nil, err
	}
	return &pb.ArticleDetailResponse{
		Article: &pb.ArticleItem{
			Id:          article.Id,
			Title:       article.Title,
			Content:     article.Content,
			Description: article.Description,
			Cover:       article.Cover,
			AuthorId:    article.AuthorId,
			LikeCount:   article.LikeNum,
			PublishTime: article.PublishTime.Unix(),
		},
	}, nil
}
