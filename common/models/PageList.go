package models

import "gorm.io/gorm"

type PageList struct {
	List       any `json:"list"`
	PageNo     int `json:"pageNo"`
	PageSize   int `json:"pageSize"`
	PageTotal  int `json:"pageTotal"`
	TotalCount int `json:"totalCount"`
}

func QueryPageList[T any](tx *gorm.DB, pageNo int, pageSize int, value T) PageList {
	var totalCount int64
	db := tx.Offset((pageNo - 1) * pageSize).Limit(pageSize)
	db.Count(&totalCount)
	db.Find(&value)
	return PageList{
		List:       value,
		PageNo:     pageNo,
		PageSize:   pageSize,
		PageTotal:  int(totalCount / int64(pageSize)),
		TotalCount: int(totalCount),
	}

}
