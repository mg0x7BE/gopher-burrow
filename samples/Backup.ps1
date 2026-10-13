<#
.SYNOPSIS
    Copies files older than a given age to the backup drive.
#>
param(
    [Parameter(Mandatory)]
    [string]$Source,
    [string]$Destination = 'E:\Backup',
    [int]$MaxAgeInDays = 7
)

$ErrorActionPreference = 'Stop'
$threshold = [DateTime]::Today.AddDays(-$MaxAgeInDays)

function Write-Log([string]$Message) {
    Write-Host "[$(Get-Date -Format 'HH:mm:ss')] $Message"
}

$files = Get-ChildItem -Path $Source -File -Recurse |
    Where-Object { $_.LastWriteTime -lt $threshold }

foreach ($file in $files) {
    $target = Join-Path $Destination $file.Name
    if (-not (Test-Path $target)) {
        Copy-Item -Path $file.FullName -Destination $target
        Write-Log "copied $($file.Name) ($([math]::Round($file.Length / 1KB, 1)) KB)"
    }
}

Write-Log "done, $($files.Count) files checked"
