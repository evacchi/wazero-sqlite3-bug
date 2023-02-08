package main

import (
	"log"
	"os"

	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func main() {
	sqlite3.Interpreter = false
	res1, err := test()
	if err != nil {
		log.Panic(err)
	}

	sqlite3.Interpreter = true
	res2, err := test()
	if err != nil {
		log.Panic(err)
	}

	const want = "2.500000"
	log.Printf("compiled %q, interpreted %q, want %q", res1, res2, want)
	if res1 != want || res2 != want {
		os.Exit(1)
	}
}

func test() (string, error) {
	db, err := sqlite3.Open(":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, _, err := db.Prepare(`SELECT printf('%f', 2.5)`)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var res string
	for stmt.Step() {
		res = stmt.ColumnText(0)
	}
	if err := stmt.Err(); err != nil {
		return "", err
	}

	err = stmt.Close()
	if err != nil {
		return "", err
	}

	err = db.Close()
	if err != nil {
		return "", err
	}
	return res, nil
}
