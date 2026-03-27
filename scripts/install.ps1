# Podutil Installation Script for Windows
# Usage: .\install.ps1 [-Version <version>]

param(
    [string]$Version = "latest",
    [string]$InstallDir = "$env:LOCALAPPDATA\podutil"
)

$ErrorActionPreference = "Stop"
$Repo = "XantX/podman-utils"

function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    if ($arch -eq "AMD64") { return "amd64" }
    if ($arch -eq "ARM64") { return "arm64" }
    throw "Arquitectura no soportada: $arch"
}

function Install-Podutil {
    $os = "windows"
    $arch = Get-Architecture

    Write-Host "Descargando podutil v$Version para $os/$arch..."

    $format = "zip"
    $downloadUrl = "https://github.com/$Repo/releases/$Version/download/podutil_${Version}_${os}_${arch}.$format"

    if ($Version -eq "latest") {
        $apiUrl = "https://api.github.com/repos/$Repo/releases/latest"
        $response = Invoke-RestMethod -Uri $apiUrl
        $asset = $response.assets | Where-Object { $_.name -match "windows.*$arch.*\.zip$" }
        if ($asset) {
            $downloadUrl = $asset.browser_download_url
        }
    }

    $tempDir = [System.IO.Path]::GetTempPathName()
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null

    $downloadPath = Join-Path $tempDir "podutil.zip"
    Write-Host "Descargando desde: $downloadUrl"
    Invoke-WebRequest -Uri $downloadUrl -OutFile $downloadPath -UseBasicParsing

    Expand-Archive -Path $downloadPath -DestinationPath $tempDir -Force

    $installDir = $InstallDir
    if (-not (Test-Path $installDir)) {
        New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    }

    $exePath = Join-Path $tempDir "podutil.exe"
    if (Test-Path $exePath) {
        Copy-Item -Path $exePath -Destination (Join-Path $installDir "podutil.exe") -Force
    }

    Remove-Item -Path $tempDir -Recurse -Force

    Write-Host "`nInstalado en: $installDir\podutil.exe" -ForegroundColor Green
    Write-Host "`nAgrega al PATH:" -ForegroundColor Yellow
    Write-Host "  `$env:PATH += `";$installDir`""
    Write-Host "`nO para hacer permanente:" -ForegroundColor Yellow
    Write-Host "  [Environment]::SetEnvironmentVariable(`"PATH`", `$env:PATH + `";$installDir`", `"User`")"
}

Install-Podutil

Write-Host "`nInstalación completada!" -ForegroundColor Green