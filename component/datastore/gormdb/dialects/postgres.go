package dialects

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Ex: postgresql://username:password@localhost:3306/db_name?sslmode=disable
func PostgresDB(source string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(source), &gorm.Config{})
}
