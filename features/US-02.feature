# language: es
# Dependencia de T-02: el tipo y el formato concreto del ID los define T-02.
# Los escenarios se refieren al ID de forma neutral y no usan el nombre como identificador.

@US-02
Característica: Modificar proyecto
  Como integrante
  Quiero modificar un proyecto existente con nombre, fecha de inicio y fecha de fin
  Para mantener sus datos actualizados

  @CA1
  Escenario: Modificación exitosa con datos válidos
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se modifica el proyecto con el nombre "Sistema de Métricas v2", la fecha de inicio "15 de marzo de 2026" y la fecha de fin "31 de julio de 2026"
    Entonces la modificación se realiza exitosamente
    Y el proyecto devuelto queda con el nombre "Sistema de Métricas v2", la fecha de inicio "15 de marzo de 2026" y la fecha de fin "31 de julio de 2026"
    Y el proyecto original conserva el nombre "Sistema de Métricas", la fecha de inicio "1 de marzo de 2026" y la fecha de fin "30 de junio de 2026"

  @CA1
  Escenario: Modificar con valores iguales a los actuales es válido
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se modifica el proyecto con el nombre "Sistema de Métricas", la fecha de inicio "1 de marzo de 2026" y la fecha de fin "30 de junio de 2026"
    Entonces la modificación se realiza exitosamente
    Y el proyecto devuelto conserva el nombre "Sistema de Métricas", la fecha de inicio "1 de marzo de 2026" y la fecha de fin "30 de junio de 2026"
    Y el proyecto devuelto conserva el estado "En curso"

  @CA2
  Escenario: Nombre vacío devuelve error de validación al modificar
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se modifica el proyecto con el nombre "", la fecha de inicio "1 de marzo de 2026" y la fecha de fin "30 de junio de 2026"
    Entonces se devuelve un error de validación de nombre
    Y la modificación no se completa exitosamente
    Y el proyecto original conserva el nombre "Sistema de Métricas", la fecha de inicio "1 de marzo de 2026" y la fecha de fin "30 de junio de 2026"

  @CA2
  Escenario: Nombre compuesto únicamente por espacios devuelve error de validación al modificar
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se modifica el proyecto con el nombre "   ", la fecha de inicio "1 de marzo de 2026" y la fecha de fin "30 de junio de 2026"
    Entonces se devuelve un error de validación de nombre
    Y la modificación no se completa exitosamente
    Y el proyecto original conserva el nombre "Sistema de Métricas", la fecha de inicio "1 de marzo de 2026" y la fecha de fin "30 de junio de 2026"

  @CA1 @CA2
  Escenario: Nombre con espacios externos se conserva al modificar
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se modifica el proyecto con el nombre "  Proyecto con espacios  ", la fecha de inicio "15 de marzo de 2026" y la fecha de fin "31 de julio de 2026"
    Entonces la modificación se realiza exitosamente
    Y el proyecto devuelto queda con el nombre "  Proyecto con espacios  "
    Y el proyecto original conserva el nombre "Sistema de Métricas"

  @CA2 @CA3
  Escenario: Nombre y fechas inválidos a la vez: prevalece el error de nombre
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se modifica el proyecto con el nombre "   ", la fecha de inicio "30 de junio de 2026" y la fecha de fin "1 de marzo de 2026"
    Entonces se devuelve un error de validación de nombre
    Y no se devuelve un error de validación de fechas
    Y la modificación no se completa exitosamente
    Y el proyecto original conserva sus datos originales

  @CA3
  Escenario: Fecha de fin anterior a la fecha de inicio devuelve error al modificar
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se modifica el proyecto con el nombre "Sistema de Métricas v2", la fecha de inicio "30 de junio de 2026" y la fecha de fin "1 de marzo de 2026"
    Entonces se devuelve un error de validación de fechas
    Y la modificación no se completa exitosamente
    Y el proyecto original conserva sus datos originales

  @CA1 @CA3
  Escenario: Fecha de fin igual a la fecha de inicio es válida al modificar
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se modifica el proyecto con el nombre "Sistema de Métricas v2", la fecha de inicio "1 de marzo de 2026" y la fecha de fin "1 de marzo de 2026"
    Entonces la modificación se realiza exitosamente
    Y el proyecto devuelto queda con la fecha de inicio "1 de marzo de 2026" y la fecha de fin "1 de marzo de 2026"

  @CA4
  Escenario: Modificar un proyecto inexistente devuelve error
    Dado que no existe un proyecto identificado por el ID del proyecto a modificar
    Cuando se intenta modificar el proyecto con el nombre "Proyecto actualizado", la fecha de inicio "1 de marzo de 2026" y la fecha de fin "30 de junio de 2026"
    Entonces se devuelve un error indicando que el proyecto no existe
    Y la modificación no se completa exitosamente

  @CA4
  Escenario: Proyecto inexistente prevalece sobre datos inválidos
    Dado que no existe un proyecto identificado por el ID del proyecto a modificar
    Cuando se intenta modificar el proyecto con el nombre "", la fecha de inicio "30 de junio de 2026" y la fecha de fin "1 de marzo de 2026"
    Entonces se devuelve un error indicando que el proyecto no existe
    Y el error no es de validación de nombre ni de fechas
    Y la modificación no se completa exitosamente

  @CA5
  Escenario: La modificación conserva el estado actual del proyecto
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se modifica el proyecto con el nombre "Sistema de Métricas v2", la fecha de inicio "15 de marzo de 2026" y la fecha de fin "31 de julio de 2026"
    Entonces la modificación se realiza exitosamente
    Y el proyecto devuelto conserva el estado "En curso"
    Y el proyecto no queda con estado "Planificado"

  @CA5
  Escenario: Se puede modificar un proyecto en estado Finalizado y se conserva el estado
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "Finalizado"
    Cuando se modifica el proyecto con el nombre "Sistema de Métricas v2", la fecha de inicio "15 de marzo de 2026" y la fecha de fin "31 de julio de 2026"
    Entonces la modificación se realiza exitosamente
    Y el proyecto devuelto conserva el estado "Finalizado"
