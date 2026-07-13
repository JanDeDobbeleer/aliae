<#
.SYNOPSIS
    Prepares the environment for a release build of aliae.

.DESCRIPTION
    Tags the current commit and sets up the Azure Trusted Signing tooling used
    to code sign the produced Windows binaries. The resulting SIGNTOOL and
    SIGNTOOLDLIB paths are written to $env:GITHUB_ENV so later workflow steps
    (goreleaser) can use them.

.PARAMETER Tag
    The git tag to create for this release (e.g. "v1.2.3").

.PARAMETER SDKVersion
    The Windows SDK version to use for signtool.exe. Defaults to "10.0.26100.0".

.EXAMPLE
    ./pre.ps1 -Tag "v1.2.3" -SDKVersion "10.0.26100.0"

.OUTPUTS
    None. Appends SIGNTOOLDLIB and SIGNTOOL to $env:GITHUB_ENV.

.NOTES
    Requires nuget.exe and runs on windows-latest in GitHub Actions.
#>

[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$Tag,

    [Parameter()]
    [string]$SDKVersion = "10.0.26100.0"
)

$ErrorActionPreference = 'Stop'
$PSNativeCommandUseErrorActionPreference = $true
$PSDefaultParameterValues['Out-File:Encoding'] = 'UTF8'

git config --global user.name "GitHub Actions"
git config --global user.email "41898282+github-actions[bot]@users.noreply.github.com"
git tag $Tag --force

# install code signing dlib
nuget.exe install Microsoft.Trusted.Signing.Client -Version 1.0.92 -ExcludeVersion -OutputDirectory $env:RUNNER_TEMP
Write-Output "SIGNTOOLDLIB=$env:RUNNER_TEMP/Microsoft.Trusted.Signing.Client/bin/x64/Azure.CodeSigning.Dlib.dll" | Out-File -FilePath $env:GITHUB_ENV -Encoding utf8 -Append

# requires Windows Dev Kit 10.0.26100.0
$signtool = "C:/Program Files (x86)/Windows Kits/10/bin/$SDKVersion/x64/signtool.exe"
Write-Output "SIGNTOOL=$signtool" | Out-File -FilePath $env:GITHUB_ENV -Encoding utf8 -Append
