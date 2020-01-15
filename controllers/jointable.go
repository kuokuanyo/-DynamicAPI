package controllers

import (
	models "DynamicAPI/model"
	"DynamicAPI/repository"
	"DynamicAPI/utils"
	"database/sql"
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
//@Param sql1 path string true "資料庫引擎"
//@Param table1 path string true "合併資料表"
//@Param sql2 path string true "資料庫引擎"
//@Param table2 path string true "合併資料表"
//@Param join query array true "合併條件"
//@Param table1 query array false "挑選欄位"
//@Param table2 query array false "挑選欄位"
//@Success 200 {object} models.object "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/jointable/{sql1}/{table1}/{sql2}/{table2} [get]
func (c Controller) JoinTable() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			datas     []map[string]interface{} //資料放置
			params    = mux.Vars(r)
			table1Col = r.URL.Query()["table1"]
			table2Col = r.URL.Query()["table2"]
			message   models.Error
			col       []string //挑選欄位
			colType   []string
			join      = r.URL.Query()["join"] //合併欄位
			Repo      repository.Repository
			err       error
		)

		//處理第一個資料表
		switch strings.ToLower(params["sql1"]) {
		case "mysql":
			//檢查資料庫是否連接
			if MysqlDB == nil {
				message.Message = ("資料庫未連接，請連接資料庫")
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}

			describetable1 := fmt.Sprintf("select COLUMN_NAME, DATA_TYPE, IS_NULLABLE,COLUMN_DEFAULT from INFORMATION_SCHEMA.COLUMNS where TABLE_NAME='%s' and TABLE_SCHEMA='%s'",
				params["table1"], mysqlinformation.Database) //執行資料庫命令
			//處理table1資料表
			rows, err := Repo.RawData(MysqlDB, describetable1)
			if err != nil {
				message.Message = "取得資料表資訊時發生錯誤"
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}
			for rows.Next() {
				var result models.Result
				err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Default)
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
		case "mssql":
			//檢查資料庫是否連接
			if MssqlDB == nil {
				message.Message = ("資料庫未連接，請連接資料庫")
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}

			describetable1 := fmt.Sprintf(`select COLUMN_NAME, DATA_TYPE, IS_NULLABLE,COLUMN_DEFAULT from %s.INFORMATION_SCHEMA.COLUMNS where TABLE_NAME='%s'`,
				mssqlinformation.Database, params["table1"]) //執行資料庫命令
			//處理table1資料表
			rows, err := Repo.RawData(MssqlDB, describetable1)
			if err != nil {
				message.Message = "取得資料表資訊時發生錯誤"
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}
			for rows.Next() {
				var result models.Result
				err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Default)
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
		}

		//處理第二個資料表
		switch strings.ToLower(params["sql2"]) {
		case "mysql":
			if MysqlDB == nil {
				message.Message = ("資料庫未連接，請連接資料庫")
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}

			describetable2 := fmt.Sprintf("select COLUMN_NAME, DATA_TYPE, IS_NULLABLE,COLUMN_DEFAULT from INFORMATION_SCHEMA.COLUMNS where TABLE_NAME='%s' and TABLE_SCHEMA='%s'",
				params["table2"], mysqlinformation.Database)
			rows, err := Repo.RawData(MysqlDB, describetable2)
			if err != nil {
				message.Message = "取得資料表資訊時發生錯誤"
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}
			for rows.Next() {
				var result models.Result
				err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Default)
				if err != nil {
					message.Message = "Scan時發生錯誤"
					utils.SendError(w, http.StatusInternalServerError, message, err)
					return
				}

				//匿名函式(檢查是否有重複的欄位)
				firstcheck := utils.Duplicate(col, result.Field)
				//第二道檢查(檢查是否出現在join欄位，可能出現有一樣的欄位卻不是join條件，必須改欄位名)
				secondcheck := utils.Duplicate(join, result.Field)

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
		case "mssql":
			if MssqlDB == nil {
				message.Message = ("資料庫未連接，請連接資料庫")
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}

			describetable2 := fmt.Sprintf(`select COLUMN_NAME, DATA_TYPE, IS_NULLABLE,COLUMN_DEFAULT from %s.INFORMATION_SCHEMA.COLUMNS where TABLE_NAME='%s'`,
				mssqlinformation.Database, params["table2"])
			rows, err := Repo.RawData(MssqlDB, describetable2)
			if err != nil {
				message.Message = "取得資料表資訊時發生錯誤"
				utils.SendError(w, http.StatusInternalServerError, message, err)
				return
			}
			for rows.Next() {
				var result models.Result
				err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Default)
				if err != nil {
					message.Message = "Scan時發生錯誤"
					utils.SendError(w, http.StatusInternalServerError, message, err)
					return
				}

				//匿名函式(檢查是否有重複的欄位)
				firstcheck := utils.Duplicate(col, result.Field)
				//第二道檢查(檢查是否出現在join欄位，可能出現有一樣的欄位卻不是join條件，必須改欄位名)
				secondcheck := utils.Duplicate(join, result.Field)

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
		}

		//設定變數
		var (
			value            = make([]string, len(col))
			valuePtrs        = make([]interface{}, len(col)) //scan時必須使用指針
			jointable1       = make([]string, len(join))
			jointable2       = make([]string, len(join))
			sliceTostringCol = strings.Join(col, ", ") //slice convert string
		)

		for i := 0; i < len(col); i++ {
			valuePtrs[i] = &value[i] //因Scan需要使用指標(valuePtrs)
		}

		for i := 0; i < len(join); i++ {
			jointable1[i] = params["table1"] + "." + join[i]
			jointable2[i] = params["table2"] + "." + join[i]
		}

		var getdata string
		if params["sql1"] == "mysql" && params["sql2"] == "mysql" {
			getdata = fmt.Sprintf("select %s from %s join %s on ",
				sliceTostringCol, params["table1"], params["table2"])
		} else if params["sql1"] == "mssql" && params["sql2"] == "mssql" {
			getdata = fmt.Sprintf("select %s from %s.dbo.%s join %s.dbo.%s on ",
				sliceTostringCol, mssqlinformation.Database, params["table1"], mssqlinformation.Database, params["table2"])
		} else if params["sql1"] != params["sql2"] {
			if params["sql1"] == "mssql" {
				getdata = fmt.Sprintf("select %s from openquery(MYSQL, 'select * from %s.%s') %s join %s.dbo.%s on ",
					sliceTostringCol, mysqlinformation.Database, params["table2"], params["table2"],
					mssqlinformation.Database, params["table1"])
			} else if params["sql2"] == "mssql" {
				getdata = fmt.Sprintf("select %s from openquery(MYSQL, 'select * from %s.%s') %s join %s.dbo.%s on ",
					sliceTostringCol, mysqlinformation.Database, params["table1"], params["table1"],
					mssqlinformation.Database, params["table2"])
			}
		}

		for i := 0; i < len(join); i++ {
			if i == len(join)-1 {
				getdata += fmt.Sprintf("%s=%s", jointable1[i], jointable2[i])
			} else {
				getdata += fmt.Sprintf("%s=%s and ", jointable1[i], jointable2[i])
			}
		}

		var rows *sql.Rows
		if params["sql1"] == "mysql" && params["sql2"] == "mysql" {
			//執行命令
			rows, err = Repo.RawData(MysqlDB, getdata)
		} else {
			rows, err = Repo.RawData(MssqlDB, getdata)
		}

		if err != nil {
			message.Message = "合併資料表時發生錯誤"
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
