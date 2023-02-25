# Shinymas Auto Login

Automatically log in to the web game “THE iDOLM@STER SHINY COLORS” to receive the daily login bonus.

## Requirements

- Go 1.20 or later to build the program
- PowerShell 7.0 or later to create the scheduled task
- Chromium-based browser (Microsoft Edge, Google Chrome, Chromium, Brave) that able to run the game and support app mode using the command line option `--app=<URL>`

## Installation

1. Run `install.bat` in Command Prompt or PowerShell to build and install the program.
2. Run `create-task.ps1` in PowerShell to create a scheduled task that trigger the program every day.

## Usage

- The program supports the following command line options:
  - `-h, --help`: Show the help message.
  - `-bn`: Specify the browser name, available values are `edge`, `chrome`, `chromium`, `brave`. Default value is `edge`.
  - `-ri`: Specify the retry interval in milliseconds. Default value is `100`.
  - `-rt`: Specify the retry times. Default value is `100`.
  - `-wt`: Specify the waiting time in seconds. Default value is `20`.

- `create-task.ps1` supports the following command line options:
  - `-tn`: Specify the task name. Default value is `ShinymasAutoLogin`.
  - `-ri`: Specify the repetition interval. Default value is `00:02:00`.
  - `-rd`: Specify the repetition duration. Default value is `00:10:00`.
  - `tt`: Specify the trigger time in +09:00 time zone. Default value is `05:00:00`.
  - `td`: Specify the trigger random delay. Default value is `00:00:05`.
  - `-pp`: Specify the program path. Default value is `$env:LOCALAPPDATA\ShinymasAutoLogin\ShinymasAutoLogin.exe`.
  - `-pa`: Specify the program arguments. Default value is `-bn edge -ri 100 -rt 100 -wt 20`.
