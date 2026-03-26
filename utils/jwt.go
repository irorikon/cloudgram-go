package utils

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/irorikon/cloudgram-go/config"
	"github.com/irorikon/cloudgram-go/logger"
	systemReq "github.com/irorikon/cloudgram-go/model/request"
)

type JWT struct {
	SecretKey []byte
}

var (
	ErrTokenValid            = errors.New("未知错误")
	ErrTokenExpired          = errors.New("token已过期")
	ErrTokenNotValidYet      = errors.New("token尚未激活")
	ErrTokenMalformed        = errors.New("这不是一个token")
	ErrTokenSignatureInvalid = errors.New("无效签名")
	ErrTokenInvalid          = errors.New("无法处理此token")
)

// NewJWT 创建JWT实例
func NewJWT() *JWT {
	return &JWT{
		SecretKey: []byte(config.JwtSecretKey),
	}
}

// CreateClaims 创建自定义声明
func (j *JWT) CreateClaims(baseClaims systemReq.BaseClaims) systemReq.CustomClaims {
	// 确保配置被正确加载
	bf, _ := time.ParseDuration(config.GlobalJWTConfig.BufferTime)
	ep, _ := time.ParseDuration(config.GlobalJWTConfig.ExpiresTime)

	claims := systemReq.CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(bf / time.Second), // 缓冲时间（秒）
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{},
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)),
			Issuer:    config.GlobalJWTConfig.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return claims
}

// GenerateToken 生成token
func (j *JWT) GenerateToken(claims systemReq.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SecretKey)
}

// ParseToken 解析token
func (j *JWT) ParseToken(tokenString string) (*systemReq.CustomClaims, error) {
	// 检查token是否为空
	if tokenString == "" {
		return nil, errors.New("token为空")
	}

	token, err := jwt.ParseWithClaims(tokenString, &systemReq.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return j.SecretKey, nil
	})

	if err != nil {
		// 记录详细的错误信息
		logger.SysErrorF("解析token失败: %v", err)

		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, ErrTokenSignatureInvalid
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		}
		return nil, err
	}

	if token != nil && token.Valid {
		if claims, ok := token.Claims.(*systemReq.CustomClaims); ok {
			return claims, nil
		}
	}

	return nil, ErrTokenInvalid
}
