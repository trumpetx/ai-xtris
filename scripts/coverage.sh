#!/bin/bash

# Coverage tracking script for Xtris Clone
# This script generates coverage reports and badges

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}📊 Generating Test Coverage Report${NC}"

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...

# Get coverage percentage
COVERAGE=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')

echo -e "${BLUE}📈 Coverage: ${COVERAGE}%${NC}"

# Generate detailed report
echo -e "\n${BLUE}📋 Detailed Coverage Report:${NC}"
go tool cover -func=coverage.out

# Generate HTML report
echo -e "\n${BLUE}🌐 Generating HTML Coverage Report...${NC}"
go tool cover -html=coverage.out -o coverage.html
echo -e "${GREEN}✅ HTML report generated: coverage.html${NC}"

# Check coverage threshold
THRESHOLD=80
if (( $(echo "$COVERAGE >= $THRESHOLD" | bc -l) )); then
    echo -e "${GREEN}✅ Coverage meets ${THRESHOLD}% threshold${NC}"
    BADGE_COLOR="green"
else
    echo -e "${YELLOW}⚠️  Coverage below ${THRESHOLD}% threshold${NC}"
    BADGE_COLOR="orange"
fi

# Generate coverage badge
echo -e "\n${BLUE}🏷️  Generating Coverage Badge...${NC}"
BADGE_URL="https://img.shields.io/badge/coverage-${COVERAGE}%25-${BADGE_COLOR}?style=flat-square&logo=go"
echo "Coverage Badge URL: $BADGE_URL"

# Save badge URL to file for CI/CD
echo "$BADGE_URL" > .coverage-badge-url

# Generate summary for README
cat > .coverage-summary << EOF
# Test Coverage Summary

- **Current Coverage**: ${COVERAGE}%
- **Target Coverage**: ${THRESHOLD}%
- **Status**: $([ $(echo "$COVERAGE >= $THRESHOLD" | bc -l) -eq 1 ] && echo "✅ PASS" || echo "⚠️  BELOW THRESHOLD")

## Coverage by Function

$(go tool cover -func=coverage.out | grep -E '^xtris-clone' | sed 's/^xtris-clone\/main.go:[0-9]*:\s*//' | sed 's/\s\+/ | /')

## Generated Reports

- **HTML Report**: [coverage.html](coverage.html)
- **Coverage Badge**: ![Coverage]($BADGE_URL)
EOF

echo -e "\n${GREEN}✅ Coverage analysis complete!${NC}"
echo -e "${BLUE}📄 Summary saved to: .coverage-summary${NC}"
echo -e "${BLUE}🏷️  Badge URL saved to: .coverage-badge-url${NC}" 