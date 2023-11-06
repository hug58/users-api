package data

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	_ "github.com/lib/pq"
)

var (
	lock       = &sync.Mutex{}
	DbInstance *sql.DB
)

func getConnection() (*sql.DB, error) {
	//uri := "postgres://hug58:grosss213@127.0.0.1:5432/microvideogame?sslmode=disable"
	pgConnString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
	)

	return sql.Open("postgres", pgConnString)
}

func makeMigrations(db *sql.DB, dir string) error {
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".sql" {
			sqlFile, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			_, err = db.Exec(string(sqlFile))
			if err != nil && !strings.Contains(err.Error(), "already exists") {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func init() {
	if DbInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		if DbInstance == nil {
			fmt.Println("Creating single instance now.")
			db, err := getConnection()
			if err != nil {
				log.Panic("Error connecting to database: ", err)
				return
			}

			DbInstance = db

			if err := makeMigrations(DbInstance, "./internal/data/models"); err != nil {
				log.Panic("Error creating database: ", err)
				return
			}

		} else {
			log.Println("Single instance already created.")
		}

	}
}
