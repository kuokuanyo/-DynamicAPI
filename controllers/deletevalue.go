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

//DeleteValue 刪除資料表數值
//@Summary 刪除資料表數值
//@Tags Table CRUD
//@Accept json
//@Produce json
//@Param sql path string true "資料庫引擎"
//@Param tablename path string true "資料表名稱"
//@Param where query array false "被刪除的欄位條件"
//@Success 200 {string} string "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/{sql}/table/{tablename} [delete]
func (c Controller) DeleteValue() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			message       models.Error
			params        = mux.Vars(r)
			mysqldescribe = fmt.Sprintf("select COLUMN_NAME, DATA_TYPE, IS_NULLABLE,COLUMN_DEFAULT from INFORMATION_SCHEMA.COLUMNS where TABLE_NAME='%s' and TABLE_SCHEMA='%s'",
				params["tablename"], mysqlinformation.Database) //執行資料庫命令
			mssqldescribe = fmt.Sprintf(`select COLUMN_NAME, DATA_TYPE, IS_NULLABLE,COLUMN_DEFAULT from %s.INFORMATION_SCHEMA.COLUMNS where TABLE_NAME='%s'`,
				mssqlinformation.Database, params["tablename"])
			where       = r.URL.Query()["where"]
			whereindex  []string
			wherevalue  []interface{}
			Repo        repository.Repository
			err         error
			rows        *sql.Rows
			deletevalue string
		)

		switch strings.ToLower(params["sql"]) {
		case "mysql":
			//檢查資料庫是否連接
			if MysqlDB == nil {
				message.Message = ("資料庫未連接，請連接資料庫")
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}
			//取得資料表欄位資訊
			rows, err = Repo.RawData(MysqlDB, mysqldescribe)
		case "mssql":
			//檢查資料庫是否連接
			if MssqlDB == nil {
				message.Message = ("資料庫未連接，請連接資料庫")
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}
			//取得資料表欄位資訊
			rows, err = Repo.RawData(MssqlDB, mssqldescribe)
		}
		if err != nil {
			message.Message = "取得欄位資訊時發生錯誤"
			utils.SendError(w, http.StatusInternalServerError, message, err)
			return
		}
		for rows.Next() {
			var result models.Result
			err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Default)
			if err != nil {
				message.Message = "Scan資料時發生錯誤"
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}

			if len(where) > 0 {
				for _, y := range where {
					new := strings.Split(y, ",")
					if new[0] == result.Field {
						whereindex = append(whereindex, new[0])
						wherevalue = append(wherevalue, new[1])
					}
				}
			}
		}
		switch strings.ToLower(params["sql"]) {
		case "mysql":
			deletevalue = fmt.Sprintf("DELETE FROM %s WHERE ", params["tablename"])
		case "mssql":
			deletevalue = fmt.Sprintf("DELETE FROM %s.dbo.%s WHERE ", mssqlinformation.Database, params["tablename"])
		}
		//處理sql語法
		for i := 0; i < len(whereindex); i++ {
			if i == len(whereindex)-1 {
				deletevalue += fmt.Sprintf(`%s='%s'`, whereindex[i], wherevalue[i])
			} else {
				deletevalue += fmt.Sprintf(`%s='%s' AND `, whereindex[i], wherevalue[i])
			}
		}
		switch strings.ToLower(params["sql"]) {
		case "mysql":
			err = Repo.Exec(MysqlDB, deletevalue)
		case "mssql":
			err = Repo.Exec(MssqlDB, deletevalue)
		}
		if err != nil {
			message.Message = "刪除資料時發生錯誤"
			utils.SendError(w, http.StatusInternalServerError, message, err)
			return
		}
		utils.SendSuccess(w, "Successfully Delete Value")
	}
}
