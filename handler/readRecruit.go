package handler

import (
	"strconv"

	"github.com/McaxDev/Salary/config"
	"github.com/McaxDev/Salary/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 记录阅读了招聘信息的用户
func ReadRecruit(c *gin.Context) {

	// 读取用户信息
	user, _, err := util.BindJwt(c)
	if err != nil {
		util.Error(c, 500, "读取用户信息失败", err)
		return
	}

	if id, _ := strconv.Atoi(c.Query("id")); id != 0 {

		// 查询对应ID的信息
		var recruit config.Recruit
		if err := DB.First(
			&recruit, "id = ?", id,
		).Error; err != nil {
			util.DbQueryError(c, err, "无法查询到对应ID的信息")
			return
		}

		// 查看次数自增
		if err := DB.Model(new(config.Recruit)).Where("id = ?", id).Update(
			"click", gorm.Expr("click + ?", 1),
		).Error; err != nil {
			util.DbQueryError(c, err, "无法找到对应的记录")
			return
		}

		// 修改是否已阅读的字段
		if err := DB.Create(&config.RecruitReaded{
			Username:  user.Username,
			Telephone: user.Telephone,
			Title:     recruit.Title,
		}).Error; err != nil {
			util.Error(c, 500, "无法将阅读信息存入数据库", err)
			return
		}

		// 返回客户端响应
		c.AbortWithStatusJSON(200, util.Resp("查询成功", &recruit))

	} else if page, err := strconv.Atoi(c.Query("page")); page != 0 {

		// 获取总页数
		query := DB.Model(new(config.Recruit)).Where("deleted_at IS NULL")
		var totalrecruits int64
		if err := query.Count(&totalrecruits).Error; err != nil {
			util.Error(c, 500, "获取总页数失败", err)
			return
		}

		// 查询特定页数的概要
		var recruitList []struct {
			ID      uint   `json:"id"`
			Title   string `json:"title"`
			Summary string `json:"summary"`
			Click   int    `json:"click"`
		}
		if err := query.Order("id desc").Limit(10).Offset(
			(page - 1) * 10,
		).Find(&recruitList).Error; err != nil {
			util.DbQueryError(c, err, "无法查询到满足条件的记录")
			return
		}

		// 返回客户端响应
		c.AbortWithStatusJSON(200, gin.H{
			"msg":        "查询成功",
			"totalpages": (totalrecruits-1)/10 + 1,
			"data":       recruitList,
		})

	} else {
		util.Error(c, 400, "你提供的查询条件无效", err)
		return
	}
}
