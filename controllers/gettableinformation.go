package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	models "DynamicAPI/model"
	"DynamicAPI/repository"
	"DynamicAPI/utils"

	"github.com/gorilla/mux"
)

//GetTableInformation 取得資料表資訊
//@Summary 取得資料表資訊
//@Tags Table Information
//@Accept json
//@Produce json
//@Param sql path string true "資料庫引擎"
//@Param tablename path string true "資料庫名稱"
//@Success 200 {object} models.object "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/{sql}/tableinformation/{tablename} [get]
func (c Controller) GetTableInformation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			message models.Error
			results []models.Result
			params  = mux.Vars(r) //印出url參數
			Repo    repository.Repository
			err     error
			rows    *sql.Rows
		)

		switch strings.ToLower(params["sql"]) {
		case "mysql":
			//檢查資料庫是否連接
			if MysqlDB == nil {
				message.Message = ("資料庫未連接，請連接資料庫")
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}
			//指令
			mysqldescribe := fmt.Sprintf("select COLUMN_NAME, DATA_TYPE, IS_NULLABLE,COLUMN_DEFAULT from INFORMATION_SCHEMA.COLUMNS where TABLE_NAME='%s' and TABLE_SCHEMA='%s'",
				params["tablename"], mysqlinformation.Database) //執行資料庫命令

			rows, err = Repo.RawData(MysqlDB, mysqldescribe)
		case "mssql":
			//檢查資料庫是否連接
			if MssqlDB == nil {
				message.Message = "資料庫未連接，請連接資料庫"
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}
			//指令
			execute := fmt.Sprintf(`select COLUMN_NAME, DATA_TYPE, IS_NULLABLE,COLUMN_DEFAULT from %s.INFORMATION_SCHEMA.COLUMNS where TABLE_NAME='%s'`, mssqlinformation.Database, params["tablename"])
			rows, err = Repo.RawData(MssqlDB, execute)
		}
		if err != nil {
			message.Message = "取得資料表資訊時發生錯誤"
			utils.SendError(w, http.StatusInternalServerError, message, err)
			return
		}
		for rows.Next() {
			var result models.Result
			rows.Scan(&result.Field, &result.Type, &result.Null, &result.Default)
			results = append(results, result)
		}
		utils.SendSuccess(w, results)
	}
}
