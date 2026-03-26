package system

import (
	"github.com/gin-gonic/gin"
	"github.com/irorikon/cloudgram-go/config"
	"github.com/irorikon/cloudgram-go/model/request"
	"github.com/irorikon/cloudgram-go/model/response"
	"github.com/irorikon/cloudgram-go/utils"
)

type SystemAuthApi struct{}

// Login 登录
func (a *SystemAuthApi) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("Request body must be valid JSON", c)
		return
	}
	// 验证用户名密码
	err := utils.Verify(req, utils.LoginVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.Username != config.AuthUser {
		response.FailWithMessage("Invalid username or password", c)
		return
	}
	if req.Password != utils.Base64Encode([]byte(config.AuthPassword)) {
		response.FailWithMessage("Invalid username or password", c)
		return
	}
	// 生成token
	token, claims, err := utils.LoginToken(req.Username)
	if err != nil {
		response.FailWithMessage("Invalid username or password", c)
		return
	}
	response.OKWithDetailed(response.LoginResponse{
		ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		Token:     token,
		User: response.User{
			Username: claims.Username,
		},
	}, "登录成功", c)
}

// refresh 刷新 Token
func (a *SystemAuthApi) Refresh(c *gin.Context) {
	// 获取当前的claims
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("Token refresh failed", c)
		return
	}

	// 使用当前用户名生成新的token
	token, newClaims, err := utils.LoginToken(claims.Username)
	if err != nil {
		response.FailWithMessage("Token refresh failed", c)
		return
	}

	response.OKWithDetailed(response.LoginResponse{
		ExpiresAt: newClaims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		Token:     token,
		User: response.User{
			Username: claims.Username,
		},
	}, "Token refresh successful", c)
}
