package logic

import (
	"context"
	"fmt"
	"github.com/LinCrayon/beyond/application/follow/rpc/internal/code"
	"github.com/LinCrayon/beyond/application/follow/rpc/internal/model"
	"github.com/LinCrayon/beyond/application/follow/rpc/internal/types"
	"gorm.io/gorm"
	"strconv"
	"time"

	"github.com/LinCrayon/beyond/application/follow/rpc/internal/svc"
	"github.com/LinCrayon/beyond/application/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Follow 关注
func (l *FollowLogic) Follow(in *pb.FollowRequest) (*pb.FollowResponse, error) {
	//基本参数校验
	if in.UserId == 0 {
		return nil, code.FollowUserIdEmpty
	}
	if in.FollowedUserId == 0 {
		return nil, code.FollowedUserIdEmpty
	}
	if in.UserId == in.FollowedUserId { //自己不能关注自己
		return nil, code.CannotFollowSelf
	}
	//根据 userid 和 follow_user_id 查询用户
	follow, err := l.svcCtx.FollowModel.FindByUserIDAndFollowedUserID(l.ctx, in.UserId, in.FollowedUserId)
	if err != nil {
		l.Logger.Errorf("[Follow] FollowModel.FindByUserIDAndFollowedUserID err: %v req: %v", err, in)
		return nil, err
	}
	//用户不存在 或 已经关注该该用户
	if follow != nil && follow.FollowStatus == types.FollowStatusFollow {
		return &pb.FollowResponse{}, nil
	}
	//事务  （tx ：事务上下文）
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if follow != nil {
			// 更新关注状态   UpdateFields: 更新指定字段
			err = model.NewFollowModel(tx).UpdateFields(l.ctx, follow.ID, map[string]interface{}{
				"follow_status": types.FollowStatusFollow,
			})
		} else {
			err = model.NewFollowModel(tx).Insert(l.ctx, &model.Follow{
				UserID:         in.UserId,
				FollowedUserID: in.FollowedUserId,
				FollowStatus:   types.FollowStatusFollow,
				CreateTime:     time.Now(),
				UpdateTime:     time.Now(),
			})
		}
		if err != nil {
			return err
		}
		//增加关注计数
		err = model.NewFollowCountModel(tx).IncrFollowCount(l.ctx, in.UserId)
		if err != nil {
			return err
		}
		//增加粉丝计数
		return model.NewFollowCountModel(tx).IncrFansCount(l.ctx, in.FollowedUserId)
	})
	if err != nil {
		l.Logger.Errorf("[Follow] Transaction error: %v", err)
		return nil, err
	}
	//查缓存看是否存在 关注列表
	followExist, err := l.svcCtx.BizRedis.ExistsCtx(l.ctx, userFollowKey(in.UserId))
	if err != nil {
		l.Logger.Errorf("[Follow] Redis Exists error: %v", err)
		return nil, err
	}

	if followExist { // Zadd 添加到有序集合并指定分数
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, userFollowKey(in.UserId), time.Now().Unix(), strconv.FormatInt(in.FollowedUserId, 10))
		if err != nil {
			l.Logger.Errorf("[Follow] Redis Zadd error: %v", err)
			return nil, err
		}
		//关注数不用太多，缓存只保留一千条 （最新的在最上面，老的删除） Zremrangebyrank：按排名移除有序集合中的成员
		_, err = l.svcCtx.BizRedis.ZremrangebyrankCtx(l.ctx, userFollowKey(in.UserId), 0, -(types.CacheMaxFollowCount + 1))
		if err != nil {
			l.Logger.Errorf("[Follow] Redis Zremrangebyrank error: %v", err)
		}
	}
	//缓存是否存在 粉丝列表
	fansExist, err := l.svcCtx.BizRedis.ExistsCtx(l.ctx, userFansKey(in.FollowedUserId))
	if err != nil {
		l.Logger.Errorf("[Follow] Redis Exists error: %v", err)
		return nil, err
	}
	if fansExist {
		_, err = l.svcCtx.BizRedis.ZaddCtx(l.ctx, userFansKey(in.FollowedUserId), time.Now().Unix(), strconv.FormatInt(in.UserId, 10))
		if err != nil {
			l.Logger.Errorf("[Follow] Redis Zadd error: %v", err)
			return nil, err
		}
		_, err = l.svcCtx.BizRedis.ZremrangebyrankCtx(l.ctx, userFansKey(in.FollowedUserId), 0, -(types.CacheMaxFansCount + 1))
		if err != nil {
			l.Logger.Errorf("[Follow] Redis Zremrangebyrank error: %v", err)
		}
	}

	return &pb.FollowResponse{}, nil
}

func userFollowKey(userId int64) string {
	return fmt.Sprintf("biz#user#follow#%d", userId)
}

func userFansKey(userId int64) string {
	return fmt.Sprintf("biz#user#fans#%d", userId)
}
