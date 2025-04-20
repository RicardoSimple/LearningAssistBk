package util

import (
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"

	"ar-app-api/conf"
	"ar-app-api/model"
	"ar-app-api/util/log"
)

type JwtC struct {
	Secret        string `json:"secret"`
	Expire        int    `json:"expire"` // 超时时间(小时)
	RefreshSecret string `json:"refresh_secret"`
	RefreshExpire int    `json:"refresh_expire"` // refresh_token 超时时间(小时)
	AesKey        string `json:"aes_key"`
}

// 定义自定义的Claims结构
type CustomClaims struct {
	// user info
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Gender   string `json:"gender"`
	jwt.StandardClaims
}
type RefreshTokenClaims struct {
	ID uint `json:"id"`
	jwt.StandardClaims
}

var jc *JwtC

func InitJwtC(ctx context.Context) {
	log.Info("[INIT] jwt config init")
	jwtc := conf.GetConfig().Jwt
	jc = &JwtC{
		Secret:        jwtc.Secret,
		Expire:        jwtc.Expire,
		RefreshSecret: jwtc.RefreshSecret,
		RefreshExpire: jwtc.RefreshExpire,
		AesKey:        jwtc.AesKey,
	}
}

func GenerateTokens(user *model.User) (*model.TokenInfo, error) {
	// 访问令牌
	accessTokenExpiration := time.Now().Add(time.Duration(jc.Expire) * time.Hour)
	accessTokenClaims := &CustomClaims{
		ID:       user.ID,
		UserName: user.Username,
		Email:    user.Email,
		Gender:   user.Gender,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExpiration.Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(jc.Secret))
	if err != nil {
		return nil, err
	}

	encryptedAccessToken, err := encryptAES(jc.AesKey, accessTokenString)
	if err != nil {
		return nil, err
	}

	// 刷新令牌
	refreshTokenExpiration := time.Now().Add(time.Duration(jc.RefreshExpire) * time.Hour)
	refreshTokenClaims := &RefreshTokenClaims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpiration.Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(jc.Secret))
	if err != nil {
		return nil, err
	}

	encryptedRefreshToken, err := encryptAES(jc.AesKey, refreshTokenString)
	if err != nil {
		return nil, err
	}

	return &model.TokenInfo{
		Token:        encryptedAccessToken,
		RefreshToken: encryptedRefreshToken,
		ID:           user.ID,
		UserName:     user.Username,
	}, nil
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	decryptedToken, err := decryptAES(jc.AesKey, tokenString)
	if err != nil {
		return nil, err
	}

	claims := &CustomClaims{}
	// 解析 Token
	token, err := jwt.ParseWithClaims(decryptedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jc.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
func RefreshToken(oldRefreshToken string) (uint, error) {
	claims := &RefreshTokenClaims{}
	decryptedToken, err := decryptAES(jc.AesKey, oldRefreshToken)
	if err != nil {
		return 0, err
	}

	token, err := jwt.ParseWithClaims(decryptedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jc.Secret), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid refresh token")
	}

	return claims.ID, nil
}
