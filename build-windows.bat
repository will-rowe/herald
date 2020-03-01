@echo off
go generate
go build -ldflags "-H windowsgui -X main.dbLocation=%APPDATA%/herald" -o herald.exe
