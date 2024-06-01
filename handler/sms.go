package handler

import (
	"fmt"
	"time"

	"github.com/McaxDev/Salary/config"
	"github.com/McaxDev/Salary/util"
	unisms "github.com/apistd/uni-go-sdk/sms"
	"github.com/gin-gonic/gin"
)

// 已经发送的验证码，键为验证码，值为过期时间和用户手机号
var SentSMS = make(map[string]*SentSmsValue)

type SentSmsValue struct {
	Telephone string
	Expire    time.Time
}

// 请求验证码的允许时间，键为客户IP，值为允许时间
var AllowTimeMap = make(map[string]time.Time)

// 请求短信验证码
func SendSMS(c *gin.Context) {

	// 从查询字符串参数获取手机号
	telephone := c.Query("telephone")

	// 检查上次请求的时间
	ip := c.ClientIP()
	if allowTime, ok := AllowTimeMap[ip]; ok {
		if timeLeft := allowTime.Sub(time.Now()).Seconds(); timeLeft > 0 {
			util.Error(c, 400, fmt.Sprintf("请 %.0f 秒后重试", timeLeft), nil)
			return
		}
		delete(AllowTimeMap, ip)
	}

	// 创建客户端
	conf := config.Config.SMS
	client := unisms.NewClient(conf["ID"])
	if config.Config.SMS["Key"] != "" {
		client = unisms.NewClient(conf["ID"], conf["Key"])
	}

	// 构造短信内容
	message := unisms.BuildMessage()
	message.SetTo(telephone)
	message.SetSignature(config.Config.SMS["signature"])
	message.SetTemplateId(config.Config.SMS["template"])
	authcode := util.RandStr(6)
	message.SetTemplateData(map[string]string{
		"code": authcode,
		"ttl":  "10",
	})

	// 发送短信
	resp, err := client.Send(message)
	if err != nil {
		util.Error(c, 500, "短信验证码发送失败", err)
		return
	}

	// 更新内存
	SentSMS[authcode] = &SentSmsValue{
		Telephone: telephone,
		Expire:    time.Now().Add(time.Minute * 10),
	}
	AllowTimeMap[ip] = time.Now().Add(time.Minute)

	// 发送响应
	c.AbortWithStatusJSON(200, util.Resp("验证码发送成功", resp.Message))
}

// 验证用户提交的验证码
func AuthSMS(telephone string, authcode string) bool {
	if value, ok := SentSMS[authcode]; ok {
		if value.Telephone == telephone && time.Now().Before(value.Expire) {
			return true
		}
	}
	return false
}

// 清理已经发送的SMS内存
func ClearSmsMemory() {
	for key, value := range SentSMS {
		if time.Now().After(value.Expire) {
			delete(SentSMS, key)
		}
	}
	for key, value := range AllowTimeMap {
		if time.Now().After(value) {
			delete(AllowTimeMap, key)
		}
	}
}
