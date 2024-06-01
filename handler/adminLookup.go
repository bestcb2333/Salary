package handler

import (
	"strconv"

	"github.com/McaxDev/Salary/config"
	"github.com/McaxDev/Salary/util"
	"github.com/gin-gonic/gin"
)

// 管理员用户查看工资阅读或招聘信息阅读的记录
func AdminLookup(c *gin.Context) {

	// 验证用户身份是否是管理员
	if _, admin, err := util.BindJwt(c); err != nil || !admin {
		util.Error(c, 400, "只有管理员才能进行此操作", err)
		return
	}

	// 根据查询字符串的参数构造查询条件
	query := DB.Limit(10).Order("id desc")
	var records any
	if c.Request.URL.Path == "/allsalary" { // 查工资条
		query = query.Model(new(config.SalaryRecord))
		records = make([]config.SalaryRecord, 0)
	} else if c.Request.URL.Path == "/recruitreader" { // 查招聘阅读者
		query = query.Model(new(config.RecruitReaded))
		if title := c.Query("title"); title != "" {
			query = query.Where("title = ?", title)
		}
		records = make([]config.RecruitReaded, 0)
	} else { // 查工资阅读者
		query = query.Model(new(config.SalaryReaded))
		records = make([]config.SalaryReaded, 0)
	}
	if name := c.Query("name"); name != "" {
		query = query.Where("username = ?", name)
	}
	if date := c.Query("date"); date != "" {
		query = query.Where("LEFT(created_at, 7) = ?", date)
	}
	if telephone := c.Query("telephone"); telephone != "" {
		query = query.Where("telephone = ?", telephone)
	}

	// 获取数据条目数量
	var totalrecords int64
	if err := query.Count(&totalrecords).Error; err != nil {
		util.Error(c, 500, "查询数据数量失败", err)
		return
	}

	// 获取查询页数
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		util.Error(c, 400, "你提供的页数无效", err)
		return
	}

	// 执行查询
	if err := query.Offset((page - 1) * 10).Find(&records).Error; err != nil {
		util.Error(c, 500, "查询失败，请联系管理员", err)
		return
	}
	c.AbortWithStatusJSON(200, gin.H{
		"msg":   "查询成功",
		"total": totalrecords,
		"data":  records,
	})
}
