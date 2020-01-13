package controllers

import (
	models "DynamicAPI/model"
	"DynamicAPI/repository"
	"DynamicAPI/utils"
	"strconv"

	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

//GetAllData 取得所有資料
//@Summary 取得所有資料
//@Tags Table CRUD
//@Accept json
//@Produce json
//@Param tablename path string true "資料表名稱"
//@Param col query array false "挑選欄位"
//@Success 200 {object} models.object "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/getall/{tablename} [get]
func (c Controller) GetAllData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			datas         []map[string]interface{} //資料放置
			message       models.Error
			index         []string                                          //欄位名稱
			coltype       []string                                          //欄位類型
			params        = mux.Vars(r)                                     //印出url參數
			describetable = fmt.Sprintf("DESCRIBE %s", params["tablename"]) //取得資料表資訊的指令
			Repo          repository.Repository
			err           error
		)

		//檢查資料庫是否連接
		if DB == nil {
			message.Message = ("資料庫未連接，請連接資料庫")
			utils.SendError(w, http.StatusInternalServerError, message, err)
			return
		}

		//取得資料表資訊
		rows, err := Repo.RawData(DB, describetable)
		if err != nil {
			message.Message = "取得資料表資訊時發生錯誤"
			utils.SendError(w, http.StatusInternalServerError, message, err)
			return
		}
		//取得欄位的資訊
		for rows.Next() {
			var result models.MysqlResult
			err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
			if err != nil {
				message.Message = "Scan時發生錯誤"
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}
			queryvalue := r.URL.Query()["col"]
			if len(queryvalue) > 0 {
				for x := 0; x < len(queryvalue); x++ {
					if queryvalue[x] == result.Field {
						index = append(index, result.Field)
						coltype = append(coltype, result.Type)
					}
				}
			} else if len(queryvalue) == 0 {
				index = append(index, result.Field)
				coltype = append(coltype, result.Type)
			}
		}

		//設定變數
		var (
			value     = make([]string, len(index))
			valuePtrs = make([]interface{}, len(index)) //因Scan需要使用指標(valuePtrs)
		)

		for i := 0; i < len(index); i++ {
			valuePtrs[i] = &value[i] //因Scan需要使用指標(valuePtrs)
		}

		//處理sql命令
		slicetostringIndex := strings.Join(index, ", ")
		getdata := fmt.Sprintf("select %s from %s", slicetostringIndex, params["tablename"])

		//取得所有資料
		rows, err = Repo.RawData(DB, getdata)
		if err != nil {
			message.Message = "取得資料時發生錯誤"
			utils.SendError(w, http.StatusInternalServerError, message, err)
			return
		}
		for rows.Next() {
			var data = make(map[string]interface{})
			rows.Scan(valuePtrs...)
			for i := range index {
				//辨別資料庫欄位類別
				if strings.Contains(coltype[i], "varchar") {
					data[index[i]] = value[i] //欄位型態為string
				} else if strings.Contains(coltype[i], "int") {
					data[index[i]], err = strconv.Atoi(value[i]) //欄位型態為int
					if err != nil {
						message.Message = "數值轉換時發生錯誤"
						utils.SendError(w, http.StatusInternalServerError, message, err)
						return
					}
				} else {
					data[index[i]] = value[i]
				}
			}
			datas = append(datas, data)
		}
		utils.SendSuccess(w, datas)
	}
}
