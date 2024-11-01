package core

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

func loadMeshInfo(db *sql.DB) (*Mesh, error) {
	var vertexCount, faceCount int
	err := db.QueryRow("SELECT vertices_count, faces_count FROM info WHERE id = 1").Scan(&vertexCount, &faceCount)
	if err != nil {
		return nil, err
	}
	mesh := &Mesh{
		Vertices:    make([][]float64, 0, vertexCount),
		Faces:       make([][]int, 0, faceCount),
		VertexCount: vertexCount,
		FaceCount:   faceCount,
	}

	return mesh, nil
}

func loadVertices(db *sql.DB, vertexCount int) ([][]float64, error) {
	vertices := make([][]float64, vertexCount)

	stmt, err := db.Prepare("SELECT id, x, y, z FROM vertices")
	if err != nil {
		return nil, fmt.Errorf("prepare vertices stmt: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("query vertices: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		vertex := make([]float64, 3)
		err := rows.Scan(&id, &vertex[0], &vertex[1], &vertex[2])
		if err != nil {
			return nil, fmt.Errorf("scan vertex: %w", err)
		}
		vertices[id-1] = vertex
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("vertices rows: %w", err)
	}

	return vertices, nil
}

func loadFaces(db *sql.DB, faceCount int) ([][]int, error) {
	faces := make([][]int, faceCount)

	stmt, err := db.Prepare("SELECT id, vertex1, vertex2, vertex3 FROM faces")
	if err != nil {
		return nil, fmt.Errorf("prepare faces stmt: %w", err)
	}
	defer stmt.Close()

	// face4Stmt, err := db.Prepare("SELECT id, vertex1, vertex2, vertex3, vertex4 FROM faces")
	// if err != nil {
	// 	return nil, err
	// }
	// defer face4Stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("query faces: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		face := make([]int, 3)
		err := rows.Scan(&id, &face[0], &face[1], &face[2])
		if err != nil {
			return nil, fmt.Errorf("scan face: %w", err)
		}
		face[0]--
		face[1]--
		face[2]--
		faces[id-1] = face
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("faces rows: %w", err)
	}

	return faces, nil
}

// It opens a mesh from an ODX file
func Open(filename string) (*Mesh, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	// Enable performance optimizations
	_, err = db.Exec("PRAGMA synchronous = OFF")
	if err != nil {
		return nil, fmt.Errorf("set synchronous mode: %w", err)
	}
	_, err = db.Exec("PRAGMA journal_mode = MEMORY")
	if err != nil {
		return nil, fmt.Errorf("set journal mode: %w", err)
	}

	// Read mesh info
	mesh, err := loadMeshInfo(db)
	if err != nil {
		return nil, fmt.Errorf("load mesh info: %w", err)
	}

	// Create error channel and wait group
	errChan := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		vertices, err := loadVertices(db, mesh.VertexCount)
		if err != nil {
			errChan <- err
			return
		}
		mesh.Vertices = vertices
	}()

	// Load faces concurrently
	go func() {
		defer wg.Done()
		faces, err := loadFaces(db, mesh.FaceCount)
		if err != nil {
			errChan <- err
			return
		}
		mesh.Faces = faces
	}()

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Check for any errors from goroutines
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return mesh, nil
}
