package dbmigration

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/sqlserver"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/spf13/pflag"
)

type opt struct {
	migrationURL string
	dbSource     string
}

type dbMigrator struct {
	id        string
	migration *migrate.Migrate
	logger    appctx.Logger
	*opt
}

func NewDBMigration(id string) *dbMigrator {
	return &dbMigrator{
		id:  id,
		opt: new(opt),
	}
}

func (m *dbMigrator) ID() string {
	return m.id
}

func (m *dbMigrator) InitFlags() {
	pflag.StringVar(&m.dbSource,
		"db-migration-source",
		"",
		"Database connection string",
	)

	pflag.StringVar(&m.migrationURL,
		"db-migration-url",
		"",
		"Database migration url - Default: file://migration",
	)
}

func (m *dbMigrator) isDisabled() bool {
	return m.dbSource == ""
}

func (m *dbMigrator) Run(ac appctx.AppContext) error {
	if m.isDisabled() {
		return nil
	}

	m.logger = ac.Logger(m.id)

	migration, err := migrate.New(m.migrationURL, m.dbSource)
	if err != nil {
		m.logger.Fatal(err, "Can not create new database migrate instance")
	}

	m.migration = migration
	return nil
}

func (m *dbMigrator) Stop() error {
	return nil
}

func (m *dbMigrator) MigrateUp() {
	if err := m.migration.Up(); err != nil && err != migrate.ErrNoChange {
		m.logger.Fatal(err, "Failed to run migrate up")
	}

	m.logger.Print("DB migrated successfully")
}

func (m *dbMigrator) MigrateDown() {
	if err := m.migration.Down(); err != nil && err != migrate.ErrNoChange {
		m.logger.Fatal(err, "Failed to run migrate down")
	}

	m.logger.Print("DB migrated successfully")
}
