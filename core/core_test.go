package core

import "testing"

func BenchmarkOpen(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	mesh, err := Open("../data/bunny.odx")
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
