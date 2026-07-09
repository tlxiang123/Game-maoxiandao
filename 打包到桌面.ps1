$ErrorActionPreference = 'Stop'

$Project = Split-Path -Parent $MyInvocation.MyCommand.Path
$ShellApkRoot = 'C:\Users\Administrator\Desktop\auto\2_Win10'
$ShellApk = Get-ChildItem -Path $ShellApkRoot -Recurse -Filter app-release.apk -File |
    Select-Object -First 1 -ExpandProperty FullName
if (-not $ShellApk) {
    throw "找不到壳 APK：$ShellApkRoot"
}
$ArmApp = Join-Path $Project 'build\arm64-v8a-release'
$X64App = Join-Path $Project 'build\x86_64-release'
$Temp = Join-Path $env:TEMP ('autogo_apk_' + [guid]::NewGuid().ToString('N'))
$Unsigned = Join-Path $Project 'build\通用-unsigned.apk'
$BaseZip = Join-Path $Project 'build\通用-base.zip'
$Signed = Join-Path $Project 'build\app-release.apk'
$DesktopCn = 'C:\Users\Administrator\Desktop\通用.apk'
$DesktopEn = 'C:\Users\Administrator\Desktop\autogo-test.apk'
$KeyStore = Join-Path $Project 'build\autogo-debug.keystore'
$Keytool = 'D:\Install\JDK21\bin\keytool.exe'
$Javac = 'D:\Install\JDK21\bin\javac.exe'
$Java = 'D:\Install\JDK21\bin\java.exe'
$SevenZip = 'D:\Install\7-Zip\7z.exe'
$JadxJar = 'D:\Install\Jadx\lib\jadx-1.5.5-all.jar'
$SignJava = Join-Path $Project 'build\SignApkWithApksig.java'

foreach ($Path in @($ShellApk, $ArmApp, $X64App, $Keytool, $Javac, $Java, $SevenZip, $JadxJar, $SignJava)) {
    if (!(Test-Path $Path)) {
        throw "缺少文件：$Path"
    }
}

Remove-Item $Temp -Recurse -Force -ErrorAction SilentlyContinue
Remove-Item $Unsigned, $BaseZip, $Signed, $DesktopCn, $DesktopEn -Force -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Path $Temp | Out-Null

Add-Type -AssemblyName System.IO.Compression.FileSystem
[System.IO.Compression.ZipFile]::ExtractToDirectory($ShellApk, $Temp)
Remove-Item (Join-Path $Temp 'META-INF') -Recurse -Force -ErrorAction SilentlyContinue

$AssetRoot = Join-Path $Temp 'assets'
foreach ($Abi in @('arm64-v8a', 'x86_64')) {
    $AbiDir = Join-Path $AssetRoot $Abi
    New-Item -ItemType Directory -Path $AbiDir -Force | Out-Null

    $AppSource = if ($Abi -eq 'arm64-v8a') { $ArmApp } else { $X64App }
    Copy-Item -LiteralPath $AppSource -Destination (Join-Path $AbiDir 'app') -Force

    $LibDir = Join-Path $Project "resources\libs\$Abi"
    if (Test-Path $LibDir) {
        Get-ChildItem -LiteralPath $LibDir -File | ForEach-Object {
            Copy-Item -LiteralPath $_.FullName -Destination (Join-Path $AbiDir $_.Name) -Force
        }
    }
}

# targetSdk 30+ requires resources.arsc to be stored, and extractNativeLibs=false
# requires native libraries under lib/ to be stored and page-aligned.
$Fs = [System.IO.File]::Open($BaseZip, [System.IO.FileMode]::CreateNew)
try {
    $Archive = [System.IO.Compression.ZipArchive]::new($Fs, [System.IO.Compression.ZipArchiveMode]::Create)
    try {
        Get-ChildItem -LiteralPath $Temp -Recurse -File |
            Where-Object { $_.FullName -notlike "*\lib\*" -and $_.Name -ne 'resources.arsc' } |
            ForEach-Object {
                $Rel = $_.FullName.Substring($Temp.Length).TrimStart('\', '/') -replace '\\', '/'
                [System.IO.Compression.ZipFileExtensions]::CreateEntryFromFile(
                    $Archive,
                    $_.FullName,
                    $Rel,
                    [System.IO.Compression.CompressionLevel]::Optimal
                ) | Out-Null
            }
    } finally {
        $Archive.Dispose()
    }
} finally {
    $Fs.Dispose()
}

Copy-Item -LiteralPath $BaseZip -Destination $Unsigned -Force
Push-Location $Temp
try {
    & $SevenZip a -tzip -mx=0 $Unsigned '.\resources.arsc' '.\lib\*' | Out-Null
    if ($LASTEXITCODE -ne 0) {
        throw "7z 添加 Store 条目失败，exit=$LASTEXITCODE"
    }
} finally {
    Pop-Location
}

if (!(Test-Path $KeyStore)) {
    & $Keytool -genkeypair -v -keystore $KeyStore -storepass android -keypass android -alias androiddebugkey -keyalg RSA -keysize 2048 -validity 10000 -dname 'CN=Android Debug,O=Android,C=US' | Out-Null
    if ($LASTEXITCODE -ne 0) {
        throw "生成 debug keystore 失败，exit=$LASTEXITCODE"
    }
}

& $Javac -encoding UTF-8 -cp $JadxJar $SignJava
if ($LASTEXITCODE -ne 0) {
    throw "编译 SignApkWithApksig 失败，exit=$LASTEXITCODE"
}

$SignOutput = & $Java -cp "$JadxJar;$Project\build" SignApkWithApksig $Unsigned $Signed $KeyStore android androiddebugkey android 2>&1
$SignOutput | ForEach-Object { $_.ToString() }
if ($LASTEXITCODE -ne 0) {
    throw "apksig 签名失败，exit=$LASTEXITCODE"
}

Copy-Item -LiteralPath $Signed -Destination $DesktopCn -Force
Copy-Item -LiteralPath $Signed -Destination $DesktopEn -Force
Remove-Item $Temp -Recurse -Force -ErrorAction SilentlyContinue

Get-Item $Signed, $DesktopCn, $DesktopEn | Select-Object FullName, Length, LastWriteTime
