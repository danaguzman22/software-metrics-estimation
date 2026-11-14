# language: es
# Estado: Propuesta
# La fuente de alcance es la Issue #11. Se reutilizan las métricas de Sprint,
# SPPlanificados y ListarPorProyecto incorporados en US-09.

@US-16
Característica: Consultar métricas de Story Points y velocidad
  Como Product Owner
  Quiero ver SP planificados, SP completados y velocidad
  Para predecir cuánto puede entregar el equipo

  @CA1
  Escenario: Calcular los Story Points planificados de un Sprint activo
    Dado un proyecto existente identificado por "p1"
    Y que el Sprint número 1 del proyecto "p1" está Activo
    Y que las historias asignadas al Sprint número 1 tienen 3 y 5 Story Points
    Cuando se consultan los Story Points planificados del Sprint número 1 del proyecto "p1"
    Entonces los Story Points planificados son 8

  @CA1
  Escenario: Consultar los Story Points planificados guardados de un Sprint cerrado
    Dado un proyecto existente identificado por "p1"
    Y que el Sprint número 1 del proyecto "p1" está Cerrado y tiene 8 SPPlanificados guardados
    Cuando se consultan los Story Points planificados del Sprint número 1 del proyecto "p1"
    Entonces los Story Points planificados son 8

  @CA2
  Escenario: Calcular los Story Points completados de un Sprint
    Dado un proyecto existente identificado por "p1"
    Y que el Sprint número 1 del proyecto "p1" tiene historias completadas de 3 y 5 Story Points
    Y que el Sprint número 1 del proyecto "p1" tiene una historia pendiente de 8 Story Points
    Cuando se consultan los Story Points completados del Sprint número 1 del proyecto "p1"
    Entonces los Story Points completados son 8

  @CA3
  @P3
  Escenario: Promediar únicamente Sprints cerrados e incluir el Sprint cerrado con velocidad cero
    Dado un proyecto existente identificado por "p1"
    Y que tiene dos Sprints cerrados con velocidades 0 y 5
    Y que tiene un Sprint activo con velocidad actual 8
    Cuando se consulta la velocidad promedio del proyecto "p1"
    Entonces la velocidad promedio es 2.5

  @CA4
  Escenario: Devolver velocidad cero cuando el proyecto no tiene Sprints cerrados
    Dado un proyecto existente identificado por "p1"
    Y que tiene un Sprint activo y ningún Sprint cerrado
    Cuando se consulta la velocidad promedio del proyecto "p1"
    Entonces la velocidad promedio es 0
    Y no se devuelve error

  @P4
  Escenario: Devolver error al consultar planificados de Sprint cerrado sin valor guardado
    Dado un proyecto existente identificado por "p1"
    Y que tiene un Sprint cerrado sin SPPlanificados guardados
    Cuando se consultan los Story Points planificados del Sprint cerrado del proyecto "p1"
    Entonces se devuelve un error indicando que faltan Story Points planificados guardados

  @CA5
  Escenario: Calcular cero por ciento cuando el backlog está vacío
    Dado un proyecto existente identificado por "p1" sin historias en el backlog
    Cuando se consulta el porcentaje de historias completadas del proyecto "p1"
    Entonces el porcentaje de historias completadas es 0
    Y no se devuelve error

  @P6
  Escenario: Calcular el porcentaje sobre todas las historias del backlog del proyecto
    Dado un proyecto existente identificado por "p1"
    Y que el backlog del proyecto "p1" tiene cuatro historias, una Completada
    Y que las historias no completadas ya no están asignadas a Sprints cerrados
    Cuando se consulta el porcentaje de historias completadas del proyecto "p1"
    Entonces el porcentaje de historias completadas es 25
