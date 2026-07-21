#!/usr/bin/env bash

set -u

# Performance Test Runner
# Runs all benchmarks and generates a report

PROJECT_ROOT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
REPORT_DIR="${PROJECT_ROOT}/dev/reports"
REPORT_FILE="${REPORT_DIR}/benchmark_output.txt"

cd "${PROJECT_ROOT}"

echo "Starting x-uentity performance tests..."

# Create reports directory
mkdir -p "${REPORT_DIR}"

# Run benchmarks
echo "Running benchmarks..."
if go test -bench=. -benchmem -benchtime=1s ./dev/benchmarks -run='^$' > "${REPORT_FILE}" 2>&1; then
    echo "✓ Benchmarks completed successfully"
    echo "  Report saved to: dev/reports/benchmark_output.txt"
else
    echo "✗ Benchmarks failed"
    cat "${REPORT_FILE}"
    exit 1
fi

echo ""
echo "Performance test completed at $(date)"
echo "Check dev/reports/ for results"
