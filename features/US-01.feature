# language: es

@US-01
Característica: Crear proyecto
  Como integrante
  Quiero crear un proyecto con nombre, fecha de inicio y fecha de fin
  Para empezar a gestionarlo

  @CA1
  Escenario: Creación exitosa con datos válidos
    Dado un integrante que quiere crear un proyecto
    Cuando ingresa el nombre "Sistema de Métricas", la fecha de inicio "1 de marzo de 2026" y la fecha de fin "30 de junio de 2026"
    Entonces el proyecto se crea exitosamente
    Y el proyecto conserva el nombre "Sistema de Métricas", la fecha de inicio "1 de marzo de 2026" y la fecha de fin "30 de junio de 2026"

  @CA2
  Escenario: Nombre vacío devuelve error de validación
    Dado un integrante que quiere crear un proyecto
    Cuando ingresa el nombre "", la fecha de inicio "1 de marzo de 2026" y la fecha de fin "30 de junio de 2026"
    Entonces se devuelve un error de validación
    Y la creación del proyecto no se completa exitosamente

  @CA2
  Escenario: Nombre compuesto únicamente por espacios devuelve error de validación
    Dado un integrante que quiere crear un proyecto
    Cuando ingresa el nombre "   ", la fecha de inicio "1 de marzo de 2026" y la fecha de fin "30 de junio de 2026"
    Entonces se devuelve un error de validación
    Y la creación del proyecto no se completa exitosamente

  @CA3
  Escenario: Fecha de fin anterior a la fecha de inicio devuelve error
    Dado un integrante que quiere crear un proyecto
    Cuando ingresa el nombre "Sistema de Métricas", la fecha de inicio "30 de junio de 2026" y la fecha de fin "1 de marzo de 2026"
    Entonces se devuelve un error
    Y la creación del proyecto no se completa exitosamente

  @CA4
  Escenario: El proyecto creado queda con estado inicial Planificado
    Dado un integrante que quiere crear un proyecto
    Cuando ingresa el nombre "Sistema de Métricas", la fecha de inicio "1 de marzo de 2026" y la fecha de fin "30 de junio de 2026"
    Entonces el proyecto se crea exitosamente
    Y el proyecto queda con estado "Planificado"

  @CA1 @CA3
  Escenario: Fecha de fin igual a la fecha de inicio es válida
    Dado un integrante que quiere crear un proyecto
    Cuando ingresa el nombre "Sistema de Métricas", la fecha de inicio "1 de marzo de 2026" y la fecha de fin "1 de marzo de 2026"
    Entonces el proyecto se crea exitosamente
    Y el proyecto conserva la fecha de inicio "1 de marzo de 2026" y la fecha de fin "1 de marzo de 2026"
