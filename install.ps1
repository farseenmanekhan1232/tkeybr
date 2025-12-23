# tkeybr installer for Windows
# Usage: iwr -useb https://raw.githubusercontent.com/farseenmanekhan1232/tkeybr/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

$REPO = "farseenmanekhan1232/tkeybr"
$BINARY_NAME = "tkeybr"

Write-Host ""
Write-Host "  ╔═══════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "  ║        tkeybr Installer               ║" -ForegroundColor Cyan
Write-Host "  ║   Terminal-based typing practice      ║" -ForegroundColor Cyan
Write-Host "  ╚═══════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# Detect architecture
$ARCH = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
Write-Host "[INFO] Detected: windows/$ARCH" -ForegroundColor Green

# Get latest version
Write-Host "[INFO] Fetching latest release..." -ForegroundColor Green
$release = Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest"
$VERSION = $release.tag_name
Write-Host "[INFO] Latest version: $VERSION" -ForegroundColor Green

# Download URL - GoReleaser format: tkeybr_0.5.0_windows_amd64.zip
$VERSION_NUM = $VERSION.TrimStart('v')
$FILENAME = "${BINARY_NAME}_${VERSION_NUM}_windows_${ARCH}.zip"
$DOWNLOAD_URL = "https://github.com/$REPO/releases/download/$VERSION/$FILENAME"

# Install directory
$INSTALL_DIR = "$env:LOCALAPPDATA\tkeybr"
if (-not (Test-Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
}

# Download
Write-Host "[INFO] Downloading $FILENAME..." -ForegroundColor Green
$TMP_FILE = "$env:TEMP\$FILENAME"
Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $TMP_FILE

# Extract
Write-Host "[INFO] Extracting..." -ForegroundColor Green
Expand-Archive -Path $TMP_FILE -DestinationPath $INSTALL_DIR -Force
Remove-Item $TMP_FILE

# Add to PATH if not already
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$INSTALL_DIR*") {
    Write-Host "[INFO] Adding to PATH..." -ForegroundColor Green
    [Environment]::SetEnvironmentVariable("Path", "$currentPath;$INSTALL_DIR", "User")
    $env:Path = "$env:Path;$INSTALL_DIR"
}

# Verify
Write-Host ""
Write-Host "[INFO] ✅ Installation complete!" -ForegroundColor Green
Write-Host ""
Write-Host "  Run 'tkeybr' to start practicing!" -ForegroundColor White
Write-Host "  (You may need to restart your terminal)" -ForegroundColor Yellow
Write-Host ""

& "$INSTALL_DIR\$BINARY_NAME.exe" --version 2>$null
