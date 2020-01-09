package controllers

import (
	models "DynamicAPI/model"
	"DynamicAPI/repository"
	"DynamicAPI/utils"

	"fmt"
	"net/http"
	"strconv"
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
			index         []string //欄位名稱
			coltype       []string //欄位類型
			getdata       = "select "
			params        = mux.Vars(r)                                     //印出url參數
			describetable = fmt.Sprintf("DESCRIBE %s", params["tablename"]) //取得資料表資訊的指令
			Repo          repository.Repository
		)

		//檢查資料庫是否連接
		if DB == nil {
			message.Message = ("資料庫未連接，請連接資料庫")
			utils.SendError(w, http.StatusInternalServerError, message)
			return
		}

		//取得資料表資訊
		rows, err := Repo.RawData(DB, describetable)
		if err != nil {
			message.Message = "取得資料表資訊時發生錯誤"
			utils.SendError(w, http.StatusInternalServerError, message)
			return
		}
		//取得欄位的資訊
		for rows.Next() {
			var result models.Result
			err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
			if err != nil {
				message.Message = "Scan時發生錯誤"
				utils.SendError(w, http.StatusInternalServerError, message)
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
			valuePtrs = make([]interface{}, len(index))
		)

		for x := range index {
			if x == len(index)-1 {
				getdata += fmt.Sprintf("%s from %s", index[x], params["tablename"])
				valuePtrs[x] = &value[x] //因Scan需要使用指標(valuePtrs)
			} else {
				getdata += fmt.Sprintf("%s, ", index[x])
				valuePtrs[x] = &value[x] //因Scan需要使用指標(valuePtrs)
			}
		}

		//取得所有資料
		rows, err = Repo.RawData(DB, getdata)
		if err != nil {
			message.Message = "取得資料時發生錯誤"
			fmt.Println(err)
			utils.SendError(w, http.StatusInternalServerError, message)
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
						utils.SendError(w, http.StatusInternalServerError, message)
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
