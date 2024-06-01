package handler

import (
	"github.com/McaxDev/Salary/config"
	"gorm.io/gorm"
)

var conf *config.ConfigStcuct
var DB *gorm.DB

func Init() {
	conf = &config.Config
	DB = config.DB
}
