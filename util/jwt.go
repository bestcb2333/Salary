package util

import (
	"errors"
	"time"

	"github.com/McaxDev/Salary/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// 通过用户id生成jwt
func GetJwt(telephone string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"telephone": telephone,
		"exp":       time.Now().Add(time.Hour).Unix(),
	}).SignedString([]byte(conf.JwtKey))
}

// 验证JWT的handler
func AuthJwt(c *gin.Context) {

	// 解析jwt的格式
	JwtToken, err := jwt.Parse(
		c.GetHeader("Authorization")[7:],
		func(token *jwt.Token) (any, error) { return []byte(conf.JwtKey), nil },
	)
	if err != nil {
		Error(c, 500, "解析JWT失败", err)
		return
	}

	// 检查jwt是否通过
	claims, ok := JwtToken.Claims.(jwt.MapClaims)
	if !ok || !JwtToken.Valid {
		Error(c, 400, "token身份信息有误", nil)
		return
	}

	// 将jwt传递给后续的业务逻辑函数
	telephone, ok := claims["telephone"].(string)
	if !ok {
		Error(c, 400, "你的jwt无效", nil)
		return
	}
	c.Set("telephone", telephone)

	// 生成新的jwt并设置
	newJwt, err := GetJwt(telephone)
	if err != nil {
		Error(c, 400, "生成新的jwt失败", err)
		return
	}
	c.Header("Authorization", "Bearer "+newJwt)
}

// 通过请求里的jwt将用户数据绑定到结构体对象
func BindJwt(c *gin.Context) (*config.User, bool, error) {

	// 从请求里读取用户ID
	telephoneAny, exist := c.Get("telephone")
	if !exist {
		return nil, false, errors.New("读取用户数据失败")
	}
	telephone, ok := telephoneAny.(string)
	if !ok {
		return nil, false, errors.New("用户数据类型不正确")
	}

	// 判断用户是否是管理员
	if _, ok := conf.Admin[telephone]; ok {
		return &config.User{
			Telephone: telephone,
			Username:  "管理员",
		}, true, nil
	}

	// 判断用户是否是预注册用户
	if _, ok := conf.User[telephone]; ok {
		return &config.User{
			Telephone: telephone,
			Username:  "普通用户",
		}, false, nil
	}

	// 根据用户ID在数据库里查找用户并返回
	var user config.User
	if err := config.DB.First(
		&user, "telephone = ?", telephone,
	).Error; err != nil {
		return nil, false, err
	}
	return &user, false, nil
}
