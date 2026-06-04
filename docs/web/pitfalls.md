# Errores comunes — Frontend

## Doble request al refrescar lista paginada

**Problema:** Invalidar caché antes de cambiar página → stale query → fetch reactivo + refetch() = 2 requests.

**Solución:** Patrón de 3 pasos:

```ts
async function refresh() {
  if (refreshing) return
  refreshing = true
  try {
    page.value = 1                     // 1. Cambiar página SIN invalidar
    await nextTick()                   // 2. Esperar watcher interno
    await queryClient.invalidateQueries({
      queryKey: ['entidad'],
      refetchType: 'active',           // 3. Refetch solo observer activo
    })
  } finally {
    refreshing = false
  }
}
```

## Servicios que NO usan el HTTP compartido

`healthCatalogService` y `setupService` tienen su propia función `request()` interna. No importan `get/post/put/patch/del` desde `@/services/http`. Verificar antes de modificar.
