package handler

import (
	"github.com/McaxDev/Salary/config"
	"github.com/McaxDev/Salary/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SignUp(c *gin.Context) {

	// 将请求体内容保存到结构体对象
	var req struct {
		Certificate string `json:"certificate"`
		Username    string `json:"username"`
		Telephone   string `json:"telephone"`
		AuthCode    string `json:"authcode"`
	}
	if err := c.BindJSON(&req); err != nil {
		util.Error(c, 400, "无法读取你的请求体", err)
		return
	}

	// 验证验证码是否正确
	if !AuthSMS(req.Telephone, req.AuthCode) {
		util.Error(c, 400, "验证码不正确", nil)
		return
	}

	// 验证用户是否已经注册过了
	if err := DB.First(
		new(config.User), "telephone = ?", req.Telephone,
	).Error; err == nil {
		util.Error(c, 400, "你已经注册过了", nil)
		return
	} else if err != gorm.ErrRecordNotFound {
		util.Error(c, 500, "服务器错误，请联系管理员", err)
		return
	}

	// 创建用户
	user := config.User{
		Username:    req.Username,
		Telephone:   req.Telephone,
		Certificate: req.Certificate,
	}
	if err := DB.Create(&user).Error; err != nil {
		util.Error(c, 500, "注册失败，请联系管理员", err)
		return
	}

	// 获取jwt
	userJwt, err := util.GetJwt(user.Telephone)
	if err != nil {
		util.Error(c, 500, "创建用户信息失败，请联系管理员", err)
		return
	}

	// 将请求发送给客户端
	c.AbortWithStatusJSON(200, util.Resp("注册成功", gin.H{
		"token":     userJwt,
		"username":  req.Username,
		"telephone": req.Telephone,
		"admin":     false,
	}))

}
