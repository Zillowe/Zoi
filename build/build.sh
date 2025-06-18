#!/bin/bash
set -e

CYAN='\033[0;36m'
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

mkdir -p build/compiled

COMMIT=$(git rev-parse --short=10 HEAD 2>/dev/null || echo "dev")

echo -e "${CYAN}Building GCT for $(uname -s)...${NC}"
go build -o "./build/compiled/gct" \
    -ldflags "-X main.VerCommit=$COMMIT" \
    ./src

if [ $? -eq 0 ]; then
    chmod +x ./build/compiled/gct
    echo -e "${GREEN}Build successful! Commit: $COMMIT${NC}"
else
    echo -e "${RED}Build failed${NC}"
    exit 1
fi