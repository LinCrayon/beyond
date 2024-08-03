package logic

import (
	"context"
	"github.com/LinCrayon/beyond/application/follow/rpc/internal/code"
	"github.com/LinCrayon/beyond/application/follow/rpc/internal/model"
	"github.com/LinCrayon/beyond/application/follow/rpc/internal/types"
	"github.com/zeromicro/go-zero/core/threading"
	"strconv"
	"time"

	"github.com/LinCrayon/beyond/application/follow/rpc/internal/svc"
	"github.com/LinCrayon/beyond/application/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

const userFollowExpireTime = 3600 * 24 * 2

type FollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowListLogic {
	return &FollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FollowList 关注列表 （首先从缓存里获取，获取不到从数据库获取，获取一千条，取出当前页，被关注数的id都取出来，去count表查他的粉丝数，粉丝数赋值，返回Item,最后写缓存（关注时间倒叙显示））
func (l *FollowListLogic) FollowList(in *pb.FollowListRequest) (*pb.FollowListResponse, error) {
	if in.UserId == 0 {
		return nil, code.UserIdEmpty
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.Cursor == 0 {
		in.Cursor = time.Now().Unix() //时间戳作为游标,可以确保获取的关注列表是从最新的记录开始，这样用户会看到最近的关注记录
	}

	var (
		err             error
		isCache, isEnd  bool
		lastId, cursor  int64
		followedUserIds []int64
		follows         []*model.Follow
		curPage         []*pb.FollowItem
	)
	followUserIds, _ := l.cacheFollowUserIds(l.ctx, in.UserId, in.Cursor, in.PageSize)
	if len(followUserIds) > 0 {
		isCache = true
		if followUserIds[len(followUserIds)-1] == -1 { //检查最后一个元素是否为 -1
			followUserIds = followUserIds[:len(followUserIds)-1] //移除最后一个元素 ,忽略这个特殊标识
			isEnd = true
		}
		if len(followUserIds) == 0 {
			return &pb.FollowListResponse{}, nil
		}
		follows, err = l.svcCtx.FollowModel.FindByFollowedUserIds(l.ctx, in.UserId, followUserIds)
		if err != nil {
			l.Logger.Errorf("[FollowList] FollowModel.FindByFollowedUserIds error: %v req: %v", err, in)
			return nil, err
		}
		for _, follow := range follows {
			followedUserIds = append(followedUserIds, follow.FollowedUserID)
			curPage = append(curPage, &pb.FollowItem{
				Id:             follow.ID,
				FollowedUserId: follow.FollowedUserID,
				CreateTime:     follow.CreateTime.Unix(),
			})
		}
	} else {
		follows, err = l.svcCtx.FollowModel.FindByUserId(l.ctx, in.UserId, types.CacheMaxFollowCount)
		if err != nil {
			l.Logger.Errorf("[FollowList] FollowModel.FindByUserId error: %v req: %v", err, in)
			return nil, err
		}
		if len(follows) == 0 {
			return &pb.FollowListResponse{}, nil
		}
		var firstPageFollows []*model.Follow
		if len(follows) > int(in.PageSize) {
			firstPageFollows = follows[:in.PageSize] //截取当前页面的数据
		} else {
			firstPageFollows = follows //不需要分页
			isEnd = true
		}
		for _, follow := range firstPageFollows {
			followedUserIds = append(followedUserIds, follow.FollowedUserID)
			curPage = append(curPage, &pb.FollowItem{
				Id:             follow.ID,
				FollowedUserId: follow.FollowedUserID,
				CreateTime:     follow.CreateTime.Unix(),
			})
		}
	}
	if len(curPage) > 0 {
		pageLast := curPage[len(curPage)-1]
		lastId = pageLast.Id
		cursor = pageLast.CreateTime
		if cursor < 0 {
			cursor = 0
		}
		for k, follow := range curPage {
			if follow.CreateTime == in.Cursor && follow.Id == in.Id {
				curPage = curPage[k:]
				break
			}
		}
	}
	fc, err := l.svcCtx.FollowCountModel.FindByUserIds(l.ctx, followedUserIds)
	if err != nil {
		l.Logger.Errorf("[FollowList] FollowCountModel.FindByUserIds error: %v followedUserIds: %v", err, followedUserIds)
	}
	uidFansCount := make(map[int64]int)
	for _, f := range fc {
		uidFansCount[f.UserID] = f.FansCount
	}
	for _, cur := range curPage {
		cur.FansCount = int64(uidFansCount[cur.FollowedUserId])
	}
	ret := &pb.FollowListResponse{
		IsEnd:  isEnd,
		Cursor: cursor,
		Id:     lastId,
		Items:  curPage,
	}

	if !isCache {
		threading.GoSafe(func() {
			if len(follows) < types.CacheMaxFollowCount && len(follows) > 0 {
				follows = append(follows, &model.Follow{FollowedUserID: -1})
			}
			err = l.addCacheFollow(context.Background(), in.UserId, follows)
			if err != nil {
				logx.Errorf("addCacheFollow error: %v", err)
			}
		})
	}

	return ret, nil
}

func (l *FollowListLogic) cacheFollowUserIds(ctx context.Context, userId, cursor, pageSize int64) ([]int64, error) {
	key := userFollowKey(userId) //关注列表缓存 key
	b, err := l.svcCtx.BizRedis.ExistsCtx(ctx, key)
	if err != nil {
		logx.Errorf("[cacheFollowUserIds] BizRedis.ExistsCtx error: %v", err)
	}
	if b {
		err = l.svcCtx.BizRedis.ExpireCtx(ctx, key, userFollowExpireTime)
		if err != nil {
			logx.Errorf("[cacheFollowUserIds] BizRedis.ExpireCtx error: %v", err)
		}
	}
	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, 0, cursor, 0, int(pageSize))
	if err != nil {
		logx.Errorf("[cacheFollowUserIds] BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx error: %v", err)
		return nil, err
	}
	var uids []int64
	for _, pair := range pairs {
		uid, err := strconv.ParseInt(pair.Key, 10, 64)
		if err != nil {
			logx.Errorf("[cacheFollowUserIds] strconv.ParseInt error: %v", err)
			continue
		}
		uids = append(uids, uid)
	}

	return uids, nil
}

func (l *FollowListLogic) addCacheFollow(ctx context.Context, userId int64, follows []*model.Follow) error {
	//判断关注列表是否为空
	if len(follows) == 0 {
		return nil
	}
	key := userFollowKey(userId)
	//遍历关注列表并添加到 有序集合
	for _, follow := range follows {
		var score int64
		if follow.FollowedUserID == -1 {
			score = 0
		} else {
			//分数设置为 CreateTime 的 Unix 时间戳，以便按时间倒序显示
			score = follow.CreateTime.Unix() //关注时间倒叙显示
		}
		//添加到有序集合
		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, key, score, strconv.FormatInt(follow.FollowedUserID, 10))
		if err != nil {
			logx.Errorf("[addCacheFollow] BizRedis.ZaddCtx error: %v", err)
			return err
		}
	}

	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, userFollowExpireTime)
}
