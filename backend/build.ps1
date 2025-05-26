# PowerShell Build Script
# Exit immediately if a command fails
$ErrorActionPreference = "Stop"

# --- Configuration Variables ---
# Path to the main Go package
$MAIN_PACKAGE_PATH = "./cmd/protection/main.go"

# Directory where the build output will be stored
$BUILD_DIR = "./build"

# Name of the final executable file
$OUTPUT_NAME = "backend"

# Target Operating System and Architecture
$TARGET_OS = "linux"   # Routers typically run Linux. Change if needed.
$TARGET_ARCH = "arm64"

# --- Create Build Directory ---
Write-Host "Creating build directory: $BUILD_DIR" -ForegroundColor Green
if (!(Test-Path -Path $BUILD_DIR)) {
    New-Item -ItemType Directory -Path $BUILD_DIR -Force | Out-Null
}

# --- Go Module Path Check (Optional, adjust if go.mod is inside the backend folder) ---
# If your go.mod file is in the project root, this section is not needed.
# If go.mod is inside the 'backend' folder, you might need to 'Set-Location backend' before building.
# This script assumes go.mod is in the project root.

# --- Start Build ---
Write-Host ""
Write-Host "Starting build for ARM64..." -ForegroundColor Yellow
Write-Host "Target OS: $TARGET_OS" -ForegroundColor Cyan
Write-Host "Target Architecture: $TARGET_ARCH" -ForegroundColor Cyan
Write-Host "Main Package: $MAIN_PACKAGE_PATH" -ForegroundColor Cyan
Write-Host "Output File: $BUILD_DIR/$OUTPUT_NAME" -ForegroundColor Cyan

# Set environment variables for Go cross-compilation
$env:GOOS = $TARGET_OS
$env:GOARCH = $TARGET_ARCH
$env:CGO_ENABLED = "0"

# CGO_ENABLED=0 is useful for creating static binaries without C library dependencies.
# This can be beneficial for router environments to avoid dynamic library issues.
# If your project has CGO dependencies (e.g., some networking libraries),
# you might need to adjust this and ensure you have the correct C cross-compiler toolchain.
# -ldflags="-s -w" strips debugging symbols and DWARF information to reduce binary size.

try {
    go build -ldflags="-s -w" -o "$BUILD_DIR/$OUTPUT_NAME" "$MAIN_PACKAGE_PATH"
}
catch {
    Write-Error "Build failed: $_"
    exit 1
}

# --- Build Complete ---
Write-Host ""
Write-Host "Build complete!" -ForegroundColor Green
Write-Host "Executable: $BUILD_DIR/$OUTPUT_NAME" -ForegroundColor Green

# --- Display File Information (Optional) ---
$outputPath = Join-Path $BUILD_DIR $OUTPUT_NAME

if (Test-Path $outputPath) {
    Write-Host ""
    Write-Host "File information:" -ForegroundColor Yellow

    # Get file information
    $fileInfo = Get-Item $outputPath
    Write-Host "File size: $([math]::Round($fileInfo.Length / 1MB, 2)) MB ($($fileInfo.Length) bytes)"
    Write-Host "Created: $($fileInfo.CreationTime)"
    Write-Host "Modified: $($fileInfo.LastWriteTime)"

    # Check if it's executable (on Windows, this is less relevant, but we can check the extension)
    if ($TARGET_OS -eq "windows" -and $OUTPUT_NAME -notlike "*.exe") {
        Write-Host "Note: For Windows targets, consider adding .exe extension" -ForegroundColor Yellow
    }
} else {
    Write-Warning "Output file not found at: $outputPath"
}

Write-Host ""
Write-Host "Script finished." -ForegroundColor Green

# Reset environment variables (optional, as they're only set for this session)
Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue