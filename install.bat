@echo off
REM tkeybr installer for Windows CMD
REM Usage: Download and run this script, or:
REM curl -fsSL https://raw.githubusercontent.com/farseenmanekhan1232/tkeybr/main/install.bat -o install.bat && install.bat

setlocal enabledelayedexpansion

set REPO=farseenmanekhan1232/tkeybr
set BINARY_NAME=tkeybr
set INSTALL_DIR=%LOCALAPPDATA%\tkeybr

echo.
echo   ========================================
echo        tkeybr Installer
echo    Terminal-based typing practice
echo   ========================================
echo.

REM Detect architecture
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    set ARCH=amd64
) else if "%PROCESSOR_ARCHITECTURE%"=="x86" (
    set ARCH=386
) else (
    echo [ERROR] Unsupported architecture: %PROCESSOR_ARCHITECTURE%
    exit /b 1
)

echo [INFO] Detected: Windows/%ARCH%

REM Get latest version using PowerShell (required for HTTPS)
echo [INFO] Fetching latest release...
for /f "delims=" %%i in ('powershell -Command "(Invoke-RestMethod -Uri 'https://api.github.com/repos/%REPO%/releases/latest').tag_name"') do set VERSION=%%i

if "%VERSION%"=="" (
    echo [ERROR] Failed to get latest version
    exit /b 1
)

echo [INFO] Latest version: %VERSION%

REM Construct download URL - GoReleaser format: tkeybr_0.5.0_windows_amd64.zip
set VERSION_NUM=%VERSION:~1%
set FILENAME=%BINARY_NAME%_%VERSION_NUM%_windows_%ARCH%.zip
set DOWNLOAD_URL=https://github.com/%REPO%/releases/download/%VERSION%/%FILENAME%

REM Create install directory
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

REM Download using PowerShell
echo [INFO] Downloading %FILENAME%...
set TMP_FILE=%TEMP%\%FILENAME%
powershell -Command "Invoke-WebRequest -Uri '%DOWNLOAD_URL%' -OutFile '%TMP_FILE%'"

if not exist "%TMP_FILE%" (
    echo [ERROR] Download failed
    exit /b 1
)

REM Extract using PowerShell
echo [INFO] Extracting...
powershell -Command "Expand-Archive -Path '%TMP_FILE%' -DestinationPath '%INSTALL_DIR%' -Force"
del "%TMP_FILE%"

REM Add to PATH if not already there
echo %PATH% | find /i "%INSTALL_DIR%" >nul
if errorlevel 1 (
    echo [INFO] Adding to PATH...
    setx PATH "%PATH%;%INSTALL_DIR%" >nul 2>&1
    set PATH=%PATH%;%INSTALL_DIR%
)

echo.
echo [INFO] Installation complete!
echo.
echo   Run 'tkeybr' to start practicing!
echo   (You may need to restart your terminal)
echo.

"%INSTALL_DIR%\%BINARY_NAME%.exe" --version 2>nul

endlocal
