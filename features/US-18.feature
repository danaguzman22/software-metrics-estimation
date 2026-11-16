# language: es
# Estado: Propuesta
# Fuente: Issue #21. Las decisiones derivadas D1-D5 se documentan en
# specs/US-18.md; fueron confirmadas por Victoria y esperan revisión de Dana.

@US-18
Característica: Consultar métricas de defectos
  Como Product Owner
  Quiero ver defectos detectados y resueltos
  Para conocer la calidad del producto

  @CA1
  @D1
  @D2
  Escenario: Contar defectos por Sprint de detección y de resolución
    Dado un proyecto con los Sprints 1 y 2
    Y un defecto detectado en el Sprint 1 y resuelto en el Sprint 2
    Cuando se consultan las métricas del Sprint 1
    Entonces el Sprint 1 informa 1 defecto detectado y 0 resueltos
    Cuando se consultan las métricas del Sprint 2
    Entonces el Sprint 2 informa 0 defectos detectados y 1 resuelto

  @CA1
  @D3
  Escenario: Consultar los totales de defectos del proyecto
    Dado un proyecto con 3 defectos registrados
    Y que 2 defectos están en estado "Resuelto" con Sprint de resolución informado
    Cuando se consultan las métricas del proyecto
    Entonces el proyecto informa 3 defectos detectados y 2 resueltos

  @CA2
  @D4
  Escenario: Devolver cero para un Sprint sin defectos
    Dado un proyecto con un Sprint sin defectos detectados ni resueltos
    Cuando se consultan las métricas de ese Sprint
    Entonces el Sprint informa 0 defectos detectados y 0 resueltos

  @CA2
  @D3
  Escenario: Devolver cero para un proyecto sin defectos
    Dado un proyecto sin defectos registrados
    Cuando se consultan las métricas del proyecto
    Entonces el proyecto informa 0 defectos detectados y 0 resueltos

  @D4
  @D5
  Escenario: Mostrar todos los Sprints con sus métricas en la página de defectos
    Dado un proyecto con dos Sprints y ningún defecto
    Cuando se abre la página de defectos del proyecto
    Entonces se muestran las métricas de ambos Sprints en orden numérico
    Y cada Sprint informa 0 defectos detectados y 0 resueltos
    Y se muestran los totales del proyecto con valor 0
