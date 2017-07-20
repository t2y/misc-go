package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/naoina/genmai"
)

func newDB(dsn string) (db *genmai.DB, err error) {
	dialect := &genmai.SQLite3Dialect{}
	db, err = genmai.New(dialect, dsn)
	return
}

type MyTable struct {
	Key   string
	Value string
}

func createAndInsertMyTable(db *genmai.DB) {
	table := &MyTable{
		Key:   "test1",
		Value: "ttt",
	}

	if err := db.CreateTableIfNotExists(table); err != nil {
		log.Fatal("create table: ", err)
	}

	if _, err := db.Insert(table); err != nil {
		log.Fatal("insert: ", err)
	}
}

type MyMap struct {
	m map[string]string
}

func createAndInsertMap(db *genmai.DB) {
	m := MyMap{
		m: map[string]string{
			"key":   "test2",
			"value": "sss",
		},
	}

	if err := db.CreateTableIfNotExists(m); err != nil {
		log.Fatal("create table: ", err)
	}

	if _, err := db.Insert(m); err != nil {
		log.Fatal("insert: ", err)
	}
}

func main() {
	db, err := newDB("test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createAndInsertMyTable(db)
	createAndInsertMap(db)
}
