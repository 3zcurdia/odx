package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open database
	db, err := sql.Open("sqlite3", "data/bunny.odb")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create tables
	createTables := `
	CREATE TABLE IF NOT EXISTS vertices (
		id INTEGER PRIMARY KEY,
		x INTEGER NOT NULL,
		y INTEGER NOT NULL,
		z INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS faces (
		id INTEGER PRIMARY KEY,
		vertex1 INTEGER NOT NULL,
		vertex2 INTEGER NOT NULL,
		vertex3 INTEGER NOT NULL,
		vertex4 INTEGER
	);
	CREATE TABLE IF NOT EXISTS info (
		name STRING,
		vertices_count INTEGER NOT NULL,
		faces_count INTEGER NOT NULL,
		enhancements STRING,
		extensions STRING
	);
	CREATE UNIQUE INDEX IF NOT EXISTS vertices_xyz ON vertices (x, y, z);
	CREATE UNIQUE INDEX IF NOT EXISTS faces_cuad ON faces (vertex1, vertex2, vertex3, vertex4);
	`

	_, err = db.Exec(createTables)
	if err != nil {
		panic(err)
	}

	// Read file
	file, err := os.Open("data/bunny.ply")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var vertices [][]float64
	var faces [][]int
	var vertexCount, faceCount int
	readingHeader := true

	// Read file line by line
	for scanner.Scan() {
		line := scanner.Text()

		if readingHeader {
			if strings.HasPrefix(line, "element vertex") {
				fields := strings.Fields(line)
				vertexCount, _ = strconv.Atoi(fields[2])
			}
			if strings.HasPrefix(line, "element face") {
				fields := strings.Fields(line)
				faceCount, _ = strconv.Atoi(fields[2])
			}
			if strings.HasPrefix(line, "end_header") {
				readingHeader = false
			}
		} else {
			if len(vertices) < vertexCount {
				fields := strings.Fields(line)
				vertex := make([]float64, 3)
				for i := 0; i < 3; i++ {
					vertex[i], _ = strconv.ParseFloat(fields[i], 64)
				}
				vertices = append(vertices, vertex)
			} else {
				fields := strings.Fields(line)
				face := make([]int, len(fields)-1)
				for i := 1; i < len(fields); i++ {
					val, _ := strconv.Atoi(fields[i])
					face[i-1] = val + 1
				}
				faces = append(faces, face)
			}
		}
	}

	// Verify counts
	if len(vertices) != vertexCount {
		panic("Vertex count mismatch")
	}
	if len(faces) != faceCount {
		panic("Face count mismatch")
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	// Insert info
	_, err = tx.Exec("INSERT INTO info (name, vertices_count, faces_count) VALUES (?, ?, ?)",
		"bunny", vertexCount, faceCount)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	// Insert vertices
	stmtVertex, err := tx.Prepare("INSERT INTO vertices (x, y, z) VALUES (?, ?, ?)")
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	defer stmtVertex.Close()

	for _, v := range vertices {
		_, err = stmtVertex.Exec(v[0], v[1], v[2])
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	// Insert faces
	stmtFace4, err := tx.Prepare("INSERT INTO faces (vertex1, vertex2, vertex3, vertex4) VALUES (?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	defer stmtFace4.Close()

	stmtFace3, err := tx.Prepare("INSERT INTO faces (vertex1, vertex2, vertex3) VALUES (?, ?, ?)")
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	defer stmtFace3.Close()

	for _, f := range faces {
		if len(f) == 4 {
			_, err = stmtFace4.Exec(f[0], f[1], f[2], f[3])
		} else {
			_, err = stmtFace3.Exec(f[0], f[1], f[2])
		}
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	fmt.Println("Done!")
}
