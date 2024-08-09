package logic

import (
	"context"
	"fmt"
	"github.com/LinCrayon/beyond/application/user/rpc/user"
	"github.com/LinCrayon/beyond/pkg/util"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
	"time"

	"github.com/LinCrayon/beyond/application/applet/internal/svc"
	"github.com/LinCrayon/beyond/application/applet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	prefixVerificationCount = "biz#verification#count#%s"
	verificationLimitPerDay = 10
	expireActivation        = 60 * 30
)

type VerificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerificationLogic {
	return &VerificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerificationLogic) Verification(req *types.VerificationRequest) (resp *types.VerificationResponse, err error) {
	//获取缓存中手机号对应的验证码发送次数
	count, err := l.getVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("getVerificationCount mobile: %s error: %v", req.Mobile, err)
	}
	//当日不能超过10次,一个手机号产生10个验证码
	if count > verificationLimitPerDay {
		return nil, err
	}
	//获取缓存中的验证码(验证码的过期时间是30分钟)
	code, err := getActivationCache(req.Mobile, l.svcCtx.BizRedis)
	if err != nil {
		logx.Errorf("getActivationCache mobile: %s error: %v", req.Mobile, err)
	}
	if len(code) == 0 {
		code = util.RandomNumeric(6)
	}
	//调用发送验证码的rpc
	_, err = l.svcCtx.UserRPC.SendSms(l.ctx, &user.SendSmsRequest{Mobile: req.Mobile})
	if err != nil {
		return nil, err
	}
	//保存验证码到缓存，30分钟过期
	err = saveActivationCache(req.Mobile, code, l.svcCtx.BizRedis)
	if err != nil {
		logx.Errorf("saveActivationCache mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}
	err = l.incrVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("incrVerificationCount mobile: %s error: %v", req.Mobile, err)
	}

	return &types.VerificationResponse{}, nil
}

func (l *VerificationLogic) getVerificationCount(mobile string) (int, error) {
	key := fmt.Sprintf(prefixVerificationCount, mobile)
	val, err := l.svcCtx.BizRedis.Get(key)
	if err != nil {
		return 0, err
	}
	if len(val) == 0 {
		return 0, err
	}
	return strconv.Atoi(val)
}

func (l *VerificationLogic) incrVerificationCount(mobile string) error {
	key := fmt.Sprintf(prefixVerificationCount, mobile)
	_, err := l.svcCtx.BizRedis.Incr(key)
	if err != nil {
		return err
	}
	return l.svcCtx.BizRedis.Expireat(key, util.EndOfDay(time.Now()).Unix()) //设置过期时间
}

func getActivationCache(mobile string, rds *redis.Redis) (string, error) {
	key := fmt.Sprintf(prefixActivation, mobile)
	return rds.Get(key)
}

// 保存验证码 验证码的过期时间30分钟
func saveActivationCache(mobile, code string, rds *redis.Redis) error {
	key := fmt.Sprintf(prefixActivation, mobile)
	return rds.Setex(key, code, expireActivation)
}
func delActivationCache(mobile string, rds *redis.Redis) error {
	key := fmt.Sprintf(prefixActivation, mobile)
	_, err := rds.Del(key)
	return err
}
