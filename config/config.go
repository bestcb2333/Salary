package config

type ConfigStcuct struct {
	Port   string
	JwtKey string
	MySQL  map[string]string
	SMS    map[string]string
	Admin  map[string]string
	User   map[string]string
	Tip    map[string]any
}

var Config = ConfigStcuct{
	Port:   "8080",
	JwtKey: "5201314",
	MySQL: map[string]string{
		"user":     "salary",
		"password": "salary",
		"database": "salary",
		"address":  "localhost",
		"port":     "3306",
	},
	SMS: map[string]string{
		"ID":        "",
		"Key":       "",
		"signature": "实习生小助手",
		"template":  "pub_verif_ttl3",
	},
	Admin: map[string]string{
		"86100000000": "admin1",
	},
	User: map[string]string{
		"88800000000": "user1",
	},
	Tip: map[string]any{
		"Port":   "后端端口，保持不变",
		"JwtKey": "密钥，不变",
		"MySQL": map[string]string{
			"address":  "数据库地址，不变",
			"database": "数据库名称，不变",
			"password": "数据库密码",
			"port":     "数据库端口，不变",
			"user":     "数据库用户名",
		},
		"SMS": map[string]string{
			"ID":        "发信ID，不变",
			"Key":       "发信KEY，不填",
			"signature": "签名",
			"template":  "发信内容模板，不变",
		},
		"Admin": map[string]string{
			"管理员账号": "管理员密码",   
			"管理员账号2": "管理员密码2",
			"可无限制添加": "可无限制添加",
		},
		"User": map[string]string{
			"测试用户名": "测试用户密码",  
			"测试用户名2": "测试用户密码2",
			"可无限制添加": "可无限制添加",
		},
	},
}
