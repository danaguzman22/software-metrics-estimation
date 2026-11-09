# T-02: Persistencia en archivos JSON

- **Versión:** 0.1
- **Estado:** Borrador
- **Autora:** Victoria

## Objetivo
Proveer un almacenamiento genérico en archivos JSON, desacoplado del dominio, para guardar y recuperar datos de la aplicación.

## Entradas
| Campo | Tipo | Obligatorio | Validación |
|---|---|---|---|
| dir | string | Sí | No vacío. Se crea si no existe |
| name | string | Sí | No vacío, sin `/`, `\`, ni ser `.` o `..` |
| v | any | Sí | Serializable a JSON (Save) / puntero (Load) |

## Salidas esperadas
- Save: crea o sobrescribe `<dir>/<name>.json` con JSON indentado.
- Load: completa el valor destino con el contenido del archivo.

## Reglas de negocio
1. Un nombre lógico corresponde a un archivo `<name>.json` dentro del directorio base.
2. La escritura es atómica: se escribe un archivo temporal y luego se renombra, para no dejar archivos a medias.
3. Guardar sobre un nombre existente lo sobrescribe.
4. Si falla el guardado, no queda archivo temporal.

## Restricciones
- Solo biblioteca estándar de Go.
- No depende de paquetes de dominio.

## Casos límite
- Directorio base inexistente (anidado): se crea.
- Sobrescribir un archivo existente.
- Valor no serializable (por ejemplo un canal): error y no se crea el archivo.
- Nombres con separadores de ruta o `..` (path traversal).

## Condiciones de error
| Condición | Error |
|---|---|
| Directorio vacío | error descriptivo |
| Nombre inválido | `ErrInvalidName` |
| Archivo inexistente al cargar | `ErrNotFound` |
| JSON malformado al cargar | `ErrCorrupt` |
| Destino de carga no es puntero | error de json (no `ErrCorrupt`) |

## Criterios de aceptación
- [ ] CA-1: guardar y cargar devuelve el mismo valor.
- [ ] CA-2: cargar un archivo inexistente devuelve `ErrNotFound`.
- [ ] CA-3: cargar JSON malformado devuelve `ErrCorrupt`.
- [ ] CA-4: nombres inválidos devuelven `ErrInvalidName` y no escriben nada.
- [ ] CA-5: guardar crea el directorio base si no existe.
- [ ] CA-6: guardar sobrescribe y no deja archivos temporales.
- [ ] CA-7: un valor no serializable devuelve error y no crea archivo.

## Trazabilidad
- Escenarios BDD: `features/persistencia.feature`
- Tests: `internal/storage/storage_test.go`
- Código: `internal/storage/storage.go`
