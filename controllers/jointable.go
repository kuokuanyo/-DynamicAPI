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
	uuid "github.com/satori/go.uuid"
)

//Identity 利用uuid建立新的api
var Identity = make(map[string][]map[string]interface{})

//JoinTable 合併表
//@Summary 合併資料表
//@Tags Table JoinTable
//@Accept json
//@Produce json
//@Param table1 path string true "合併資料表"
//@Param table2 path string true "合併資料表"
//@Param join query array true "合併條件"
//@Param table1 query array false "挑選欄位"
//@Param table2 query array false "挑選欄位"
//@Success 200 {object} models.object "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/getall/{tablename} [get]
func (c Controller) JoinTable() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			datas          []map[string]interface{} //資料放置
			params         = mux.Vars(r)
			describetable1 = fmt.Sprintf("DESCRIBE %s", params["table1"]) //取得資料表資訊的指令
			describetable2 = fmt.Sprintf("DESCRIBE %s", params["table2"]) //取得資料表資訊的指令
			table1Col      = r.URL.Query()["table1"]
			table2Col      = r.URL.Query()["table2"]
			message        models.Error
			col            []string //挑選欄位
			colType        []string
			join           = r.URL.Query()["join"] //合併欄位
			Repo           repository.Repository
			err            error
		)

		//檢查資料庫是否連接
		if DB == nil {
			message.Message = ("資料庫未連接，請連接資料庫")
			utils.SendError(w, http.StatusInternalServerError, message, err)
			return
		}

		//處理table1資料表
		rows, err := Repo.RawData(DB, describetable1)
		if err != nil {
			message.Message = "取得資料表資訊時發生錯誤"
			utils.SendError(w, http.StatusInternalServerError, message, err)
			return
		}
		for rows.Next() {
			var result models.MysqlResult
			err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
			if err != nil {
				message.Message = "Scan時發生錯誤"
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}

			if len(table1Col) > 0 {
				for x := range table1Col {
					if table1Col[x] == result.Field {
						result.Field = params["table1"] + "." + result.Field
						col = append(col, result.Field)
						colType = append(colType, result.Type)
					}
				}
			} else if len(table1Col) == 0 {
				result.Field = params["table1"] + "." + result.Field
				col = append(col, result.Field)
				colType = append(colType, result.Type)
			}
		}

		//處理table2資料表
		rows, err = Repo.RawData(DB, describetable2)
		if err != nil {
			message.Message = "取得資料表資訊時發生錯誤"
			utils.SendError(w, http.StatusInternalServerError, message, err)
			return
		}
		for rows.Next() {
			var result models.MysqlResult
			err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
			if err != nil {
				message.Message = "Scan時發生錯誤"
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}

			//匿名函式(檢查是否有重複的欄位)
			firstcheck := duplicate(col, result.Field)
			//第二道檢查(檢查是否出現在join欄位，可能出現有一樣的欄位卻不是join條件，必須改欄位名)
			secondcheck := duplicate(join, result.Field)

			if firstcheck {
				if len(table2Col) > 0 {
					for x := range table2Col {
						if table2Col[x] == result.Field {
							result.Field = params["table2"] + "." + result.Field
							col = append(col, result.Field)
							colType = append(colType, result.Type)
						}
					}
				} else if len(table2Col) == 0 {
					result.Field = params["table2"] + "." + result.Field
					col = append(col, result.Field)
					colType = append(colType, result.Type)
				}
			} else {
				if secondcheck {
					for x := range table2Col {
						if table2Col[x] == result.Field {
							result.Field = params["table2"] + "." + result.Field
							col = append(col, result.Field)
							colType = append(colType, result.Type)
						}
					}
				}
			}
		}

		//設定變數
		var (
			value     = make([]string, len(col))
			valuePtrs = make([]interface{}, len(col)) //scan時必須使用指針
		)
		for i := 0; i < len(col); i++ {
			valuePtrs[i] = &value[i] //因Scan需要使用指標(valuePtrs)
		}

		//處理sql命令
		//slice convert string
		sliceTostringCol := strings.Join(col, ", ")
		sliceTostringJoin := strings.Join(join, ", ")

		getdata := fmt.Sprintf("select %s from %s join %s using (%s)",
			sliceTostringCol, params["table1"], params["table2"], sliceTostringJoin)

		//執行命令
		rows, err = Repo.RawData(DB, getdata)
		if err != nil {
			message.Message = "合併資料表時發生錯誤"
			fmt.Println(err)
			utils.SendError(w, http.StatusInternalServerError, message, err)
			return
		}
		for rows.Next() {
			data := make(map[string]interface{})
			rows.Scan(valuePtrs...)
			for i := range col {
				if strings.Contains(colType[i], "varchar") { //欄位為string
					data[col[i]] = value[i]
				} else if strings.Contains(colType[i], "int") {
					data[col[i]], err = strconv.Atoi(value[i]) //欄位為int
					if err != nil {
						message.Message = "數值轉換時發生錯誤"
						utils.SendError(w, http.StatusInternalServerError, message, err)
						return
					}
				} else {
					data[col[i]] = value[i]
				}
			}
			datas = append(datas, data)
		}

		//建立一組驗證的uuid
		uuid := uuid.Must(uuid.NewV4())
		uuidtostring := uuid.String() //convert to string
		Identity[uuidtostring] = datas

		getuuid := fmt.Sprintf("New uuid is %s", uuidtostring)
		utils.SendSuccess(w, getuuid)
		utils.SendSuccess(w, datas)
	}
}

//duplicate 檢查是否有重複字串
func duplicate(col []string, s string) bool {
	b := true
	for x := range col {
		if strings.Contains(col[x], s) {
			b = false
			break
		} else {
			continue
		}
	}
	return b
}
