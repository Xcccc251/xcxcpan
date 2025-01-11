package test

import (
	"XcxcPan/common/helper"
	"XcxcPan/common/imageCode"
	"fmt"
	"os"
	"testing"
)

func TestCreateImageCode(t *testing.T) {
	// 创建新的验证码图片
	imgCode := imageCode.NewCreateImageCode()

	// 将图片保存到文件
	f, err := os.Create("test_captcha.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// 写入图片
	err = imgCode.Write(f)
	if err != nil {
		t.Fatal(err)
	}

	// 获取验证码文本
	code := imgCode.GetCode()
	if len(code) != 4 {
		t.Errorf("Expected code length 4, got %d", len(code))
	}

	t.Logf("Generated captcha code: %s", code)
}

func TestPassword(t *testing.T) {
	fmt.Println(helper.AnalysisBcryptPassword("$2a$10$dNxL5960nAAgvvWVEXJaqOwll2nudPBsEBdIEpje2L7zCAe9WO7Tu", "test123456"))
}
