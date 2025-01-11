package Enum

// RegisterOrLoginType 定义结构体来模拟枚举
type RegisterOrLoginType struct {
	RegisterType int
	Description  string
	Strategy     string
}

// 定义常量，模拟枚举值
var (
	EMAIL  = RegisterOrLoginType{RegisterType: 0, Description: "邮箱登录", Strategy: "email"}
	GITEE  = RegisterOrLoginType{RegisterType: 1, Description: "Gitee登录", Strategy: "gitee"}
	GITHUB = RegisterOrLoginType{RegisterType: 2, Description: "Github登录", Strategy: "github"}
)

// 枚举集合，用于查找和遍历
var registerOrLoginTypes = []RegisterOrLoginType{EMAIL, GITEE, GITHUB}

func GetType(r RegisterOrLoginType) int {
	return r.RegisterType
}

// 根据 RegisterType 查找枚举值
func GetRegisterOrLoginTypeByRegisterType(registerType int) (RegisterOrLoginType, bool) {
	for _, v := range registerOrLoginTypes {
		if v.RegisterType == registerType {
			return v, true
		}
	}
	return RegisterOrLoginType{}, false
}

// 根据策略查找枚举值
func GetRegisterOrLoginTypeByStrategy(strategy string) (RegisterOrLoginType, bool) {
	for _, v := range registerOrLoginTypes {
		if v.Strategy == strategy {
			return v, true
		}
	}
	return RegisterOrLoginType{}, false
}
