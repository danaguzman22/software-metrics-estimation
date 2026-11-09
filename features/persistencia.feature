# language: es
Característica: Persistencia en archivos JSON (T-02)
  Como Product Builder
  Quiero guardar y recuperar datos en archivos JSON
  Para no perder la información al cerrar la aplicación

  Escenario: Guardar y cargar un dato
    Dado un almacenamiento en un directorio vacío
    Cuando guardo el dato "proyectos" con id 1 y nombre "Demo"
    Y cargo el dato "proyectos"
    Entonces obtengo id 1 y nombre "Demo"

  Escenario: Sobrescribir un dato existente
    Dado un almacenamiento con el dato "proyectos" guardado con nombre "Viejo"
    Cuando guardo el dato "proyectos" con id 1 y nombre "Nuevo"
    Y cargo el dato "proyectos"
    Entonces obtengo id 1 y nombre "Nuevo"

  Escenario: Cargar un dato que no existe
    Dado un almacenamiento en un directorio vacío
    Cuando cargo el dato "inexistente"
    Entonces obtengo el error "archivo no encontrado"

  Escenario: Cargar un archivo con JSON malformado
    Dado un almacenamiento con el archivo "roto" con contenido inválido
    Cuando cargo el dato "roto"
    Entonces obtengo el error "contenido JSON inválido"

  Esquema del escenario: Rechazar nombres inválidos
    Dado un almacenamiento en un directorio vacío
    Cuando guardo el dato "<nombre>" con id 1 y nombre "X"
    Entonces obtengo el error "nombre inválido"

    Ejemplos:
      | nombre |
      | ../x   |
      | a/b    |
      | ..     |
