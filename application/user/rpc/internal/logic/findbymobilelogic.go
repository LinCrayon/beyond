package logic

import (
	"context"

	"github.com/LinCrayon/beyond/application/user/rpc/internal/svc"
	"github.com/LinCrayon/beyond/application/user/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByMobileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByMobileLogic {
	return &FindByMobileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByMobileLogic) FindByMobile(in *service.FindByMobileRequest) (*service.FindByMobileResponse, error) {
	user, err := l.svcCtx.UserModel.FindByMobile(l.ctx, in.Mobile)
	if err != nil {
		logx.Errorf("FindByMobile mobile: %s error: %v", in.Mobile, err)
		return nil, err
	}

	//用户不存在
	if user == nil {
		return &service.FindByMobileResponse{}, nil
	}

	return &service.FindByMobileResponse{
		UserId:   int64(user.Id),
		Username: user.Username,
		Avatar:   user.Avatar,
	}, nil
}
