package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type (
	// TokenOptions 用于配置生成令牌的选项，包括密钥、过期时间和自定义字段
	TokenOptions struct {
		AccessSecret string
		AccessExpire int64
		Fields       map[string]interface{} //存储在token中的附加信息
	}
	// Token : 表示生成的令牌及其过期时间
	Token struct {
		AccessToken  string `json:"access_token"`
		AccessExpire int64  `json:"access_expire"`
	}
)

func BuildTokens(opt TokenOptions) (Token, error) {
	var token Token
	now := time.Now().Add(-time.Minute).Unix() //当前时间减一分钟
	accessToken, err := genToken(now, opt.AccessSecret, opt.Fields, opt.AccessExpire)
	if err != nil {
		return token, err
	}
	token.AccessToken = accessToken
	token.AccessExpire = now + opt.AccessExpire

	return token, nil
}

// 带有自定义声明（claims）的 JWT 令牌
func genToken(iat int64, secretKey string, payloads map[string]interface{}, seconds int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds //iat: 签发时间（issued at），Unix 时间戳。
	claims["iat"] = iat
	for k, v := range payloads { //payloads: 要包含在令牌中的自定义声明。
		claims[k] = v //遍历 payloads map，将所有键值对添加到 claims 中
	}
	token := jwt.New(jwt.SigningMethodHS256) //创建 JWT 令牌使用hs256加密
	token.Claims = claims
	//签名并生成令牌字符串
	return token.SignedString([]byte(secretKey))
}
