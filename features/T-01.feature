# language: es

@T-01
Característica: Base técnica del proyecto Go, módulo y CI
  Como equipo de desarrollo
  Quiero un repositorio Go compilable, con capas separadas y CI en GitHub Actions
  Para comenzar el desarrollo con verificación automática y trazabilidad

  @CA1
  Esquema del escenario: El repositorio compila con go build
    Dado el módulo Go "github.com/danaguzman22/software-metrics-estimation"
    Y la directiva "go 1.24" en go.mod
    Y los paquetes "cmd/server", "internal/domain", "internal/service", "internal/storage" e "internal/web"
    Cuando se ejecuta "<comando>"
    Entonces el resultado de la compilación es "<resultado>"

    Ejemplos:
      | comando                                                  | resultado |
      | go build ./...                                           | exitosa   |
      | go build ./cmd/server                                    | exitosa   |
      | go build ./internal/domain                               | exitosa   |
      | go build ./... con un import inexistente en internal/web | fallida   |

  @CA2
  Esquema del escenario: GitHub Actions ejecuta las pruebas en cada push
    Dado el workflow de GitHub Actions con los pasos "go build ./..." y "go test ./... -cover"
    Y el workflow usa la versión de Go declarada en go.mod
    Cuando ocurre "<evento>"
    Entonces GitHub Actions inicia el workflow
    Y se ejecuta el comando "go test ./... -cover"
    Y el resultado del workflow es "<resultado>"

    Ejemplos:
      | evento                                                                   | resultado |
      | push a la rama feature/T-01-estructura-proyecto-go                       | exitoso   |
      | pull_request hacia develop                                               | exitoso   |
      | pull_request hacia main                                                  | exitoso   |
      | push a la rama feature/T-01-estructura-proyecto-go con un test que falla | fallido   |

  @CA3
  Esquema del escenario: La lógica de negocio se mantiene separada de la capa web
    Dado el paquete "internal/domain" con "<condición>"
    Y el paquete "internal/web"
    Cuando se inspeccionan los imports de "internal/domain"
    Entonces la verificación estructural reporta "<resultado>"

    Ejemplos:
      | condición                                 | resultado                     |
      | no importa "internal/web"                 | separación correcta           |
      | importa "fmt" y no importa "internal/web" | separación correcta           |
      | no declara imports                        | separación correcta           |
      | importa "internal/web"                    | violación de separación       |
      | importa "net/http"                        | se incumple la guarda técnica |
