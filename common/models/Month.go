package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

type MyMonth time.Time

func (t *MyMonth) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")
	t1, err := time.Parse("2006-01", timeStr)
	*t = MyMonth(t1)
	return err
}

func (t MyMonth) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%v\"", time.Time(t).Format("2006-01"))
	return []byte(formatted), nil
}

func (t MyMonth) Value() (driver.Value, error) {
	// MyTime 转换成 time.Time 类型
	tTime := time.Time(t)
	return tTime.Format("2006-01"), nil
}

func (t *MyMonth) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		// 字符串转成 time.Time 类型
		*t = MyMonth(vt)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

func (t *MyMonth) String() string {
	return fmt.Sprintf("hhh:%s", time.Time(*t).String())
}
