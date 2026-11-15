# Manual breve de usuario

Software Metrics & Estimation. Guía rápida para usar la aplicación web.

## 1. Iniciar la aplicación

```powershell
docker compose up --build -d --wait
```

Abrí <http://localhost:8080>. Para detenerla: `docker compose down` (conserva los datos) o `docker compose down -v` (borra la base).

## 2. Proyectos

1. En **Proyectos** elegí **Nuevo proyecto**.
2. Completá nombre, fecha de inicio y fecha de fin (la de fin no puede ser anterior a la de inicio).
3. Guardá. El proyecto aparece en el listado con su estado (se calcula por las fechas: por comenzar, en curso o finalizado).
4. Entrá al proyecto para ver su detalle y acceder al Backlog y a los Sprints.

## 3. Product Backlog

1. En el proyecto abrí **Backlog** y elegí **Nuevo ítem**.
2. Completá título, descripción, prioridad y criterios de aceptación. El identificador (`HU-01`, `HU-02`, ...) se asigna solo.
3. Usá **Editar** para cambiar los datos o la prioridad de una historia y **Estimar** para asignarle Story Points.
4. Si un dato no es válido, el formulario vuelve a mostrarse con lo que escribiste y el motivo del error.

## 4. Sprints

1. En **Sprints** elegí **Nuevo Sprint** e ingresá el objetivo (Sprint Goal) y las fechas. Solo puede haber **un Sprint activo** a la vez.
2. En la fila del Sprint, elegí una historia del backlog y **Asignar**. La columna **Story Points planificados** se actualiza.
3. Durante el Sprint, en cada historia pendiente usá **Marcar completada**.
4. Al terminar, **Cerrar Sprint**. Las historias que no se completaron **vuelven al backlog** y el Sprint queda cerrado con sus métricas.

## 5. Métricas

En la pantalla de Sprints se ven, por cada Sprint, los **Story Points planificados**, los **completados** y la **velocidad**. A nivel proyecto se muestran la **velocidad promedio** (solo Sprints cerrados; 0.00 si todavía no hay ninguno) y el **porcentaje de historias completadas**.

## 6. Mensajes de error frecuentes

| Mensaje | Qué significa |
|---|---|
| El Sprint ya está cerrado | No se pueden asignar ni completar historias en un Sprint cerrado |
| Ya existe un Sprint activo | Cerrá el activo antes de crear otro |
| El ítem no existe en el proyecto | La historia elegida no pertenece a este proyecto |
| Ocurrió un error inesperado | Error interno; reintentá y avisá al equipo (el detalle queda en el log) |

## 7. Qué viene

Planning Poker en la interfaz, registro de esfuerzo, gestión de defectos completa, dashboard con gráficos y exportación a PDF (Sprints 2 a 4).
