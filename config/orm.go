package config

import (
	"database/sql/driver"
	"errors"
	"time"

	"gorm.io/gorm"
)

var DB *gorm.DB

type Time string

// 从数据库里读取Date字段的方法
func (d *Date) Scan(value any) error {
	if datetime, ok := value.([]uint8); ok {
		*d = Date(datetime)
		return nil
	}
	return errors.New("类型断言失败Date")
}

// 将年月存入数据库Date字段的方法
func (d Date) Value() (driver.Value, error) {
	return time.Now().Format("2006-01"), nil
}

// 从数据库里读取Time字段的方法
func (d *Time) Scan(value any) error {
	if datetime, ok := value.([]uint8); ok {
		*d = Time(datetime[:19])
		return nil
	}
	return errors.New("类型断言失败Time")
}

// 将年月存入数据库Time字段的方法
func (d Time) Value() (driver.Value, error) {
	return time.Now().Format("2006-01-02 15:04:05.000"), nil
}

type Date string

type User struct {
	gorm.Model
	Username    string // 用户名
	Certificate string // 身份证
	Telephone   string // 电话号码
}

type SalaryRecord struct {
	ID                 uint   `gorm:"primarykey"` // 用户UID
	CreatedAt          Time   // 记录创建时间
	Username           string // 姓名
	Certificate        string // 身份证号
	Telephone          string // 电话号码
	CreditCard         string // 银行卡号
	BasicSalary        string // 基本工资
	AttendanceRequired string // 应出勤天数
	AttendanceActual   string // 实出勤天数
	WorkHour           string // 总工时
	Performance        string // 绩效
	Allowance          string // 津贴
	Subsidy            string // 补助
	OvertimeSalary     string // 加班工资
	Excitation         string // 正负激励
	Discipline         string // 违纪扣款
	Withholding        string // 代扣部分
	BackPayment        string // 补发扣
	ShouldSalary       string // 应发工资
	UtilitiesFee       string // 水电物业费
	Tax                string // 个税
	AdvanceSalary      string // 预支工资
	ActualSalary       string // 实发工资
	Date                      // 对应年份
}

type RecruitReaded struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt Time
	Username  string
	Telephone string
	Title     string
}

type SalaryReaded struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt Time
	Username  string
	Telephone string
}

// 招聘信息列表
type Recruit struct {
	gorm.Model
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Content string `json:"content"`
	Click   int    `json:"click"`
}
