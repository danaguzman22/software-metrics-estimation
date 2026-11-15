# language: es
# Estado: Propuesta
# Fuente: Issue #18. Las decisiones derivadas P1-P7 están documentadas en
# specs/US-17.md para revisión humana.

@US-17
Característica: Consultar horas estimadas, reales y desviación
  Como Product Owner
  Quiero ver horas estimadas, reales y desviación
  Para evaluar la precisión de las estimaciones

  @CA1
  Escenario: Sumar horas reales de una historia
    Dado un proyecto existente identificado por "p1"
    Y que la historia "HU-01" tiene 4 horas estimadas
    Y que la historia "HU-01" tiene registros de esfuerzo de 1.5 y 2 horas
    Cuando se consultan las métricas de esfuerzo de la historia "HU-01" del proyecto "p1"
    Entonces las horas reales de la historia son 3.5

  @CA2
  Escenario: Calcular la desviación porcentual de una historia
    Dado un proyecto existente identificado por "p1"
    Y que la historia "HU-01" tiene 4 horas estimadas
    Y que la historia "HU-01" tiene registros de esfuerzo que suman 5 horas
    Cuando se consultan las métricas de esfuerzo de la historia "HU-01" del proyecto "p1"
    Entonces las horas reales son 5
    Y la desviación porcentual es 25

  @CA3
  Escenario: No definir desviación cuando las horas estimadas son cero
    Dado un proyecto existente identificado por "p1"
    Y que la historia "HU-01" tiene 0 horas estimadas
    Y que la historia "HU-01" tiene registros de esfuerzo que suman 2 horas
    Cuando se consultan las métricas de esfuerzo de la historia "HU-01" del proyecto "p1"
    Entonces las horas reales son 2
    Y la desviación no está disponible
    Y no se produce división por cero

  @CA3
  @P1
  @P4
  Escenario: Devolver reales y marcar como incompletas las métricas sin estimación
    Dado un proyecto existente identificado por "p1"
    Y que la historia "HU-01" no tiene horas estimadas
    Y que la historia "HU-01" tiene registros de esfuerzo que suman 2 horas
    Cuando se consultan las métricas de esfuerzo de la historia "HU-01" del proyecto "p1"
    Entonces las horas reales son 2
    Y las estimaciones están incompletas
    Y la desviación no está disponible

  @CA2
  @CA4
  @P7
  Escenario: Calcular la desviación de un Sprint sobre sus totales
    Dado un proyecto existente identificado por "p1"
    Y que el Sprint número 1 tiene historias con 10 horas estimadas en total
    Y que los registros atribuibles al Sprint número 1 suman 12 horas reales
    Cuando se consultan las métricas de esfuerzo del Sprint número 1 del proyecto "p1"
    Entonces las horas estimadas del Sprint son 10
    Y las horas reales del Sprint son 12
    Y la desviación porcentual del Sprint es 20
    Y la desviación se calcula una vez sobre los totales del Sprint

  @CA2
  @CA4
  @P7
  Escenario: Calcular la desviación del proyecto sobre sus totales y no promediar historias
    Dado un proyecto existente identificado por "p1"
    Y que sus historias tienen respectivamente 1 y 9 horas estimadas
    Y que sus registros de esfuerzo suman respectivamente 2 y 9 horas reales
    Cuando se consultan las métricas de esfuerzo del proyecto "p1"
    Entonces las horas estimadas del proyecto son 10
    Y las horas reales del proyecto son 11
    Y la desviación porcentual del proyecto es 10
    Y la desviación se calcula una vez sobre los totales del proyecto

  @CA1
  @CA4
  @P5
  @P6
  Escenario: Conservar métricas de Sprint después de devolver historias al backlog
    Dado que el Sprint número 1 del proyecto "p1" cerró con historias completadas y no completadas
    Y que las asignaciones actuales de las historias no completadas se eliminaron al cerrar
    Cuando se consultan las métricas históricas de esfuerzo del Sprint número 1 del proyecto "p1"
    Entonces se incluyen las historias y estimaciones que pertenecían al Sprint al cerrarse
    Y se incluyen los registros de esfuerzo atribuidos históricamente al Sprint

  @P5
  Escenario: No reconstruir métricas de un Sprint cerrado antes de US-17
    Dado que el Sprint número 1 del proyecto "p1" se cerró antes de crear snapshots de US-17
    Cuando se consultan las métricas de esfuerzo del Sprint número 1 del proyecto "p1"
    Entonces se devuelve el error de métricas históricas de Sprint no disponibles

  @P1
  Escenario: Persistir las horas estimadas de una historia
    Dado un proyecto existente identificado por "p1"
    Y que existe la historia "HU-01" en el proyecto "p1"
    Cuando se establecen 4 horas estimadas para la historia "HU-01"
    Entonces la consulta de métricas de "HU-01" devuelve 4 horas estimadas
