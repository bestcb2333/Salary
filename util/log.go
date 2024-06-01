package util

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 向客户端输出错误信息，并将错误输出到服务器
func Error(c *gin.Context, status int, msg string, err error) {
	c.AbortWithStatusJSON(status, Resp(msg, nil))
	if err != nil {
		fmt.Println("出现了错误：" + err.Error())
	}
}

func DbQueryError(c *gin.Context, err error, message string) {
	if err == gorm.ErrRecordNotFound {
		Error(c, 400, message, nil)
	} else {
		Error(c, 500, "服务器错误，请联系管理员", err)
	}
}

// 制作响应体
func Resp(msg string, data any) gin.H {
	return gin.H{"msg": msg, "data": data}
}
