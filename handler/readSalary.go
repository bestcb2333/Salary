package handler

import (
	"github.com/McaxDev/Salary/config"
	"github.com/McaxDev/Salary/util"
	"github.com/gin-gonic/gin"
)

// 用户查询工资
func ReadSalary(c *gin.Context) {

	// 从查询字符串参数获取身份证和月份
	certificate := c.Query("certificate")
	date := c.Query("date")

	// 查询工资信息
	var record config.SalaryRecord
	if err := DB.First(&record,
		"certificate = ? AND date = ?", certificate, date,
	).Error; err != nil {
		util.DbQueryError(c, err, "无法查询到符合条件的记录")
		return
	}

	// 存储用户的阅读记录
	if err := DB.Model(new(config.SalaryReaded)).Create(&config.SalaryReaded{
		Username:  record.Username,
		Telephone: record.Telephone,
	}).Error; err != nil {
		util.Error(c, 500, "无法存储你的阅读记录", err)
		return
	}

	// 将工资信息发送给客户端
	c.AbortWithStatusJSON(200, util.Resp("查询成功", &record))
}
