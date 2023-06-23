package gormdb

import (
	"errors"
	"fmt"
	"strings"
	"time"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/component/datastore/gormdb/dialects"
	"github.com/spf13/pflag"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type gormDBDriver int

const (
	gormDBDriverNotSupported gormDBDriver = iota
	gormDBDriverMySQL
	gormDBDriverPostgres
	gormDBDriverSQLite
	gormDBDriverMSSQL
)

type gormOpt struct {
	prefix          string
	source          string
	dbDriver        string
	maxOpenConns    int
	maxIdleConns    int
	connMaxIdleTime int
}

type gormDB struct {
	id     string
	logger appctx.Logger
	db     *gorm.DB
	*gormOpt
}

func NewGormDB(id, prefix string) *gormDB {
	return &gormDB{
		id: id,
		gormOpt: &gormOpt{
			prefix: strings.TrimSpace(prefix),
		},
	}
}

func (gdb *gormDB) ID() string {
	return gdb.id
}

func (gdb *gormDB) InitFlags() {
	prefix := gdb.prefix
	if prefix != "" {
		prefix += "-"
	}

	pflag.StringVar(
		&gdb.source,
		fmt.Sprintf("%sdb-source", prefix),
		"",
		"Database connection string",
	)

	pflag.StringVar(
		&gdb.dbDriver,
		fmt.Sprintf("%sdb-driver", prefix),
		"mysql",
		"Database driver (mysql | postgres | sqlite | mssql) - Default: mysql",
	)

	pflag.IntVar(
		&gdb.maxOpenConns,
		fmt.Sprintf("%sdb-max-open-conns", prefix),
		30,
		"Maximum number of open connections to the database - Default: 30",
	)

	pflag.IntVar(
		&gdb.maxIdleConns,
		fmt.Sprintf("%sdb-max-ide-conns", prefix),
		10,
		"Maximum number of database connections in the idle - Default: 10",
	)

	pflag.IntVar(
		&gdb.connMaxIdleTime,
		fmt.Sprintf("%sdb-max-conn-ide-time", prefix),
		3600,
		"Maximum amount of time a connection may be idle in seconds - Default: 3600",
	)
}

func (gdb *gormDB) isDisabled() bool {
	return gdb.source == ""
}

func (gdb *gormDB) Run(ac appctx.AppContext) error {
	if gdb.isDisabled() {
		return nil
	}

	gdb.logger = ac.Logger(gdb.id)

	dbDriver := getDBDriver(gdb.dbDriver)
	if dbDriver == gormDBDriverNotSupported {
		return errors.New("database driver not supported")
	}

	gdb.logger.Info("Connect to Gorm DB at ", gdb.source, " ...")

	var err error
	gdb.db, err = gdb.getDBConn(dbDriver)
	if err != nil {
		gdb.logger.Error(err, "Cannot connect to database")
		return err
	}

	return nil
}

func (gdb *gormDB) Stop() error {
	return nil
}

func (gdb *gormDB) GetDB() *gorm.DB {
	if gdb.logger.GetLevel() == "debug" || gdb.logger.GetLevel() == "trace" {
		return gdb.db.Session(&gorm.Session{NewDB: true}).Debug()
	}

	newSession := gdb.db.Session(
		&gorm.Session{
			NewDB:  true,
			Logger: gdb.db.Logger.LogMode(logger.Silent),
		},
	)

	if db, err := newSession.DB(); err != nil {
		db.SetMaxOpenConns(gdb.maxOpenConns)
		db.SetMaxIdleConns(gdb.maxIdleConns)
		db.SetConnMaxIdleTime(time.Second * time.Duration(gdb.connMaxIdleTime))
	}

	return newSession
}

func getDBDriver(driver string) gormDBDriver {
	switch strings.ToLower(driver) {
	case "mysql":
		return gormDBDriverMySQL
	case "postgres":
		return gormDBDriverPostgres
	case "sqlite":
		return gormDBDriverSQLite
	case "mssql":
		return gormDBDriverMSSQL
	}

	return gormDBDriverNotSupported
}

func (gdb *gormDB) getDBConn(driver gormDBDriver) (dbConn *gorm.DB, err error) {
	switch driver {
	case gormDBDriverMySQL:
		return dialects.MySqlDB(gdb.source)
	case gormDBDriverPostgres:
		return dialects.PostgresDB(gdb.source)
	case gormDBDriverSQLite:
		return dialects.SQLiteDB(gdb.source)
	case gormDBDriverMSSQL:
		return dialects.MSSqlDB(gdb.source)
	}

	return nil, nil
}
