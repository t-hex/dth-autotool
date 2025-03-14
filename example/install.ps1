Param(
    [Parameter()]
    [Alias("sil")]
    [Switch]$SkipInstallLatest = $False
)

Function Get-IsAdminMode {
    ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

Function Create-Path {
    Param(
        [Parameter(Mandatory, ValueFromPipeline)]
        [ValidateNotNullOrEmpty()]
        [String]$Path
    )

    $AbsPath = [System.IO.Path]::GetFullPath($Path)
    Write-Host "Checking path: '$AbsPath'."
    If(-not $(Get-Item $AbsPath -ErrorAction 'SilentlyContinue')) {
        Write-Warning "Path does not exist! Creating path: '$AbsPath'."
        New-Item -Path $AbsPath -ItemType Directory -Force | Out-Null
    }
}

Function Create-Link {
    Param(
        [Parameter(Mandatory, ValueFromPipeline)]
        [ValidateNotNullOrEmpty()]
        [String]$TargetPath,

        [Parameter(Mandatory, Position = 0)]
        [ValidateNotNullOrEmpty()]
        [String]$LinkPath
    )

    $AbsLinkPath = [System.IO.Path]::GetFullPath($LinkPath)
    $AbsTargetPath = [System.IO.Path]::GetFullPath($TargetPath)

    $LinkItem = Get-Item $AbsLinkPath -ErrorAction 'SilentlyContinue'
    If($LinkItem) {
        Write-Host "Removing previous symbolic link: '$AbsLinkPath'."
        $LinkItem.Delete()
    }
    If(Get-Item $AbsTargetPath -ErrorAction 'SilentlyContinue') {
        Write-Host "Creating new symbolic link: '$AbsLinkPath' to '$AbsTargetPath'."
        New-Item -Path $AbsLinkPath -ItemType SymbolicLink -Value $AbsTargetPath | Out-Null
    } Else {
        Write-Warning "Path '$AbsTargetPath' does not exist."
        Write-Warning "New link '$AbsLinkPath' will not be created!"
    }
}

If (-not $(Get-IsAdminMode)) {
    Start-Process powershell.exe -ArgumentList ( `
        '-NoExit', `
        '-ExecutionPolicy', 'Bypass', `
        '-Command', "cd '$PSScriptRoot'; $($MyInvocation.Line -Replace '`"','\"')" `
    ) -Verb RunAs # -WorkingDirectory $PSScriptRoot
    Exit
}

$config = @{
    ProjectConfigFileName = "project.config.json"
    DthAutoToolAppConfigFileName = "app.config.json"
    AlembicExporterConfigFileName = "abc.config.sagan"

    AppPath = "$PSScriptRoot\app"
    ConfigPath = "$PSScriptRoot\config"
    TempPath = "$PSScriptRoot\temp"
    ExportPath = "$PSScriptRoot\exports"

    FbxTempExportPath = "$($env:USERPROFILE)\Documents\DAZ 3D\Bridges\Daz To Maya\Exports"
    AbcTempExportPath = "$($env:TEMP)\DTH-AutoTool\Exports\ABC-Temp"

    FbxTempExportLinkName = "fbx"
    AbcTempExportLinkName = "abc"

    DownloadPageUrl = 'https://e.pcloud.link/publink/show?code=XZptNlZabWBSUzGL3JsE2W6OfOLtJGX32Dy'
}

# $config.FbxTempExportPath | Create-Path # should exist already
$config.AbcTempExportPath | Create-Path

$config.AppPath | Create-Path
$config.ConfigPath | Create-Path
$config.TempPath | Create-Path
$config.ExportPath | Create-Path

If($SkipInstallLatest) {
    Write-Warning 'Skipping latest DTH-AT download & installation'
} Else {
    Get-ChildItem -Path $config.AppPath -Recurse | Remove-Item -Recurse -Force

    $downloadPageResponse = Invoke-WebRequest $config.DownloadPageUrl
    $pubLinkVarDeclOffset = $downloadPageResponse.Content.IndexOf('publinkData')
    $pubLinkVarDefStartOffset = $downloadPageResponse.Content.IndexOf('{', $pubLinkVarDeclOffset)
    $pubLinkVarDefEndOffset = $downloadPageResponse.Content.IndexOf('};', $pubLinkVarDefStartOffset) + 1
    $pubLinkDataString = $downloadPageResponse.Content.Substring($pubLinkVarDefStartOffset, $pubLinkVarDefEndOffset-$pubLinkVarDefStartOffset)
    $pubLinkDataJsonData = $pubLinkDataString  | ConvertFrom-Json

    Invoke-WebRequest -Uri $pubLinkDataJsonData.downloadlink -OutFile "$($config.AppPath)\dth-at-latest.zip"
    Expand-Archive -Path "$($config.AppPath)\dth-at-latest.zip" -DestinationPath $config.AppPath -Force
    Get-Item -Path "$($config.AppPath)\dth-at-latest.zip" | Remove-Item -Force
}

$config.FbxTempExportPath | Create-Link "$($config.TempPath)/$($config.FbxTempExportLinkName)"
$config.AbcTempExportPath | Create-Link "$($config.TempPath)/$($config.AbcTempExportLinkName)"

"$($config.ConfigPath)/$($config.ProjectConfigFileName)" | Create-Link "$($config.AppPath)/$($config.ProjectConfigFileName)"
"$($config.ConfigPath)/$($config.DthAutoToolAppConfigFileName)" | Create-Link "$($config.AppPath)/$($config.DthAutoToolAppConfigFileName)"
"$($config.ConfigPath)/$($config.AlembicExporterConfigFileName)" | Create-Link "$($config.AppPath)/$($config.AlembicExporterConfigFileName)"
