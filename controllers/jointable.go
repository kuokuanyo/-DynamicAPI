package controllers

import (
	models "DynamicAPI/model"
	"DynamicAPI/repository"
	"DynamicAPI/utils"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

//JoinTable 合併表
func (c Controller) JoinTable() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			//datas          []map[string]interface{} //資料放置
			params         = mux.Vars(r)
			describetable1 = fmt.Sprintf("DESCRIBE %s", params["table1"]) //取得資料表資訊的指令
			describetable2 = fmt.Sprintf("DESCRIBE %s", params["table2"]) //取得資料表資訊的指令
			//join           = r.URL.Query()["join"]
			table1Col = r.URL.Query()["table1"]
			table2Col = r.URL.Query()["table2"]
			message   models.Error
			//getdata        = "select "
			col     []string //挑選欄位
			colType []string
			//join    = r.URL.Query()["join"] //合併欄位
			Repo repository.Repository
		)

		//檢查資料庫是否連接
		if DB == nil {
			message.Message = ("資料庫未連接，請連接資料庫")
			utils.SendError(w, http.StatusInternalServerError, message)
			return
		}

		//處理table1資料表
		rows, err := Repo.RawData(DB, describetable1)
		if err != nil {
			message.Message = "取得資料表資訊時發生錯誤"
			utils.SendError(w, http.StatusInternalServerError, message)
			return
		}
		for rows.Next() {
			var result models.Result
			err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
			if err != nil {
				message.Message = "Scan時發生錯誤"
				utils.SendError(w, http.StatusInternalServerError, message)
				return
			}

			if len(table1Col) > 0 {
				for x := range table1Col {
					if table1Col[x] == result.Field {
						col = append(col, result.Field)
						colType = append(colType, result.Type)
					}
				}
			} else if len(table1Col) == 0 {
				col = append(col, result.Field)
				colType = append(colType, result.Type)
			}
		}

		//處理table2資料表
		rows, err = Repo.RawData(DB, describetable2)
		if err != nil {
			message.Message = "取得資料表資訊時發生錯誤"
			utils.SendError(w, http.StatusInternalServerError, message)
			return
		}
		for rows.Next() {
			var result models.Result
			err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
			if err != nil {
				message.Message = "Scan時發生錯誤"
				utils.SendError(w, http.StatusInternalServerError, message)
				return
			}

			if len(table2Col) > 0 {
				for x := range table2Col {
					if table2Col[x] == result.Field {
						col = append(col, result.Field)
						colType = append(col, result.Type)
					}
				}
			} else if len(table2Col) == 0 {
				col = append(col, result.Field)
				colType = append(colType, result.Type)
			}
		}

		fmt.Println(col)
		//fmt.Println(colType)
	}
}
