param (
    [Parameter()][Alias("tn")][string]
    $taskName = "ShinymasAutoLogin",

    [Parameter()][Alias("ri")][ValidatePattern("\d{2}:\d{2}:\d{2}")][string]
    $repetitionInterval = "00:02:00",

    [Parameter()][Alias("rd")][ValidatePattern("\d{2}:\d{2}:\d{2}")][string]
    $repetitionDuration = "00:10:00",

    [Parameter()][Alias("tt")][ValidatePattern("\d{2}:\d{2}:\d{2}")][string]
    $triggerTime = "05:00:00",

    [Parameter()][Alias("td")][ValidatePattern("\d{2}:\d{2}:\d{2}")][string]
    $triggerRandomDelay = "00:00:05",

    [Parameter()][Alias("pp")][string]
    $programPath = "$env:LOCALAPPDATA\ShinymasAutoLogin\ShinymasAutoLogin.exe",

    [Parameter()][Alias("pa")][string]
    $programArguments = "-bn edge -ri 100 -rt 100 -wt 20"
)

$taskDescription = "Automatically log in to the web game `“THE iDOLM@STER SHINY COLORS`” to receive the daily login bonus."
$taskAuthor = "yujinlin0224"

$triggerDate = (Get-Date -Date "2018-04-24T00:00:00+09:00")

$taskPrincipal = New-ScheduledTaskPrincipal `
    -UserId ([System.Security.Principal.WindowsIdentity]::GetCurrent().Name) `
    -LogonType Interactive `
    -RunLevel Limited
$taskTriggerRepetition = (
    New-ScheduledTaskTrigger `
        -Once `
        -At (Get-Date) `
        -RepetitionInterval ([System.TimeSpan]::Parse($repetitionInterval)) `
        -RepetitionDuration ([System.TimeSpan]::Parse($repetitionDuration))
).Repetition
$taskTrigger = New-ScheduledTaskTrigger `
    -Daily `
    -At ($triggerDate + [System.TimeSpan]::Parse($triggerTime)) `
    -DaysInterval 1 `
    -RandomDelay ([System.TimeSpan]::Parse($triggerRandomDelay))
$taskTrigger.Repetition = $taskTriggerRepetition
$taskAction = New-ScheduledTaskAction -Execute $programPath
$taskAction.Arguments = $programArguments
$task = New-ScheduledTask `
    -Description $taskDescription `
    -Principal $taskPrincipal `
    -Trigger $taskTrigger `
    -Action $taskAction
$task.Author = $taskAuthor
$task | Register-ScheduledTask -TaskName $taskName
