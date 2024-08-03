package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ArticleModel = (*customArticleModel)(nil)

type (
	// ArticleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customArticleModel.
	ArticleModel interface {
		articleModel
		ArticlesByUserId(ctx context.Context, userId int64, status int, likeNum int64, pubTime, sortField string, limit int) ([]*Article, error)
		UpdateArticleStatus(ctx context.Context, id int64, status int) error
	}

	customArticleModel struct {
		*defaultArticleModel
	}
)

// NewArticleModel returns a model for the database table.
func NewArticleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ArticleModel {
	return &customArticleModel{
		defaultArticleModel: newArticleModel(conn, c, opts...),
	}
}

func (m *customArticleModel) ArticlesByUserId(ctx context.Context, userId int64, status int, likeNum int64, pubTime, sortField string, limit int) ([]*Article, error) {
	var (
		err      error
		sql      string
		anyField any
		articles []*Article
	)
	if sortField == "like_num" {
		anyField = likeNum
		sql = fmt.Sprintf("select "+articleRows+" from "+m.table+" where author_id=? and status=? and like_num < ? order by %s desc limit ?", sortField)
	} else {
		anyField = pubTime
		sql = fmt.Sprintf("select "+articleRows+" from "+m.table+" where author_id=? and status=? and publish_time < ? order by %s desc limit ?", sortField)
	}
	err = m.QueryRowsNoCacheCtx(ctx, &articles, sql, userId, status, anyField, limit) //从数据库中查询数据，并将结果存储到 articles 切片中
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (m *customArticleModel) UpdateArticleStatus(ctx context.Context, id int64, status int) error {
	beyondArticleArticleIdKey := fmt.Sprintf("%s%v", cacheBeyondArticleArticleIdPrefix, id)
	// 执行数据库更新操作
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		// 创建SQL更新语句
		query := fmt.Sprintf("update %s set status = ? where `id` = ?", m.table)
		// 执行更新操作
		return conn.ExecCtx(ctx, query, status, id)
	}, beyondArticleArticleIdKey)

	return err
}
