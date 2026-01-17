$ErrorActionPreference = "Stop"

$AppName = "orchastration"
$Version = $env:VERSION; if ([string]::IsNullOrEmpty($Version)) { $Version = "dev" }
$Commit = $env:COMMIT; if ([string]::IsNullOrEmpty($Commit)) { $Commit = "unknown" }
$BuildTime = $env:BUILDTIME; if ([string]::IsNullOrEmpty($BuildTime)) { $BuildTime = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ") }

New-Item -ItemType Directory -Force -Path dist | Out-Null

$ldflags = "-X 'orchastration/internal/version.Version=$Version' -X 'orchastration/internal/version.Commit=$Commit' -X 'orchastration/internal/version.BuildTime=$BuildTime'"

$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -trimpath -ldflags $ldflags -o ("dist/" + $AppName + ".exe") ("./cmd/" + $AppName)
