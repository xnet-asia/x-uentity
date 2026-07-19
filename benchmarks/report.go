package benchmarks

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ReportGenerator generates performance test reports
type ReportGenerator struct {
	Results []RepositoryBenchmarkResult
}

// RepositoryBenchmarkResult represents a single benchmark result
type RepositoryBenchmarkResult struct {
	Operation      string        `json:"operation"`
	Iterations     int           `json:"iterations"`
	TotalTime      string        `json:"total_time"`
	AvgTime        string        `json:"avg_time_ns"`
	OpsPerSecond   float64       `json:"ops_per_second"`
	Timestamp      string        `json:"timestamp"`
}

// NewReportGenerator creates a new report generator
func NewReportGenerator() *ReportGenerator {
	return &ReportGenerator{
		Results: make([]RepositoryBenchmarkResult, 0),
	}
}

// AddResult adds a benchmark result
func (rg *ReportGenerator) AddResult(operation string, iterations int, totalTime time.Duration, opsPerSecond float64) {
	result := RepositoryBenchmarkResult{
		Operation:    operation,
		Iterations:   iterations,
		TotalTime:    totalTime.String(),
		AvgTime:      fmt.Sprintf("%d ns", totalTime.Nanoseconds()/int64(iterations)),
		OpsPerSecond: opsPerSecond,
		Timestamp:    time.Now().Format(time.RFC3339),
	}
	rg.Results = append(rg.Results, result)
}

// SaveReport saves the report to a JSON file
func (rg *ReportGenerator) SaveReport(filename string) error {
	// Create reports directory if it doesn't exist
	reportDir := "reports"
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(reportDir, filename)

	// Create report structure
	report := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"total_benchmarks": len(rg.Results),
		"results":        rg.Results,
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	fmt.Printf("Report saved to: %s\n", filePath)
	return nil
}

// PrintReport prints the report to stdout
func (rg *ReportGenerator) PrintReport() {
	fmt.Println("\n" + string([]byte{0x1b, 0x5b, 0x31, 0x3b, 0x33, 0x36, 0x6d}) + "=== Performance Report ===" + string([]byte{0x1b, 0x5b, 0x30, 0x6d}))
	fmt.Printf("Generated: %s\n\n", time.Now().Format(time.RFC3339))

	fmt.Printf("%-30s | %-12s | %-15s | %-15s | %-15s\n", "Operation", "Iterations", "Total Time", "Avg Time (ns)", "Ops/Sec")
	fmt.Println(string(make([]byte, 100)))

	for _, result := range rg.Results {
		fmt.Printf("%-30s | %-12d | %-15s | %-15s | %-15.2f\n",
			result.Operation,
			result.Iterations,
			result.TotalTime,
			result.AvgTime,
			result.OpsPerSecond,
		)
	}

	fmt.Println()
}
