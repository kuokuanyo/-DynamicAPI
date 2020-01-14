package models

//MysqlDBinformation 資料庫資訊
type MysqlDBinformation struct {
	UserName string
	Password string
	Host     string
	Port     string
	Database string
	MaxIdle  int
	MaxOpen  int
}

//MssqlDBinformation 資料庫資訊
type MssqlDBinformation struct {
	UserName string
	Password string
	Host     string
	Port     string
	Database string
	MaxIdle  int
	MaxOpen  int
}