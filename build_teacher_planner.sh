#!/bin/bash
#
# build_teacher_planner.sh
# Script para generar un planificador docente automáticamente
# Período académico: Septiembre - Junio
#
# Uso:
#   ./build_teacher_planner.sh 2025     # Genera planner para 2025-2026
#   ./build_teacher_planner.sh 2026     # Genera planner para 2026-2027
#

set -e  # Exit on error

# Colores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Verificar argumentos
if [ $# -eq 0 ]; then
    echo -e "${YELLOW}Uso: ./build_teacher_planner.sh <AÑO>${NC}"
    echo "Ejemplo: ./build_teacher_planner.sh 2025"
    echo ""
    echo "Esto genera un planner de septiembre 2025 - junio 2026"
    exit 1
fi

START_YEAR=$1
END_YEAR=$((START_YEAR + 1))
SCHOOL_YEAR="${START_YEAR}-${END_YEAR: -2}"

echo -e "${BLUE}════════════════════════════════════════${NC}"
echo -e "${BLUE}  📚 Generador de Planificador Docente${NC}"
echo -e "${BLUE}════════════════════════════════════════${NC}"
echo ""
echo -e "${GREEN}Período académico:${NC} $SCHOOL_YEAR"
echo -e "${GREEN}Año de generación:${NC} $START_YEAR (sept 2025 - junio 2026)"
echo ""

# Verificar que estamos en la carpeta correcta
if [ ! -f "single.sh" ]; then
    echo -e "${YELLOW}Error: No se encuentra single.sh${NC}"
    echo "Asegúrate de estar en la carpeta raíz de latex-yearly-planner"
    exit 1
fi

# Verificar que exista teacher_base.yaml
if [ ! -f "cfg/teacher_base.yaml" ]; then
    echo -e "${YELLOW}Error: No se encuentra cfg/teacher_base.yaml${NC}"
    echo "Asegúrate de haber copiado teacher_base.yaml en la carpeta cfg/"
    exit 1
fi

# Generar el planificador
echo -e "${BLUE}Generando planificador...${NC}"
echo ""

PLANNER_YEAR=$START_YEAR \
PASSES=1 \
CFG="cfg/base.yaml,cfg/teacher_base.yaml,cfg/template_breadcrumb.yaml" \
NAME="teacher_planner_${SCHOOL_YEAR}" \
./single.sh

echo ""
echo -e "${GREEN}✅ ¡Listo!${NC}"
echo ""
echo -e "Archivo generado: ${GREEN}teacher_planner_${SCHOOL_YEAR}.pdf${NC}"
echo ""
echo "Próximos pasos:"
echo "  1️⃣  Revisar el PDF generado"
echo "  2️⃣  Ajustar configuración si es necesario"
echo "  3️⃣  Cuando esté listo, proceder a FASE 2 (horario semanal)"
echo ""
