# language: es
# Estado: Fase 2 implementada (GREEN); el ajuste Fase 2b (desviación por
# Sprint calculada por el servicio, R23) está especificado y pendiente de
# RED/GREEN.
# Los escenarios describen el comportamiento acordado para CA1–CA3 y las
# decisiones D1–D8 (una página por proyecto GET /proyectos/{id}/dashboard,
# enlazada desde el detalle y desde la navegación del proyecto; contenido con
# el estado calculado por fechas, Story Points totales, velocidad promedio,
# porcentaje de historias completadas, horas totales con su desviación y la
# cantidad de Sprints cerrados; dos gráficos SVG dibujados en el servidor, sin
# JavaScript ni CDN; proyecto sin datos responde 200 con ceros y un mensaje en
# cada gráfico).
# FASE 1 incluye el servicio que arma los datos del dashboard
# (ServicioDashboard.ObtenerDashboard), reutilizando los servicios existentes,
# sin métodos nuevos de repositorio ni migraciones. FASE 2 incluye la página y
# los gráficos; se implementa en un PR aparte, después de FASE 1.
# FASE 2b (ajuste tras la revisión) calcula en el servicio la desviación de
# horas por Sprint (PuntoHorasSprint.Desviacion) con la misma política de
# nil que el total; la web solo la formatea con signo o muestra "—".
# Los defectos (US-15 y US-18) quedan fuera de alcance.
# Los errores se describen por su significado; sus identificadores y códigos
# HTTP exactos están decididos en la spec (sección 8).
# Los escenarios de FASE 1 están marcados con @Fase1, los de FASE 2 con
# @Fase2 y los del ajuste tras la revisión con @Fase2b.
# No se escenifican el HTTP 500/405, el patrón de composición, la
# verificación estructural de imports (sin <script> ni recursos externos),
# la no recalculación de métricas en la web ni la ausencia de sección de
# defectos: se verifican con tests sin escenario BDD propio.

@US-19
Característica: Dashboard con gráficos
  Como equipo
  Quiero un dashboard por proyecto con el estado, las métricas clave y gráficos
  Para ver la situación del proyecto de un vistazo

  @CA1 @Fase1
  Escenario: El servicio arma el dashboard con las métricas del proyecto
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 cerrado con 8 Story Points completados
    Y que el proyecto "PRJ-001" tiene el Sprint 2 cerrado con 13 Story Points completados
    Y que el proyecto "PRJ-001" tiene 5 historias completadas de 10 registradas
    Cuando se consulta el dashboard del proyecto "PRJ-001"
    Entonces se devuelven los Story Points planificados y completados totales del proyecto
    Y se devuelve la velocidad promedio de los Sprints cerrados
    Y se devuelve el porcentaje de historias completadas
    Y se devuelven las horas estimadas, las horas reales y su desviación
    Y se devuelve la cantidad de Sprints cerrados

  @CA1 @CA2 @Fase1
  Escenario: El servicio informa solo los Sprints cerrados en orden ascendente
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 cerrado con 8 Story Points completados
    Y que el proyecto "PRJ-001" tiene el Sprint 2 activo
    Y que el proyecto "PRJ-001" tiene el Sprint 3 cerrado con 5 Story Points completados
    Cuando se consulta el dashboard del proyecto "PRJ-001"
    Entonces la lista de velocidad por Sprint incluye los Sprints "1" y "3" en ese orden
    Y la lista de horas por Sprint incluye los Sprints "1" y "3" en ese orden
    Y la lista no incluye el Sprint "2"

  @CA2 @Fase1
  Escenario: Un Sprint cerrado sin métricas de esfuerzo se informa como métricas no disponibles
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 cerrado con 8 Story Points completados
    Y que el Sprint 1 está cerrado sin snapshot de esfuerzo
    Cuando se consulta el dashboard del proyecto "PRJ-001"
    Entonces el Sprint 1 se informa con métricas no disponibles
    Y el resto del dashboard se devuelve sin error

  @CA1 @Fase1
  Escenario: Un proyecto con estimaciones incompletas no informa desviación de horas
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" tiene al menos una historia sin estimar en horas
    Y que el proyecto "PRJ-001" tiene horas estimadas y horas reales registradas
    Cuando se consulta el dashboard del proyecto "PRJ-001"
    Entonces las horas estimadas suman las estimaciones disponibles
    Y la desviación de horas no está disponible

  @CA1 @Fase1
  Escenario: Un proyecto con la suma de estimaciones en cero no informa desviación de horas
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" tiene horas reales registradas
    Y que las horas estimadas del proyecto suman cero
    Cuando se consulta el dashboard del proyecto "PRJ-001"
    Entonces las horas estimadas son cero
    Y la desviación de horas no está disponible

  @CA3 @Fase1
  Escenario: Un proyecto sin Sprints ni historias devuelve el dashboard en cero
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" no tiene Sprints
    Y que el proyecto "PRJ-001" no tiene historias
    Cuando se consulta el dashboard del proyecto "PRJ-001"
    Entonces los totales de Story Points son cero
    Y la velocidad promedio es cero
    Y el porcentaje de historias completadas es cero
    Y las horas estimadas y reales son cero
    Y la desviación de horas no está disponible
    Y la cantidad de Sprints cerrados es cero

  @CA1 @Fase1
  Escenario: Consultar el dashboard de un proyecto inexistente devuelve proyecto no encontrado
    Cuando se consulta el dashboard del proyecto "PRJ-999"
    Entonces se devuelve el error de proyecto no encontrado

  @CA1 @Fase2
  Escenario: El dashboard muestra el estado del proyecto y las métricas clave
    Dado que hoy es "1 de febrero de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-03-01" y fecha de fin "2026-06-30"
    Cuando se visita el dashboard del proyecto "PRJ-001"
    Entonces se muestra el estado "Planificado"
    Y se muestran los Story Points planificados y completados totales
    Y se muestra la velocidad promedio
    Y se muestra el porcentaje de historias completadas
    Y se muestran las horas estimadas y reales totales con su desviación
    Y se muestra la cantidad de Sprints cerrados

  @CA1 @Fase2
  Escenario: El detalle del proyecto y la navegación enlazan al dashboard
    Dado que hoy es "1 de febrero de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-03-01" y fecha de fin "2026-06-30"
    Cuando se visita el detalle del proyecto "PRJ-001"
    Entonces se ofrece un enlace para ver el dashboard
    Y la navegación ofrece un enlace al dashboard del proyecto

  @CA2 @Fase2
  Escenario: El gráfico de velocidad por Sprint muestra las barras y sus datos
    Dado que hoy es "1 de febrero de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 cerrado con 8 Story Points completados
    Y que el proyecto "PRJ-001" tiene el Sprint 2 cerrado con 13 Story Points completados
    Cuando se visita el dashboard del proyecto "PRJ-001"
    Entonces se muestra el gráfico "Velocidad por Sprint" con una barra por Sprint cerrado
    Y cada barra muestra su valor sobre ella
    Y el gráfico muestra un eje vertical con etiquetas y un eje horizontal con los Sprints
    Y debajo del gráfico se muestra una tabla con el Sprint y sus Story Points completados

  @CA2 @Fase2
  Escenario: El gráfico de horas estimadas vs reales por Sprint muestra barras agrupadas y sus datos
    Dado que hoy es "1 de febrero de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 cerrado con 40 horas estimadas y 35 horas reales
    Cuando se visita el dashboard del proyecto "PRJ-001"
    Entonces se muestra el gráfico "Horas estimadas vs reales por Sprint" con barras agrupadas por Sprint cerrado
    Y se distingue una barra para las horas estimadas y otra para las horas reales
    Y debajo del gráfico se muestra una tabla con el Sprint, las horas estimadas, las horas reales y la desviación

  @CA2 @Fase2
  Escenario: Un Sprint cerrado sin métricas de esfuerzo se muestra como métricas no disponibles
    Dado que hoy es "1 de febrero de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 cerrado sin snapshot de esfuerzo
    Cuando se visita el dashboard del proyecto "PRJ-001"
    Entonces el gráfico de horas no dibuja barras para el Sprint 1
    Y la tabla de horas muestra "métricas no disponibles" para el Sprint 1
    Y la página responde sin error

  @CA2 @Fase2b
  Escenario: La desviación por Sprint se muestra con signo o como no disponible
    Dado que hoy es "1 de febrero de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 cerrado con 40 horas estimadas y 35 horas reales
    Y que el proyecto "PRJ-001" tiene el Sprint 2 cerrado con al menos una historia sin estimar en horas
    Cuando se visita el dashboard del proyecto "PRJ-001"
    Entonces la tabla de horas muestra una desviación de "-12.50%" para el Sprint 1
    Y la tabla de horas muestra la desviación del Sprint 2 como "—"
    Y la página responde sin error

  @CA3 @Fase2
  Escenario: Un proyecto sin datos muestra el dashboard con ceros y los mensajes en los gráficos
    Dado que hoy es "1 de febrero de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-03-01" y fecha de fin "2026-06-30"
    Y que el proyecto "PRJ-001" no tiene Sprints
    Y que el proyecto "PRJ-001" no tiene historias
    Cuando se visita el dashboard del proyecto "PRJ-001"
    Entonces se muestran Story Points totales en cero
    Y se muestra la velocidad promedio en "0.00"
    Y se muestra el porcentaje de historias completadas en "0.00%"
    Y se muestran las horas estimadas y reales en "0.00 h"
    Y se muestra la cantidad de Sprints cerrados en cero
    Y el gráfico de velocidad muestra el mensaje "Todavía no hay Sprints cerrados"
    Y el gráfico de horas muestra el mensaje "Todavía no hay Sprints cerrados"
    Y la página responde sin error

  @CA3 @Fase2
  Escenario: Un proyecto con Sprints solo activos no muestra cerrados en los gráficos
    Dado que hoy es "1 de febrero de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 activo
    Cuando se visita el dashboard del proyecto "PRJ-001"
    Entonces se muestra la cantidad de Sprints cerrados en cero
    Y el gráfico de velocidad muestra el mensaje "Todavía no hay Sprints cerrados"
    Y el gráfico de horas muestra el mensaje "Todavía no hay Sprints cerrados"
    Y la página responde sin error

  @CA1 @Fase2
  Escenario: La desviación se muestra con signo cuando hay datos
    Dado que hoy es "1 de febrero de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene horas estimadas y horas reales registradas
    Cuando se visita el dashboard del proyecto "PRJ-001"
    Entonces se muestra la desviación de horas con signo y dos decimales y el sufijo "%"

  @CA1 @Fase2
  Escenario: La desviación se muestra como no disponible cuando las estimaciones están incompletas
    Dado que hoy es "1 de febrero de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene al menos una historia sin estimar en horas
    Cuando se visita el dashboard del proyecto "PRJ-001"
    Entonces se muestra la desviación de horas como "—" con una nota que explica que hay historias sin estimar en horas

  @CA1 @Fase2
  Escenario: El dashboard de un proyecto inexistente responde 404
    Dado que hoy es "1 de febrero de 2026"
    Cuando se visita el dashboard del proyecto "PRJ-999"
    Entonces se muestra un mensaje de proyecto no encontrado
    Y la respuesta tiene el código de error 404
