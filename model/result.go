package models

import "database/sql"

//Result 資料表結構
type Result struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default sql.NullString
	Extra   string
}
