package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func main() {
	os.Remove("./stocks.db")

	db, err := sql.Open("sqlite3", "./stocks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sql := `
	create table stocks (id integer not null primary key, symbol text);
	delete from stocks;
	`
	_, err = db.Exec(sql)
	if err != nil {
		log.Printf("%q: %s\n", err, sql)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into stocks(id, symbol) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		//_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
		_, err = stmt.Exec(i, fmt.Sprintf("poppy poopy%03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()

	rows, err := db.Query("select id, symbol from stocks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var symbol string
		rows.Scan(&id, &symbol)
		fmt.Println(id, symbol)
	}
	rows.Close()

	stmt, err = db.Prepare("select symbol from stocks where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var symbol string
	err = stmt.QueryRow("3").Scan(&symbol)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(symbol)

	_, err = db.Exec("delete from stocks")
	if err != nil {
		log.Fatal(err)
	}

	rows, err = db.Query("select id, symbol from stocks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var symbol string
		rows.Scan(&id, &symbol)
		fmt.Println(id, symbol)
	}
}
