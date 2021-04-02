package db

import (
	"encoding/json"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"raspberrypi-gpio-manager-backend/model"
)

var (
	Connection *gorm.DB
)

func ConnectDatabase(db model.Database) {

	dsn := db.User + ":" + db.Pass + "@/" + db.Name + "?charset=utf8&parseTime=True&loc=Local"
	con, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	Connection = con
}

func LoadDatabaseConfig() model.Database {
	configFile, err := os.Open("config/config.db.json")
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	db := model.Database{}
	err = decoder.Decode(&db)
	if err != nil {
		panic(err)
	}

	return db
}
