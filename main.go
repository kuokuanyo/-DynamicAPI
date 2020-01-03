//@title Restful API
//@version 1.0.0
//@description Define an API
//@Schemes http
//@host localhost:8080
//@BasePath /v1
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//DB 資料庫引擎
var DB *gorm.DB

//Result 資料表結構
type Result struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
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
	router.HandleFunc("/getall/{tablename}", GetAllData).Methods("GET")

	//伺服器連線
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}

}

//GetAllData 取得所有資料
func GetAllData(w http.ResponseWriter, r *http.Request) {
	//資料放置
	var data = make(map[string]interface{})
	var results []Result
	var message Error

	//印出url參數
	params := mux.Vars(r)

	//取得資料表資訊的指令
	describetable := fmt.Sprintf("DESCRIBE %s", params["tablename"])
	//將執行的命令
	getalldata := fmt.Sprintf("select * from %s", params["tablename"])

	//取得資料表資訊
	rows, err := DB.Raw(describetable).Rows()
	if err != nil {
		message.Message = "Server(database) error!"
		SendError(w, http.StatusInternalServerError, message)
		return
	}
	for rows.Next() {
		var result Result
		rows.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &result.Default, &result.Extra)
		results = append(results, result)
	}

	//取得所有資料
	rows, err = DB.Raw(getalldata).Rows()
	for rows.Next() {

		var index []string
		var value = make([]interface{}, len(results))

		//取得資料表
		for _, result := range results {
			index = append(index, result.Field)
		}

		rows.Scan(value...)
		for i := 0; i < len(index); i++ {
			data[index[i]] = value[i]
		}
		fmt.Println(data)
		SendSuccess(w, data)
	}
}

//GetTableInformation 取得資料表資訊
//@Summary 取得資料表資訊
//@Tags Table Information
//@Accept json
//@Produce json
//@Success 200 {object} []string "Successfully"
//@Failure 500 {object} models.Error "Internal Server Error"
//@Router /v1/tableinformation/{tablename} [get]
func GetTableInformation(w http.ResponseWriter, r *http.Request) {
	var (
		message Error
		results []Result
	)

	//印出url參數
	params := mux.Vars(r)

	//指令
	describetable := fmt.Sprintf("DESCRIBE %s", params["tablename"])

	rows, err := DB.Raw(describetable).Rows()
	if err != nil {
		message.Message = "Server(database) error!"
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

	//查詢單行
	//func (s *DB) Pluck(column string, value interface{}) *DB
	if err := DB.Raw("show tables").Pluck("Tables Names", &tablenames).Error; err != nil {
		message.Message = "Server(database) error!"
		SendError(w, http.StatusInternalServerError, message)
		return
	}

	SendSuccess(w, tablenames)
}

//ConnectDb 連接資料庫
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
		err         error
		message     Error
	)

	//decode
	json.NewDecoder(r.Body).Decode(&information)

	//印出url參數
	params := mux.Vars(r)

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
			message.Message = "Server(database) error!"
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

		DB, err = gorm.Open("mssql", MssqlDataSourceName)
		if err != nil {
			message.Message = "Server(database) error!"
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
