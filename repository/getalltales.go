package repository

import (
	"github.com/jinzhu/gorm"
)

//GetAlltables 列出所有資料表
func (r Repository) GetAlltables(DB *gorm.DB, tablenames *[]string) error {
	if err := DB.Raw("show tables").Pluck("Tables Names", tablenames).Error; err != nil {
		return err
	}
	return nil
}
