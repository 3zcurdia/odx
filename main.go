package main

import (
	"fmt"

	core "github.com/3zcurdia/odb/odb/core"
	parsers "github.com/3zcurdia/odb/odb/parsers"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := core.Init("data/bunny.odb")
	if err != nil {
		panic(fmt.Sprintf("error initializing database: %v", err))
	}
	defer db.Close()

	fmt.Println("Reading PLY file... ")
	mesh, err := parsers.LoadPLY("data/bunny.ply")
	if err != nil {
		panic(fmt.Sprintf("error loading PLY file: %v", err))
	}

	fmt.Println("Inserting data into database... ")
	if err := core.Insert(db, mesh); err != nil {
		panic(fmt.Sprintf("error inserting data: %v", err))
	}

	fmt.Println("Done!")
}
