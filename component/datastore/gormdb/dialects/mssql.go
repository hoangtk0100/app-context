package dialects

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// Ex: sqlserver://username:password@localhost:1433?database=db_name
func MSSqlDB(source string) (db *gorm.DB, err error) {
	return gorm.Open(sqlserver.Open(source), &gorm.Config{})
}
