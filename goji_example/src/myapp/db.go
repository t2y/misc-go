package main

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func isExistFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func InitDatabase(dbfile string) (db *sql.DB) {
	isExist := isExistFile(dbfile)

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		Logger.Fatal(err)
	}

	if isExist {
		return
	}

	sqlStmt := `
		create table users (
			id integer not null primary key,
			name text,
			age	integer
		);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		Logger.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	InsertUser(db, 1, "user1", 23)
	return
}

func InsertUser(db *sql.DB, id int, name string, age int) {
	if db == nil {
		return
	}
	s := "insert into users(id, name, age) values (?, ?, ?)"
	ExecSql(db, s, id, name, age)
}

func UpdateUser(db *sql.DB, oldName string, newName string, age int) {
	if db == nil {
		return
	}
	s := "update users set name=?, age=? where name=?"
	ExecSql(db, s, newName, age, oldName)
}

func DeleteUser(db *sql.DB, name string) {
	if db == nil {
		return
	}
	s := "delete from users where name=?"
	ExecSql(db, s, name)
}

func ExecSql(db *sql.DB, sql string, params ...interface{}) {
	tx, err := db.Begin()
	if err != nil {
		Logger.Fatal(err)
	}

	stmt, err := tx.Prepare(sql)
	if err != nil {
		Logger.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		Logger.Fatal(err)
	}

	tx.Commit()
}

func QueryRowSql(db *sql.DB, sqlStmt string, params ...interface{}) (row *sql.Row) {
	if db == nil {
		return
	}
	row = db.QueryRow(sqlStmt, params...)
	return
}

func QuerySql(db *sql.DB, sqlStmt string, params ...interface{}) (rows *sql.Rows) {
	if db == nil {
		return
	}
	rows, _ = db.Query(sqlStmt, params...)
	return
}
