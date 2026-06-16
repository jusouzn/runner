#!/usr/bin/env bash
set -euo pipefail

PLANTUML_JAR="plantuml.jar"
PLANTUML_URL="https://github.com/plantuml/plantuml/releases/latest/download/plantuml.jar"
DIAGRAMAS_DIR="diagramas"
IMAGENS_DIR="${DIAGRAMAS_DIR}/imagens"

if [ ! -f "${PLANTUML_JAR}" ]; then
  echo "Baixando plantuml.jar..."
  curl -sL "${PLANTUML_URL}" -o "${PLANTUML_JAR}"
fi

mkdir -p "${IMAGENS_DIR}"

echo "Gerando diagramas..."
java -jar "${PLANTUML_JAR}" -tsvg -o "$(pwd)/${IMAGENS_DIR}" "${DIAGRAMAS_DIR}"/*.puml
echo "Diagramas gerados em ${IMAGENS_DIR}/"
