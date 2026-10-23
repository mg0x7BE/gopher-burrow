# Prints the 16 ANSI colors as text and as backgrounds, plus typical messages
$e = [char]27
$names = 'black', 'red', 'green', 'yellow', 'blue', 'magenta', 'cyan', 'white'

foreach ($bright in $false, $true) {
    $text = ''
    $back = ''
    for ($i = 0; $i -lt 8; $i++) {
        $fg = if ($bright) { 90 + $i } else { 30 + $i }
        $label = if ($bright) { "br-$($names[$i])" } else { $names[$i] }
        $text += "${e}[${fg}m{0,-12}${e}[0m" -f $label
        $back += "${e}[$($fg + 10)m{0,-12}${e}[0m" -f $label
    }
    Write-Host $text
    Write-Host $back
}

Write-Host ''
Write-Host 'Typical output:'
Write-Host 'ok   inventory/store   0.42s' -ForegroundColor Green
Write-Host 'FAIL inventory/api     1.07s' -ForegroundColor Red
Write-Host 'warning: deprecated flag --legacy' -ForegroundColor Yellow
Write-Host 'info: listening on :8080' -ForegroundColor Cyan
Write-Host 'debug: cache miss for key items/42' -ForegroundColor DarkGray
Write-Host "${e}[1mbold${e}[0m ${e}[2mdim${e}[0m ${e}[3mitalic${e}[0m ${e}[4munderline${e}[0m ${e}[7minverse${e}[0m"
