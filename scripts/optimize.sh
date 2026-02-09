#!/bin/bash
# Binary optimization script for search CLI
# This script builds and compresses the binary using various optimization techniques

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Project info
VERSION=${VERSION:-"0.1.0"}
BINARY_NAME="search"
BUILD_DIR="bin"
ORIGINAL_SIZE=0
OPTIMIZED_SIZE=0
COMPRESSED_SIZE=0

echo -e "${BLUE}=== Search CLI Binary Optimization ===${NC}"
echo ""

# Function to print file size
print_size() {
    local file=$1
    local label=$2
    if [ -f "$file" ]; then
        local size=$(ls -lh "$file" | awk '{print $5}')
        local bytes=$(stat -f%z "$file" 2>/dev/null || stat -c%s "$file" 2>/dev/null)
        echo -e "  ${label}: ${GREEN}${size}${NC} (${bytes} bytes)"
        echo "$bytes"
    else
        echo -e "  ${label}: ${RED}File not found${NC}"
        echo "0"
    fi
}

# Check if upx is installed
check_upx() {
    if ! command -v upx &> /dev/null; then
        echo -e "${YELLOW}Warning: UPX not found${NC}"
        echo ""
        echo "UPX is required for binary compression."
        echo "Install with:"
        echo "  - Debian/Ubuntu: sudo apt-get install upx"
        echo "  - macOS: brew install upx"
        echo "  - Fedora/RHEL: sudo dnf install upx"
        echo "  - Arch: sudo pacman -S upx"
        echo ""
        return 1
    fi
    return 0
}

# Create build directory
mkdir -p "$BUILD_DIR"

# Step 1: Build standard binary
echo -e "${BLUE}Step 1: Building standard binary...${NC}"
go build -ldflags="-X github.com/mule-ai/search/pkg/version.Version=$VERSION" \
    -o "$BUILD_DIR/${BINARY_NAME}-standard" ./cmd/search

ORIGINAL_SIZE=$(print_size "$BUILD_DIR/${BINARY_NAME}-standard" "Standard build")
echo ""

# Step 2: Build optimized binary (with -ldflags "-s -w" and -trimpath)
echo -e "${BLUE}Step 2: Building optimized binary (-trimpath -ldflags \"-s -w\")...${NC}"
go build -trimpath \
    -ldflags="-s -w -X github.com/mule-ai/search/pkg/version.Version=$VERSION" \
    -o "$BUILD_DIR/${BINARY_NAME}-optimized" ./cmd/search

OPTIMIZED_SIZE=$(print_size "$BUILD_DIR/${BINARY_NAME}-optimized" "Optimized build")

# Calculate size reduction
if [ "$ORIGINAL_SIZE" -gt 0 ] && [ "$OPTIMIZED_SIZE" -gt 0 ]; then
    reduction=$((ORIGINAL_SIZE - OPTIMIZED_SIZE))
    percent=$((reduction * 100 / ORIGINAL_SIZE))
    echo -e "  ${GREEN}Reduced by: ${reduction} bytes (${percent}%)${NC}"
fi
echo ""

# Step 3: Compress with UPX
if check_upx; then
    echo -e "${BLUE}Step 3: Compressing with UPX (--best --lzma)...${NC}"
    
    # Copy optimized binary for compression
    cp "$BUILD_DIR/${BINARY_NAME}-optimized" "$BUILD_DIR/${BINARY_NAME}"
    
    # Compress with UPX
    upx --best --lzma --force --no-color "$BUILD_DIR/${BINARY_NAME}" 2>/dev/null || true
    
    COMPRESSED_SIZE=$(print_size "$BUILD_DIR/${BINARY_NAME}" "UPX compressed")
    
    # Calculate compression ratio
    if [ "$OPTIMIZED_SIZE" -gt 0 ] && [ "$COMPRESSED_SIZE" -gt 0 ]; then
        compression=$((OPTIMIZED_SIZE - COMPRESSED_SIZE))
        percent=$((compression * 100 / OPTIMIZED_SIZE))
        echo -e "  ${GREEN}Compressed by: ${compression} bytes (${percent}%)${NC}"
    fi
    
    # Calculate total reduction
    if [ "$ORIGINAL_SIZE" -gt 0 ] && [ "$COMPRESSED_SIZE" -gt 0 ]; then
        total_reduction=$((ORIGINAL_SIZE - COMPRESSED_SIZE))
        total_percent=$((total_reduction * 100 / ORIGINAL_SIZE))
        echo -e "  ${GREEN}Total reduction: ${total_reduction} bytes (${total_percent}%)${NC}"
    fi
    echo ""
fi

# Step 4: Show comparison table
echo -e "${BLUE}=== Size Comparison ===${NC}"
echo ""
printf "%-20s %-15s\n" "Build Type" "Size"
printf "%-20s %-15s\n" "--------------------" "---------------"
if [ "$ORIGINAL_SIZE" -gt 0 ]; then
    printf "%-20s %-15s\n" "Standard" "$(ls -lh "$BUILD_DIR/${BINARY_NAME}-standard" | awk '{print $5}')"
fi
if [ "$OPTIMIZED_SIZE" -gt 0 ]; then
    printf "%-20s %-15s\n" "Optimized" "$(ls -lh "$BUILD_DIR/${BINARY_NAME}-optimized" | awk '{print $5}')"
fi
if [ "$COMPRESSED_SIZE" -gt 0 ]; then
    printf "%-20s %-15s\n" "UPX Compressed" "$(ls -lh "$BUILD_DIR/${BINARY_NAME}" | awk '{print $5}')"
fi
echo ""

# Step 5: Test the compressed binary
if [ -f "$BUILD_DIR/${BINARY_NAME}" ]; then
    echo -e "${BLUE}Step 4: Testing compressed binary...${NC}"
    if "$BUILD_DIR/${BINARY_NAME}" --version &> /dev/null; then
        echo -e "  ${GREEN}✓ Binary works correctly${NC}"
        "$BUILD_DIR/${BINARY_NAME}" --version
    else
        echo -e "  ${RED}✗ Binary test failed${NC}"
        exit 1
    fi
    echo ""
fi

# Step 6: Recommendations
echo -e "${BLUE}=== Recommendations ===${NC}"
echo ""
echo "For development:"
echo "  - Use standard build: go build -o bin/search ./cmd/search"
echo ""
echo "For distribution:"
echo "  - Use optimized build: go build -trimpath -ldflags=\"-s -w\" -o bin/search ./cmd/search"
echo "  - For even smaller size: upx --best --lzma bin/search"
echo ""
echo "For releases:"
echo "  - Use GoReleaser (includes UPX compression automatically)"
echo "  - make snapshot  # Test release build"
echo ""

echo -e "${GREEN}=== Optimization Complete ===${NC}"
echo ""
echo "Binaries available in: $BUILD_DIR/"
echo "  - ${BINARY_NAME}-standard"
echo "  - ${BINARY_NAME}-optimized"
if [ -f "$BUILD_DIR/${BINARY_NAME}" ]; then
    echo "  - ${BINARY_NAME} (UPX compressed)"
fi
