package handler

import (
	"github.com/McaxDev/Salary/config"
	"github.com/McaxDev/Salary/util"
	"github.com/gin-gonic/gin"
)

// 上传招聘信息
func AdminEdit(c *gin.Context) {

	// 验证用户是否是管理员
	if _, admin, err := util.BindJwt(c); err != nil || !admin {
		util.Error(c, 400, "无法验证你的管理员身份", err)
		return
	}

	if c.Request.URL.Path == "/upload/recruit" {

		// 读取请求体
		var data config.Recruit
		if err := c.ShouldBindJSON(&data); err != nil {
			util.Error(c, 400, "无法读取你的请求体", err)
			return
		}
		data.Click = 0

		// 将信息存储到数据库
		if err := DB.Create(&data).Error; err != nil {
			util.Error(c, 500, "无法将信息存储到数据库", err)
			return
		}

	} else if c.Request.URL.Path == "/edit/recruit" {

		// 读取请求体
		var data struct {
			ID      uint   `json:"id"`
			Title   string `json:"title"`
			Summary string `json:"summary"`
			Content string `json:"content"`
		}
		if err := c.ShouldBindJSON(&data); err != nil {
			util.Error(c, 400, "无法读取你的请求体", err)
			return
		}

		// 在数据库里更新信息
		if err := DB.Model(new(config.Recruit)).Where(
			"id = ?", data.ID,
		).Updates(&data).Error; err != nil {
			util.Error(c, 500, "更新数据失败", err)
			return
		}

	} else if c.Request.URL.Path == "/delete/recruit" {

		// 检查对应的数据是否存在
		id := c.Query("id")
		if err := DB.First(
			new(config.Recruit), "id = ?", id,
		).Error; err != nil {
			util.DbQueryError(c, err, "对应的数据并不存在")
			return
		}

		// 删除数据
		if err := DB.Delete(
			new(config.Recruit), "id = ?", id,
		).Error; err != nil {
			util.Error(c, 500, "删除数据失败", err)
			return
		}
	} else if c.Request.URL.Path == "/addsalary" {

		// 从请求体读取工资数据到对象
		var salary []config.SalaryRecord
		if err := c.BindJSON(&salary); err != nil {
			util.Error(c, 400, "无法读取你的请求", err)
			return
		}

		// 将工资信息添加到数据库
		if err := DB.Create(&salary).Error; err != nil {
			util.Error(c, 500, "数据添加失败，请联系管理员", err)
			return
		}

	} else {
		util.Error(c, 500, "服务器不存在这个路径", nil)
		return
	}

	// 返回客户端响应
	c.AbortWithStatusJSON(200, util.Resp("请求成功", nil))
}
