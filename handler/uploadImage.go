package handler

import (
	"fmt"
	"path/filepath"

	"github.com/McaxDev/Salary/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 手动上传图像
func UploadImage(c *gin.Context) {

	// 从中间件读取用户信息
	if _, admin, err := util.BindJwt(c); err != nil || !admin {
		ImgUploadErr(c, 400, "只有管理员才能添加新图片", err)
		return
	}

	// 从请求体读取用户的图片
	file, err := c.FormFile("file")
	if err != nil || file == nil {
		ImgUploadErr(c, 400, "无法读取你的图片", err)
		return
	}

	// 将文件保存在文件系统
	fileName := uuid.New().String()
	fp := filepath.Join("image", fileName)
	if err := c.SaveUploadedFile(
		file, fp,
	); err != nil {
		ImgUploadErr(c, 500, "无法将文件上传到服务器", err)
		return
	}

	// 获取客户端请求URL
	imageURL := fmt.Sprintf(
		"%s://%s/%s",
		c.Request.URL.Scheme, c.Request.URL.Host, fp,
	)

	// 返回响应
	c.AbortWithStatusJSON(200, gin.H{
		"errno": 0,
		"data": gin.H{
			"url":  imageURL,
			"alt":  fileName,
			"href": imageURL,
		},
	})
}

// 返回错误响应并记录日志
func ImgUploadErr(c *gin.Context, status int, msg string, err error) {
	c.AbortWithStatusJSON(status, gin.H{"errno": 1, "message": msg})
	if err != nil {
		fmt.Println("上传文件时出现了错误："+ err.Error())                       
	}
}
