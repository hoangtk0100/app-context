package dialects

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Ex: username:password@tcp(localhost:3306)/db_name?charset=utf8mb4&parseTime=True&loc=Local
func MySqlDB(source string) (db *gorm.DB, err error) {
	return gorm.Open(mysql.Open(source), &gorm.Config{})
}
