//@title Restful API
//@version 1.0.0
//@description Define an API
//@Schemes http
//@host localhost:8080
//@BasePath /v1
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//DB 資料庫
var DB *gorm.DB

//Result 資料表結構
type Result struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default sql.NullString
	Extra   string
}

//DBinformation 資料庫資訊
type DBinformation struct {
	UserName string
	Password string
	Host     string
	Port     string
	Database string
	MaxIdle  int
	MaxOpen  int
}

//Error 錯誤回傳
type Error struct {
	Message string
}

func main() {
	//路由器
	router := mux.NewRouter()

	//func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *Route
	//func (r *Router) Methods(methods ...string) *Route
	//連接資料庫
	router.HandleFunc("/v1/opendb/{sql}", ConnectDb).Methods("POST")
	router.HandleFunc("/v1/getalltables", GetAlltables).Methods("GET")
	router.HandleFunc("/v1/tableinformation/{tablename}", GetTableInformation).Methods("GET")
	router.HandleFunc("/v1/getall/{tablename}", GetAllData).Methods("GET")
	router.HandleFunc("/v1/getsome/{tablename}", GetSomeData).Methods("GET")
	router.HandleFunc("/v1/addvalue/{tablename}", AddValue).Methods("POST")

	//伺服器連線
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

//AddValue 加入數值至資料表
func AddValue(w http.ResponseWriter, r *http.Request) {
	var (
		message       Error
		params        = mux.Vars(r)
		describetable = fmt.Sprintf("DESCRIBE %s", params["tablename"]) //執行資料庫命令
		queryvalue    = r.URL.Query()["value"]
		index         []string //欄位名稱
		value         []interface{}
	)

	//檢查資料庫是否連接
	if DB == nil {
		message.Message = ("資料庫未連接，請連接資料庫")
		SendError(w, http.StatusInternalServerError, message)
		return
	}

	//取得資料表欄位資訊
	rows, err := DB.Raw(describetable).Rows()
	if err != nil {
		message.Message = "取得欄位資訊時發生錯誤"
		SendError(w, http.StatusInternalServerError, message)
		return
	}
	for rows.Next() {
		var result Result
		err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
		if err != nil {
			message.Message = "Scan資料時發生錯誤"
			SendError(w, http.StatusInternalServerError, message)
			return
		}

		//選擇欄位
		if len(queryvalue) > 0 {
			for x, y := range queryvalue {
				new := strings.Split(queryvalue[i], "=")
			}

		}
	}
}

//GetSomeData 取得部分資料
//@Summary 取得部分資料
//@Tags Table CRUD
//@Accept json
//@Produce json
//@Param tablename path string true "資料表名稱"
//@Param col query array false "挑選欄位"
//@Param where query array false "選擇條件"
//@Success 200 {object} []map[string]interface{} "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/getsome/{tablename} [get]
func GetSomeData(w http.ResponseWriter, r *http.Request) {
	//設定變數
	var (
		datas          []map[string]interface{} //資料放置
		message        Error
		params         = mux.Vars(r)
		describetable  = fmt.Sprintf("DESCRIBE %s", params["tablename"]) //執行資料庫命令
		condition      = make(map[string]interface{})
		getdata        = "select "
		index          []string //欄位名稱
		coltype        []string //欄位類型
		queryvalue     = r.URL.Query()["col"]
		conditionvalue = r.URL.Query()["where"]
	)
	//檢查資料庫是否連接
	if DB == nil {
		message.Message = ("資料庫未連接，請連接資料庫")
		SendError(w, http.StatusInternalServerError, message)
		return
	}

	//取得資料表欄位資訊
	rows, err := DB.Raw(describetable).Rows()
	if err != nil {
		message.Message = "取得欄位資訊時發生錯誤"
		SendError(w, http.StatusInternalServerError, message)
		return
	}
	for rows.Next() {
		var result Result
		err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
		if err != nil {
			message.Message = "Scan資料時發生錯誤"
			SendError(w, http.StatusInternalServerError, message)
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
			new := strings.Split(y, "=")
			condition[new[0]] = new[1]
		}
	}

	//設定變數
	var (
		value     = make([]string, len(index))
		valuePtrs = make([]interface{}, len(index))
	)

	//sql命令
	for i := range index {
		if i == len(index)-1 {
			getdata += fmt.Sprintf("%s from %s ", index[i], params["tablename"])
			valuePtrs[i] = &value[i] //因Scan需要使用指標(valuePtrs)
		} else {
			getdata += fmt.Sprintf("%s, ", index[i])
			valuePtrs[i] = &value[i] //因Scan需要使用指標(valuePtrs)
		}
	}
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
	rows, err = DB.Raw(getdata).Rows()
	if err != nil {
		message.Message = "取資料時發生錯誤"
		SendError(w, http.StatusInternalServerError, message)
		return
	}
	for rows.Next() {
		var data = make(map[string]interface{})
		err = rows.Scan(valuePtrs...)
		if err != nil {
			message.Message = "Scan資料時發生錯誤"
			SendError(w, http.StatusInternalServerError, message)
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
					SendError(w, http.StatusInternalServerError, message)
					return
				}
			} else {
				data[index[i]] = value[i]
			}
		}
		datas = append(datas, data)
	}
	SendSuccess(w, datas)
}

//GetAllData 取得所有資料
//@Summary 取得所有資料
//@Tags Table CRUD
//@Accept json
//@Produce json
//@Param tablename path string true "資料表名稱"
//@Param col query array false "挑選欄位"
//@Success 200 {object} []map[string]interface{} "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/getsll/{tablename} [get]
func GetAllData(w http.ResponseWriter, r *http.Request) {
	var (
		datas         []map[string]interface{} //資料放置
		message       Error
		index         []string //欄位名稱
		coltype       []string //欄位類型
		getdata       = "select "
		params        = mux.Vars(r)                                     //印出url參數
		describetable = fmt.Sprintf("DESCRIBE %s", params["tablename"]) //取得資料表資訊的指令
	)

	//檢查資料庫是否連接
	if DB == nil {
		message.Message = ("資料庫未連接，請連接資料庫")
		SendError(w, http.StatusInternalServerError, message)
		return
	}

	//取得資料表資訊
	rows, err := DB.Raw(describetable).Rows()
	if err != nil {
		message.Message = "取得資料表資訊時發生錯誤"
		SendError(w, http.StatusInternalServerError, message)
		return
	}
	//取得欄位的資訊
	for rows.Next() {
		var result Result
		err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
		if err != nil {
			message.Message = "Scan時發生錯誤"
			SendError(w, http.StatusInternalServerError, message)
			fmt.Println(err)
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
	rows, err = DB.Raw(getdata).Rows()
	if err != nil {
		message.Message = "取得資料時發生錯誤"
		fmt.Println(err)
		SendError(w, http.StatusInternalServerError, message)
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
					SendError(w, http.StatusInternalServerError, message)
					return
				}
			} else {
				data[index[i]] = value[i]
			}
		}
		datas = append(datas, data)
	}
	SendSuccess(w, datas)
}

//GetTableInformation 取得資料表資訊
//@Summary 取得資料表資訊
//@Tags Table Information
//@Accept json
//@Produce json
//@Param tablename path string true "資料庫名稱"
//@Success 200 {object} []string "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/tableinformation/{tablename} [get]
func GetTableInformation(w http.ResponseWriter, r *http.Request) {
	var (
		message Error
		results []Result
		params  = mux.Vars(r) //印出url參數
	)

	//檢查資料庫是否連接
	if DB == nil {
		message.Message = ("資料庫未連接，請連接資料庫")
		SendError(w, http.StatusInternalServerError, message)
		return
	}

	//指令
	describetable := fmt.Sprintf("DESCRIBE %s", params["tablename"])

	rows, err := DB.Raw(describetable).Rows()
	if err != nil {
		message.Message = "取得資料表資訊時發生錯誤"
		SendError(w, http.StatusInternalServerError, message)
		return
	}
	for rows.Next() {
		var result Result
		rows.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
		results = append(results, result)
	}

	SendSuccess(w, results)
}

//GetAlltables 取得所有資料表
//@Summary 取得所有資料表
//@Tags Table Information
//@Accept json
//@Produce json
//@Success 200 {object} []string "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/getalltables [get]
func GetAlltables(w http.ResponseWriter, r *http.Request) {
	var (
		tablenames []string
		message    Error
	)

	//檢查資料庫是否連接
	if DB == nil {
		message.Message = ("資料庫未連接，請連接資料庫")
		SendError(w, http.StatusInternalServerError, message)
		return
	}

	//查詢單行
	//func (s *DB) Pluck(column string, value interface{}) *DB
	if err := DB.Raw("show tables").Pluck("Tables Names", &tablenames).Error; err != nil {
		message.Message = "取得資料表時發生錯誤"
		SendError(w, http.StatusInternalServerError, message)
		return
	}

	SendSuccess(w, tablenames)
}

//ConnectDb 測試是否能連接資料庫
//@Summary 連接資料庫
//@Tags Connect Database(Must be connected first)
//@Accept json
//@Produce json
//@Param sql path string true "資料庫引擎"
//@Param information body DBinformation false "資料庫資訊"
//@Success 200 {string} string "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/opendb/{sql} [post]
func ConnectDb(w http.ResponseWriter, r *http.Request) {
	var (
		information DBinformation
		message     Error
		params      = mux.Vars(r) //印出url參數
		err         error
	)

	//decode
	json.NewDecoder(r.Body).Decode(&information)

	switch strings.ToLower(params["sql"]) {
	case "mysql":
		//完整的資料格式: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
		MysqlDataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			information.UserName,
			information.Password,
			information.Host,
			information.Port,
			information.Database)

		//開啟資料庫連線
		DB, err = gorm.Open("mysql", MysqlDataSourceName)
		if err != nil {
			message.Message = "Connect Database failed"
			SendError(w, http.StatusInternalServerError, message)
			return
		}

		DB.DB().SetMaxIdleConns(information.MaxIdle)
		DB.DB().SetMaxOpenConns(information.MaxOpen)

		SendSuccess(w, "Successfully Connect Database")
	case "mssql":
		MssqlDataSourceName := fmt.Sprintf("sqlserver://%s:%s@%s:%s? database=%s",
			information.UserName,
			information.Password,
			information.Host,
			information.Port,
			information.Database)

		DB, err := gorm.Open("mssql", MssqlDataSourceName)
		if err != nil {
			message.Message = "Connect Database failed"
			SendError(w, http.StatusInternalServerError, message)
			return
		}

		DB.DB().SetMaxIdleConns(information.MaxIdle)
		DB.DB().SetMaxOpenConns(information.MaxOpen)

		SendSuccess(w, "Successfully Connect Database")
	}
}

//SendError response error
func SendError(w http.ResponseWriter, status int, error Error) {
	w.WriteHeader(status)
	//encode
	json.NewEncoder(w).Encode(error)
}

//SendSuccess response success
func SendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	//encode
	json.NewEncoder(w).Encode(data)
}
