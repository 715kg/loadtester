@echo off
echo ========================================
echo    Load Tester - Cross Platform Build
echo ========================================
echo.

REM Получаем версию из git или используем дефолтную
for /f "tokens=*" %%i in ('git describe --tags --always 2^>nul') do set VERSION=%%i
if "%VERSION%"=="" set VERSION=v1.0.0

echo Building version: %VERSION%
echo.

REM Создаем директорию для сборок
if not exist "dist" mkdir dist
if not exist "dist\windows" mkdir dist\windows
if not exist "dist\linux" mkdir dist\linux
if not exist "dist\macos" mkdir dist\macos

REM Получаем дополнительную информацию для сборки
for /f "tokens=*" %%i in ('git rev-parse --short HEAD 2^>nul') do set GIT_COMMIT=%%i
if "%GIT_COMMIT%"=="" set GIT_COMMIT=unknown

set BUILD_USER=%USERNAME%
if "%BUILD_USER%"=="" set BUILD_USER=unknown

set BUILD_TIME=%DATE% %TIME%

REM Устанавливаем расширенные флаги сборки
set LDFLAGS=-ldflags "-s -w -X main.version=%VERSION% -X main.gitCommit=%GIT_COMMIT% -X 'main.buildTime=%BUILD_TIME%' -X main.buildUser=%BUILD_USER%"
set CGO_ENABLED=0

echo Building for Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
go build %LDFLAGS% -o dist\windows\loadtester-%VERSION%-windows-amd64.exe .\cmd\loadtester
if %ERRORLEVEL% neq 0 (
    echo ERROR: Windows build failed
    exit /b 1
)

echo Building for Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
go build %LDFLAGS% -o dist\linux\loadtester-%VERSION%-linux-amd64 .\cmd\loadtester
if %ERRORLEVEL% neq 0 (
    echo ERROR: Linux build failed
    exit /b 1
)

echo Building for Linux (arm64)...
set GOOS=linux
set GOARCH=arm64
go build %LDFLAGS% -o dist\linux\loadtester-%VERSION%-linux-arm64 .\cmd\loadtester
if %ERRORLEVEL% neq 0 (
    echo ERROR: Linux ARM64 build failed
    exit /b 1
)

echo Building for macOS (amd64 - Intel)...
set GOOS=darwin
set GOARCH=amd64
go build %LDFLAGS% -o dist\macos\loadtester-%VERSION%-darwin-amd64 .\cmd\loadtester
if %ERRORLEVEL% neq 0 (
    echo ERROR: macOS Intel build failed
    exit /b 1
)

echo Building for macOS (arm64 - Apple Silicon)...
set GOOS=darwin
set GOARCH=arm64
go build %LDFLAGS% -o dist\macos\loadtester-%VERSION%-darwin-arm64 .\cmd\loadtester
if %ERRORLEVEL% neq 0 (
    echo ERROR: macOS Apple Silicon build failed
    exit /b 1
)

echo.
echo ========================================
echo           Build Complete!
echo ========================================
echo.
echo Built files:
dir /b dist\windows\
dir /b dist\linux\
dir /b dist\macos\
echo.
echo Total size:
for /r dist %%f in (*) do echo %%~nxf - %%~zf bytes

REM Сбрасываем переменные окружения
set GOOS=
set GOARCH=
set CGO_ENABLED=

echo.
echo Build completed successfully!
echo Files are located in the 'dist' directory.