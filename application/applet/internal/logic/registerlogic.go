package logic

import (
	"context"
	"github.com/LinCrayon/beyond/application/applet/internal/code"
	"github.com/LinCrayon/beyond/application/user/rpc/user"
	"github.com/LinCrayon/beyond/pkg/encrypt"
	"github.com/LinCrayon/beyond/pkg/jwt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strings"

	"github.com/LinCrayon/beyond/application/applet/internal/svc"
	"github.com/LinCrayon/beyond/application/applet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	prefixActivation = "biz#activation#%s"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	//参数基本校验
	req.Name = strings.TrimSpace(req.Name) //strings.TrimSpace函数用于去除字符串两端的空白字符
	req.Mobile = strings.TrimSpace(req.Mobile)
	if len(req.Mobile) == 0 {
		return nil, code.RegisterMobileEmpty
	}
	req.Password = strings.TrimSpace(req.Password)
	if len(req.Password) == 0 {
		return nil, code.RegisterPasswdEmpty
	} else {
		req.Password = encrypt.EncPassword(req.Password) //密码用md5加密
	}
	req.VerificationCode = strings.TrimSpace(req.VerificationCode)
	if len(req.VerificationCode) == 0 {
		return nil, code.VerificationCodeEmpty
	}
	//验证验证码是否有效
	err = checkVerificationCode(l.svcCtx.BizRedis, req.Mobile, req.VerificationCode)
	if err != nil {
		logx.Errorf("checkVerificationCode error: %v", err)
		return nil, err
	}
	//手机号进行AES加密 (公司中手机号一定要加密的)
	mobile, err := encrypt.EncMobile(req.Mobile)
	if err != nil {
		logx.Errorf("EncMobile mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}
	//调用UserPRC中用手机号查询用户
	u, err := l.svcCtx.UserRPC.FindByMobile(l.ctx, &user.FindByMobileRequest{Mobile: mobile})
	if err != nil {
		logx.Errorf("FindByMobile error: %v", err)
		return nil, err
	}
	//用户已经存在
	if u != nil && u.UserId > 0 {
		return nil, code.MobileHasRegistered
	}
	//根据手机号和用户名注册用户
	regRet, err := l.svcCtx.UserRPC.Register(l.ctx, &user.RegisterRequest{
		Username: req.Name,
		Mobile:   mobile,
	})
	if err != nil {
		logx.Errorf("Register error: %v", err)
		return nil, err
	}
	//生成jwt,在token中保存了自定义的数据
	token, err := jwt.BuildTokens(jwt.TokenOptions{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fields: map[string]interface{}{
			"userId": regRet.UserId, //在claim中存userId
		},
	})
	if err != nil {
		logx.Errorf("BuildTokens error: %v", err)
		return nil, err
	}
	//注册成功删除缓存中的验证码
	_ = delActivationCache(req.Mobile, l.svcCtx.BizRedis)

	return &types.RegisterResponse{
		UserId: regRet.UserId,
		Token: types.Token{
			AccessToken:  token.AccessToken,
			AccessExpire: token.AccessExpire,
		},
	}, nil
}

// 验证验证码是否正确存在
func checkVerificationCode(rds *redis.Redis, mobile, code string) error {
	cacheCode, err := getActivationCache(mobile, rds)
	if err != nil {
		return err
	}
	//验证码为空串的话==过期
	if cacheCode == "" {
		return errors.New("verification code expired")
	}
	if cacheCode != code {
		return errors.New("verification code failed")
	}
	return nil
}
