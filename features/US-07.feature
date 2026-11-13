# language: es
# Estado: Propuesta
# Los escenarios describen el comportamiento acordado para CA1–CA3 y las
# decisiones D1–D9: el Sprint pertenece a un proyecto existente; se crea con
# número, Sprint Goal y fechas; el objetivo se recorta y no puede quedar vacío;
# el número debe ser mayor que cero y es único dentro del proyecto; la fecha de
# fin no puede ser anterior a la de inicio comparando por año, mes y día; todo
# Sprint nace en estado "Activo" y no puede haber dos Sprints activos en el
# mismo proyecto. El orden de evaluación es: proyecto, número, objetivo, fechas,
# Sprint activo y número duplicado.
# Los errores se describen por su significado; sus identificadores exactos están
# decididos en la spec (sección 8).
# El escenario de número duplicado se plantea con un Sprint ya registrado que
# no está activo, porque si el Sprint existente estuviera activo prevalece el
# error de Sprint activo (D7). El cierre de Sprints y los estados distintos de
# "Activo" quedan fuera de alcance (D2).
# No se escenifica la defensa de concurrencia ni el contrato de errores de
# repositorio: se verifican con tests sin escenario BDD propio.

@US-07
Característica: Crear un Sprint
  Como equipo
  Quiero crear un Sprint con su objetivo y fechas
  Para planificar una iteración

  @CA1
  Escenario: Crear un Sprint con número, objetivo y fechas válidas
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto no tiene Sprints
    Cuando se crea el Sprint número 1 con objetivo "Planificar la primera iteración", fecha de inicio "2 de marzo de 2026" y fecha de fin "20 de marzo de 2026"
    Entonces el Sprint se crea exitosamente
    Y el Sprint queda con número 1, objetivo "Planificar la primera iteración", fecha de inicio "2 de marzo de 2026", fecha de fin "20 de marzo de 2026" y estado "Activo"
    Y el proyecto queda con un Sprint registrado

  @CA1
  Escenario: Recortar los espacios externos del objetivo
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto no tiene Sprints
    Cuando se crea el Sprint número 1 con objetivo "   Planificar la primera iteración   ", fecha de inicio "2 de marzo de 2026" y fecha de fin "20 de marzo de 2026"
    Entonces el Sprint se crea exitosamente
    Y el Sprint queda con objetivo "Planificar la primera iteración" y estado "Activo"

  @CA1
  Escenario: Número cero devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto no tiene Sprints
    Cuando se intenta crear el Sprint número 0 con objetivo "Planificar la primera iteración", fecha de inicio "2 de marzo de 2026" y fecha de fin "20 de marzo de 2026"
    Entonces se devuelve un error indicando que el número del Sprint debe ser mayor que cero
    Y el Sprint no se crea
    Y el proyecto no queda con ningún Sprint

  @CA1
  Escenario: Número negativo devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto no tiene Sprints
    Cuando se intenta crear el Sprint número -1 con objetivo "Planificar la primera iteración", fecha de inicio "2 de marzo de 2026" y fecha de fin "20 de marzo de 2026"
    Entonces se devuelve un error indicando que el número del Sprint debe ser mayor que cero
    Y el Sprint no se crea
    Y el proyecto no queda con ningún Sprint

  @CA1
  Escenario: Número duplicado en el mismo proyecto devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto ya tiene registrado el número 1 en un Sprint que no está activo
    Y el proyecto no tiene ningún Sprint activo
    Cuando se intenta crear el Sprint número 1 con objetivo "Planificar una nueva iteración", fecha de inicio "2 de marzo de 2026" y fecha de fin "20 de marzo de 2026"
    Entonces se devuelve un error indicando que el número del Sprint ya existe en el proyecto
    Y el Sprint no se crea
    Y el proyecto sigue teniendo un solo Sprint con número 1

  @CA1
  Escenario: El mismo número en otro proyecto es válido
    Dado un proyecto existente identificado por su ID, con nombre "Portal de Clientes", fecha de inicio "5 de enero de 2026", fecha de fin "30 de abril de 2026" y estado "En curso"
    Y que ese proyecto tiene un Sprint con número 1, objetivo "Planificar el portal", fecha de inicio "12 de enero de 2026", fecha de fin "30 de enero de 2026" y estado "Activo"
    Y un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto "Sistema de Métricas" no tiene Sprints
    Cuando se crea el Sprint número 1 en el proyecto "Sistema de Métricas" con objetivo "Planificar la primera iteración", fecha de inicio "2 de marzo de 2026" y fecha de fin "20 de marzo de 2026"
    Entonces el Sprint se crea exitosamente
    Y el Sprint del proyecto "Sistema de Métricas" queda con número 1 y estado "Activo"
    Y el Sprint del proyecto "Portal de Clientes" conserva su objetivo "Planificar el portal"

  @CA1
  Escenario: Fecha de fin anterior a la de inicio devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto no tiene Sprints
    Cuando se intenta crear el Sprint número 1 con objetivo "Planificar la primera iteración", fecha de inicio "20 de marzo de 2026" y fecha de fin "2 de marzo de 2026"
    Entonces se devuelve un error indicando que la fecha de fin del Sprint no puede ser anterior a la fecha de inicio
    Y el Sprint no se crea
    Y el proyecto no queda con ningún Sprint

  @CA1
  Escenario: Fecha de inicio igual a la de fin es válida
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto no tiene Sprints
    Cuando se crea el Sprint número 2 con objetivo "Realizar una iteración de un día", fecha de inicio "10 de marzo de 2026" y fecha de fin "10 de marzo de 2026"
    Entonces el Sprint se crea exitosamente
    Y el Sprint queda con fecha de inicio "10 de marzo de 2026", fecha de fin "10 de marzo de 2026" y estado "Activo"

  @CA1
  Escenario: Fechas del mismo día con distinta hora son válidas
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto no tiene Sprints
    Cuando se crea el Sprint número 2 con objetivo "Realizar una iteración de un día", fecha de inicio "10 de marzo de 2026 a las 18:00" y fecha de fin "10 de marzo de 2026 a las 09:00"
    Entonces el Sprint se crea exitosamente
    Y el Sprint queda con fecha de inicio "10 de marzo de 2026" y fecha de fin "10 de marzo de 2026"

  @CA1
  Escenario: Sprint fuera de las fechas del proyecto es válido
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto no tiene Sprints
    Cuando se crea el Sprint número 3 con objetivo "Planificar la iteración de julio", fecha de inicio "1 de julio de 2026" y fecha de fin "15 de julio de 2026"
    Entonces el Sprint se crea exitosamente
    Y el Sprint queda con fecha de inicio "1 de julio de 2026", fecha de fin "15 de julio de 2026" y estado "Activo"

  @CA1
  Escenario: Proyecto inexistente devuelve error
    Dado que no existe un proyecto identificado por el ID del proyecto donde se quiere crear el Sprint
    Cuando se intenta crear el Sprint número 1 con objetivo "Planificar la primera iteración", fecha de inicio "2 de marzo de 2026" y fecha de fin "20 de marzo de 2026"
    Entonces se devuelve un error indicando que el proyecto no existe
    Y el Sprint no se crea

  @CA2
  Escenario: Objetivo vacío devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto no tiene Sprints
    Cuando se intenta crear el Sprint número 1 con objetivo "", fecha de inicio "2 de marzo de 2026" y fecha de fin "20 de marzo de 2026"
    Entonces se devuelve un error indicando que el objetivo del Sprint no puede estar vacío
    Y el Sprint no se crea
    Y el proyecto no queda con ningún Sprint

  @CA2
  Escenario: Objetivo compuesto solo por espacios devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto no tiene Sprints
    Cuando se intenta crear el Sprint número 1 con objetivo "   ", fecha de inicio "2 de marzo de 2026" y fecha de fin "20 de marzo de 2026"
    Entonces se devuelve un error indicando que el objetivo del Sprint no puede estar vacío
    Y el Sprint no se crea
    Y el proyecto no queda con ningún Sprint

  @CA3
  Escenario: No se puede crear un segundo Sprint activo en el proyecto
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un Sprint con número 1, objetivo "Planificar la primera iteración", fecha de inicio "2 de marzo de 2026", fecha de fin "20 de marzo de 2026" y estado "Activo"
    Cuando se intenta crear el Sprint número 2 con objetivo "Planificar la segunda iteración", fecha de inicio "23 de marzo de 2026" y fecha de fin "10 de abril de 2026"
    Entonces se devuelve un error indicando que el proyecto ya tiene un Sprint activo
    Y el Sprint no se crea
    Y el proyecto sigue teniendo un solo Sprint activo

  @CA3
  Escenario: Un Sprint activo en otro proyecto no interfiere
    Dado un proyecto existente identificado por su ID, con nombre "Portal de Clientes", fecha de inicio "5 de enero de 2026", fecha de fin "30 de abril de 2026" y estado "En curso"
    Y que ese proyecto tiene un Sprint con número 7, objetivo "Planificar el portal", fecha de inicio "12 de enero de 2026", fecha de fin "30 de enero de 2026" y estado "Activo"
    Y un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto "Sistema de Métricas" no tiene Sprints
    Cuando se crea el Sprint número 2 en el proyecto "Sistema de Métricas" con objetivo "Planificar la primera iteración", fecha de inicio "2 de marzo de 2026" y fecha de fin "20 de marzo de 2026"
    Entonces el Sprint se crea exitosamente
    Y el Sprint del proyecto "Sistema de Métricas" queda con número 2 y estado "Activo"
    Y el Sprint activo del proyecto "Portal de Clientes" conserva su número 7

  @CA1
  Escenario: Proyecto inexistente con número inválido devuelve primero el error de proyecto
    Dado que no existe un proyecto identificado por el ID del proyecto donde se quiere crear el Sprint
    Cuando se intenta crear el Sprint número 0 con objetivo "Planificar la primera iteración", fecha de inicio "2 de marzo de 2026" y fecha de fin "20 de marzo de 2026"
    Entonces se devuelve un error indicando que el proyecto no existe
    Y no se devuelve un error de número inválido
    Y el Sprint no se crea

  @CA1
  Escenario: Número inválido con objetivo vacío devuelve primero el error de número
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto no tiene Sprints
    Cuando se intenta crear el Sprint número 0 con objetivo "", fecha de inicio "2 de marzo de 2026" y fecha de fin "20 de marzo de 2026"
    Entonces se devuelve un error indicando que el número del Sprint debe ser mayor que cero
    Y no se devuelve un error de objetivo vacío
    Y el Sprint no se crea

  @CA2
  Escenario: Objetivo vacío con fechas inválidas devuelve primero el error de objetivo
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto no tiene Sprints
    Cuando se intenta crear el Sprint número 1 con objetivo "   ", fecha de inicio "20 de marzo de 2026" y fecha de fin "2 de marzo de 2026"
    Entonces se devuelve un error indicando que el objetivo del Sprint no puede estar vacío
    Y no se devuelve un error de fechas inválidas
    Y el Sprint no se crea

  @CA3
  Escenario: Sprint activo existente con número duplicado devuelve primero el error de Sprint activo
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un Sprint con número 1, objetivo "Planificar la primera iteración", fecha de inicio "2 de marzo de 2026", fecha de fin "20 de marzo de 2026" y estado "Activo"
    Cuando se intenta crear el Sprint número 1 con objetivo "Planificar otra iteración", fecha de inicio "23 de marzo de 2026" y fecha de fin "10 de abril de 2026"
    Entonces se devuelve un error indicando que el proyecto ya tiene un Sprint activo
    Y no se devuelve un error de número duplicado
    Y el Sprint no se crea
