package driver

import (
	"DynamicAPI/model"
	"DynamicAPI/utils"

	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

//Driver struct
type Driver struct{}


func (d Driver) ConnectDb() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			information model.DBinformation
			message     model.Error
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
			DB, err = d.Connect("mysql", MysqlDataSourceName, information.MaxIdle, information.MaxOpen)
			if err != nil {
				message.Message = "Connect Database failed"
				utils.SendError(w, http.StatusInternalServerError, message)
				return
			}
			fmt.Println(DB)
			utils.SendSuccess(w, "Successfully Connect Database")
			/*
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
						utils.SendError(w, http.StatusInternalServerError, message)
						return
					}

					DB.DB().SetMaxIdleConns(information.MaxIdle)
					DB.DB().SetMaxOpenConns(information.MaxOpen)

					utils.SendSuccess(w, "Successfully Connect Database")
			*/
		}
	}
}

func (d Driver) Connect(engine string, SourceName string, MaxIdle int, MaxOpen int) (*gorm.DB, error) {
	//開啟資料庫連線
	DB, err := gorm.Open(engine, SourceName)
	if err != nil {
		return nil, err
	}
	DB.DB().SetMaxIdleConns(MaxIdle)
	DB.DB().SetMaxOpenConns(MaxOpen)
	return DB, nil
}
