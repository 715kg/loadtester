# Load Tester - Cross Platform Build Script (PowerShell)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "   Load Tester - Cross Platform Build" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host

# Получаем версию из git или используем дефолтную
try {
    $VERSION = git describe --tags --always 2>$null
    if (-not $VERSION) {
        $VERSION = "v1.0.0"
    }
} catch {
    $VERSION = "v1.0.0"
}

Write-Host "Building version: $VERSION" -ForegroundColor Green
Write-Host

# Создаем директории для сборок
$null = New-Item -ItemType Directory -Force -Path "dist"
$null = New-Item -ItemType Directory -Force -Path "dist\windows"
$null = New-Item -ItemType Directory -Force -Path "dist\linux"
$null = New-Item -ItemType Directory -Force -Path "dist\macos"

# Получаем дополнительную информацию для сборки
try {
    $GIT_COMMIT = git rev-parse --short HEAD 2>$null
    if (-not $GIT_COMMIT) { $GIT_COMMIT = "unknown" }
} catch {
    $GIT_COMMIT = "unknown"
}

$BUILD_USER = $env:USERNAME
if (-not $BUILD_USER) { $BUILD_USER = "unknown" }

$BUILD_TIME = Get-Date -Format "yyyy-MM-dd HH:mm:ss"

# Устанавливаем расширенные флаги сборки
$LDFLAGS = "-ldflags `"-s -w -X main.version=$VERSION -X main.gitCommit=$GIT_COMMIT -X 'main.buildTime=$BUILD_TIME' -X main.buildUser=$BUILD_USER`""
$env:CGO_ENABLED = "0"

# Функция для сборки
function Build-Platform {
    param(
        [string]$OS,
        [string]$Arch,
        [string]$OutputPath,
        [string]$Description
    )
    
    Write-Host "Building for $Description..." -ForegroundColor Yellow
    $env:GOOS = $OS
    $env:GOARCH = $Arch
    
    $buildCmd = "go build $LDFLAGS -o `"$OutputPath`" ./cmd/loadtester"
    Invoke-Expression $buildCmd
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: $Description build failed" -ForegroundColor Red
        exit 1
    }
}

# Сборка для всех платформ
Build-Platform "windows" "amd64" "dist\windows\loadtester-$VERSION-windows-amd64.exe" "Windows (amd64)"
Build-Platform "linux" "amd64" "dist\linux\loadtester-$VERSION-linux-amd64" "Linux (amd64)"
Build-Platform "linux" "arm64" "dist\linux\loadtester-$VERSION-linux-arm64" "Linux (arm64)"
Build-Platform "darwin" "amd64" "dist\macos\loadtester-$VERSION-darwin-amd64" "macOS (amd64 - Intel)"
Build-Platform "darwin" "arm64" "dist\macos\loadtester-$VERSION-darwin-arm64" "macOS (arm64 - Apple Silicon)"

Write-Host
Write-Host "========================================" -ForegroundColor Green
Write-Host "          Build Complete!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host

Write-Host "Built files:" -ForegroundColor Cyan
Get-ChildItem -Path "dist" -Recurse -File | ForEach-Object {
    $size = [math]::Round($_.Length / 1MB, 2)
    Write-Host "  $($_.Name) - $size MB" -ForegroundColor White
}

Write-Host
$totalSize = (Get-ChildItem -Path "dist" -Recurse -File | Measure-Object -Property Length -Sum).Sum
$totalSizeMB = [math]::Round($totalSize / 1MB, 2)
Write-Host "Total size: $totalSizeMB MB" -ForegroundColor Green

# Сбрасываем переменные окружения
Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue

Write-Host
Write-Host "Build completed successfully!" -ForegroundColor Green
Write-Host "Files are located in the 'dist' directory." -ForegroundColor Cyan

# Показываем команды для тестирования
Write-Host
Write-Host "Test commands:" -ForegroundColor Yellow
Write-Host "  Windows: .\dist\windows\loadtester-$VERSION-windows-amd64.exe --version" -ForegroundColor Gray
Write-Host "  Linux:   ./dist/linux/loadtester-$VERSION-linux-amd64 --version" -ForegroundColor Gray
Write-Host "  macOS:   ./dist/macos/loadtester-$VERSION-darwin-amd64 --version" -ForegroundColor Gray