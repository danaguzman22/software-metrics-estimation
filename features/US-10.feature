# language: es
# Estado: Propuesta
# Fuente: Issue #19. Las decisiones derivadas P1-P5 están documentadas en
# specs/US-10.md para revisión humana.

@US-10
Característica: Consultar Sprints anteriores
  Como equipo
  Quiero ver el histórico de Sprints cerrados
  Para analizar la evolución

  @CA1
  @P1
  @P2
  @P3
  Escenario: Consultar Sprint Goals y métricas de los Sprints cerrados
    Dado un proyecto existente con Sprints cerrados y activos
    Y los Sprints cerrados tienen Sprint Goal, SP planificados, SP completados y velocidad
    Cuando el equipo consulta el histórico de Sprints del proyecto
    Entonces la lista incluye únicamente los Sprints cerrados
    Y cada Sprint muestra su Sprint Goal, SP planificados, SP completados y velocidad
    Y los Sprints se muestran en orden ascendente por número

  @CA2
  @P4
  Escenario: Consultar el historial cuando no hay Sprints cerrados
    Dado un proyecto existente sin Sprints cerrados
    Cuando el equipo consulta el histórico de Sprints del proyecto
    Entonces se muestra una lista vacía sin error
