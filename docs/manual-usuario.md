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
3. Guardá. El proyecto aparece en el listado con su estado (se calcula por las fechas: Planificado, En curso o Finalizado).
4. Entrá al proyecto para ver su detalle y acceder a Backlog, Sprints, Defectos, Dashboard, Reporte e Integrantes.

## Navegación

La barra superior cambia según dónde estés: fuera de un proyecto muestra solo **Proyectos**; dentro de un proyecto muestra **Backlog, Sprints, Dashboard, Defectos, Reporte e Integrantes** de ese proyecto. Debajo aparecen las migas de pan (por ejemplo Proyectos › Mi proyecto › Backlog) y el botón **← Volver**, que lleva a la pantalla anterior de la jerarquía.

## Integrantes

1. En el proyecto elegí **Integrantes** (en la barra o con el botón **Ver integrantes**).
2. Completá **Nombre** y **Rol** y elegí **Registrar integrante**. No se admiten nombres repetidos dentro del mismo proyecto.
3. El listado muestra a todos los integrantes. Se necesitan para votar en Planning Poker: si el proyecto no tiene integrantes, la pantalla de Planning Poker lo avisa.

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
5. **Histórico de Sprints**: desde la pantalla de Sprints, el enlace "Histórico de Sprints" muestra los Sprints cerrados con su Sprint Goal, Story Points planificados y completados y velocidad. Si todavía no hay Sprints cerrados, indica que no hay ninguno.

## Planning Poker

1. En el Backlog, en la fila de la historia, elegí **Planning Poker**.
2. Con **Iniciar ronda** se abre la votación.
3. Cada integrante elige su nombre y su voto (1, 2, 3, 5, 8, 13 o 21) y presiona **Votar**. Puede cambiar su voto mientras la ronda esté abierta. La pantalla solo muestra quiénes ya votaron, nunca los valores.
4. Cualquier integrante del proyecto puede **Revelar**: se muestra la tabla con todos los votos. Es irreversible.
5. Con **Nueva ronda** se vota otra vez sobre la misma historia.
6. Con la ronda revelada, la pantalla muestra las **diferencias**: valor mínimo, máximo y quién votó cada uno, además del **historial** de rondas anteriores.
7. **Acordar estimación**: elegí tu nombre y el valor acordado (de la serie 1, 2, 3, 5, 8, 13, 21). Se guarda como Story Points de la historia y se ve en el backlog. No se puede acordar si la ronda no fue revelada.

## 5. Métricas

En la pantalla de Sprints se ven, por cada Sprint, los **Story Points planificados**, los **completados** y la **velocidad**. A nivel proyecto se muestran la **velocidad promedio** (solo Sprints cerrados; 0.00 si todavía no hay ninguno) y el **porcentaje de historias completadas**.

## 6. Defectos

1. En el proyecto, elegí **Ver defectos**.
2. En el formulario completá la descripción, la severidad (Alta, Media o Baja), el estado, la historia relacionada y el Sprint donde se detectó. Si el defecto ya está resuelto, indicá también el Sprint de resolución, que no puede ser anterior al de detección.
3. En el listado, cambiá el **estado** de un defecto y, al marcarlo como Resuelto, elegí el **Sprint de resolución**.

## 7. Dashboard

Desde el detalle del proyecto o desde la barra superior, elegí **Dashboard**. Muestra, en tarjetas compactas, el estado del proyecto, los Story Points planificados y completados, la velocidad promedio, el porcentaje de historias completadas, las horas estimadas y reales con su desviación y, en paneles lado a lado, dos gráficos: **velocidad por Sprint** y **horas estimadas vs. reales por Sprint** (solo Sprints cerrados), cada uno con su tabla de datos. Si el proyecto todavía no tiene Sprints cerrados, los gráficos lo indican.

## 8. Reportes

1. **Reporte del proyecto**: en el detalle del proyecto elegí **Ver reporte**. Incluye historias planificadas y completadas, estimaciones, esfuerzo, métricas y defectos.
2. **Reporte de un Sprint**: en la pantalla de Sprints, en la fila del Sprint, elegí **Ver reporte del Sprint**.
3. **Descargar PDF**: en cada reporte, el botón **Descargar PDF** baja el informe en formato PDF, generado por el servidor.
4. También podés usar **Imprimir** (o Ctrl+P): la impresión oculta la navegación y los botones.

## 9. Mensajes de error frecuentes

| Mensaje | Qué significa |
|---|---|
| El Sprint ya está cerrado | No se pueden asignar ni completar historias en un Sprint cerrado |
| Ya existe un Sprint activo | Cerrá el activo antes de crear otro |
| El ítem no existe en el proyecto | La historia elegida no pertenece a este proyecto |
| Ocurrió un error inesperado | Error interno; reintentá y avisá al equipo (el detalle queda en el log) |

## 10. Limitaciones y qué viene

- Las horas estimadas y los registros de esfuerzo todavía no tienen pantalla de carga: el dashboard y los reportes muestran las horas solo si los datos ya existen.
- No hay inicio de sesión: cualquier persona con acceso a la aplicación puede modificar sus datos.
- Un proyecto Finalizado todavía permite modificaciones (queda como trabajo futuro).
- Trabajo futuro: pantallas de esfuerzo y horas, autenticación, bloqueo de proyectos finalizados y filtros por período.
