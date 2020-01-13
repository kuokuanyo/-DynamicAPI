package controllers

import (
	models "DynamicAPI/model"
	"DynamicAPI/repository"
	"DynamicAPI/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

//MssqlDB mssql引擎
var MssqlDB *gorm.DB

//MysqlDB mysql引擎
var MysqlDB *gorm.DB

//information資料庫資訊
var information models.DBinformation

//Controller struct
type Controller struct{}

//ConnectDb 測試是否能連接資料庫
//@Summary 連接資料庫
//@Tags Connect Database(Must be connected first)
//@Accept json
//@Produce json
//@Param sql path string true "資料庫引擎"
//@Param information body models.DBinformation false "資料庫資訊"
//@Success 200 {string} string "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/opendb/{sql} [post]
func (c Controller) ConnectDb() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			message models.Error
			params  = mux.Vars(r) //印出url參數
			err     error
			Repo    repository.Repository
		)

		//decode
		json.NewDecoder(r.Body).Decode(&information)

		switch strings.ToLower(params["sql"]) {
		case "mysql":
			//完整的資料格式: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
			MysqlDataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
				information.UserName,
				information.Password,
				information.Host,
				information.Port,
				information.Database)

			//開啟資料庫連線
			MysqlDB, err = Repo.ConnectDb("mysql", MysqlDataSourceName)
			if err != nil {
				message.Message = "Connect Database failed"
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}

			MysqlDB.DB().SetMaxIdleConns(information.MaxIdle)
			MysqlDB.DB().SetMaxOpenConns(information.MaxOpen)
			utils.SendSuccess(w, "Successfully Connect Database")

		case "mssql":
			MssqlDataSourceName := fmt.Sprintf("sqlserver://%s:%s@%s:%s? database=%s",
				information.UserName,
				information.Password,
				information.Host,
				information.Port,
				information.Database)

			MssqlDB, err = Repo.ConnectDb("mssql", MssqlDataSourceName)
			if err != nil {
				message.Message = "Connect Database failed"
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}

			MssqlDB.DB().SetMaxIdleConns(information.MaxIdle)
			MssqlDB.DB().SetMaxOpenConns(information.MaxOpen)
			utils.SendSuccess(w, "Successfully Connect Database")
		}
	}
}
