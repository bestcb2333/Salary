package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/McaxDev/Salary/config"
	"github.com/McaxDev/Salary/handler"
	"github.com/McaxDev/Salary/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// 初始化程序
	if err := config.ProgramInit(); err != nil {
		log.Fatalln("程序初始化失败：" + err.Error())
	}
	util.Init()
	handler.Init()

	// 定时清理内存
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for {
			<-ticker.C
			handler.ClearSmsMemory()
			fmt.Println(time.Now().String() + " 内存清理完成")
		}
	}()

	// 创建图片文件夹
	if _, err := os.Stat("./image"); os.IsNotExist(err) {
		if err := os.MkdirAll("./image", os.ModePerm); err != nil {
			log.Fatalln("图片文件夹创建失败：" + err.Error())
		}
	} else if err != nil {
		log.Fatalln("无法判断图片文件夹是否存在：" + err.Error())
	}

	//允许CORS跨域
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Static("/image", "./image")
	r.GET("/sms", handler.SendSMS)
	r.POST("/signup", handler.SignUp)
	r.POST("/login", handler.Login)
	r.Use(util.AuthJwt)
	r.GET("/readrecruit", handler.ReadRecruit)
	r.GET("/readsalary", handler.ReadSalary)
	r.GET("/recruitreader", handler.AdminLookup)
	r.GET("/salaryreader", handler.AdminLookup)
	r.GET("/allsalary", handler.AdminLookup)
	r.POST("/upload/image", handler.UploadImage)
	r.POST("/addsalary", handler.AdminEdit)
	r.POST("/upload/recruit", handler.AdminEdit)
	r.POST("/edit/recruit", handler.AdminEdit)
	r.GET("/delete/recruit", handler.AdminEdit)
	r.Run(":" + config.Config.Port)
}
