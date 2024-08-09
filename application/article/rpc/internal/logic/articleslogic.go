package logic

import (
	"cmp"
	"context"
	"fmt"
	"github.com/LinCrayon/beyond/application/article/rpc/internal/code"
	"github.com/LinCrayon/beyond/application/article/rpc/internal/model"
	"github.com/LinCrayon/beyond/application/article/rpc/internal/types"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"
	"slices"
	"strconv"
	"time"

	"github.com/LinCrayon/beyond/application/article/rpc/internal/svc"
	"github.com/LinCrayon/beyond/application/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	prefixArticles = "biz#articles#%d#%d"
	articlesExpire = 3600 * 24 * 2
)

type ArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticlesLogic {
	return &ArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticlesLogic) Articles(in *pb.ArticlesRequest) (*pb.ArticlesResponse, error) {
	//排序类型==发布时间或者点赞数
	if in.SortType != types.SortPublishTime && in.SortType != types.SortLikeCount { //排序类型不等发布时间或者点赞数
		return nil, code.SortTypeInvalid //排序类型无效
	}
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	//设置默认的分页大小和游标
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.Cursor == 0 {
		if in.SortType == types.SortPublishTime { //发布时间排序
			in.Cursor = time.Now().Unix()
		} else {
			in.Cursor = types.DefaultSortLikeCursor //点赞数排序
		}
	}

	var (
		sortField       string //排序字段
		sortLikeNum     int64  // 1 点赞数排序
		sortPublishTime string // 0 发布时间排序
	)
	//根据排序类型设置排序字段和游标
	if in.SortType == types.SortLikeCount { //点赞数排序
		sortField = "like_num"  //设置排序字段为 like_num
		sortLikeNum = in.Cursor //将输入中的游标值赋值给sortLikeNum，表示排序的点赞数
	} else {
		sortField = "publish_time"
		sortPublishTime = time.Unix(in.Cursor, 0).Format("2006-01-02 15:04:05") //Unix时间戳转换为 time.Time 对象
	}

	var (
		err            error
		isCache, isEnd bool //标识缓存状态和是否到达数据结束
		lastId, cursor int64
		curPage        []*pb.ArticleItem //当前页
		articles       []*model.Article
	)
	// 从缓存中获取文章ID列表
	articleIds, _ := l.cacheArticles(l.ctx, in.UserId, in.Cursor, in.PageSize, in.SortType)
	if len(articleIds) > 0 { //缓存命中
		isCache = true                           //命中了缓存
		if articleIds[len(articleIds)-1] == -1 { //检查 articleIds 切片的最后一个元素是否为 -1
			isEnd = true // 到达数据结束
		}
		//TODO 根据 文章id查文章  MapReduce 操作
		articles, err = l.articleByIds(l.ctx, articleIds)
		if err != nil {
			return nil, err
		}

		// 通过sortFiled对articles切片进行排序 （根据排序字段对文章进行排序）
		var cmpFunc func(a, b *model.Article) int
		if sortField == "like_num" {
			cmpFunc = func(a, b *model.Article) int {
				return cmp.Compare(b.LikeNum, a.LikeNum) //比较 b.LikeNum 和 a.LikeNum ,返回-1，0，1
			}
		} else {
			cmpFunc = func(a, b *model.Article) int {
				return cmp.Compare(b.PublishTime.Unix(), a.PublishTime.Unix())
			}
		}
		slices.SortFunc(articles, cmpFunc) //将切片按比较函数的规则排序

		for _, article := range articles {
			curPage = append(curPage, &pb.ArticleItem{ //遍历append到当前页
				Id:           article.Id,
				Title:        article.Title,
				Content:      article.Content,
				LikeCount:    article.LikeNum,
				CommentCount: article.CommentNum,
				PublishTime:  article.PublishTime.Unix(),
			})
		}
	} else { //缓存未命中  获取文章数据，并构建当前页的文章列表
		//SingleFlightGroup.Do 方法确保在并发情况下，对于相同的键（key），只有一个请求会被实际执行，其他的请求会等待这个请求的结果 ,避免了多次重复查询数据库
		v, err := l.svcCtx.SingleFlightGroup.Do(fmt.Sprintf("ArticlesByUserId:%d:%d", in.UserId, in.SortType), func() (interface{}, error) {
			return l.svcCtx.ArticleModel.ArticlesByUserId(l.ctx, in.UserId, types.ArticleStatusVisible, sortLikeNum, sortPublishTime, sortField, types.DefaultLimit)
		})
		if err != nil {
			logx.Errorf("ArticlesByUserId userId: %d sortField: %s error: %v", in.UserId, sortField, err)
			return nil, err
		}
		if v == nil {
			return &pb.ArticlesResponse{}, nil
		}
		articles = v.([]*model.Article)        //类型断言
		var firstPageArticles []*model.Article //用于存储当前页的文章
		//分页处理
		if len(articles) > int(in.PageSize) {
			firstPageArticles = articles[:int(in.PageSize)] //取出前 PageSize 个文章作为第一页数据
		} else {
			firstPageArticles = articles //表示所有文章都在第一页
			isEnd = true
		}
		for _, article := range firstPageArticles {
			curPage = append(curPage, &pb.ArticleItem{
				Id:           article.Id,
				Title:        article.Title,
				Content:      article.Content,
				LikeCount:    article.LikeNum,
				CommentCount: article.CommentNum,
				PublishTime:  article.PublishTime.Unix(),
			})
		}
	}

	//去重逻辑
	if len(curPage) > 0 {
		pageLast := curPage[len(curPage)-1] //获取最后一篇文章 即当前页最后一篇文章的 ID
		lastId = pageLast.Id
		//更新游标(根据排序类型)
		if in.SortType == types.SortPublishTime {
			cursor = pageLast.PublishTime //cursor 设置为最后一篇文章的发布时间
		} else {
			cursor = pageLast.LikeCount
		}
		if cursor < 0 { //修正 cursor：确保 cursor 不小于 0，避免负值对后续逻辑造成影响
			cursor = 0
		}
		for k, article := range curPage {
			//根据排序类型检查文章
			if in.SortType == types.SortPublishTime {
				if article.PublishTime == in.Cursor && article.Id == in.ArticleId {
					curPage = curPage[k:] //从匹配的文章开始 （从用户指定的文章开始，避免显示重复的文章。）
					break
				}
			} else {
				if article.LikeCount == in.Cursor && article.Id == in.ArticleId {
					curPage = curPage[k:]
					break
				}
			}
		}
	}

	ret := &pb.ArticlesResponse{
		IsEnd:     isEnd,
		Cursor:    cursor,
		ArticleId: lastId,
		Articles:  curPage,
	}

	//写缓存操作
	if !isCache {
		threading.GoSafe(func() {
			if len(articles) < types.DefaultLimit && len(articles) > 0 {
				articles = append(articles, &model.Article{Id: -1})
			}
			err = l.addCacheArticles(context.Background(), articles, in.UserId, in.SortType)
			if err != nil {
				logx.Errorf("addCacheArticles error: %v", err)
			}
		})
	}

	return ret, nil
}

// 查询缓存的详情 (MapReduce 模式来并行处理文章的查询，并将结果聚合成最终的文章列表)
/*
Map:将输入数据集映射为中间的键值对集合 (每个单词映射到一个计数值（例如 (word, 1)）)
Shuffle:将 Map 阶段输出的中间结果按照键进行分组，以便于后续的 Reduce 阶段处理 (将所有相同的单词组合（例如，将所有 (word, 1) 对合并为 word 和 [1, 1, 1, ...])
Reduce:汇总和处理 Shuffle 阶段中的分组数据，生成最终结果。 (对每个单词的计数进行汇总，计算出每个单词的总出现次数，并生成结果)
*/
//过滤有效的文章ID（Map）、查找文章（Reduce）和收集文章（聚合）
func (l *ArticlesLogic) articleByIds(ctx context.Context, articleIds []int64) ([]*model.Article, error) {
	// 使用 MapReduce 并行处理库
	articles, err := mr.MapReduce[int64, *model.Article, []*model.Article](
		// TODO Map 阶段  遍历 articleIds 切片，将每个有效的 ID 发送到 source 通道, 过滤掉 ID 为 -1 的项
		func(source chan<- int64) { //Map 阶段(处理阶段)
			for _, aid := range articleIds { //遍历 articleIds 切片，将每个有效的 ID 发送到 source 通道, 过滤掉 ID 为 -1 的项
				if aid == -1 {
					continue
				}
				source <- aid
			}
		},
		//TODO Reduce 阶段  查找每个文章ID对应的文章
		func(id int64, writer mr.Writer[*model.Article], cancel func(error)) { //Reduce 阶段
			p, err := l.svcCtx.ArticleModel.FindOne(ctx, id) //FindOne 自动生成行记录的缓存
			if err != nil {
				cancel(err) //cancel 函数来处理错误，并终止处理
				return
			}
			writer.Write(p)
		},
		//TODO 聚合阶段   从 pipe 通道中接收文章，将它们添加到 articles 切片中
		func(pipe <-chan *model.Article, writer mr.Writer[[]*model.Article], cancel func(error)) { //聚合阶段 (从 pipe 通道中接收文章，将它们添加到 articles 切片中。)
			var articles []*model.Article
			for article := range pipe {
				articles = append(articles, article)
			}
			writer.Write(articles)
		})
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func articlesKey(uid int64, sortType int32) string {
	return fmt.Sprintf(prefixArticles, uid, sortType)
}

func (l *ArticlesLogic) cacheArticles(ctx context.Context, uid, cursor, ps int64, sortType int32) ([]int64, error) {
	key := articlesKey(uid, sortType)
	b, err := l.svcCtx.BizRedis.ExistsCtx(ctx, key) //判断key是否存在
	if err != nil {
		logx.Errorf("ExistsCtx key: %s error: %v", key, err)
	}
	if b { //key存在
		// 对缓存键进行续期操作
		err = l.svcCtx.BizRedis.ExpireCtx(ctx, key, articlesExpire)
		if err != nil {
			logx.Errorf("ExpireCtx key: %s error: %v", key, err)
		}
	}
	//从有序集合中按分数范围倒序获取元素及其分数。使用这种方法可以实现高效的数据检索和分页显示  ,ps 表示每页显示的元素数量
	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(l.ctx, key, 0, cursor, 0, int(ps))
	if err != nil {
		logx.Errorf("ZrevrangebyscoreWithScoresAndLimit key: %s error: %v", key, err)
		return nil, err
	}
	var ids []int64              //用于存储文章ID
	for _, pair := range pairs { //遍历键值对列表
		id, err := strconv.ParseInt(pair.Key, 10, 64) //将键值对中的键转换为 int64 类型的文章ID。
		if err != nil {
			logx.Errorf("strconv.ParseInt key: %s error: %v", pair.Key, err)
			return nil, err
		}
		ids = append(ids, id) //将转换后的文章ID添加到 ids 切片中
	}
	return ids, nil
}

func (l *ArticlesLogic) addCacheArticles(ctx context.Context, articles []*model.Article, userId int64, sortType int32) error {
	// 如果文章列表为空，则直接返回
	if len(articles) == 0 {
		return nil
	}
	key := articlesKey(userId, sortType)
	// 遍历文章列表
	for _, article := range articles {
		var score int64
		//根据排序类型设置分数
		if sortType == types.SortLikeCount {
			score = article.LikeNum
		} else if sortType == types.SortPublishTime && article.Id != -1 {
			score = article.PublishTime.Local().Unix()
		}
		//防止无效的或不合适的负数值影响程序的逻辑
		if score < 0 {
			score = 0
		}
		// 将文章ID和分数添加到 Redis 有序集合中
		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, key, score, strconv.Itoa(int(article.Id)))
		if err != nil {
			return err
		}
	}
	// 对缓存键进行续期操作
	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, articlesExpire)
}
