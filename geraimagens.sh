#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"

PLANTUML_JAR="${SCRIPT_DIR}/plantuml.jar"
PLANTUML_URL="https://github.com/plantuml/plantuml/releases/latest/download/plantuml.jar"
DIAGRAMAS_DIR="${SCRIPT_DIR}/diagramas"
IMAGENS_DIR="${DIAGRAMAS_DIR}/imagens"

if [ ! -f "${PLANTUML_JAR}" ]; then
  echo "Baixando plantuml.jar..."
  curl -fSL "${PLANTUML_URL}" -o "${PLANTUML_JAR}"
fi

mkdir -p "${IMAGENS_DIR}"

echo "Gerando diagramas..."
(
  cd "${DIAGRAMAS_DIR}"
  java -jar "${PLANTUML_JAR}" -tsvg -o "imagens" ./*.puml
)
echo "Diagramas gerados em ${IMAGENS_DIR}/"
