# language: es
# T-14 implementada: pendiente. Pantalla de integrantes en internal/web y aviso
# en Planning Poker. Fuente: issue #87 y decisiones DECIDIDAS 1–5 y D9 de specs/T-14.md.
# Los escenarios describen CA1–CA6 y las decisiones confirmadas: una página por
# proyecto con listado y formulario; enlace "Ver integrantes" en el detalle sin
# quitar los existentes; alta con POST-Redirect-GET 303; errores de dominio 422
# con mensaje exacto y valores conservados; proyecto inexistente 404; método no
# permitido 405; error de infraestructura 500 genérico; aviso en Poker cuando no
# hay integrantes; sin reglas duplicadas y con salida escapada. Las reglas de
# negocio (nombre, rol, duplicado) siguen en el dominio de US-03.
# Los detalles de rutas, plantilla, textos y tests son DERIVADAS 6–8 y 10–13 de
# la spec; D9 (nombre del proyecto con ObtenerProyecto) está confirmada.
# No se escenifican cosas fuera de alcance: modificar o eliminar integrantes,
# catálogo de roles, autenticación, paginación, JavaScript, cambios en el
# servicio, repositorios, migraciones o reglas de Poker.

@T-14
Característica: Pantalla para registrar integrantes del proyecto
  Como integrante del equipo
  Quiero registrar y ver los integrantes de un proyecto desde la interfaz web
  Para no cargarlos por SQL y poder usar Planning Poker

  @CA1
  Escenario: La página del proyecto muestra su nombre, sus integrantes y el formulario de alta
    Dado un proyecto existente identificado por su ID y con nombre "Sistema de Métricas"
    Y que el proyecto tiene registrados a "Dana Guzmán" con rol "Desarrollador" y a "Victoria Caro" con rol "Tester"
    Cuando se abre la página de integrantes del proyecto
    Entonces la página responde con éxito
    Y muestra el nombre del proyecto "Sistema de Métricas"
    Y muestra a "Dana Guzmán" con rol "Desarrollador"
    Y muestra a "Victoria Caro" con rol "Tester"
    Y muestra un formulario con los campos "Nombre" y "Rol"
    Y muestra el botón "Registrar integrante"

  @CA1
  Escenario: Un proyecto sin integrantes muestra el estado vacío y el formulario
    Dado un proyecto existente identificado por su ID
    Y que el proyecto no tiene integrantes registrados
    Cuando se abre la página de integrantes del proyecto
    Entonces la página responde con éxito
    Y muestra el mensaje "El proyecto no tiene integrantes"
    Y muestra el formulario de alta de integrantes

  @CA1
  Escenario: El detalle del proyecto ofrece el enlace a integrantes y conserva los enlaces existentes
    Dado un proyecto existente identificado por su ID
    Cuando se abre el detalle del proyecto
    Entonces la página muestra el enlace "Ver integrantes" hacia la página de integrantes
    Y conserva los enlaces "Ver backlog", "Ver sprints", "Ver defectos", "Ver dashboard" y "Ver reporte"

  @CA2
  Escenario: Alta válida de un integrante redirige al listado
    Dado un proyecto existente identificado por su ID
    Cuando se registra en ese proyecto el integrante "Dana Guzmán" con rol "Desarrollador"
    Entonces la respuesta redirige con el código 303 a la página de integrantes del proyecto
    Y al abrir el listado, la integrante "Dana Guzmán" aparece con rol "Desarrollador"

  @CA2
  Escenario: El primer integrante se registra correctamente
    Dado un proyecto existente identificado por su ID y sin integrantes
    Cuando se registra en ese proyecto el integrante "Dana Guzmán" con rol "Desarrollador"
    Entonces el alta se realiza exitosamente
    Y el proyecto queda con la integrante "Dana Guzmán"

  @CA2
  Escenario: El listado conserva el orden devuelto por el servicio
    Dado un proyecto existente identificado por su ID
    Y que el servicio devuelve a "Dana Guzmán" y luego a "Victoria Caro"
    Cuando se abre la página de integrantes del proyecto
    Entonces la página muestra primero a "Dana Guzmán" y después a "Victoria Caro"

  @CA3
  Escenario: Nombre vacío devuelve error y conserva el rol ingresado
    Dado un proyecto existente identificado por su ID
    Cuando se intenta registrar un integrante con nombre "" y rol "Desarrollador"
    Entonces la página responde con el código 422
    Y muestra el mensaje "el nombre del integrante no puede estar vacío"
    Y el formulario conserva el rol "Desarrollador"
    Y el integrante no queda registrado

  @CA3
  Escenario: Rol vacío devuelve error y conserva el nombre ingresado
    Dado un proyecto existente identificado por su ID
    Cuando se intenta registrar el integrante "Dana Guzmán" con rol ""
    Entonces la página responde con el código 422
    Y muestra el mensaje "el rol del integrante no puede estar vacío"
    Y el formulario conserva el nombre "Dana Guzmán"
    Y el integrante no queda registrado

  @CA3
  Escenario: Integrante duplicado devuelve error y conserva los valores
    Dado un proyecto existente identificado por su ID
    Y que el proyecto tiene registrada a la integrante "Dana Guzmán" con rol "Desarrollador"
    Cuando se intenta registrar nuevamente a "Dana Guzmán" con rol "Desarrollador"
    Entonces la página responde con el código 422
    Y muestra el mensaje "el integrante ya está registrado en el proyecto"
    Y el formulario conserva el nombre "Dana Guzmán" y el rol "Desarrollador"
    Y el proyecto sigue teniendo únicamente a la integrante "Dana Guzmán"

  @CA3
  Escenario: El duplicado se reconoce ignorando mayúsculas y espacios extremos
    Dado un proyecto existente identificado por su ID
    Y que el proyecto tiene registrada a la integrante "Dana Guzmán" con rol "Desarrollador"
    Cuando se intenta registrar a "  dana guzmán  " con rol "Tester"
    Entonces la página responde con el código 422
    Y muestra el mensaje "el integrante ya está registrado en el proyecto"
    Y el integrante no queda registrado

  @CA4
  Escenario: Abrir la página de un proyecto inexistente devuelve 404
    Dado que no existe un proyecto identificado por el ID solicitado
    Cuando se abre la página de integrantes de ese proyecto
    Entonces la página responde con el código 404
    Y muestra el mensaje de proyecto no encontrado

  @CA4
  Escenario: Registrar un integrante en un proyecto inexistente devuelve 404
    Dado que no existe un proyecto identificado por el ID solicitado
    Cuando se intenta registrar en ese proyecto el integrante "Dana Guzmán" con rol "Desarrollador"
    Entonces la página responde con el código 404
    Y muestra el mensaje de proyecto no encontrado
    Y el integrante no queda registrado

  @CA4
  Escenario: Un método no permitido devuelve 405
    Dado un proyecto existente identificado por su ID
    Cuando se envía una petición con un método distinto de GET y POST a la ruta de integrantes
    Entonces la respuesta tiene el código 405
    Y el handler de integrantes no se ejecuta

  @CA4
  Escenario: Un error de infraestructura devuelve 500 con mensaje genérico
    Dado un proyecto existente identificado por su ID
    Y que el servicio de integrantes falla por infraestructura
    Cuando se abre la página de integrantes del proyecto
    Entonces la página responde con el código 500
    Y muestra el mensaje "Ocurrió un error inesperado. Intentá nuevamente."
    Y no muestra el detalle del error de infraestructura

  @CA5
  Escenario: Planning Poker sin integrantes muestra un aviso con enlace
    Dado un proyecto existente identificado por su ID y sin integrantes
    Y una historia del backlog con una ronda de Planning Poker
    Cuando se abre la página de Planning Poker de la historia
    Entonces la página muestra el aviso "Este proyecto todavía no tiene integrantes"
    Y muestra el enlace "Registrar integrantes" hacia la página de integrantes del proyecto

  @CA5
  Escenario: Planning Poker sin integrantes no muestra desplegables vacíos
    Dado un proyecto existente identificado por su ID y sin integrantes
    Y una historia del backlog con una ronda abierta de Planning Poker
    Cuando se abre la página de Planning Poker de la historia
    Entonces la página no muestra el formulario de votar
    Y la página no muestra el formulario de revelar
    Y la página no muestra el formulario de acordar estimación

  @CA5
  Escenario: Planning Poker con integrantes conserva sus formularios
    Dado un proyecto existente identificado por su ID
    Y que el proyecto tiene registrada a la integrante "Dana Guzmán"
    Y una historia del backlog con una ronda abierta de Planning Poker
    Cuando se abre la página de Planning Poker de la historia
    Entonces la página muestra el formulario de votar con la integrante "Dana Guzmán"
    Y la página muestra el formulario de revelar
    Y no muestra el aviso de proyecto sin integrantes

  @CA6
  Escenario: El nombre con markup se muestra escapado
    Dado un proyecto existente identificado por su ID
    Y que el proyecto tiene registrado un integrante cuyo nombre contiene "<script>alert(1)</script>"
    Cuando se abre la página de integrantes del proyecto
    Entonces la página no interpreta el nombre como HTML
    Y muestra el texto del nombre escapado

  @CA6
  Escenario: La composición del servidor registra las rutas de integrantes
    Dado el servidor armado con la composición real y un pool nil
    Cuando se envía un método no permitido a la ruta de integrantes de un proyecto
    Entonces la respuesta tiene el código 405
    Y la ruta inexistente de integrantes responde 404
    Y la ruta de salud sigue respondiendo con éxito
