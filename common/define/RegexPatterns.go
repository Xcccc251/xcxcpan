package define

import "regexp"

var (
	PHONE_REGEX       = "^1([38][0-9]|4[579]|5[0-3,5-9]|6[6]|7[0135678]|9[89])\\d{8}$"
	EMAIL_REGEX       = "^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$"
	PASSWORD_REGEX    = "^\\w{4,32}$"
	VERIFY_CODE_REGEX = "^[a-zA-Z\\d]{6}$"
)

// IsPhoneValid 检查字符串是否符合手机号正则表达式
func IsPhoneValid(phone string) bool {
	re := regexp.MustCompile(PHONE_REGEX)
	return re.MatchString(phone)

}

// IsEmailValid 检查字符串是否符合邮箱正则表达式
func IsEmailValid(email string) bool {
	re := regexp.MustCompile(EMAIL_REGEX)
	return re.MatchString(email)
}

// IsPasswordValid 检查字符串是否符合密码正则表达式
func IsPasswordValid(password string) bool {
	re := regexp.MustCompile(PASSWORD_REGEX)
	return re.MatchString(password)
}

// IsVerifyCodeValid 检查字符串是否符合验证码正则表达式
func IsVerifyCodeValid(verifyCode string) bool {
	re := regexp.MustCompile(VERIFY_CODE_REGEX)
	return re.MatchString(verifyCode)
}
