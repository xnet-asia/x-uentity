#!/bin/bash

# Performance Test Runner
# Runs all benchmarks and generates a report

echo "Starting x-uentity performance tests..."
echo ""

# Create reports directory
mkdir -p reports

# Run benchmarks
echo "Running benchmarks..."
go test -bench=. -benchmem -benchtime=1s ./benchmarks -run=^$ > reports/benchmark_output.txt 2>&1

if [ $? -eq 0 ]; then
    echo "✓ Benchmarks completed successfully"
    echo "  Report saved to: reports/benchmark_output.txt"
else
    echo "✗ Benchmarks failed"
    exit 1
fi

echo ""
echo "Performance test completed at $(date)"
echo "Check reports/ directory for results"
