package parsers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	core "github.com/3zcurdia/odb/odb/core"
)

// LoadPLY reads and parses the PLY file into a core Mesh struct
func LoadPLY(filename string) (*core.Mesh, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	data := &core.Mesh{
		Vertices: make([][]float64, 0),
		Faces:    make([][]int, 0),
	}

	scanner := bufio.NewScanner(file)
	readingHeader := true

	for scanner.Scan() {
		line := scanner.Text()

		if readingHeader {
			if strings.HasPrefix(line, "element vertex") {
				fields := strings.Fields(line)
				data.VertexCount, _ = strconv.Atoi(fields[2])
			}
			if strings.HasPrefix(line, "element face") {
				fields := strings.Fields(line)
				data.FaceCount, _ = strconv.Atoi(fields[2])
			}
			if strings.HasPrefix(line, "end_header") {
				readingHeader = false
			}
			continue
		}

		if len(data.Vertices) < data.VertexCount {
			vertex, err := parseVertex(line)
			if err != nil {
				return nil, fmt.Errorf("error parsing vertex: %v", err)
			}
			data.Vertices = append(data.Vertices, vertex)
		} else {
			face, err := parseFace(line)
			if err != nil {
				return nil, fmt.Errorf("error parsing face: %v", err)
			}
			data.Faces = append(data.Faces, face)
		}
	}
	return data, nil
}

func parseVertex(line string) ([]float64, error) {
	fields := strings.Fields(line)
	vertex := make([]float64, 3)
	for i := 0; i < 3; i++ {
		val, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing vertex coordinate: %v", err)
		}
		vertex[i] = val
	}
	return vertex, nil
}

// parseFace parses a face line into integer values
func parseFace(line string) ([]int, error) {
	fields := strings.Fields(line)
	face := make([]int, len(fields)-1)
	for i := 1; i < len(fields); i++ {
		val, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil, fmt.Errorf("error parsing face index: %v", err)
		}
		face[i-1] = val + 1
	}
	return face, nil
}
