#!/bin/bash
#
# build_teacher_planner.sh
# Script para generar un planificador docente automáticamente
# Período académico: Septiembre - Agosto (12 meses, 4 trimestres iguales)
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
    echo "Esto genera un planner de septiembre 2025 - agosto 2026"
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
echo -e "${GREEN}Año de generación:${NC} $START_YEAR (sept $START_YEAR - agosto $END_YEAR)"
echo ""

# Verificar que estamos en la carpeta correcta
if [ ! -f "single.sh" ]; then
    echo -e "${YELLOW}Error: No se encuentra single.sh${NC}"
    echo "Asegúrate de estar en la carpeta raíz de latex-yearly-planner"
    exit 1
fi

for f in cfg/teacher_base.yaml cfg/teacher_schedule.yaml; do
    if [ ! -f "$f" ]; then
        echo -e "${YELLOW}Error: No se encuentra $f${NC}"
        exit 1
    fi
done

# Calendario escolar del curso (festivos, vacaciones, inicio/fin de curso).
# Es un archivo por curso: cfg/teacher_calendar_2026-27.yaml, etc. Si no
# existe para el curso pedido, se genera sin festivos (con aviso).
CALENDAR="cfg/teacher_calendar_${SCHOOL_YEAR}.yaml"
if [ -f "$CALENDAR" ]; then
    CALENDAR_CFG=",${CALENDAR}"
else
    echo -e "${YELLOW}Aviso: no existe $CALENDAR; el planner no llevará festivos.${NC}"
    echo "       Copia cfg/teacher_calendar_2026-27.yaml y adapta las fechas."
    CALENDAR_CFG=""
fi

# Cadena de configs ("el último gana"):
#   base.yaml                      valores por defecto del proyecto
#   rm2.base.yaml                  papel y márgenes de ReMarkable 2
#   template_months_on_side.yaml   layout "months on side"
#   rm2.mos.default.yaml           ajustes de RM2 para ese layout
#   teacher_base.yaml              curso académico (sept-agosto), líneas, lunes
#   teacher_schedule.yaml          horario de clases y horas lectivas (8-15)
#   teacher_calendar_AAAA-AA.yaml  calendario escolar del curso (festivos)
#
# Si tu RM2 usa firmware DDVK, cambia cfg/rm2.base.yaml por
# cfg/rm2_ddvk.base.yaml (o cfg/rm2_ddvk_lh.base.yaml para zurdos).
# Para generar SIN horario (columna de horas clásica) basta con quitar
# cfg/teacher_schedule.yaml de la lista.
echo -e "${BLUE}Generando planificador...${NC}"
echo ""

# PASSES=2 es necesario: las pestañas laterales (marginnote) solo se colocan
# bien en la segunda pasada de xelatex. Con una sola pasada aparecen
# desplazadas encima del contenido.
# TRANSLATION=spanish traduce los textos (meses, días, "Horario", etc.)
# con translations/spanish.json; quita la línea para dejarlo en inglés.
PLANNER_YEAR=$START_YEAR \
PASSES=2 \
TRANSLATION=spanish \
CFG="cfg/base.yaml,cfg/rm2.base.yaml,cfg/template_months_on_side.yaml,cfg/rm2.mos.default.yaml,cfg/teacher_base.yaml,cfg/teacher_schedule.yaml${CALENDAR_CFG}" \
NAME="teacher_planner_${SCHOOL_YEAR}" \
./single.sh

echo ""
echo -e "${GREEN}✅ ¡Listo!${NC}"
echo ""
echo -e "Archivo generado: ${GREEN}teacher_planner_${SCHOOL_YEAR}.pdf${NC}"
echo ""
echo "Para cambiar el horario de clases: edita cfg/teacher_schedule.yaml"
echo "Para cambiar festivos y vacaciones: edita ${CALENDAR}"
echo ""
