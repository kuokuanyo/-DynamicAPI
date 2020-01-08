package controllers

import (
	models "DynamicAPI/model"
	"DynamicAPI/repository"
	"DynamicAPI/utils"

	"net/http"
)

//GetAlltables 取得所有資料表
//@Summary 取得所有資料表
//@Tags Table Information
//@Accept json
//@Produce json
//@Success 200 {object} models.object "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/getalltables [get]
func (c Controller) GetAlltables() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			tablenames []string
			message    models.Error
			Repo       repository.Repository
		)

		//檢查資料庫是否連接
		if DB == nil {
			message.Message = ("資料庫未連接，請連接資料庫")
			utils.SendError(w, http.StatusInternalServerError, message)
			return
		}

		//查詢單行
		//func (s *DB) Pluck(column string, value interface{}) *DB
		if err := Repo.GetAlltables(DB, &tablenames); err != nil {
			message.Message = "取得資料表時發生錯誤"
			utils.SendError(w, http.StatusInternalServerError, message)
			return
		}
		utils.SendSuccess(w, tablenames)
	}
}
