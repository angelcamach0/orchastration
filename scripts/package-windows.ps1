$ErrorActionPreference = "Stop"

$AppName = "orchastration"
$Version = $env:VERSION; if ([string]::IsNullOrEmpty($Version)) { $Version = "dev" }

New-Item -ItemType Directory -Force -Path dist/package/windows | Out-Null
Copy-Item -Force ("dist/" + $AppName + ".exe") "dist/package/windows/"
Copy-Item -Force "configs/config.example.toml" ("dist/package/windows/" + $AppName + ".toml")

$zipPath = "dist/" + $AppName + "-" + $Version + "-windows-amd64.zip"
if (Test-Path $zipPath) { Remove-Item -Force $zipPath }
Compress-Archive -Path "dist/package/windows/*" -DestinationPath $zipPath
