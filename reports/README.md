# Performance Test Reports

This directory contains performance test reports for x-uentity.

## Benchmark Reports

- `benchmark_output.txt` - Raw benchmark output from Go testing
- `performance_test_report.json` - Structured JSON report with detailed metrics
- `sample_performance_report.json` - Example performance report

## Running Benchmarks

### Quick Start

```bash
# Run all benchmarks
go test -bench=. -benchmem ./benchmarks
```

### Using the Provided Script

```bash
bash scripts/run_benchmarks.sh
```

## Benchmark Operations

1. **Create** - Insert new entities into the repository
2. **Query** - Filter entities using predicates
3. **GetByIdentifier** - Direct key-based lookups
4. **Update** - Modify existing entities
5. **Delete** - Remove entities from storage
6. **ConcurrentOperations** - Mixed operations under concurrent load

## Metrics

Each benchmark measures:
- **Iterations** - Number of operations performed
- **Total Time** - Wall-clock time for all iterations
- **Avg Time** - Average time per operation (nanoseconds)
- **Ops/Sec** - Operations per second throughput

## Interpreting Results

- Lower average time = better performance
- Higher ops/sec = better throughput
- Concurrent operations test validates thread-safety
