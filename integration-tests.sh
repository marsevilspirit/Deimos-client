#!/bin/bash
set -e

# Function to clean up resources
cleanup() {
  echo "Stopping and removing Deimos cluster..."
  docker compose down
}

# Trap EXIT signal to ensure cleanup is always performed
trap cleanup EXIT

echo "Starting Deimos cluster..."
docker compose up -d

echo "Waiting for Deimos cluster to start..."
for i in {1..30}; do
  if curl -s -f http://127.0.0.1:4001/machines && \
     curl -s -f http://127.0.0.1:4002/machines && \
     curl -s -f http://127.0.0.1:4003/machines; then
    echo "Deimos cluster is up!"
    break
  fi
  sleep 2
done

if [ $i -eq 30 ]; then
  echo "Deimos cluster failed to start"
  docker compose logs
  exit 1
fi

echo "Running integration tests..."

# Track test results
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

for dir in example/*/; do
  if [ -f "${dir}main.go" ]; then
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo "==> Running test in ${dir}"
    
    # Run the test and capture exit code
    if (cd "${dir}" && go run main.go); then
      echo "<== ✓ Test PASSED in ${dir}"
      PASSED_TESTS=$((PASSED_TESTS + 1))
    else
      echo "<== ✗ Test FAILED in ${dir}"
      FAILED_TESTS=$((FAILED_TESTS + 1))
      echo "ERROR: Integration test failed in ${dir}"
      echo "Test Results Summary:"
      echo "  Total: ${TOTAL_TESTS}"
      echo "  Passed: ${PASSED_TESTS}"
      echo "  Failed: ${FAILED_TESTS}"
      exit 1
    fi
  fi
done

echo ""
echo "=== Integration Test Results ==="
echo "Total tests: ${TOTAL_TESTS}"
echo "Passed: ${PASSED_TESTS}"
echo "Failed: ${FAILED_TESTS}"

if [ ${FAILED_TESTS} -eq 0 ]; then
  echo "✓ All integration tests passed!"
  exit 0
else
  echo "✗ Some integration tests failed!"
  exit 1
fi