package controllers

import (
	models "DynamicAPI/model"
	"DynamicAPI/repository"
	"DynamicAPI/utils"
	"database/sql"

	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

//AddValue 加入數值至資料表
//@Summary 加入數值至資料表
//@Tags Table CRUD
//@Accept json
//@Produce json
//@Param sql path string true "資料庫引擎"
//@Param tablename path string true "資料表名稱"
//@Param value query array false "插入數值"
//@Success 200 {string} string "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/{sql}/addvalue/{tablename} [post]
func (c Controller) AddValue() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			message       models.Error
			params        = mux.Vars(r)
			mysqldescribe = fmt.Sprintf("select COLUMN_NAME, DATA_TYPE, IS_NULLABLE,COLUMN_DEFAULT from INFORMATION_SCHEMA.COLUMNS where TABLE_NAME='%s' and TABLE_SCHEMA='%s'",
				params["tablename"], mysqlinformation.Database) //執行資料庫命令
			mssqldescribe = fmt.Sprintf(`select COLUMN_NAME, DATA_TYPE, IS_NULLABLE,COLUMN_DEFAULT from %s.INFORMATION_SCHEMA.COLUMNS where TABLE_NAME='%s'`,
				mssqlinformation.Database, params["tablename"])
			queryvalue = r.URL.Query()["value"]
			index      []string //欄位名稱
			value      []string
			Repo       repository.Repository
			err        error
			rows       *sql.Rows
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

			//選擇欄位
			if len(queryvalue) > 0 {
				for _, y := range queryvalue {
					new := strings.Split(y, ",")
					if new[0] == result.Field {
						index = append(index, new[0])
						value = append(value, new[1])
					}
				}
			}
		}
		//處理sql語法
		slicetostringIndex := strings.Join(index, ", ")
		slicetostringValue := strings.Join(value, `', ' `)
		switch strings.ToLower(params["sql"]) {
		case "mysql":
			insertvalue := fmt.Sprintf(`INSERT INTO %s(%s) VALUES('%s')`, params["tablename"],
				slicetostringIndex, slicetostringValue)
			err = Repo.Exec(MysqlDB, insertvalue)
		case "mssql":
			insertvalue := fmt.Sprintf(`INSERT INTO %s.dbo.%s(%s) VALUES('%s')`,
				mssqlinformation.Database, params["tablename"], slicetostringIndex, slicetostringValue)
			err = Repo.Exec(MssqlDB, insertvalue)
		}
		if err != nil {
			message.Message = "插入資料時發生錯誤"
			utils.SendError(w, http.StatusInternalServerError, message, err)
			return
		}
		utils.SendSuccess(w, "Successfully Add Value")
	}
}
