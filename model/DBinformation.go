package models

//DBinformation 資料庫資訊
type DBinformation struct {
	UserName string
	Password string
	Host     string
	Port     string
	Database string
	MaxIdle  int
	MaxOpen  int
}
