package controllers

import (
	models "DynamicAPI/model"
	"DynamicAPI/repository"
	"DynamicAPI/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

//UpdateValue 更新資料表數值
//@Summary 更新資料表數值
//@Tags Table CRUD
//@Accept json
//@Produce json
//@Param sql path string true "資料庫引擎"
//@Param tablename path string true "資料表名稱"
//@Param set query array false "更新的欄位數值"
//@Param where query array false "被更新的欄位條件"
//@Success 200 {string} string "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/{sql}/update/{tablename} [put]
func (c Controller) UpdateValue() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			message       models.Error
			params        = mux.Vars(r)
			mysqldescribe = fmt.Sprintf("select COLUMN_NAME, DATA_TYPE, IS_NULLABLE,COLUMN_DEFAULT from INFORMATION_SCHEMA.COLUMNS where TABLE_NAME='%s' and TABLE_SCHEMA='%s'",
				params["tablename"], mysqlinformation.Database) //執行資料庫命令
			mssqldescribe = fmt.Sprintf(`select COLUMN_NAME, DATA_TYPE, IS_NULLABLE,COLUMN_DEFAULT from %s.INFORMATION_SCHEMA.COLUMNS where TABLE_NAME='%s'`,
				mssqlinformation.Database, params["tablename"])
			mysqlupdatevalue = fmt.Sprintf("UPDATE %s SET ", params["tablename"])
			mssqlupdatevalue = fmt.Sprintf("UPDATE %s.dbo.%s SET ", mssqlinformation.Database, params["tablename"])
			set              = r.URL.Query()["set"]
			where            = r.URL.Query()["where"]
			setindex         []string
			setvalue         []interface{}
			whereindex       []string
			wherevalue       []interface{}
			Repo             repository.Repository
			err              error
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
			rows, err := Repo.RawData(MysqlDB, mysqldescribe)
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
				if len(set) > 0 {
					for _, y := range set {
						new := strings.Split(y, ",")
						if new[0] == result.Field {
							setindex = append(setindex, new[0])
							setvalue = append(setvalue, new[1])
						}
					}
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

			//處理sql語法
			for i := 0; i < len(setindex); i++ {
				if i == len(setindex)-1 {
					mysqlupdatevalue += fmt.Sprintf(`%s="%s" WHERE `, setindex[i], setvalue[i])
				} else {
					mysqlupdatevalue += fmt.Sprintf(`%s="%s" AND `, setindex[i], setvalue[i])
				}
			}
			for i := 0; i < len(whereindex); i++ {
				if i == len(whereindex)-1 {
					mysqlupdatevalue += fmt.Sprintf("%s=%s", whereindex[i], wherevalue[i])
				} else {
					mysqlupdatevalue += fmt.Sprintf("%s=%s AND ", whereindex[i], wherevalue[i])
				}
			}

			//執行更新數值命令
			if err = Repo.Exec(MysqlDB, mysqlupdatevalue); err != nil {
				message.Message = "資料更新時發生錯誤"
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}

		case "mssql":
			//檢查資料庫是否連接
			if MssqlDB == nil {
				message.Message = ("資料庫未連接，請連接資料庫")
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}

			//取得資料表欄位資訊
			rows, err := Repo.RawData(MssqlDB, mssqldescribe)
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
				if len(set) > 0 {
					for _, y := range set {
						new := strings.Split(y, ",")
						if new[0] == result.Field {
							setindex = append(setindex, new[0])
							setvalue = append(setvalue, new[1])
						}
					}
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

			//處理sql語法
			for i := 0; i < len(setindex); i++ {
				if i == len(setindex)-1 {
					mssqlupdatevalue += fmt.Sprintf(`%s='%s' WHERE `, setindex[i], setvalue[i])
				} else {
					mssqlupdatevalue += fmt.Sprintf(`%s='%s', `, setindex[i], setvalue[i])
				}
			}
			for i := 0; i < len(whereindex); i++ {
				if i == len(whereindex)-1 {
					mssqlupdatevalue += fmt.Sprintf(`%s='%s'`, whereindex[i], wherevalue[i])
				} else {
					mssqlupdatevalue += fmt.Sprintf(`%s='%s' AND `, whereindex[i], wherevalue[i])
				}
			}

			//執行更新數值命令
			if err = Repo.Exec(MssqlDB, mssqlupdatevalue); err != nil {
				message.Message = "資料更新時發生錯誤"
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}
		}
		utils.SendSuccess(w, "Successfully Update Value")
	}
}
