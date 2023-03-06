param (
    [Parameter()][Alias("tn")][string]
    $TaskName = "ShinymasAutoLogin",

    [Parameter()][Alias("ri")][ValidatePattern("\d{2}:\d{2}:\d{2}")][string]
    $RepetitionInterval = "00:02:00",

    [Parameter()][Alias("rd")][ValidatePattern("\d{2}:\d{2}:\d{2}")][string]
    $RepetitionDuration = "00:10:00",

    [Parameter()][Alias("tt")][ValidatePattern("\d{2}:\d{2}:\d{2}")][string]
    $TriggerTime = "05:00:00",

    [Parameter()][Alias("td")][ValidatePattern("\d{2}:\d{2}:\d{2}")][string]
    $TriggerRandomDelay = "00:00:05",

    [Parameter()][Alias("pp")][string]
    $ProgramPath = "$Env:LocalAppData\ShinymasAutoLogin\ShinymasAutoLogin.exe",

    [Parameter()][Alias("pa")][string]
    $ProgramArguments = "-bn edge -ri 100 -rt 100 -wt 20"
)

$TaskDescription = "Automatically log in to the web game `“THE iDOLM@STER SHINY COLORS`” to receive the daily login bonus."
$TaskAuthor = "yujinlin0224"

$TriggerDate = (Get-Date -Date "2018-04-24T00:00:00+09:00")

$TaskPrincipal = New-ScheduledTaskPrincipal `
    -UserId ([System.Security.Principal.WindowsIdentity]::GetCurrent().Name) `
    -LogonType Interactive `
    -RunLevel Limited
$TaskTriggerRepetition = (
    New-ScheduledTaskTrigger `
        -Once `
        -At (Get-Date) `
        -RepetitionInterval ([System.TimeSpan]::Parse($RepetitionInterval)) `
        -RepetitionDuration ([System.TimeSpan]::Parse($RepetitionDuration))
).Repetition
$TaskTrigger = New-ScheduledTaskTrigger `
    -Daily `
    -At ($TriggerDate + [System.TimeSpan]::Parse($TriggerTime)) `
    -DaysInterval 1 `
    -RandomDelay ([System.TimeSpan]::Parse($TriggerRandomDelay))
$TaskTrigger.Repetition = $TaskTriggerRepetition
$TaskAction = New-ScheduledTaskAction -Execute $ProgramPath
$TaskAction.Arguments = $ProgramArguments
$Task = New-ScheduledTask `
    -Description $TaskDescription `
    -Principal $TaskPrincipal `
    -Trigger $TaskTrigger `
    -Action $TaskAction
$Task.Author = $TaskAuthor
$Task | Register-ScheduledTask -TaskName $TaskName
