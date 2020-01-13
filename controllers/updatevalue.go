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
//@Param tablename path string true "資料表名稱"
//@Param set query array false "更新的欄位數值"
//@Param where query array false "被更新的欄位條件"
//@Success 200 {string} string "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/update/{tablename} [put]
func (c Controller) UpdateValue() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			message       models.Error
			params        = mux.Vars(r)
			describetable = fmt.Sprintf("DESCRIBE %s", params["tablename"]) //執行資料庫命令
			updatevalue   = fmt.Sprintf("UPDATE %s SET ", params["tablename"])
			set           = r.URL.Query()["set"]
			where         = r.URL.Query()["where"]
			setindex      []string
			setvalue      []interface{}
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
				updatevalue += fmt.Sprintf(`%s="%s" WHERE `, setindex[i], setvalue[i])
			} else {
				updatevalue += fmt.Sprintf(`%s="%s" AND `, setindex[i], setvalue[i])
			}
		}
		for i := 0; i < len(whereindex); i++ {
			if i == len(whereindex)-1 {
				updatevalue += fmt.Sprintf("%s=%s", whereindex[i], wherevalue[i])
			} else {
				updatevalue += fmt.Sprintf("%s=%s AND ", whereindex[i], wherevalue[i])
			}
		}

		//執行更新數值命令
		if err = Repo.Exec(DB, updatevalue); err != nil {
			message.Message = "資料更新時發生錯誤"
			utils.SendError(w, http.StatusInternalServerError, message, err)
			return
		}
		utils.SendSuccess(w, "Successfully Update Value")
	}
}
