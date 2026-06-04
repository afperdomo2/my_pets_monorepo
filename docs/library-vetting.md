# Validación de nuevas librerías

Reglas a seguir **antes** de agregar cualquier dependencia nueva al proyecto (npm, pnpm, Go modules).

## Checklist de validación

- [ ] **Reputación**: revisar GitHub stars, download count (npm trends), fecha del último commit
- [ ] **No legacy/obsoleta**: si no recibe actualizaciones desde hace >1 año, buscar alternativa activa
- [ ] **Preferir oficial**: dar prioridad a librerías mantenidas por la organización del framework (ej: `@tanstack/vue-query` sobre alternativas third-party)
- [ ] **Bundle size** (frontend): consultar bundlephobia.com antes de agregar una lib JS
- [ ] **Licencia**: compatible con el proyecto (MIT, Apache-2.0, etc.)
- [ ] **Documentación**: usar Context7 MCP para obtener documentación actualizada antes de integrar
- [ ] **Justificación**: si la elección no es obvia, agregar un comentario breve explicando por qué se eligió esa lib sobre otras

## Proceso recomendado

1. Identificar la necesidad y buscar 2-3 candidatos
2. Ejecutar `context7_resolve-library-id` + `context7_query-docs` para cada candidato
3. Comparar según los criterios arriba
4. Instalar con `pnpm add <lib>` o `go get <lib>`
5. Si aplica, verificar tipos con `pnpm type-check`
