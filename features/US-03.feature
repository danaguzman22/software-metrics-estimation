# language: es
# US-03 implementada: dominio en internal/domain/integrante.go, servicio en
# internal/service/integrantes.go, tests en internal/domain/integrante_test.go e
# internal/service/integrantes_test.go; evidencias en docs/evidencias/US-03-red.txt y
# docs/evidencias/US-03-green.txt. Los escenarios describen el comportamiento acordado
# para CA1–CA4
# y las decisiones confirmadas (proyecto existente identificado por ID; duplicado por nombre
# ignorando mayúsculas y espacios al inicio y al final, sin que el rol participe y con los
# espacios internos como distintos; mismo nombre en otro proyecto permitido; nombre o rol
# solo espacios = vacío; nombre y rol se guardan sin espacios extremos; orden de errores:
# proyecto, nombre, rol, duplicado).
# El rol es texto libre: cualquier texto no vacío se acepta (no hay catálogo de roles).
# Los errores se describen por su significado; sus identificadores exactos están decididos
# en la spec (D6).
# No quedan pendientes abiertos. Las decisiones cerradas "no se decide / no se valida"
# (ID propio, límites de longitud/formato, estado del proyecto, cantidad de integrantes)
# no se escenifican, como tampoco nada fuera de alcance (modificar o eliminar integrantes,
# esfuerzo, votos, interfaz web).

@US-03
Característica: Registrar integrantes del proyecto
  Como integrante
  Quiero registrar a los miembros del equipo con su rol
  Para asignarles esfuerzo y votos

  @CA1
  Escenario: Registro exitoso de un integrante con nombre y rol
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se registra en ese proyecto la integrante "Dana Guzmán" con el rol "Desarrollador"
    Entonces el registro se realiza exitosamente
    Y la integrante "Dana Guzmán" queda registrada en el proyecto "Sistema de Métricas" con el rol "Desarrollador"

  @CA1
  Escenario: Se puede registrar un segundo integrante en el mismo proyecto
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene registrada a la integrante "Dana Guzmán" con el rol "Desarrollador"
    Cuando se registra en ese proyecto la integrante "Victoria Caro" con el rol "Revisor"
    Entonces el registro se realiza exitosamente
    Y el proyecto "Sistema de Métricas" queda con las integrantes "Dana Guzmán" y "Victoria Caro"

  @CA2
  Escenario: Integrante duplicado exacto devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene registrada a la integrante "Dana Guzmán" con el rol "Desarrollador"
    Cuando se intenta registrar en ese proyecto la integrante "Dana Guzmán" con el rol "Desarrollador"
    Entonces se devuelve un error indicando que la integrante ya está registrada en el proyecto
    Y el registro no se completa exitosamente
    Y el proyecto "Sistema de Métricas" sigue teniendo únicamente a la integrante "Dana Guzmán" con el rol "Desarrollador"

  @CA2
  Escenario: Duplicado ignorando mayúsculas
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene registrada a la integrante "Dana Guzmán" con el rol "Desarrollador"
    Cuando se intenta registrar en ese proyecto la integrante "dana guzmán" con el rol "Revisor"
    Entonces se devuelve un error indicando que la integrante ya está registrada en el proyecto
    Y el registro no se completa exitosamente
    Y el proyecto "Sistema de Métricas" sigue teniendo únicamente a la integrante "Dana Guzmán" con el rol "Desarrollador"

  @CA2
  Escenario: Duplicado ignorando espacios al inicio y al final
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene registrada a la integrante "Dana Guzmán" con el rol "Desarrollador"
    Cuando se intenta registrar en ese proyecto la integrante "  Dana Guzmán  " con el rol "Desarrollador"
    Entonces se devuelve un error indicando que la integrante ya está registrada en el proyecto
    Y el registro no se completa exitosamente
    Y el proyecto "Sistema de Métricas" sigue teniendo únicamente a la integrante "Dana Guzmán" con el rol "Desarrollador"

  @CA2
  Escenario: El mismo nombre con distinto rol también es duplicado
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene registrada a la integrante "Dana Guzmán" con el rol "Desarrollador"
    Cuando se intenta registrar en ese proyecto la integrante "Dana Guzmán" con el rol "Product Owner"
    Entonces se devuelve un error indicando que la integrante ya está registrada en el proyecto
    Y el registro no se completa exitosamente
    Y el proyecto "Sistema de Métricas" sigue teniendo únicamente a la integrante "Dana Guzmán" con el rol "Desarrollador"

  @CA2
  Escenario: Nombres con espacios internos distintos no son duplicado
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene registrada a la integrante "Ana Paz" con el rol "Revisor"
    Cuando se registra en ese proyecto la integrante "AnaPaz" con el rol "Desarrollador"
    Entonces el registro se realiza exitosamente
    Y el proyecto "Sistema de Métricas" queda con las integrantes "Ana Paz" y "AnaPaz"

  @CA2
  Escenario: El mismo nombre en otro proyecto sí se permite
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene registrada a la integrante "Dana Guzmán" con el rol "Desarrollador"
    Y un proyecto existente identificado por su ID, con nombre "App Móvil", fecha de inicio "1 de abril de 2026", fecha de fin "31 de diciembre de 2026" y estado "Planificado"
    Cuando se registra en el proyecto "App Móvil" la integrante "Dana Guzmán" con el rol "Desarrollador"
    Entonces el registro se realiza exitosamente
    Y la integrante "Dana Guzmán" queda registrada en el proyecto "App Móvil"
    Y el proyecto "Sistema de Métricas" sigue teniendo a la integrante "Dana Guzmán" con el rol "Desarrollador"

  @CA3
  Escenario: Nombre vacío devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se intenta registrar en ese proyecto la integrante con el nombre "" y el rol "Desarrollador"
    Entonces se devuelve un error indicando que el nombre no puede estar vacío
    Y el registro no se completa exitosamente
    Y el proyecto "Sistema de Métricas" no queda con integrantes registrados

  @CA3
  Escenario: Nombre compuesto únicamente por espacios devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se intenta registrar en ese proyecto la integrante con el nombre "   " y el rol "Desarrollador"
    Entonces se devuelve un error indicando que el nombre no puede estar vacío
    Y el registro no se completa exitosamente
    Y el proyecto "Sistema de Métricas" no queda con integrantes registrados

  @CA4
  Escenario: Rol vacío devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se intenta registrar en ese proyecto la integrante "Dana Guzmán" con el rol ""
    Entonces se devuelve un error indicando que el rol no puede estar vacío
    Y el registro no se completa exitosamente
    Y el proyecto "Sistema de Métricas" no queda con integrantes registrados

  @CA4
  Escenario: Rol compuesto únicamente por espacios devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se intenta registrar en ese proyecto la integrante "Dana Guzmán" con el rol "   "
    Entonces se devuelve un error indicando que el rol no puede estar vacío
    Y el registro no se completa exitosamente
    Y el proyecto "Sistema de Métricas" no queda con integrantes registrados

  @CA1 @CA3
  Escenario: Nombre con espacios externos no se considera vacío
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se registra en ese proyecto la integrante "  Ana Paz  " con el rol "Revisor"
    Entonces el registro se realiza exitosamente
    Y la integrante queda registrada en el proyecto "Sistema de Métricas" con el rol "Revisor"
    Y su nombre queda almacenado como "Ana Paz" (sin espacios al inicio y al final)

  @CA1
  Escenario: Proyecto inexistente devuelve error
    Dado que no existe un proyecto identificado por el ID del proyecto al que se quiere registrar la integrante
    Cuando se intenta registrar en ese proyecto la integrante "Dana Guzmán" con el rol "Desarrollador"
    Entonces se devuelve un error indicando que el proyecto no existe
    Y el registro no se completa exitosamente

  @CA1
  Escenario: Proyecto inexistente con nombre vacío devuelve error de proyecto inexistente
    Dado que no existe un proyecto identificado por el ID del proyecto al que se quiere registrar la integrante
    Cuando se intenta registrar en ese proyecto la integrante con el nombre "" y el rol "Desarrollador"
    Entonces se devuelve un error indicando que el proyecto no existe
    Y el registro no se completa exitosamente
