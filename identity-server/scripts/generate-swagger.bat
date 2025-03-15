@echo off
setlocal

REM Set default values for environment variables if not provided
if "%API_HOST%"=="" set API_HOST=localhost:8000

REM Generate swagger docs
swag init -g main.go

REM Replace placeholder values in swagger.json and swagger.yaml
powershell -Command "(Get-Content docs\swagger.json) -replace '\${API_HOST}', '%API_HOST%' | Set-Content docs\swagger.json"
powershell -Command "(Get-Content docs\swagger.yaml) -replace '\${API_HOST}', '%API_HOST%' | Set-Content docs\swagger.yaml"

echo Swagger documentation generated successfully

endlocal