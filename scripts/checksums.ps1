# Скрипт для генерации контрольных сумм файлов

Write-Host "Генерация контрольных сумм для Load Tester..." -ForegroundColor Cyan
Write-Host

$checksumFile = "CHECKSUMS.txt"
$distPath = "dist"

if (-not (Test-Path $distPath)) {
    Write-Host "Директория dist не найдена. Сначала выполните сборку." -ForegroundColor Red
    exit 1
}

# Создаем файл с контрольными суммами
$content = @"
# Load Tester - Контрольные суммы (SHA256)
# Дата генерации: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')
# Версия: v1.0.0

"@

Write-Host "Вычисление SHA256 хешей..." -ForegroundColor Yellow

Get-ChildItem -Path $distPath -Recurse -File | ForEach-Object {
    $hash = Get-FileHash -Path $_.FullName -Algorithm SHA256
    $relativePath = $_.FullName.Replace((Get-Location).Path + "\", "")
    $size = [math]::Round($_.Length / 1MB, 2)
    
    Write-Host "  $($_.Name) - $size MB" -ForegroundColor White
    
    $content += "$($hash.Hash.ToLower())  $relativePath`n"
}

# Сохраняем контрольные суммы
$content | Out-File -FilePath $checksumFile -Encoding UTF8

Write-Host
Write-Host "Контрольные суммы сохранены в файл: $checksumFile" -ForegroundColor Green

# Показываем содержимое файла
Write-Host
Write-Host "Содержимое файла контрольных сумм:" -ForegroundColor Cyan
Get-Content $checksumFile

Write-Host
Write-Host "Для проверки целостности используйте:" -ForegroundColor Yellow
Write-Host "  Windows: certutil -hashfile <файл> SHA256" -ForegroundColor Gray
Write-Host "  Linux:   sha256sum <файл>" -ForegroundColor Gray
Write-Host "  macOS:   shasum -a 256 <файл>" -ForegroundColor Gray