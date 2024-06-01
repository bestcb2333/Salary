package handler

import (
	"github.com/McaxDev/Salary/config"
	"github.com/McaxDev/Salary/util"
	"github.com/gin-gonic/gin"
)

// 登录请求处理函数
func Login(c *gin.Context) {

	// 将请求体绑定到结构体对象
	var req struct {
		Telephone string `json:"telephone"`
		AuthCode  string `json:"authcode"`
	}
	if err := c.BindJSON(&req); err != nil {
		util.Error(c, 400, "你的请求不合法", err)
		return
	}

	// 验证用户提供的验证码
	admin := false
	var user config.User
	if req.AuthCode == conf.Admin[req.Telephone] {
		admin = true
		user.Telephone = req.Telephone
	} else if req.AuthCode == conf.User[req.Telephone] {
		user.Telephone = req.Telephone
	} else if !AuthSMS(req.Telephone, req.AuthCode) {
		util.Error(c, 400, "验证码错误", nil)
		return
	} else if err := DB.First(
		&user, "telephone = ?", req.Telephone,
	).Error; err != nil {
		util.DbQueryError(c, err, "你尚未注册")
		return
	}

	// 生成jwt
	userJwt, err := util.GetJwt(user.Telephone)
	if err != nil {
		util.Error(c, 500, "生成jwt失败", err)
		return
	}

	// 将请求发送给客户端
	c.AbortWithStatusJSON(200, util.Resp("登录成功", gin.H{
		"token":     userJwt,
		"userId":    user.ID,
		"username":  user.Username,
		"telephone": user.Telephone,
		"admin":     admin,
	}))
}
