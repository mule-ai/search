#!/bin/bash
# Coverage report generation script

set -e

echo "Generating coverage report..."

# Run tests with coverage
go test -coverprofile=coverage.out ./... -v 2>&1 | grep -E "^(ok|FAIL|coverage:)" > coverage_summary.txt

# Generate total coverage
TOTAL=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}')
echo ""
echo "=== Total Coverage: $TOTAL ==="
echo ""

# Generate per-package coverage
echo "=== Per-Package Coverage ==="
go tool cover -func=coverage.out | grep -E "github.com/mule-ai/search" | awk '
{
    # Extract package name (everything before :)
    split($1, parts, ":")
    pkg = parts[1]
    
    # Extract coverage percentage
    cov = $3
    
    # Store last coverage for each package
    pkg_cov[pkg] = cov
}

END {
    for (p in pkg_cov) {
        print p " " pkg_cov[p]
    }
}' | sort

echo ""

# Generate HTML report
echo "Generating HTML coverage report..."
go tool cover -html=coverage.out -o coverage.html
echo "HTML report generated: coverage.html"

# Check if we meet the threshold
THRESHOLD=80
echo ""
echo "Checking against threshold of ${THRESHOLD}%..."

# Remove the % sign and compare
COVERAGE_NUM=$(echo $TOTAL | sed 's/%//')

if (( $(echo "$COVERAGE_NUM >= $THRESHOLD" | bc -l) )); then
    echo "✓ Coverage threshold met: ${TOTAL} >= ${THRESHOLD}%"
    exit 0
else
    echo "✗ Coverage threshold not met: ${TOTAL} < ${THRESHOLD}%"
    exit 1
fi