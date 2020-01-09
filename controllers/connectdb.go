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

//DB 資料庫引擎
var DB *gorm.DB

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
			information models.DBinformation
			message     models.Error
			params      = mux.Vars(r) //印出url參數
			err         error
			Repo        repository.Repository
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
			DB, err = Repo.ConnectDb("mysql", MysqlDataSourceName)
			if err != nil {
				message.Message = "Connect Database failed"
				utils.SendError(w, http.StatusInternalServerError, message)
				return
			}

			utils.SendSuccess(w, "Successfully Connect Database")

		case "mssql":
			MssqlDataSourceName := fmt.Sprintf("sqlserver://%s:%s@%s:%s? database=%s",
				information.UserName,
				information.Password,
				information.Host,
				information.Port,
				information.Database)

			DB, err = gorm.Open("mssql", MssqlDataSourceName)
			if err != nil {
				message.Message = "Connect Database failed"
				utils.SendError(w, http.StatusInternalServerError, message)
				return
			}

			DB.DB().SetMaxIdleConns(information.MaxIdle)
			DB.DB().SetMaxOpenConns(information.MaxOpen)

			utils.SendSuccess(w, "Successfully Connect Database")

		}
	}
}
