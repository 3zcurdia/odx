package parsers

import "testing"

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	_, err := LoadPLY("../data/bunny.ply")
	if err != nil {
		b.Fatal(err)
	}

	b.StopTimer()
}
