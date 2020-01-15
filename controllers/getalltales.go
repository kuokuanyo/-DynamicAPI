package controllers

import (
	"fmt"
	"net/http"
	"strings"

	models "DynamicAPI/model"
	"DynamicAPI/repository"
	"DynamicAPI/utils"

	"github.com/gorilla/mux"
)

//GetAlltables 取得所有資料表
//@Summary 取得所有資料表
//@Tags Table Information
//@Accept json
//@Produce json
//@Param sql path string true "資料庫引擎"
//@Success 200 {object} models.object "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/{sql}/getalltables [get]
func (c Controller) GetAlltables() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			params     = mux.Vars(r)
			tablenames []string
			message    models.Error
			Repo       repository.Repository
			err        error
		)
		switch strings.ToLower(params["sql"]) {
		case "mysql":
			//檢查資料庫是否連接
			if MysqlDB == nil {
				message.Message = ("資料庫未連接，請連接資料庫")
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}
			//查詢單行
			//func (s *DB) Pluck(column string, value interface{}) *DB
			if err = Repo.GetAlltables(MysqlDB, &tablenames); err != nil {
				message.Message = "取得資料表時發生錯誤"
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}
		case "mssql":
			execute := fmt.Sprintf("SELECT TABLE_NAME FROM %s.INFORMATION_SCHEMA.TABLES", mssqlinformation.Database)
			//檢查資料庫是否連接
			if MssqlDB == nil {
				message.Message = ("資料庫未連接，請連接資料庫")
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}
			rows, err := Repo.RawData(MssqlDB, execute)
			if err != nil {
				message.Message = "取得資料表時發生錯誤"
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}
			for rows.Next() {
				var table string
				rows.Scan(&table)
				tablenames = append(tablenames, table)
			}
		}
		utils.SendSuccess(w, tablenames)
	}
}
