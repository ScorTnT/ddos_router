#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# --- Configuration Variables ---
# Path to the main Go package
MAIN_PACKAGE_PATH="./cmd/protection/main.go"

# Directory where the build output will be stored
BUILD_DIR="./build"

# Name of the final executable file
OUTPUT_NAME="backend"

# Target Operating System and Architecture
TARGET_OS="linux"   # Routers typically run Linux. Change if needed.
TARGET_ARCH="arm64"

# --- Create Build Directory ---
echo "Creating build directory: ${BUILD_DIR}"
mkdir -p "${BUILD_DIR}"

# --- Go Module Path Check (Optional, adjust if go.mod is inside the backend folder) ---
# If your go.mod file is in the project root, this section is not needed.
# If go.mod is inside the 'backend' folder, you might need to 'cd backend' before building.
# This script assumes go.mod is in the project root.

# --- Start Build ---
echo "Starting build for ARM64..."
echo "Target OS: ${TARGET_OS}"
echo "Target Architecture: ${TARGET_ARCH}"
echo "Main Package: ${MAIN_PACKAGE_PATH}"
echo "Output File: ${BUILD_DIR}/${OUTPUT_NAME}"

# CGO_ENABLED=0 is useful for creating static binaries without C library dependencies.
# This can be beneficial for router environments to avoid dynamic library issues.
# If your project has CGO dependencies (e.g., some networking libraries),
# you might need to adjust this and ensure you have the correct C cross-compiler toolchain.
# -ldflags="-s -w" strips debugging symbols and DWARF information to reduce binary size.
env GOOS=${TARGET_OS} GOARCH=${TARGET_ARCH} CGO_ENABLED=0 go build -ldflags="-s -w" -o "${BUILD_DIR}/${OUTPUT_NAME}" "${MAIN_PACKAGE_PATH}"

# --- Build Complete ---
echo ""
echo "Build complete!"
echo "Executable: ${BUILD_DIR}/${OUTPUT_NAME}"

# --- Display File Information (Optional) ---
if command -v file &> /dev/null
then
    echo "File information:"
    file "${BUILD_DIR}/${OUTPUT_NAME}"
fi

if command -v ls &> /dev/null
then
    echo "File size:"
    ls -lh "${BUILD_DIR}/${OUTPUT_NAME}"
fi

echo ""
echo "Script finished."