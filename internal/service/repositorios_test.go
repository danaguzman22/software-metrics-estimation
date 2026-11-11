package service

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/danaguzman22/software-metrics-estimation/internal/domain"
)

func verificarFirmaT06(t *testing.T, contrato reflect.Type, nombre string, firma reflect.Type) {
	t.Helper()
	metodo, ok := contrato.MethodByName(nombre)
	if !ok {
		t.Errorf("%s no declara el método %q", contrato.Name(), nombre)
		return
	}
	if metodo.Type != firma {
		t.Errorf("%s.%s firma = %v, se esperaba %v", contrato.Name(), nombre, metodo.Type, firma)
	}
}

func TestT06_RepositorioProyectos_FirmasExactas(t *testing.T) {
	contrato := reflect.TypeOf((*RepositorioProyectos)(nil)).Elem()
	registro := reflect.TypeOf(ProyectoRegistrado{})
	proyecto := reflect.TypeOf(domain.Proyecto{})
	errType := reflect.TypeOf((*error)(nil)).Elem()
	stringType := reflect.TypeOf("")

	esperados := map[string]reflect.Type{
		"Obtener":    reflect.FuncOf([]reflect.Type{stringType}, []reflect.Type{registro, errType}, false),
		"Existe":     reflect.FuncOf([]reflect.Type{stringType}, []reflect.Type{reflect.TypeOf(false), errType}, false),
		"Listar":     reflect.FuncOf(nil, []reflect.Type{reflect.SliceOf(registro), errType}, false),
		"Actualizar": reflect.FuncOf([]reflect.Type{stringType, proyecto}, []reflect.Type{errType}, false),
	}

	if contrato.NumMethod() != len(esperados) {
		t.Errorf("RepositorioProyectos tiene %d métodos, se esperaban %d", contrato.NumMethod(), len(esperados))
	}
	for nombre, firma := range esperados {
		verificarFirmaT06(t, contrato, nombre, firma)
	}
}

func TestT06_RepositorioIntegrantes_FirmasExactas(t *testing.T) {
	contrato := reflect.TypeOf((*RepositorioIntegrantes)(nil)).Elem()
	integrante := reflect.TypeOf(domain.Integrante{})
	errType := reflect.TypeOf((*error)(nil)).Elem()
	stringType := reflect.TypeOf("")

	esperados := map[string]reflect.Type{
		"ListarPorProyecto": reflect.FuncOf([]reflect.Type{stringType}, []reflect.Type{reflect.SliceOf(integrante), errType}, false),
		"Agregar":           reflect.FuncOf([]reflect.Type{stringType, integrante}, []reflect.Type{errType}, false),
	}

	if contrato.NumMethod() != len(esperados) {
		t.Errorf("RepositorioIntegrantes tiene %d métodos, se esperaban %d", contrato.NumMethod(), len(esperados))
	}
	for nombre, firma := range esperados {
		verificarFirmaT06(t, contrato, nombre, firma)
	}
}

func TestT06_RepositorioBacklog_FirmasExactas(t *testing.T) {
	contrato := reflect.TypeOf((*RepositorioBacklog)(nil)).Elem()
	item := reflect.TypeOf(domain.ItemBacklog{})
	errType := reflect.TypeOf((*error)(nil)).Elem()
	stringType := reflect.TypeOf("")

	esperados := map[string]reflect.Type{
		"AgregarItem":       reflect.FuncOf([]reflect.Type{stringType, item}, []reflect.Type{item, errType}, false),
		"ObtenerItem":       reflect.FuncOf([]reflect.Type{stringType, stringType}, []reflect.Type{item, errType}, false),
		"ListarPorProyecto": reflect.FuncOf([]reflect.Type{stringType}, []reflect.Type{reflect.SliceOf(item), errType}, false),
		"ActualizarItem":    reflect.FuncOf([]reflect.Type{stringType, item}, []reflect.Type{errType}, false),
	}

	if contrato.NumMethod() != len(esperados) {
		t.Errorf("RepositorioBacklog tiene %d métodos, se esperaban %d", contrato.NumMethod(), len(esperados))
	}
	for nombre, firma := range esperados {
		verificarFirmaT06(t, contrato, nombre, firma)
	}
}

func TestT06_ProyectoRegistrado_CamposYTagsJSON(t *testing.T) {
	registro := ProyectoRegistrado{
		ID:       "p1",
		Proyecto: domain.Proyecto{Nombre: "Demo"},
	}

	datos, err := json.Marshal(registro)
	if err != nil {
		t.Fatalf("json.Marshal(ProyectoRegistrado): %v", err)
	}

	var crudo map[string]json.RawMessage
	if err := json.Unmarshal(datos, &crudo); err != nil {
		t.Fatalf("json.Unmarshal a mapa: %v", err)
	}
	if _, ok := crudo["id"]; !ok {
		t.Errorf("el JSON de ProyectoRegistrado no tiene la clave \"id\": %s", datos)
	}
	if _, ok := crudo["proyecto"]; !ok {
		t.Errorf("el JSON de ProyectoRegistrado no tiene la clave \"proyecto\": %s", datos)
	}

	entrada := []byte(`{"id":"p9","proyecto":{"Nombre":"Round"}}`)
	var vuelta ProyectoRegistrado
	if err := json.Unmarshal(entrada, &vuelta); err != nil {
		t.Fatalf("json.Unmarshal(ProyectoRegistrado): %v", err)
	}
	if vuelta.ID != "p9" {
		t.Errorf("round-trip ID = %q, se esperaba %q", vuelta.ID, "p9")
	}
	if vuelta.Proyecto.Nombre != "Round" {
		t.Errorf("round-trip Proyecto.Nombre = %q, se esperaba %q", vuelta.Proyecto.Nombre, "Round")
	}
}
