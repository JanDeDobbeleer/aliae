Param
(
    [parameter(Mandatory = $true)]
    [string]
    $Version
)

function Get-HashForArchitecture {
    param (
        [parameter(Mandatory = $true)]
        [string]
        $Architecture
    )
    $hash = Get-Content "../../dist/install-$Architecture.exe.sha256" -Raw
    return $hash.Trim()
}

New-Item -Path "." -Name "dist" -ItemType "directory"

$HashAmd64 = Get-HashForArchitecture -Architecture 'amd64'
$Hash386 = Get-HashForArchitecture -Architecture '386'

$content = Get-Content '.\aliae.json' -Raw
$content = $content.Replace('<VERSION>', $Version)
$content = $content.Replace('<HASH-AMD64>', $HashAmd64)
$content = $content.Replace('<HASH-386>', $Hash386)
$content | Out-File -Encoding 'UTF8' './dist/aliae.json'
