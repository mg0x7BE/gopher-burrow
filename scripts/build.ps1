$ErrorActionPreference = 'Stop'

$root = Split-Path $PSScriptRoot -Parent
$src = Join-Path $root 'src'
$buildDir = Join-Path $root 'build'
$version = ([xml](Get-Content (Join-Path $src 'META-INF/plugin.xml'))).'idea-plugin'.version
$out = Join-Path $buildDir "gopher-burrow-$version.jar"

New-Item -ItemType Directory -Force $buildDir | Out-Null
Remove-Item (Join-Path $buildDir '*.jar')

# Entries need forward slashes, otherwise the IDE can't find META-INF/plugin.xml
Add-Type -AssemblyName System.IO.Compression.FileSystem
$zip = [System.IO.Compression.ZipFile]::Open($out, 'Create')
Get-ChildItem $src -Recurse -File | ForEach-Object {
    $entry = $_.FullName.Substring($src.Length + 1).Replace('\', '/')
    [System.IO.Compression.ZipFileExtensions]::CreateEntryFromFile($zip, $_.FullName, $entry) | Out-Null
}
$zip.Dispose()

Write-Host "Built $out"
