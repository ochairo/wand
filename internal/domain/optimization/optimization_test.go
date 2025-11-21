package optimization

import "testing"

func TestDiskSpaceChecker(t *testing.T) {
	checker := NewDiskSpaceChecker(1024 * 1024 * 100)
	if !checker.HasEnoughSpace("/tmp") {
		t.Error("HasEnoughSpace should return true")
	}
}

func TestMemoryOptimizer(t *testing.T) {
	optimizer := NewMemoryOptimizer(1024 * 1024 * 1024)
	stats := optimizer.GetMemoryStats()
	if stats.Alloc == 0 {
		t.Error("Memory stats should be populated")
	}
	optimizer.ForceCleanup()
}

func TestBenchmarker(t *testing.T) {
	benchmarker := NewBenchmarker()
	benchmarker.Record("test", 1000000)
	benchmarker.Record("test", 2000000)

	results := benchmarker.GetResults()
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	avg := benchmarker.GetAverage("test")
	expected := int64(1500000)
	if avg != expected {
		t.Errorf("Expected average %d, got %d", expected, avg)
	}
}
