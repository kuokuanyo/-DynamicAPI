//@title Restful API
//@version 1.0.0
//@description Define an API
//@Schemes http
//@host localhost:8080
//@BasePath /v1
package main

import (
	"DynamicAPI/controllers"

	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	controller := controllers.Controller{}

	//路由器
	router := mux.NewRouter()

	//func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *Route
	//func (r *Router) Methods(methods ...string) *Route
	//連接資料庫
	router.HandleFunc("/v1/opendb/{sql}", controller.ConnectDb()).Methods("POST")
	//資料庫資訊
	router.HandleFunc("/v1/{sql}/getalltables", controller.GetAlltables()).Methods("GET")
	router.HandleFunc("/v1/{sql}/tableinformation/{tablename}", controller.GetTableInformation()).Methods("GET")
	//CRUD
	router.HandleFunc("/v1/{sql}/getall/{tablename}", controller.GetAllData()).Methods("GET")
	router.HandleFunc("/v1/{sql}/getsome/{tablename}", controller.GetSomeData()).Methods("GET")
	router.HandleFunc("/v1/{sql}/addvalue/{tablename}", controller.AddValue()).Methods("POST")
	router.HandleFunc("/v1/{sql}/update/{tablename}", controller.UpdateValue()).Methods("PUT")
	router.HandleFunc("/v1/{sql}/delete/{tablename}", controller.DeleteValue()).Methods("DELETE")
	//合併表
	router.HandleFunc("/v1/jointable/{sql1}/{table1}/{sql2}/{table2}", controller.JoinTable()).Methods("GET")
	router.HandleFunc("/v1/jointable/{uuid}", controller.GetTableByUUID()).Methods("GET")

	//伺服器連線
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
