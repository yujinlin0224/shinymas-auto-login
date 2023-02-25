@echo off

set name=ShinymasAutoLogin
set location=%LocalAppData%\%name%
set output=%location%\%name%.exe

if not exist %location% mkdir %location%
go build -ldflags "-H=windowsgui" -o %output%
