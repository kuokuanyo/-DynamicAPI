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

//DeleteValue 刪除資料表數值
//@Summary 刪除資料表數值
//@Tags Table CRUD
//@Accept json
//@Produce json
//@Param tablename path string true "資料表名稱"
//@Param where query array false "被刪除的欄位條件"
//@Success 200 {string} string "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/delete/{tablename} [delete]
func (c Controller) DeleteValue() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			message       models.Error
			params        = mux.Vars(r)
			describetable = fmt.Sprintf("DESCRIBE %s", params["tablename"]) //執行資料庫命令
			deletevalue   = fmt.Sprintf("DELETE FROM %s WHERE ", params["tablename"])
			where         = r.URL.Query()["where"]
			whereindex    []string
			wherevalue    []interface{}
			Repo          repository.Repository
			err           error
		)

		//檢查資料庫是否連接
		if DB == nil {
			message.Message = ("資料庫未連接，請連接資料庫")
			utils.SendError(w, http.StatusInternalServerError, message, err)
			return
		}

		//取得資料表欄位資訊
		rows, err := Repo.RawData(DB, describetable)
		if err != nil {
			message.Message = "取得欄位資訊時發生錯誤"
			utils.SendError(w, http.StatusInternalServerError, message, err)
			return
		}
		for rows.Next() {
			var result models.MysqlResult
			err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
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

		//處理sql語法
		for i := 0; i < len(whereindex); i++ {
			if i == len(whereindex)-1 {
				deletevalue += fmt.Sprintf(`%s="%s"`, whereindex[i], wherevalue[i])
			} else {
				deletevalue += fmt.Sprintf(`%s="%s" AND `, whereindex[i], wherevalue[i])
			}
		}

		//執行刪除數值命令
		if err = Repo.Exec(DB, deletevalue); err != nil {
			message.Message = "刪除資料時發生錯誤"
			utils.SendError(w, http.StatusInternalServerError, message, err)
			return
		}
		utils.SendSuccess(w, "Successfully Delete Value")
	}
}
