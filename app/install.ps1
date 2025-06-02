#!/usr/bin/env pwsh
param(
    [Switch]$NoPathUpdate = $false 
)

$ErrorActionPreference = "Stop"

$RepoOwner = "Zusty"
$RepoName = "GCT"
$BaseUrl = "https://codeberg.org/$RepoOwner/$RepoName/releases/download/latest"
$InstallDir = Join-Path $env:LOCALAPPDATA "GCT" 
$BinName = "gct.exe"                             

$Os = "windows"
$Arch = ""
$SystemType = (Get-CimInstance Win32_ComputerSystem).SystemType
if ($SystemType -match "x64-based") {
    $Arch = "amd64"
}
elseif ($SystemType -match "ARM64") {
    $Arch = "arm64" 
}
else {
    Write-Error "Install Failed: GCT currently requires a 64-bit (x64 or ARM64) Windows system. Detected: $SystemType"
    exit 1
}

$TargetBin = "gct-${Os}-${Arch}.zip"
$DownloadUrl = "$BaseUrl/$TargetBin"
$OutputPath = Join-Path $InstallDir $BinName
$TempZipPath = Join-Path $env:TEMP $TargetBin

Write-Host "Installing/Updating GCT for $Os ($Arch)..."

if (-not (Test-Path $InstallDir)) {
    Write-Host "Creating installation directory: $InstallDir"
    New-Item -Path $InstallDir -ItemType Directory -Force | Out-Null
}

Write-Host "Downloading GCT from: $DownloadUrl"
try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempZipPath -UseBasicParsing

    Write-Host "Downloaded successfully to: $TempZipPath"
}
catch {
    Write-Error "Install Failed: Could not download GCT from $DownloadUrl"
    Write-Error $_.Exception.Message
    if (Test-Path $TempZipPath) { Remove-Item $TempZipPath -Force }
    exit 1
}

if ((Get-Item $TempZipPath).Length -lt 1KB) {
    Write-Error "Install Failed: Downloaded file seems too small or corrupted."
    Remove-Item $TempZipPath -Force
    exit 1
}

if (Test-Path $OutputPath) {
    Write-Host "Removing existing binary at $OutputPath..."
    Remove-Item $OutputPath -Force -ErrorAction SilentlyContinue | Out-Null
}

Write-Host "Extracting archive..."
try {
    Expand-Archive -Path $TempZipPath -DestinationPath $InstallDir -Force

    $ExtractedExe = Get-ChildItem -Path $InstallDir -Filter "gct-${Os}-${Arch}.exe" -ErrorAction SilentlyContinue
    if ($ExtractedExe) {
        Write-Host "Renaming extracted binary to $BinName..."
        Move-Item -Path $ExtractedExe.FullName -Destination $OutputPath -Force
    }
    else {
        $AnyExe = Get-ChildItem -Path $InstallDir -Filter "*.exe" -ErrorAction SilentlyContinue
        if ($AnyExe) {
            Write-Host "Renaming extracted binary to $BinName..."
            Move-Item -Path $AnyExe.FullName -Destination $OutputPath -Force
        }
        else {
            throw "Could not find executable in extracted archive"
        }
    }

    Write-Host "Extraction successful to $InstallDir."
}
catch {
    Write-Error "Install Failed: Could not extract archive $TempZipPath"
    Write-Error $_.Exception.Message
    if (Test-Path $TempZipPath) { Remove-Item $TempZipPath -Force }
    exit 1
}
finally {
    if (Test-Path $TempZipPath) { Remove-Item $TempZipPath -Force }
}

if (-not $NoPathUpdate) {
    Write-Host "Checking user PATH environment variable..."
    try {
        $UserPath = [Environment]::GetEnvironmentVariable('Path', 'User')
        if ($UserPath -notlike "*$InstallDir*") {
            Write-Host "Adding '$InstallDir' to user PATH..."
            $Separator = ""
            if (-not ([string]::IsNullOrEmpty($UserPath)) -and (-not $UserPath.EndsWith(";"))) {
                $Separator = ";"
            }
            $NewPath = $UserPath + $Separator + $InstallDir
            [Environment]::SetEnvironmentVariable('Path', $NewPath, 'User')
            Write-Host "PATH updated. You need to restart your terminal for the change to take effect."
        }
        else {
            Write-Host "'$InstallDir' is already in the user PATH."
        }
    }
    catch {
        Write-Warning "Could not automatically update user PATH. Error: $($_.Exception.Message)"
        Write-Warning "Please add '$InstallDir' to your PATH manually."
    }
}
else {
    Write-Host "Skipping PATH update as requested. Add '$InstallDir' to your PATH manually."
}

Write-Host ""
Write-Host "GCT ($(Split-Path $TargetBin -Leaf)) installed/updated successfully to: $InstallDir"
Write-Host "Run 'gct --version' in a *new* terminal window to verify."
