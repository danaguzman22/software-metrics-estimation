# language: es

@T-02
Característica: Persistencia en archivos JSON
  Como Product Builder
  Quiero guardar y recuperar datos en archivos JSON
  Para no perder la información al cerrar la aplicación

  @CA1
  Escenario: Guardar y cargar un dato devuelve el mismo valor
    Dado un almacenamiento en un directorio vacío
    Cuando guardo el dato "proyectos" con id 1 y nombre "Demo"
    Y cargo el dato "proyectos"
    Entonces obtengo id 1 y nombre "Demo"

  @CA2
  Escenario: Cargar un dato que no existe devuelve error de archivo no encontrado
    Dado un almacenamiento en un directorio vacío
    Cuando cargo el dato "inexistente"
    Entonces obtengo el error "archivo no encontrado"

  @CA3
  Escenario: Cargar un archivo con JSON malformado devuelve error de contenido inválido
    Dado un almacenamiento con el archivo "roto" con contenido JSON inválido
    Cuando cargo el dato "roto"
    Entonces obtengo el error "contenido JSON inválido"

  @CA4
  Esquema del escenario: Rechazar nombres inválidos al guardar y al cargar
    Dado un almacenamiento en un directorio vacío
    Cuando guardo el dato "<nombre>" con id 1 y nombre "X"
    Y cargo el dato "<nombre>"
    Entonces obtengo el error "nombre inválido" en ambas operaciones
    Y no se crea ningún archivo en el directorio

    Ejemplos:
      | nombre |
      |        |
      | .      |
      | ..     |
      | ../x   |
      | a/b    |
      | a\\b   |

  @CA5
  Escenario: Crear el almacenamiento crea el directorio base si no existe
    Dado un directorio base inexistente en la ruta anidada "a/b/data"
    Cuando creo un almacenamiento en ese directorio
    Entonces el directorio base "a/b/data" existe

  @CA6
  Escenario: Crear el almacenamiento con un directorio vacío devuelve error
    Dado un directorio base compuesto únicamente por espacios "  "
    Cuando creo un almacenamiento en ese directorio
    Entonces obtengo un error descriptivo

  @CA7
  Escenario: Guardar sobre un dato existente lo sobrescribe sin dejar archivos temporales
    Dado un almacenamiento con el dato "p" guardado con id 1 y nombre "Viejo"
    Cuando guardo el dato "p" con id 1 y nombre "Nuevo"
    Y cargo el dato "p"
    Entonces obtengo id 1 y nombre "Nuevo"
    Y en el directorio solo queda el archivo "p.json"

  @CA8
  Escenario: Guardar un valor no serializable devuelve error y no crea archivo
    Dado un almacenamiento en un directorio vacío
    Cuando guardo el dato "canal" con un valor no serializable
    Entonces obtengo un error
    Y no se crea ningún archivo en el directorio

  @CA9
  Escenario: Cargar con un destino que no es puntero devuelve un error distinto de contenido inválido
    Dado un almacenamiento con el dato "x" guardado con id 1
    Cuando cargo el dato "x" en un destino que no es puntero
    Entonces obtengo un error
    Y el error no es "contenido JSON inválido"
