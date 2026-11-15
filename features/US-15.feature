# language: es
# Estado: Aprobada por Victoria; pendiente de revisión de Dana
# Fuente: Issue #20. Las decisiones derivadas D1-D5 están documentadas en
# specs/US-15.md y fueron aprobadas por Victoria.

@US-15
Característica: Registrar y seguir defectos
  Como equipo
  Quiero registrar y seguir defectos
  Para medir la calidad

  @CA1
  @D1
  @D2
  @D3
  Escenario: Registrar un defecto con todas sus referencias
    Dado un proyecto con la historia "HU-01" y los Sprints 1 y 2
    Cuando se registra un defecto con descripción "No guarda los cambios", severidad "Alta", estado "Abierto", Sprint de detección 1 y Sprint de resolución 2
    Entonces el defecto queda asociado a la historia "HU-01" y al proyecto
    Y conserva el Sprint de detección 1 y el Sprint de resolución 2

  @CA1
  @D4
  Escenario: Registrar un defecto no resuelto sin Sprint de resolución
    Dado un proyecto con la historia "HU-01" y el Sprint 1
    Cuando se registra un defecto con descripción "La vista no carga", severidad "Media", estado "Abierto", Sprint de detección 1 y sin Sprint de resolución
    Entonces el defecto se registra con Sprint de resolución no informado

  @D5
  Escenario: Consultar los defectos registrados del proyecto
    Dado un proyecto sin defectos registrados
    Cuando se consulta la lista de defectos del proyecto
    Entonces se muestra una lista vacía
    Cuando se registra un defecto en el proyecto
    Y se vuelve a consultar la lista de defectos del proyecto
    Entonces se muestra el defecto registrado en orden de creación

  @CA2
  @D3
  Escenario: Rechazar un Sprint de resolución anterior al de detección
    Dado un proyecto con la historia "HU-01" y los Sprints 1 y 2
    Cuando se registra un defecto detectado en el Sprint 2 y resuelto en el Sprint 1
    Entonces se devuelve el error de Sprint de resolución anterior a detección
    Y el defecto no se persiste

  @CA3
  @D4
  Escenario: Exigir Sprint de resolución al marcar un defecto como resuelto
    Dado un defecto abierto detectado en el Sprint 1
    Cuando se intenta cambiar su estado a "Resuelto" sin indicar Sprint de resolución
    Entonces se devuelve el error de Sprint de resolución obligatorio
    Y el defecto conserva su estado anterior

  @CA3
  @D4
  Escenario: Resolver un defecto informando un Sprint válido
    Dado un defecto abierto detectado en el Sprint 1
    Y que existe el Sprint 2 en el proyecto
    Cuando se cambia el estado del defecto a "Resuelto" e informa el Sprint de resolución 2
    Entonces el defecto queda resuelto en el Sprint 2

  @D4
  Escenario: Reabrir un defecto resuelto
    Dado un defecto resuelto en el Sprint 2
    Cuando se cambia su estado a "En curso" sin Sprint de resolución
    Entonces el defecto queda en estado "En curso"
    Y se limpia el Sprint de resolución

  @CA4
  @D1
  Escenario: Rechazar severidad no válida
    Dado un proyecto con la historia "HU-01" y el Sprint 1
    Cuando se registra un defecto con severidad "Crítica"
    Entonces se devuelve el error de severidad inválida
    Y el defecto no se persiste

  @CA4
  @D2
  Escenario: Rechazar estado no válido
    Dado un proyecto con la historia "HU-01" y el Sprint 1
    Cuando se registra un defecto con estado "Cerrado"
    Entonces se devuelve el error de estado inválido
    Y el defecto no se persiste
