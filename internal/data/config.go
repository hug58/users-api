package data

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

func makeMigration(db *sql.DB) error {
	b, err := ioutil.ReadFile("./internal/data/models/user.sql")
	if err != nil {
		return err
	}

	rows, err := db.Query(string(b))
	if err != nil {
		return err
	}

	return rows.Close()
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

			if err := makeMigration(DbInstance); err != nil {
				log.Panic("Error creating database: ", err)
				return
			}

		} else {
			log.Println("Single instance already created.")
		}

	}
}
