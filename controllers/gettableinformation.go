package controllers

import (
	models "DynamicAPI/model"
	"DynamicAPI/repository"
	"DynamicAPI/utils"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

//GetTableInformation 取得資料表資訊
//@Summary 取得資料表資訊
//@Tags Table Information
//@Accept json
//@Produce json
//@Param tablename path string true "資料庫名稱"
//@Success 200 {object} models.object "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/tableinformation/{tablename} [get]
func (c Controller) GetTableInformation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			message models.Error
			results []models.Result
			params  = mux.Vars(r) //印出url參數
			Repo    repository.Repository
		)

		//檢查資料庫是否連接
		if DB == nil {
			message.Message = ("資料庫未連接，請連接資料庫")
			utils.SendError(w, http.StatusInternalServerError, message)
			return
		}

		//指令
		describetable := fmt.Sprintf("DESCRIBE %s", params["tablename"])

		rows, err := Repo.RawData(DB, describetable)
		if err != nil {
			message.Message = "取得資料表資訊時發生錯誤"
			utils.SendError(w, http.StatusInternalServerError, message)
			return
		}
		for rows.Next() {
			var result models.Result
			rows.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
			results = append(results, result)
		}
		utils.SendSuccess(w, results)
	}
}
