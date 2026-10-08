$ErrorActionPreference = 'Stop'

$root = Split-Path $PSScriptRoot -Parent
$src = Join-Path $root 'src'
$buildDir = Join-Path $root 'build'
$palette = Get-Content (Join-Path $root 'palette.json') -Raw | ConvertFrom-Json -AsHashtable
$version = ([xml](Get-Content (Join-Path $src 'META-INF/plugin.xml'))).'idea-plugin'.version
$out = Join-Path $buildDir "gopher-burrow-$version.jar"

New-Item -ItemType Directory -Force $buildDir | Out-Null
Remove-Item (Join-Path $buildDir '*.jar')

# $name placeholders in theme files are replaced with hex values from palette.json
function Expand-Palette([string]$text, [string]$file) {
    [regex]::Replace($text, '\$([a-z][a-z0-9-]*)', {
        param($m)
        $name = $m.Groups[1].Value
        if (-not $palette.ContainsKey($name)) { throw "Unknown palette color '$name' in $file" }
        $palette[$name]
    })
}

Add-Type -AssemblyName System.IO.Compression.FileSystem
$zip = [System.IO.Compression.ZipFile]::Open($out, 'Create')
try {
    Get-ChildItem $src -Recurse -File | ForEach-Object {
        # Entries need forward slashes, otherwise the IDE can't find META-INF/plugin.xml
        $entryName = $_.FullName.Substring($src.Length + 1).Replace('\', '/')
        $entry = $zip.CreateEntry($entryName)
        $stream = $entry.Open()
        if ($_.Extension -in '.xml', '.json') {
            $bytes = [System.Text.UTF8Encoding]::new($false).GetBytes((Expand-Palette (Get-Content $_.FullName -Raw) $entryName))
        } else {
            $bytes = [System.IO.File]::ReadAllBytes($_.FullName)
        }
        $stream.Write($bytes, 0, $bytes.Length)
        $stream.Dispose()
    }
} finally {
    $zip.Dispose()
}

Write-Host "Built $out"
