package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化程序
func ProgramInit() error {
	var err error
	var data []byte

	// 检测配置文件是否存在
	if _, err = os.Stat("config.json"); os.IsNotExist(err) {
		fmt.Println("未检测到配置文件，将新建并使用默认配置")
		if data, err = json.MarshalIndent(&Config, "", "  "); err != nil {
			return errors.New("配置文件编码失败：" + err.Error())
		}
		if err = os.WriteFile("config.json", data, 0644); err != nil {
			return errors.New("配置文件写入失败：" + err.Error())
		}
	} else if err != nil {
		return errors.New("配置文件状态异常：" + err.Error())
	} else {
		if data, err = os.ReadFile("config.json"); err != nil {
			return errors.New("配置文件读取失败：" + err.Error())
		}
		if err = json.Unmarshal(data, &Config); err != nil {
			return errors.New("配置文件格式错误，请删除：" + err.Error())
		}
	}

	// 连接数据库
	if DB, err = gorm.Open(mysql.Open(fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Config.MySQL["user"],
		Config.MySQL["password"],
		Config.MySQL["address"],
		Config.MySQL["port"],
		Config.MySQL["database"],
	))); err != nil {
		return err
	}

	// 自动迁移数据库表
	err = DB.AutoMigrate(
		new(User),
		new(SalaryRecord),
		new(Recruit),
		new(RecruitReaded),
		new(SalaryReaded),
	)
	if err != nil {
		return err
	}
	return nil
}
