package core

import "testing"

func BenchmarkOpen(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	_, err := Open("../data/bunny.odx")
	if err != nil {
		b.Fatal(err)
	}

	b.StopTimer()
}
