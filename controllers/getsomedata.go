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

//GetSomeData 取得部分資料
//@Summary 取得部分資料
//@Tags Table CRUD
//@Accept json
//@Produce json
//@Param tablename path string true "資料表名稱"
//@Param col query array false "挑選欄位"
//@Param where query array false "選擇條件"
//@Success 200 {object} models.object "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/getsome/{tablename} [get]
func (c Controller) GetSomeData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//設定變數
		var (
			datas          []map[string]interface{} //資料放置
			message        models.Error
			params         = mux.Vars(r)
			describetable  = fmt.Sprintf("DESCRIBE %s", params["tablename"]) //執行資料庫命令
			condition      = make(map[string]interface{})
			index          []string //欄位名稱
			coltype        []string //欄位類型
			queryvalue     = r.URL.Query()["col"]
			conditionvalue = r.URL.Query()["where"]
			Repo           repository.Repository
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
				for x := 0; x < len(queryvalue); x++ {
					if queryvalue[x] == result.Field {
						index = append(index, result.Field)    //新增欄位名稱
						coltype = append(coltype, result.Type) //新增欄位類型
					}
				}
			} else if len(queryvalue) == 0 {
				index = append(index, result.Field)    //新增欄位名稱
				coltype = append(coltype, result.Type) //新增欄位類型
			}

		}
		//查詢條件
		if len(conditionvalue) > 0 {
			for _, y := range conditionvalue {
				new := strings.Split(y, ",")
				condition[new[0]] = new[1]
			}
		}

		//設定變數
		var (
			value     = make([]string, len(index))
			valuePtrs = make([]interface{}, len(index))
		)
		for i := range index {
			valuePtrs[i] = &value[i] //因Scan需要使用指標(valuePtrs)
		}

		//處理sql命令
		slicetostringIndex := strings.Join(index, ", ")
		getdata := fmt.Sprintf("select %s from %s ", slicetostringIndex, params["tablename"])

		if len(condition) > 0 {
			i := 0
			for x, y := range condition {
				if i == 0 {
					getdata += fmt.Sprintf(`WHERE %s="%s" `, x, y)
					i++
				} else {
					getdata += fmt.Sprintf(`AND %s="%s"`, x, y)
				}
			}
		}
		
		//取得資料
		rows, err = Repo.RawData(DB, getdata)
		if err != nil {
			message.Message = "取資料時發生錯誤"
			utils.SendError(w, http.StatusInternalServerError, message)
			return
		}
		for rows.Next() {
			var data = make(map[string]interface{})
			err = rows.Scan(valuePtrs...)
			if err != nil {
				message.Message = "Scan資料時發生錯誤"
				utils.SendError(w, http.StatusInternalServerError, message)
				return
			}
			for i := range index {
				//辨別資料庫欄位類別
				if strings.Contains(coltype[i], "varchar") {
					data[index[i]] = value[i] //欄位型態為string
				} else if strings.Contains(coltype[i], "int") {
					data[index[i]], err = strconv.Atoi(value[i]) //欄位型態為int
					if err != nil {
						message.Message = "數值轉換錯誤"
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
