# Software Metrics & Estimation

## Trabajo Práctico Integrador 2026

**Universidad Tecnológica Nacional**  
**Facultad Regional San Rafael**  
**Ingeniería en Sistemas de Información**  
**Ingeniería y Calidad de Software**

---

## Descripción del proyecto

**Software Metrics & Estimation** es una aplicación destinada a la gestión de proyectos de software, permitiendo realizar tareas relacionadas con su estimación, planificación, seguimiento y medición de calidad.

El proyecto corresponde al Trabajo Práctico Integrador de la asignatura Ingeniería y Calidad de Software 2026.

La aplicación deberá permitir administrar proyectos y obtener información relacionada con el avance, esfuerzo, estimaciones, defectos y métricas del proyecto.

---

## Objetivo

Desarrollar una aplicación que permita realizar la estimación, seguimiento y medición de proyectos de software, aplicando los conceptos y prácticas adquiridos durante la asignatura.

---

## Funcionalidades principales

La aplicación deberá contemplar, como mínimo, las siguientes funcionalidades:

### 1. Gestión de proyectos

- Crear proyectos.
- Modificar proyectos.
- Registrar integrantes.
- Registrar fecha de inicio y finalización.
- Consultar el estado de un proyecto.

### 2. Product Backlog

Para cada proyecto se deberá administrar un Product Backlog.

Cada elemento deberá contener:

- Identificador.
- Título.
- Descripción.
- Prioridad.
- Estado.
- Story Points.
- Criterios de aceptación.

### 3. Gestión de Sprints

La aplicación deberá permitir:

- Crear Sprints.
- Definir un Sprint Goal.
- Asignar historias a un Sprint.
- Registrar historias completadas.
- Cerrar un Sprint.
- Consultar Sprints anteriores.

### 4. Estimación

La aplicación deberá permitir estimar historias mediante Story Points.

También deberá implementar un mecanismo de **Planning Poker** que permita:

- Registrar estimaciones individuales.
- Mantener las estimaciones ocultas hasta finalizar la votación.
- Mostrar las estimaciones realizadas.
- Detectar diferencias entre estimaciones.
- Realizar nuevas rondas.
- Registrar la estimación acordada.

### 5. Registro de esfuerzo

Para cada historia o tarea se deberá poder registrar:

- Integrante.
- Fecha.
- Actividad realizada.
- Horas trabajadas.

Esto permitirá comparar posteriormente el esfuerzo estimado con el esfuerzo real.

### 6. Gestión de defectos

El sistema deberá permitir registrar defectos indicando:

- Descripción.
- Severidad.
- Estado.
- Historia relacionada.
- Sprint de detección.
- Sprint de resolución.

### 7. Métricas

La aplicación deberá calcular como mínimo:

- Story Points planificados.
- Story Points completados.
- Velocidad del equipo.
- Horas estimadas.
- Horas reales.
- Desviación entre esfuerzo estimado y real.
- Porcentaje de historias completadas.
- Cantidad de defectos detectados.
- Cantidad de defectos resueltos.

### 8. Dashboard

El sistema deberá presentar un Dashboard con información sobre el estado del proyecto y representaciones gráficas de algunas de las métricas obtenidas.

### 9. Reportes

La aplicación deberá generar un reporte de un proyecto o Sprint incluyendo:

- Historias planificadas y completadas.
- Estimaciones.
- Esfuerzo registrado.
- Métricas.
- Defectos.

---

## 🛠️ Tecnologías y metodologías

El proyecto deberá desarrollarse utilizando:

- **Go (Golang)** como lenguaje de programación.
- **Git** para el control de versiones.
- **GitHub** para el repositorio y la gestión del proyecto.
- **Scrum** para la organización y gestión del proyecto.
- **SDD (Specification-Driven Development)** para la especificación de funcionalidades.
- **BDD (Behavior-Driven Development)** para la definición y validación de comportamientos.
- **TDD (Test-Driven Development)** para el desarrollo y prueba de componentes.
- **Herramientas de Inteligencia Artificial** como soporte al proceso de desarrollo.

Todo resultado generado mediante Inteligencia Artificial deberá ser comprendido, revisado y validado por el equipo.

---

## Equipo

| Integrante | Rol |
|---|---|
| **Dana Guzmán** | Agile Enabler / Product Builder |
| **Victoria Caro Boldrini** | Product Builder |

### Product Architects

Los profesores de la asignatura cumplen el rol de **Product Architects**.

---

## Metodología de trabajo

El proyecto será gestionado mediante **Scrum**, utilizando un Product Backlog y diferentes Sprints.

El trabajo se organizará de acuerdo con las siguientes etapas:

### Sprint 0 — La Preparación

Objetivo:

- Formar el equipo.
- Asignar roles.
- Configurar el entorno de trabajo.
- Crear el repositorio.
- Configurar GitHub Projects.
- Elaborar el Product Backlog inicial junto al Cliente.

### Sprint 1 — El MVP

Objetivo:

- Entregar una versión funcional básica del sistema.

Se realizarán las ceremonias correspondientes:

- Sprint Planning.
- Daily.
- Sprint Review.
- Sprint Retrospective.

### Sprint 2 — La Interfaz

Objetivo:

- Dotar al software de una interfaz usable.

### Sprint 3 — Funcionalidad y Calidad

Objetivo:

- Añadir funcionalidades clave.
- Robustecer el sistema.
- Implementar la visualización.

### Sprint 4 — El Cierre y la Entrega Final

Objetivo:

- Pulir el producto.
- Completar la documentación.
- Preparar la entrega final.
- Permitir la exportación de un informe del proyecto en formato PDF.

---

## Estado actual

**Sprint 0 — La Preparación**

Actualmente el proyecto se encuentra en la etapa inicial de configuración.

### Actividades iniciales

- [x] Formación del equipo.
- [x] Asignación de roles.
- [ ] Creación y configuración del repositorio.
- [ ] Configuración de GitHub Projects.
- [ ] Elaboración del Product Backlog inicial.
- [ ] Reunión inicial con los Product Architects.

---

## Trazabilidad

Durante el desarrollo se deberá mantener trazabilidad entre:

**Historia de Usuario → Especificación SDD → Criterios de Aceptación → Escenarios BDD → Tests → Código Go**

Esta trazabilidad permitirá demostrar el proceso completo de desarrollo de las funcionalidades seleccionadas.

---

## Entregables

El proyecto deberá contemplar:

1. Repositorio Git con historial de contribuciones.
2. Tablero Scrum.
3. Product Backlog y Sprint Backlogs.
4. Especificaciones SDD.
5. Escenarios BDD.
6. Pruebas automatizadas.
7. Código fuente en Go.
8. Evidencias de aplicación de TDD.
9. Software funcional.
10. Informe de métricas y cobertura de pruebas.
11. Actas de retrospectivas.
12. Documentación técnica y manual breve de usuario.
13. Presentación y demostración final.

---

## 📍 Institución

**Universidad Tecnológica Nacional — Facultad Regional San Rafael**

**Carrera:** Ingeniería en Sistemas de Información  
**Asignatura:** Ingeniería y Calidad de Software  
**Año:** 2026
