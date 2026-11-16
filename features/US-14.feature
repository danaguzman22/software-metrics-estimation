# language: es
# Estado: Propuesta
# La fuente de alcance es la Issue #17. Los escenarios describen las reglas
# decididas en specs/US-14.md. El nombre del error de fecha futura (D28) sigue
# DERIVADO. No se escenifican R6, el CHECK de la base, la migración, los
# listados ni las horas estimadas.

@US-14
Característica: Registrar el esfuerzo realizado en una historia
  Como integrante
  Quiero registrar las horas trabajadas en una historia
  Para comparar el esfuerzo estimado con el real

  @CA1
  @R1
  Escenario: Registrar un esfuerzo válido en una historia del proyecto
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión"
    Y que el integrante "Ana Pérez" está registrado en el proyecto "p1"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-01" del proyecto "p1" con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "Análisis" y 1.5 horas
    Entonces se registra el esfuerzo con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "Análisis" y 1.5 horas
    Y el esfuerzo queda asociado al proyecto "p1" y a la historia "HU-01"

  @CA1
  @R9
  Escenario: Guardar vacía una actividad vacía
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión"
    Y que el integrante "Ana Pérez" está registrado en el proyecto "p1"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-01" del proyecto "p1" con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "" y 1.5 horas
    Entonces se registra el esfuerzo con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "" y 1.5 horas

  @CA1
  @R9
  Escenario: Guardar vacía una actividad compuesta solo por espacios
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión"
    Y que el integrante "Ana Pérez" está registrado en el proyecto "p1"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-01" del proyecto "p1" con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "   " y 1.5 horas
    Entonces se registra el esfuerzo con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "" y 1.5 horas

  @CA1
  @R9
  Escenario: Recortar los espacios externos de la actividad
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión"
    Y que el integrante "Ana Pérez" está registrado en el proyecto "p1"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-01" del proyecto "p1" con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "  Análisis  " y 1.5 horas
    Entonces se registra el esfuerzo con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "Análisis" y 1.5 horas

  @CA1
  @R1
  @R4
  Escenario: Guardar el nombre del integrante tal como está registrado aunque la entrada difiera en mayúsculas y espacios externos
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión"
    Y que el integrante "Ana Pérez" está registrado en el proyecto "p1"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-01" del proyecto "p1" con integrante "  ana pérez  ", fecha "15 de noviembre de 2026", actividad "Análisis" y 1.5 horas
    Entonces se registra el esfuerzo con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "Análisis" y 1.5 horas

  @CA1
  @R8
  Esquema del escenario: La fecha de hoy se registra sin importar el horario
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión"
    Y que el integrante "Ana Pérez" está registrado en el proyecto "p1"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-01" del proyecto "p1" con integrante "Ana Pérez", fecha "15 de noviembre de 2026 a las <hora>", actividad "Análisis" y 1.5 horas
    Entonces se registra el esfuerzo con integrante "Ana Pérez", actividad "Análisis" y 1.5 horas
    Y la fecha del esfuerzo corresponde al 15 de noviembre de 2026

    Ejemplos:
      | hora  |
      | 00:00 |
      | 12:30 |
      | 23:59 |

  @CA1
  @R8
  Esquema del escenario: La fecha de mañana devuelve el error de fecha futura sin importar el horario
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión"
    Y que el integrante "Ana Pérez" está registrado en el proyecto "p1"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-01" del proyecto "p1" con integrante "Ana Pérez", fecha "16 de noviembre de 2026 a las <hora>", actividad "Análisis" y 1.5 horas
    Entonces se devuelve un error indicando que la fecha del esfuerzo es futura
    Y no se guarda ningún registro de esfuerzo

    Ejemplos:
      | hora  |
      | 00:00 |
      | 12:30 |
      | 23:59 |

  @CA2
  @R5
  Esquema del escenario: Registrar horas válidas dentro del rango permitido
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión"
    Y que el integrante "Ana Pérez" está registrado en el proyecto "p1"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-01" del proyecto "p1" con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "Análisis" y <horas> horas
    Entonces se registra el esfuerzo con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "Análisis" y <horas> horas

    Ejemplos:
      | horas |
      | 0.01  |
      | 1.5   |
      | 24    |

  @CA2
  @R5
  Esquema del escenario: Las horas fuera del rango o con más de dos decimales devuelven el error de horas inválidas
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión"
    Y que el integrante "Ana Pérez" está registrado en el proyecto "p1"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-01" del proyecto "p1" con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "Análisis" y <horas> horas
    Entonces se devuelve un error indicando que las horas del esfuerzo son inválidas
    Y no se guarda ningún registro de esfuerzo

    Ejemplos:
      | horas |
      | 0     |
      | -1    |
      | 24.01 |
      | 25    |
      | 1.234 |

  @CA3
  @R2
  Escenario: No registrar esfuerzo en un proyecto inexistente
    Dado que no existe un proyecto identificado por "p999"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-01" del proyecto "p999" con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "Análisis" y 1.5 horas
    Entonces se devuelve un error indicando que el proyecto no existe
    Y no se guarda ningún registro de esfuerzo

  @CA3
  @R3
  Escenario: No registrar esfuerzo en una historia inexistente en el proyecto
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" no tiene la historia "HU-99"
    Y que el integrante "Ana Pérez" está registrado en el proyecto "p1"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-99" del proyecto "p1" con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "Análisis" y 1.5 horas
    Entonces se devuelve un error indicando que la historia no existe en el proyecto
    Y no se guarda ningún registro de esfuerzo

  @CA3
  @R3
  Escenario: El mismo ID de historia en otro proyecto no cuenta como historia existente
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y un proyecto existente identificado por "p2", con nombre "Otro proyecto"
    Y que el proyecto "p1" no tiene la historia "HU-01"
    Y que el proyecto "p2" tiene la historia "HU-01" con título "Iniciar sesión"
    Y que el integrante "Ana Pérez" está registrado en el proyecto "p1"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-01" del proyecto "p1" con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "Análisis" y 1.5 horas
    Entonces se devuelve un error indicando que la historia no existe en el proyecto
    Y no se guarda ningún registro de esfuerzo

  @CA3
  @R7
  Escenario: Un ID de historia malformado devuelve el error de historia inexistente
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el integrante "Ana Pérez" está registrado en el proyecto "p1"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "historia-uno" del proyecto "p1" con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "Análisis" y 1.5 horas
    Entonces se devuelve un error indicando que la historia no existe en el proyecto
    Y no se guarda ningún registro de esfuerzo

  @CA3
  @R4
  Escenario: No registrar esfuerzo con un integrante que no pertenece al proyecto
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión"
    Y que el integrante "Ana Pérez" está registrado en el proyecto "p1"
    Y que el integrante "Luis Gómez" no está registrado en el proyecto "p1"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-01" del proyecto "p1" con integrante "Luis Gómez", fecha "15 de noviembre de 2026", actividad "Análisis" y 1.5 horas
    Entonces se devuelve un error indicando que el integrante no pertenece al proyecto
    Y no se guarda ningún registro de esfuerzo

  @CA3
  @R4
  Escenario: Un integrante que existe solo en otro proyecto no pertenece al proyecto indicado
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y un proyecto existente identificado por "p2", con nombre "Otro proyecto"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión"
    Y que el integrante "Luis Gómez" está registrado en el proyecto "p2"
    Y que el integrante "Luis Gómez" no está registrado en el proyecto "p1"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-01" del proyecto "p1" con integrante "Luis Gómez", fecha "15 de noviembre de 2026", actividad "Análisis" y 1.5 horas
    Entonces se devuelve un error indicando que el integrante no pertenece al proyecto
    Y no se guarda ningún registro de esfuerzo

  @CA3
  @R10
  Escenario: No registrar esfuerzo en una historia de un Sprint cerrado
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Cerrado"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión" y estado "Completada"
    Y que la historia "HU-01" está asignada al Sprint número 1 del proyecto "p1"
    Y que el integrante "Ana Pérez" está registrado en el proyecto "p1"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-01" del proyecto "p1" con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "Análisis" y 1.5 horas
    Entonces se devuelve un error indicando que el Sprint está cerrado
    Y no se guarda ningún registro de esfuerzo

  @CA3
  @R10
  Escenario: Registrar esfuerzo en una historia completada de un Sprint activo
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión" y estado "Completada"
    Y que la historia "HU-01" está asignada al Sprint número 1 del proyecto "p1"
    Y que el integrante "Ana Pérez" está registrado en el proyecto "p1"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-01" del proyecto "p1" con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "Análisis" y 1.5 horas
    Entonces se registra el esfuerzo con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "Análisis" y 1.5 horas
    Y el esfuerzo queda asociado al proyecto "p1" y a la historia "HU-01"

  @CA3
  @R10
  Escenario: Registrar esfuerzo en una historia sin Sprint asignado
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión" y estado "Pendiente"
    Y que la historia "HU-01" no está asignada a ningún Sprint del proyecto "p1"
    Y que el integrante "Ana Pérez" está registrado en el proyecto "p1"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-01" del proyecto "p1" con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "Análisis" y 1.5 horas
    Entonces se registra el esfuerzo con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "Análisis" y 1.5 horas
    Y el esfuerzo queda asociado al proyecto "p1" y a la historia "HU-01"

  @CA3
  @D18
  Escenario: El error de proyecto prevalece sobre los errores de datos
    Dado que no existe un proyecto identificado por "p999"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-01" del proyecto "p999" con integrante "Ana Pérez", fecha "16 de noviembre de 2026", actividad "Análisis" y 0 horas
    Entonces se devuelve un error indicando que el proyecto no existe
    Y no se devuelve un error de horas inválidas
    Y no se devuelve un error de fecha futura
    Y no se guarda ningún registro de esfuerzo

  @CA3
  @D18
  Escenario: El error de datos prevalece sobre el error de historia
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" no tiene la historia "HU-99"
    Y que el integrante "Ana Pérez" está registrado en el proyecto "p1"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-99" del proyecto "p1" con integrante "Ana Pérez", fecha "15 de noviembre de 2026", actividad "Análisis" y 0 horas
    Entonces se devuelve un error indicando que las horas del esfuerzo son inválidas
    Y no se devuelve un error de historia inexistente
    Y no se guarda ningún registro de esfuerzo

  @CA3
  @D18
  Escenario: El error de Sprint prevalece sobre el error de integrante
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Cerrado"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión" y estado "Completada"
    Y que la historia "HU-01" está asignada al Sprint número 1 del proyecto "p1"
    Y que el integrante "Luis Gómez" no está registrado en el proyecto "p1"
    Y que hoy es "15 de noviembre de 2026 a las 09:00"
    Cuando se registra un esfuerzo en la historia "HU-01" del proyecto "p1" con integrante "Luis Gómez", fecha "15 de noviembre de 2026", actividad "Análisis" y 1.5 horas
    Entonces se devuelve un error indicando que el Sprint está cerrado
    Y no se devuelve un error de integrante que no pertenece al proyecto
    Y no se guarda ningún registro de esfuerzo
