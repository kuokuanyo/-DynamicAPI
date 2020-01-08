package repository

import (
	"github.com/jinzhu/gorm"
)

//Repository struct
type Repository struct{}

//ConnectDb 連接資料庫
func (r Repository) ConnectDb(engine string, SourceName string) (*gorm.DB, error) {
	DB, err := gorm.Open(engine, SourceName)
	if err != nil {
		return nil, err
	}
	return DB, nil
}
