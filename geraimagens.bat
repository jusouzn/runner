@echo off
setlocal

set "SCRIPT_DIR=%~dp0"
set "PLANTUML_JAR=%SCRIPT_DIR%plantuml.jar"
set "PLANTUML_URL=https://github.com/plantuml/plantuml/releases/latest/download/plantuml.jar"
set "DIAGRAMAS_DIR=%SCRIPT_DIR%diagramas"
set "IMAGENS_DIR=%DIAGRAMAS_DIR%\imagens"

if not exist "%PLANTUML_JAR%" (
    echo Baixando plantuml.jar...
    powershell -NoProfile -Command "$ProgressPreference='SilentlyContinue'; try { Invoke-WebRequest -Uri '%PLANTUML_URL%' -OutFile '%PLANTUML_JAR%' -UseBasicParsing -ErrorAction Stop } catch { exit 1 }"
    if errorlevel 1 (
        echo ERRO: falha ao baixar %PLANTUML_URL%
        exit /b 1
    )
)

if not exist "%IMAGENS_DIR%" mkdir "%IMAGENS_DIR%"

echo Gerando diagramas...
pushd "%DIAGRAMAS_DIR%"
java -jar "%PLANTUML_JAR%" -tsvg -o "imagens" *.puml
set "RC=%ERRORLEVEL%"
popd
if not "%RC%"=="0" exit /b %RC%

echo Diagramas gerados em %IMAGENS_DIR%\
