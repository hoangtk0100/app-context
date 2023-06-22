package dialects

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Ex: /tmp/gorm.db
func SQLiteDB(source string) (db *gorm.DB, err error) {
	return gorm.Open(sqlite.Open(source), &gorm.Config{})
}
