<#
.SYNOPSIS
    Builds the Inno Setup installer for aliae.

.DESCRIPTION
    Copies the goreleaser-built executable and the repo's LICENSE into a
    staging folder, renders the Inno Setup script for the requested version,
    and invokes ISCC to produce a signed installer. Signing uses Azure
    Trusted Signing via signtool's /dlib mechanism, authenticated through
    AZURE_CLIENT_ID / AZURE_CLIENT_SECRET / AZURE_TENANT_ID set by the caller.

.PARAMETER Architecture
    The target architecture. Must be 'amd64', 'arm64' or '386'.

.PARAMETER Version
    The version number to assign to the installer (e.g. "1.2.3", no leading "v").

.PARAMETER SDKVersion
    The Windows SDK version to use for signtool.exe. Defaults to "10.0.26100.0".

.EXAMPLE
    ./build.ps1 -Architecture amd64 -Version "1.2.3"

.OUTPUTS
    Creates Output/install-<Architecture>.exe and Output/install-<Architecture>.exe.sha256.

.NOTES
    Expects the goreleaser artifact at ../../dist/aliae-windows-<Architecture>.exe
    and the license file at ../../LICENSE. Run from the packages/inno directory.
#>

[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [ValidateSet('amd64', 'arm64', '386')]
    [string]$Architecture,

    [Parameter(Mandatory = $true)]
    [ValidateNotNullOrEmpty()]
    [string]$Version,

    [Parameter()]
    [ValidateNotNullOrEmpty()]
    [string]$SDKVersion = "10.0.26100.0"
)

$ErrorActionPreference = 'Stop'
$PSNativeCommandUseErrorActionPreference = $true
$PSDefaultParameterValues['Out-File:Encoding'] = 'UTF8'

function Initialize-SigningEnvironment {
    <#
    .SYNOPSIS
        Sets up the Azure Trusted Signing environment and returns signing tool paths.

    .PARAMETER SDKVersion
        The Windows SDK version to use.

    .OUTPUTS
        Hashtable containing SignTool and SignToolDlib paths.
    #>
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)]
        [string]$SDKVersion
    )

    Write-Verbose "Setting up signing environment" -Verbose

    # Install Microsoft.Trusted.Signing.Client
    nuget.exe install Microsoft.Trusted.Signing.Client -Version 1.0.92 -x | Out-Null

    $signtoolDlib = "$PWD/Microsoft.Trusted.Signing.Client/bin/x64/Azure.CodeSigning.Dlib.dll" -replace '\\', '/'
    $signtool = "C:/Program Files (x86)/Windows Kits/10/bin/$SDKVersion/x64/signtool.exe" -replace '\\', '/'

    if (-not (Test-Path $signtool)) {
        throw "signtool.exe not found at: $signtool"
    }
    if (-not (Test-Path $signtoolDlib)) {
        throw "Azure.CodeSigning.Dlib.dll not found at: $signtoolDlib"
    }

    [hashtable]$result = @{
        SignTool     = $signtool
        SignToolDlib = $signtoolDlib
    }

    return $result
}

New-Item -Path "." -Name "bin" -ItemType Directory -ErrorAction SilentlyContinue | Out-Null

# copy the executable produced by goreleaser
$file = "aliae-windows-$Architecture.exe"
$name = "aliae.exe"
$sourcePath = "../../dist/$file"

if (-not (Test-Path $sourcePath)) {
    throw "Source file not found: $sourcePath"
}

Copy-Item -Path $sourcePath -Destination "./bin/$name" -Force

# license, taken from the local checkout
Copy-Item -Path "../../LICENSE" -Destination "./bin/LICENSE.txt" -Force

$content = Get-Content '.\aliae.iss' -Raw
$content = $content.Replace('<VERSION>', $Version)
$ISSName = ".aliae-$Architecture-$Version.iss"
$content | Out-File -Encoding 'UTF8' $ISSName

# set up Azure Trusted Signing
$signingTools = Initialize-SigningEnvironment -SDKVersion $SDKVersion
$signtool = $signingTools.SignTool
$signtoolDlib = $signingTools.SignToolDlib
$metadataPath = (Resolve-Path "../../src/metadata.json").Path -replace '\\', '/'

# package content
$installer = "install-$Architecture"
ISCC.exe /F$installer "/Ssigntool=$signtool sign /v /debug /fd SHA256 /tr http://timestamp.acs.microsoft.com /td SHA256 /dlib $signtoolDlib /dmdf $metadataPath `$f" $ISSName

# get hash
$installerHash = Get-FileHash "Output/$installer.exe" -Algorithm SHA256
$installerHash.Hash | Out-File -Encoding 'UTF8' "Output/$installer.exe.sha256"
