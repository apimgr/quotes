# install-windows.ps1 - Windows installer with NSSM service
# Requires: PowerShell 5.0+ and Administrator privileges

param(
    [switch]$User = $false
)

$ErrorActionPreference = "Stop"

$PROJECTNAME = "quotes"
$PROJECT_DISPLAY = "Quotes"
$GITHUB_REPO = "apimgr/quotes"
$VERSION = "latest"

Write-Host "=== ${PROJECT_DISPLAY} Installer for Windows ===" -ForegroundColor Green

# Detect architecture
$ARCH = if ([Environment]::Is64BitOperatingSystem) {
    if ([Environment]::GetEnvironmentVariable("PROCESSOR_ARCHITEW6432")) { "arm64" } else { "amd64" }
} else { "amd64" }

Write-Host "Architecture: $ARCH"

# Check admin privileges
$isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if ($isAdmin) {
    $BIN_DIR = "C:\Program Files\${PROJECT_DISPLAY}"
    $CONFIG_DIR = "${env:ProgramData}\${PROJECT_DISPLAY}\config"
    $DATA_DIR = "${env:ProgramData}\${PROJECT_DISPLAY}\data"
    $LOG_DIR = "${env:ProgramData}\${PROJECT_DISPLAY}\logs"
} else {
    $BIN_DIR = "${env:LOCALAPPDATA}\${PROJECT_DISPLAY}\bin"
    $CONFIG_DIR = "${env:APPDATA}\${PROJECT_DISPLAY}\config"
    $DATA_DIR = "${env:APPDATA}\${PROJECT_DISPLAY}\data"
    $LOG_DIR = "${env:APPDATA}\${PROJECT_DISPLAY}\logs"
}

Write-Host "Install mode: $(if ($isAdmin) { 'System (Administrator)' } else { 'User' })"

# Create directories
New-Item -ItemType Directory -Force -Path $BIN_DIR | Out-Null
New-Item -ItemType Directory -Force -Path $CONFIG_DIR | Out-Null
New-Item -ItemType Directory -Force -Path $DATA_DIR | Out-Null
New-Item -ItemType Directory -Force -Path "$DATA_DIR\db" | Out-Null
New-Item -ItemType Directory -Force -Path $LOG_DIR | Out-Null

# Download binary
$BINARY_NAME = "${PROJECTNAME}-windows-${ARCH}.exe"
$DOWNLOAD_URL = "https://github.com/${GITHUB_REPO}/releases/${VERSION}/download/${BINARY_NAME}"
$BINARY_PATH = "${BIN_DIR}\${PROJECTNAME}.exe"

Write-Host "Downloading ${BINARY_NAME}..."
Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $BINARY_PATH

Write-Host "✓ Binary installed to ${BINARY_PATH}" -ForegroundColor Green

# Download and install NSSM if admin
if ($isAdmin) {
    $NSSM_URL = "https://nssm.cc/release/nssm-2.24.zip"
    $NSSM_ZIP = "${env:TEMP}\nssm.zip"
    $NSSM_DIR = "${env:TEMP}\nssm-2.24"

    Write-Host "Downloading NSSM..."
    Invoke-WebRequest -Uri $NSSM_URL -OutFile $NSSM_ZIP
    Expand-Archive -Path $NSSM_ZIP -DestinationPath $env:TEMP -Force

    $NSSM_EXE = if ($ARCH -eq "amd64") {
        "${NSSM_DIR}\win64\nssm.exe"
    } else {
        "${NSSM_DIR}\win64\nssm.exe"
    }

    # Install service
    & $NSSM_EXE install $PROJECTNAME $BINARY_PATH
    & $NSSM_EXE set $PROJECTNAME AppDirectory $BIN_DIR
    & $NSSM_EXE set $PROJECTNAME AppEnvironmentExtra "CONFIG_DIR=${CONFIG_DIR}" "DATA_DIR=${DATA_DIR}" "LOGS_DIR=${LOG_DIR}"
    & $NSSM_EXE set $PROJECTNAME DisplayName "${PROJECT_DISPLAY} API"
    & $NSSM_EXE set $PROJECTNAME Description "${PROJECT_DISPLAY} API Server - 27,500 quotes and jokes"
    & $NSSM_EXE set $PROJECTNAME Start SERVICE_AUTO_START
    & $NSSM_EXE start $PROJECTNAME

    Write-Host "✓ Service installed and started" -ForegroundColor Green
    Write-Host "  sc query ${PROJECTNAME}"
    Write-Host "  net start ${PROJECTNAME}"
    Write-Host "  net stop ${PROJECTNAME}"
} else {
    Write-Host "User install: Run manually with: ${BINARY_PATH}" -ForegroundColor Yellow
}

# Add to PATH
if ($isAdmin) {
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
    if ($currentPath -notlike "*${BIN_DIR}*") {
        [Environment]::SetEnvironmentVariable("Path", "$currentPath;$BIN_DIR", "Machine")
        Write-Host "✓ Added to system PATH" -ForegroundColor Green
    }
} else {
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($currentPath -notlike "*${BIN_DIR}*") {
        [Environment]::SetEnvironmentVariable("Path", "$currentPath;$BIN_DIR", "User")
        Write-Host "✓ Added to user PATH" -ForegroundColor Green
    }
}

Write-Host ""
Write-Host "✓ Installation complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Admin credentials: ${CONFIG_DIR}\admin-credentials.txt"
Write-Host "Configuration: ${CONFIG_DIR}\"
Write-Host "Data: ${DATA_DIR}\"
Write-Host "Logs: ${LOG_DIR}\"
