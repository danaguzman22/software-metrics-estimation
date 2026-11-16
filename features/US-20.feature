# language: es
# Estado: Propuesta (pendiente de revisión de la persona; nada se implementa aún).
# Los escenarios describen el comportamiento acordado para CA1–CA2 y las
# decisiones D1–D7 del enunciado (dos páginas HTML imprimibles, sin PDF ni
# librerías: GET /proyectos/{id}/reporte y GET /proyectos/{id}/sprints/{numero}/reporte;
# botón Imprimir y estilos @media print; enlazadas desde el detalle del proyecto
# y desde cada Sprint; el reporte del proyecto incluye nombre y estado, historias
# planificadas y completadas con sus Story Points, totales de SP, velocidad
# promedio y porcentaje completado reutilizando ObtenerDashboard, horas estimadas
# y reales (la columna por historia se titula "Horas reales (acumuladas)") y
# desviación, Sprints cerrados con velocidad, y defectos por severidad y estado
# con su lista; el reporte del Sprint incluye Sprint Goal, fechas y estado,
# historias asignadas con SP/estado/horas, SP planificados, completados y
# velocidad, horas estimadas vs reales y defectos detectados y resueltos en el
# Sprint; los errores devuelven 404 para proyecto o Sprint inexistente o número
# de Sprint inválido, y 500 genérico con el detalle solo en el log; sin datos
# la página responde 200 con mensajes).
# FASE 1 incluye el servicio que arma los reportes (ServicioReportes.ReporteProyecto
# y ReporteSprint), reutilizando los servicios existentes, sin métodos nuevos de
# repositorio ni migraciones. FASE 2 incluye las dos páginas web; ambas fases
# van en un único PR y un único ciclo RED/GREEN/REFACTOR, con GREEN por tandas
# internas (primero el servicio, luego la web y la composición; D8). Los totales
# de defectos Detectados/Resueltos se reutilizan de US-18
# (MetricasDefectosPorProyecto/MetricasDefectosPorSprint); los conteos por
# severidad y por estado y las listas salen de ListarDefectos (US-18 no los
# expone) (D6).
# Los errores se describen por su significado; sus identificadores y códigos
# HTTP exactos están decididos en la spec (sección 8).
# Los escenarios de FASE 1 están marcados con @Fase1 y los de FASE 2 con @Fase2.
# No se escenifican el HTTP 500/405, la verificación estructural de imports, la
# no recalculación de métricas en la web, la impresión @media print, el botón
# Imprimir ni la ausencia de enlaces en la navegación: se verifican con tests
# sin escenario BDD propio.

@US-20
Característica: Reporte de proyecto o Sprint
  Como Product Owner
  Quiero un reporte imprimible de un proyecto o de un Sprint
  Para presentarlo al cliente

  @CA1 @Fase1
  Escenario: El servicio arma el reporte del proyecto con historias, métricas, esfuerzo y defectos
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 cerrado con 8 Story Points completados
    Y que el proyecto "PRJ-001" tiene el Sprint 2 cerrado con 13 Story Points completados
    Y que el proyecto "PRJ-001" tiene una historia completada con 5 Story Points
    Y que el proyecto "PRJ-001" tiene una historia pendiente con 3 Story Points
    Y que el proyecto "PRJ-001" tiene horas estimadas y horas reales registradas
    Y que el proyecto "PRJ-001" tiene 2 defectos de severidad Alta y 1 de severidad Baja
    Y que el proyecto "PRJ-001" tiene 2 defectos en estado Abierto y 1 defecto en estado Resuelto
    Cuando se consulta el reporte del proyecto "PRJ-001"
    Entonces se devuelven las historias planificadas y las completadas con sus Story Points
    Y se devuelven los totales de Story Points planificados y completados
    Y se devuelve la velocidad promedio de los Sprints cerrados
    Y se devuelve el porcentaje de historias completadas
    Y se devuelven las horas estimadas, las horas reales y su desviación
    Y se devuelven los Sprints cerrados con su velocidad
    Y se devuelven los conteos de defectos por severidad y por estado y su lista completa

  @CA1 @Fase1
  Escenario: Consultar el reporte de un proyecto inexistente devuelve proyecto no encontrado
    Cuando se consulta el reporte del proyecto "PRJ-999"
    Entonces se devuelve el error de proyecto no encontrado

  @CA2 @Fase1
  Escenario: El servicio arma el reporte del Sprint con historias, esfuerzo y defectos
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 cerrado con objetivo "Entregar el núcleo del sistema", fecha de inicio "2026-02-02" y fecha de fin "2026-02-27"
    Y que el Sprint 1 tiene dos historias asignadas: una completada con 5 Story Points y una pendiente con 3 Story Points
    Y que el Sprint 1 tiene horas estimadas y horas reales registradas
    Y que el proyecto tiene 1 defecto detectado en el Sprint 1
    Y que el proyecto tiene 1 defecto resuelto en el Sprint 1
    Y que el proyecto tiene 1 defecto detectado en el Sprint 2
    Cuando se consulta el reporte del Sprint 1 del proyecto "PRJ-001"
    Entonces se devuelve el objetivo, las fechas y el estado del Sprint
    Y se devuelven las historias asignadas con sus Story Points, su estado y sus horas (columna Horas reales acumuladas)
    Y se devuelven los Story Points planificados, los completados y la velocidad del Sprint
    Y se devuelven las horas estimadas y las horas reales del Sprint
    Y se devuelven los defectos detectados y los resueltos en el Sprint 1
    Y no se incluyen los defectos del Sprint 2

  @CA2 @Fase1
  Escenario: Consultar el reporte de un Sprint inexistente devuelve Sprint no encontrado
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Cuando se consulta el reporte del Sprint 9 del proyecto "PRJ-001"
    Entonces se devuelve el error de Sprint no encontrado

  @CA2 @Fase1
  Escenario: Un Sprint cerrado sin snapshot de esfuerzo no es un error
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 cerrado sin snapshot de esfuerzo
    Y que el Sprint 1 tiene una historia asignada con 5 Story Points
    Cuando se consulta el reporte del Sprint 1 del proyecto "PRJ-001"
    Entonces las horas del Sprint se informan como no disponibles
    Y el resto del reporte se devuelve sin error

  @CA1 @Fase2
  Escenario: La página del proyecto muestra todas las secciones del reporte
    Dado que hoy es "1 de junio de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene una historia completada con 5 Story Points
    Y que el proyecto "PRJ-001" tiene una historia pendiente con 3 Story Points
    Y que el proyecto "PRJ-001" tiene el Sprint 1 cerrado con 8 Story Points completados
    Y que el proyecto "PRJ-001" tiene 2 defectos de severidad Alta y 1 defecto en estado Resuelto
    Cuando se visita el reporte del proyecto "PRJ-001"
    Entonces se muestra el nombre y el estado del proyecto
    Y se muestran las historias planificadas y las completadas con sus Story Points y sus horas reales (acumuladas)
    Y se muestran los totales de Story Points, la velocidad promedio y el porcentaje de historias completadas
    Y se muestran las horas estimadas, las horas reales y su desviación
    Y se muestra la tabla de Sprints cerrados con su velocidad
    Y se muestran los defectos contabilizados por severidad y por estado y su lista
    Y la página responde sin error

  @CA2 @Fase2
  Escenario: La página del Sprint muestra todas las secciones del reporte
    Dado que hoy es "1 de junio de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 cerrado con objetivo "Entregar el núcleo del sistema", fecha de inicio "2026-02-02" y fecha de fin "2026-02-27"
    Y que el Sprint 1 tiene una historia completada con 5 Story Points y 40 horas estimadas
    Y que el Sprint 1 tiene horas reales registradas
    Y que el proyecto tiene 1 defecto detectado y 1 defecto resuelto en el Sprint 1
    Cuando se visita el reporte del Sprint 1 del proyecto "PRJ-001"
    Entonces se muestra el Sprint Goal, las fechas y el estado del Sprint
    Y se muestran las historias asignadas con sus Story Points, su estado y sus horas (columna Horas reales acumuladas)
    Y se muestran los Story Points planificados, los completados y la velocidad del Sprint
    Y se muestran las horas estimadas y las horas reales del Sprint
    Y se muestran los defectos detectados y los resueltos en el Sprint
    Y la página responde sin error

  @CA1 @Fase2
  Escenario: El detalle del proyecto enlaza al reporte
    Dado que hoy es "1 de junio de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Cuando se visita el detalle del proyecto "PRJ-001"
    Entonces se ofrece un enlace al reporte del proyecto

  @CA2 @Fase2
  Escenario: Cada Sprint enlaza a su reporte
    Dado que hoy es "1 de junio de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 activo
    Cuando se visita el listado de Sprints del proyecto "PRJ-001"
    Entonces se ofrece un enlace al reporte del Sprint 1

  @CA2 @Fase2
  Escenario: El reporte de un Sprint inexistente responde 404
    Dado que hoy es "1 de junio de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Cuando se visita el reporte del Sprint 9 del proyecto "PRJ-001"
    Entonces se muestra un mensaje de Sprint no encontrado
    Y la respuesta tiene el código de error 404

  @CA2 @Fase2
  Escenario: Un número de Sprint inválido responde 404
    Dado que hoy es "1 de junio de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Cuando se visita el reporte del Sprint "abc" del proyecto "PRJ-001"
    Entonces la respuesta tiene el código de error 404

  @CA1 @Fase2
  Escenario: El reporte de un proyecto inexistente responde 404
    Dado que hoy es "1 de junio de 2026"
    Cuando se visita el reporte del proyecto "PRJ-999"
    Entonces se muestra un mensaje de proyecto no encontrado
    Y la respuesta tiene el código de error 404

  @CA2 @Fase2
  Escenario: Un Sprint cerrado sin snapshot de esfuerzo muestra las métricas no disponibles
    Dado que hoy es "1 de junio de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 cerrado sin snapshot de esfuerzo
    Y que el Sprint 1 tiene una historia asignada con 5 Story Points
    Cuando se visita el reporte del Sprint 1 del proyecto "PRJ-001"
    Entonces el bloque de horas del Sprint muestra el mensaje "métricas no disponibles"
    Y la página responde sin error

  @CA1 @Fase2
  Escenario: Un proyecto sin historias muestra el mensaje en el reporte
    Dado que hoy es "1 de junio de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" no tiene historias
    Cuando se visita el reporte del proyecto "PRJ-001"
    Entonces se muestra el mensaje "No hay historias" en la lista de historias planificadas
    Y se muestra el mensaje "No hay historias" en la lista de historias completadas
    Y la página responde sin error

  @CA1 @Fase2
  Escenario: Un proyecto sin defectos muestra el mensaje en el reporte
    Dado que hoy es "1 de junio de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" no tiene defectos
    Cuando se visita el reporte del proyecto "PRJ-001"
    Entonces se muestra el mensaje "No hay defectos"
    Y los conteos de defectos se muestran en cero
    Y la página responde sin error

  @CA2 @Fase2
  Escenario: Un Sprint sin historias muestra el mensaje en el reporte
    Dado que hoy es "1 de junio de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 activo sin historias asignadas
    Cuando se visita el reporte del Sprint 1 del proyecto "PRJ-001"
    Entonces se muestra el mensaje "No hay historias asignadas"
    Y los Story Points planificados, completados y la velocidad se muestran en cero
    Y la página responde sin error

  @CA2 @Fase2
  Escenario: Un Sprint sin defectos muestra el mensaje en el reporte
    Dado que hoy es "1 de junio de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 cerrado con la fecha de inicio "2026-02-02" y la fecha de fin "2026-02-27"
    Y que el proyecto "PRJ-001" no tiene defectos
    Cuando se visita el reporte del Sprint 1 del proyecto "PRJ-001"
    Entonces se muestra el mensaje "No hay defectos"
    Y los defectos detectados y los resueltos en el Sprint se muestran en cero
    Y la página responde sin error

  @CA1 @Fase2
  Escenario: La desviación se muestra con signo o como no disponible
    Dado que hoy es "1 de junio de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene horas estimadas y horas reales registradas
    Cuando se visita el reporte del proyecto "PRJ-001"
    Entonces se muestra la desviación de horas con signo y dos decimales y el sufijo "%"
    Y la página responde sin error