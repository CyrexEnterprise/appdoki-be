package main

import (
	"appdoki-be/config"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
)

type migrationsLogger struct {
	verbose bool
}

func (ml *migrationsLogger) Printf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

func (ml *migrationsLogger) Verbose() bool {
	return ml.verbose
}

func runMigrations(db *sql.DB, conf *config.DatabaseConfig) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(conf.MigrationsDir, "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	m.Log = &migrationsLogger{
		verbose: conf.MigrationsLogVerbose,
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}
