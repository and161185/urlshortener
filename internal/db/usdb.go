package usstorage

import (
	"database/sql"
	"errors"
	"os"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

type dbdriver struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewUSStorage(log *logrus.Logger, dbdrivername string, dbname string) *dbdriver {
	log.Info(dbdrivername, " ", dbname)
	if dbdrivername == "sqlite3" {
		return newUSStorageSqlite3(log, dbdrivername, dbname)
	}
	if dbdrivername == "postgres" {
		return newUSStoragePostgres(log, dbdrivername, dbname)
	}
	log.Fatalf("Work with driver '%s' not defined", dbdrivername)
	return nil
}

func newUSStorageSqlite3(log *logrus.Logger, dbdrivername string, dbname string) *dbdriver {

	dir := "./database"
	filename := "./database/" + dbname
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatal("can't create database directory", err)
	}

	if _, err = os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		_, err = os.Create(filename)
		if err != nil {
			log.Fatal("can't create database file", err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open(dbdrivername, filename)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	CreateUrlsTableSqlite3(db, log)
	CreateClicksTableSqlite3(db, log)

	return &dbdriver{
		db:  db,
		log: log,
	}
}

func newUSStoragePostgres(log *logrus.Logger, dbdrivername string, dbname string) *dbdriver {

	db, err := sql.Open(dbdrivername, dbname)
	if err != nil {
		log.Fatal("Cant connect to database", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("ping")

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	CreateUrlsTablePostgres(db, log)
	CreateClicksTablePostgres(db, log)

	return &dbdriver{
		db:  db,
		log: log,
	}
}
