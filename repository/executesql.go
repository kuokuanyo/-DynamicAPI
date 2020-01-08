package repository

import (
	"github.com/jinzhu/gorm"
)

//Execbyvalue 執行sql指令(有值)
func (r Repository) Execbyvalue(DB *gorm.DB, sql string, value []interface{}) error {
	if err := DB.Exec(sql, value...).Error; err != nil {
		return err
	}
	return nil
}

//Exec 執行sql指令
func (r Repository) Exec(DB *gorm.DB, sql string) error {
	if err := DB.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}
