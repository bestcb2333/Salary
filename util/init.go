package util

import "github.com/McaxDev/Salary/config"

var conf *config.ConfigStcuct

func Init() {
	conf = &config.Config
}
