package logic

import (
	"context"
	"encoding/json"
	"github.com/LinCrayon/beyond/application/user/rpc/user"

	"github.com/LinCrayon/beyond/application/applet/internal/svc"
	"github.com/LinCrayon/beyond/application/applet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (resp *types.UserInfoResponse, err error) {
	//Value方法返回的是 interface{} 类型，因此需要通过类型断言将其转换为具体的类型，这里是 json.Number
	userId, err := l.ctx.Value(types.UserIdKey).(json.Number).Int64()
	if err != nil {
		return nil, err
	}
	if userId == 0 {
		return &types.UserInfoResponse{}, nil
	}
	//RPC userRPC 根据id查询用户
	u, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{UserId: userId})
	if err != nil {
		logx.Errorf("FindById userId: %d error: %v", userId, err)
		return nil, err
	}

	return &types.UserInfoResponse{
		UserId:   u.UserId,
		Username: u.Username,
		Avatar:   u.Avatar,
	}, nil
}
