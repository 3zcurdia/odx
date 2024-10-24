package core

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	OutputFile string
	SkipIndex  bool
}

type Mesh struct {
	Vertices    [][]float64
	Faces       [][]int
	VertexCount int
	FaceCount   int
}

// Init initializes a new ODB file with the core tables
func Init(config Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", config.OutputFile)
	if err != nil {
		return nil, err
	}

	if err := createTables(db, config.SkipIndex); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB, skipIndexes bool) error {
	query := `
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
	  id INTEGER PRIMARY KEY,
		vertices_count INTEGER NOT NULL,
		faces_count INTEGER NOT NULL,
		enhancements STRING,
		extensions STRING,
		inserted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	if !skipIndexes {
		_, err = db.Exec(`
			CREATE UNIQUE INDEX IF NOT EXISTS vertices_xyz ON vertices (x, y, z);
			CREATE UNIQUE INDEX IF NOT EXISTS faces_cuad ON faces (vertex1, vertex2, vertex3, vertex4);
			`)
		if err != nil {
			return err
		}
	}
	return nil
}

// Insert inserts a mesh into an ODB file
func Insert(db *sql.DB, mesh *Mesh) error {
	validate(mesh)
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = insertInfo(tx, mesh); err != nil {
		return err
	}
	if err = insertVertices(tx, mesh.Vertices); err != nil {
		return err
	}
	if err = insertFaces(tx, mesh.Faces); err != nil {
		return err
	}

	return tx.Commit()
}

func validate(mesh *Mesh) error {
	if len(mesh.Vertices) != mesh.VertexCount {
		return fmt.Errorf("vertex count mismatch: expected %d, got %d",
			mesh.VertexCount, len(mesh.Vertices))
	}
	if len(mesh.Faces) != mesh.FaceCount {
		return fmt.Errorf("face count mismatch: expected %d, got %d",
			mesh.FaceCount, len(mesh.Faces))
	}
	return nil
}

func insertInfo(tx *sql.Tx, mesh *Mesh) error {
	_, err := tx.Exec("INSERT INTO info (vertices_count, faces_count) VALUES (?, ?)", mesh.VertexCount, mesh.FaceCount)
	return err
}

func insertVertices(tx *sql.Tx, vertices [][]float64) error {
	stmt, err := tx.Prepare("INSERT INTO vertices (x, y, z) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range vertices {
		if _, err := stmt.Exec(v[0], v[1], v[2]); err != nil {
			return err
		}
	}
	return nil
}

func insertFaces(tx *sql.Tx, faces [][]int) error {
	stmt3, err := tx.Prepare("INSERT INTO faces (vertex1, vertex2, vertex3) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt3.Close()

	stmt4, err := tx.Prepare("INSERT INTO faces (vertex1, vertex2, vertex3, vertex4) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt4.Close()

	for _, f := range faces {
		if len(f) == 4 {
			if _, err := stmt4.Exec(f[0], f[1], f[2], f[3]); err != nil {
				return err
			}
		} else {
			if _, err := stmt3.Exec(f[0], f[1], f[2]); err != nil {
				return err
			}
		}
	}
	return nil
}
