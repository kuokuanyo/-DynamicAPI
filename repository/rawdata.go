package repository

import (
	"database/sql"

	"github.com/jinzhu/gorm"
)

//RawData 取得資料表欄位資訊
func (r Repository) RawData(DB *gorm.DB, describetable string) (*sql.Rows, error) {
	rows, err := DB.Raw(describetable).Rows()
	if err != nil {
		return nil, err
	}
	return rows, nil
}
