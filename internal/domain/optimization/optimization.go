// Package optimization provides system optimization strategies.
package optimization

import "runtime"

// DiskSpaceChecker checks available disk space
type DiskSpaceChecker struct {
	minFreeBytes int64
}

// NewDiskSpaceChecker creates a new disk space checker
func NewDiskSpaceChecker(minFreeBytes int64) *DiskSpaceChecker {
	return &DiskSpaceChecker{
		minFreeBytes: minFreeBytes,
	}
}

// HasEnoughSpace checks if there's enough free disk space
func (d *DiskSpaceChecker) HasEnoughSpace(_ string) bool {
	// Platform-specific implementation would go here
	// This is a placeholder that always returns true for now
	return true
}

// MemoryOptimizer provides memory optimization utilities
type MemoryOptimizer struct {
	maxMemoryPercent float64
}

// NewMemoryOptimizer creates a new memory optimizer
func NewMemoryOptimizer(maxMemoryPercent float64) *MemoryOptimizer {
	return &MemoryOptimizer{
		maxMemoryPercent: maxMemoryPercent,
	}
}

// GetMemoryStats returns current memory statistics
func (m *MemoryOptimizer) GetMemoryStats() runtime.MemStats {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	return stats
}

// ShouldCleanup determines if cleanup should run based on memory usage
func (m *MemoryOptimizer) ShouldCleanup() bool {
	stats := m.GetMemoryStats()
	// If memory usage exceeds threshold, suggest cleanup
	return float64(stats.Alloc) > m.maxMemoryPercent
}

// ForceCleanup forces garbage collection
func (m *MemoryOptimizer) ForceCleanup() {
	runtime.GC()
}

// Benchmarker provides benchmarking utilities
type Benchmarker struct {
	results map[string]BenchmarkResult
}

// BenchmarkResult holds benchmark data
type BenchmarkResult struct {
	Name     string
	Duration int64 // nanoseconds
	Calls    int64
}

// NewBenchmarker creates a new benchmarker
func NewBenchmarker() *Benchmarker {
	return &Benchmarker{
		results: make(map[string]BenchmarkResult),
	}
}

// Record adds a benchmark result
func (b *Benchmarker) Record(name string, durationNs int64) {
	result, exists := b.results[name]
	if exists {
		result.Duration += durationNs
		result.Calls++
	} else {
		result = BenchmarkResult{
			Name:     name,
			Duration: durationNs,
			Calls:    1,
		}
	}
	b.results[name] = result
}

// GetResults returns all benchmark results
func (b *Benchmarker) GetResults() map[string]BenchmarkResult {
	return b.results
}

// GetAverage returns average duration for a benchmark
func (b *Benchmarker) GetAverage(name string) int64 {
	result, exists := b.results[name]
	if !exists || result.Calls == 0 {
		return 0
	}
	return result.Duration / result.Calls
}
