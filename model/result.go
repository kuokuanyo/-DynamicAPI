package models

import "database/sql"

//MysqlResult 資料表結構
type MysqlResult struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default sql.NullString
	Extra   string
}

//MssqlResult 資料表結構
type MssqlResult struct {
	Field   string
	Type    string
	Null    string
	Default sql.NullString
}
