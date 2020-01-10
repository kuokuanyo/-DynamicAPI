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

//AddValue 加入數值至資料表
//@Summary 加入數值至資料表
//@Tags Table CRUD
//@Accept json
//@Produce json
//@Param tablename path string true "資料表名稱"
//@Param value query array false "插入數值"
//@Success 200 {string} string "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/addvalue/{tablename} [post]
func (c Controller) AddValue() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			message       models.Error
			params        = mux.Vars(r)
			describetable = fmt.Sprintf("DESCRIBE %s", params["tablename"]) //執行資料庫命令
			queryvalue    = r.URL.Query()["value"]
			index         []string //欄位名稱
			value         []string
			Repo          repository.Repository
		)

		//檢查資料庫是否連接
		if DB == nil {
			message.Message = ("資料庫未連接，請連接資料庫")
			utils.SendError(w, http.StatusInternalServerError, message)
			return
		}

		//取得資料表欄位資訊
		rows, err := Repo.RawData(DB, describetable)
		if err != nil {
			message.Message = "取得欄位資訊時發生錯誤"
			utils.SendError(w, http.StatusInternalServerError, message)
			return
		}
		for rows.Next() {
			var result models.Result
			err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
			if err != nil {
				message.Message = "Scan資料時發生錯誤"
				utils.SendError(w, http.StatusInternalServerError, message)
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
		//加入欄位
		slicetostringIndex := strings.Join(index, ", ")
		slicetostringValue := strings.Join(value, `"," `)

		insertvalue := fmt.Sprintf(`INSERT INTO %s(%s) VALUES("%s")`, params["tablename"], slicetostringIndex, slicetostringValue)
		fmt.Println(insertvalue)

		if err = Repo.Exec(DB, insertvalue); err != nil {
			message.Message = "插入資料時發生錯誤"
			utils.SendError(w, http.StatusInternalServerError, message)
			return
		}
		utils.SendSuccess(w, "Successfully Add Value")

	}
}
