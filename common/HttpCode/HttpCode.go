package HttpCode

type AppHttpCodeEnum int

const (
	SUCCESS                    AppHttpCodeEnum = 200
	FAILURE                    AppHttpCodeEnum = 500
	USERNAME_OR_PASSWORD_ERROR AppHttpCodeEnum = 1001
	NOT_LOGIN                  AppHttpCodeEnum = 1002
	NO_PERMISSION              AppHttpCodeEnum = 1003
	REQUEST_FREQUENTLY         AppHttpCodeEnum = 1004
	VERIFY_CODE_ERROR          AppHttpCodeEnum = 1005
	USERNAME_OR_EMAIL_EXIST    AppHttpCodeEnum = 1006
	PARAM_ERROR                AppHttpCodeEnum = 1007
	OTHER_ERROR                AppHttpCodeEnum = 1008
	SESSION_LIMIT              AppHttpCodeEnum = 1009
	NO_DELETE_CHILD_MENU       AppHttpCodeEnum = 1010
	FILE_UPLOAD_ERROR          AppHttpCodeEnum = 1011
	BLACK_LIST_ERROR           AppHttpCodeEnum = 1012
)

// 定义一个map来存储错误码对应的消息
var AppHttpCodeEnumMsg = map[AppHttpCodeEnum]string{
	SUCCESS:                    "success",
	FAILURE:                    "failure",
	USERNAME_OR_PASSWORD_ERROR: "用户名或密码错误",
	NOT_LOGIN:                  "请先登录",
	NO_PERMISSION:              "无权限",
	REQUEST_FREQUENTLY:         "请求频繁",
	VERIFY_CODE_ERROR:          "验证码错误",
	USERNAME_OR_EMAIL_EXIST:    "用户名或邮箱已存在",
	PARAM_ERROR:                "参数错误",
	OTHER_ERROR:                "其他错误",
	SESSION_LIMIT:              "会话数量已达上限",
	NO_DELETE_CHILD_MENU:       "存在子菜单，无法删除",
	FILE_UPLOAD_ERROR:          "文件上传失败",
	BLACK_LIST_ERROR:           "账号被封禁",
}

// 获取错误码对应的消息
func (code AppHttpCodeEnum) GetMsg() string {
	return AppHttpCodeEnumMsg[code]
}

// 获取错误码
func (code AppHttpCodeEnum) GetCode() int {
	return int(code)
}
