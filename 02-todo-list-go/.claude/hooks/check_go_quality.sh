#!/bin/bash

# Colores ANSI
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Archivo de log
LOG_FILE="./logs/go_quality_check.log"

# 1. Verificar si se recibió la ruta del archivo
if [ -z "$1" ]; then
    echo -e "${RED}❌ Error: No se proporcionó la ruta del archivo${NC}"
    echo "Uso: $0 <ruta_completa_archivo.go>"
    exit 1
fi

FILE_PATH="$1"

# 2. Verificar si el archivo existe
if [ ! -f "$FILE_PATH" ]; then
    echo -e "${RED}❌ Error: El archivo '$FILE_PATH' no existe${NC}"
    exit 1
fi

# 3. Verificar si el archivo termina en .go
if [[ ! "$FILE_PATH" =~ \.go$ ]]; then
    echo -e "${RED}❌ Error: El archivo debe terminar en .go${NC}"
    echo "Archivo recibido: $FILE_PATH"
    exit 1
fi

echo -e "${BLUE}🔍 Analizando archivo: $FILE_PATH${NC}"
echo ""

# Ejecutar golangci-lint y capturar la salida
LINT_OUTPUT=$(golangci-lint run "$FILE_PATH" 2>&1)
LINT_EXIT_CODE=$?

# Contar el número de errores/warnings
# golangci-lint muestra cada problema en una línea que contiene el nombre del archivo
ERROR_COUNT=$(echo "$LINT_OUTPUT" | grep -c "\.go:" || echo "0")

# Si golangci-lint retorna 0, no hay errores
if [ $LINT_EXIT_CODE -eq 0 ]; then
    ERROR_COUNT=0
fi

# Determinar el tipo de código según la cantidad de errores
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')
if [ $ERROR_COUNT -eq 0 ]; then
    # Tipo A: 0 errores
    QUALITY_TYPE="A"
    EMOJI="🟢"
    COLOR=$GREEN
    MESSAGE="Código excelente - Sin errores"
elif [ $ERROR_COUNT -ge 1 ] && [ $ERROR_COUNT -le 5 ]; then
    # Tipo B: 1-5 errores
    QUALITY_TYPE="B"
    EMOJI="🟡"
    COLOR=$YELLOW
    MESSAGE="Código aceptable - $ERROR_COUNT error(es) encontrado(s)"
else
    # Tipo C: más de 5 errores
    QUALITY_TYPE="C"
    EMOJI="🔴"
    COLOR=$RED
    MESSAGE="Código necesita mejoras - $ERROR_COUNT error(es) encontrado(s)"
fi

# Mostrar resultado en consola
echo -e "${COLOR}${EMOJI} TIPO $QUALITY_TYPE: $MESSAGE${NC}"
echo ""

# Guardar en el log
echo "========================================" >> "$LOG_FILE"
echo "Timestamp: $TIMESTAMP" >> "$LOG_FILE"
echo "Archivo: $FILE_PATH" >> "$LOG_FILE"
echo "Tipo: $QUALITY_TYPE $EMOJI" >> "$LOG_FILE"
echo "Errores encontrados: $ERROR_COUNT" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

if [ $ERROR_COUNT -gt 0 ]; then
    echo "Detalles de los errores:" >> "$LOG_FILE"
    echo "$LINT_OUTPUT" >> "$LOG_FILE"
    echo "" >> "$LOG_FILE"

    # Mostrar los errores en consola también
    echo -e "${YELLOW}📋 Detalles de los errores:${NC}"
    echo "$LINT_OUTPUT"
    echo ""
fi

echo "========================================" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

# Mensaje final
echo -e "${GREEN}✅ Resultado guardado en: $LOG_FILE${NC}"

# Mostrar resumen
echo ""
echo -e "${BLUE}📊 RESUMEN:${NC}"
echo -e "  Archivo: $FILE_PATH"
echo -e "  Calidad: ${COLOR}$EMOJI TIPO $QUALITY_TYPE${NC}"
echo -e "  Errores: $ERROR_COUNT"

exit 0
