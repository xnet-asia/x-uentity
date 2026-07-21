package benchmarks

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestPerformanceReportGeneration tests the report generation functionality
func TestPerformanceReportGeneration(t *testing.T) {
	gen := NewReportGenerator()
	const reportName = "performance_test_report.json"
	t.Cleanup(func() {
		_ = os.Remove(filepath.Join("reports", reportName))
		_ = os.Remove("reports")
	})

	// Add sample results
	gen.AddResult("Create", 1000000, 500*time.Millisecond, 2000000)
	gen.AddResult("Query", 100000, 200*time.Millisecond, 500000)
	gen.AddResult("GetByIdentifier", 5000000, 100*time.Millisecond, 50000000)
	gen.AddResult("Update", 1000000, 600*time.Millisecond, 1666667)
	gen.AddResult("Delete", 2000000, 300*time.Millisecond, 6666667)
	gen.AddResult("ConcurrentOperations", 500000, 1000*time.Millisecond, 500000)

	// Save report
	err := gen.SaveReport(reportName)
	if err != nil {
		t.Fatalf("Failed to save report: %v", err)
	}

	// Print report
	gen.PrintReport()

	// Verify results
	if len(gen.Results) != 6 {
		t.Fatalf("Expected 6 results, got %d", len(gen.Results))
	}
}
