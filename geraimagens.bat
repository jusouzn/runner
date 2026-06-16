@echo off
setlocal

set PLANTUML_JAR=plantuml.jar
set PLANTUML_URL=https://github.com/plantuml/plantuml/releases/latest/download/plantuml.jar
set DIAGRAMAS_DIR=diagramas
set IMAGENS_DIR=%DIAGRAMAS_DIR%\imagens

if not exist "%PLANTUML_JAR%" (
    echo Baixando plantuml.jar...
    powershell -Command "Invoke-WebRequest -Uri '%PLANTUML_URL%' -OutFile '%PLANTUML_JAR%'"
)

if not exist "%IMAGENS_DIR%" mkdir "%IMAGENS_DIR%"

echo Gerando diagramas...
java -jar "%PLANTUML_JAR%" -tsvg -o "%CD%\%IMAGENS_DIR%" "%DIAGRAMAS_DIR%\*.puml"
echo Diagramas gerados em %IMAGENS_DIR%\
