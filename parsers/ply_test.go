package parsers

import "testing"

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	mesh, err := LoadPLY("../data/bunny.ply")
	if err != nil {
		b.Fatal(err)
	}

	b.StopTimer()

	if len(mesh.Vertices) != mesh.VertexCount {
		b.Fatalf("vertex count mismatch: expected %d, got %d",
			mesh.VertexCount, len(mesh.Vertices))
	}
	if len(mesh.Faces) != mesh.FaceCount {
		b.Fatalf("face count mismatch: expected %d, got %d",
			mesh.FaceCount, len(mesh.Faces))
	}
}
