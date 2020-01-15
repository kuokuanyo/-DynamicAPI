package controllers

import (
	models "DynamicAPI/model"
	"DynamicAPI/utils"
	"net/http"

	"github.com/gorilla/mux"
)

//GetTableByUUID 利用uuid取得表
//@Summary 利用專屬uuid呼叫合併表
//@Tags Table JoinTable
//@Accept json
//@Produce json
//@Param uuid path string true "專屬uuid"
//@Success 200 {object} models.object "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/jointable/{uuid} [get]
func (c Controller) GetTableByUUID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			params  = mux.Vars(r)
			uuid    = params["uuid"] //取得url的uuid
			datas   = Identity[uuid]
			message models.Error
			err     error
		)
		//檢查資料庫是否連接
		if MysqlDB == nil {
			message.Message = "資料庫未連接，請連接資料庫"
			utils.SendError(w, http.StatusInternalServerError, message, err)
			return
		}
		if Identity[uuid] == nil {
			message.Message = "無此uuid，請建立關聯表，取得uuid"
			utils.SendError(w, http.StatusInternalServerError, message, err)
			return
		}
		utils.SendSuccess(w, datas)
	}
}
