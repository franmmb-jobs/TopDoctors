# TopDoctors Diagnostics API - Code Challenge

Este repositorio contiene la resolución de la prueba técnica para el puesto de Backend Developer en TopDoctors. La solución consiste en una API REST para gestionar diagnósticos de pacientes, incluyendo autenticación mediante JWT y documentación automática con Swagger.

## Índice
1. [Resumen de Resolución](#resumen-de-resolución)
2. [Stack Tecnológico](#stack-tecnológico)
3. [Arquitectura y Decisiones](#arquitectura-y-decisiones)
4. [Configuración y Secretos](#configuración-y-secretos)
5. [API Documentation (Swagger)](#api-documentation-swagger)
6. [Cómo ejecutar el proyecto](#cómo-ejecutar-el-proyecto)

---

## Resumen de Resolución

Se ha completado la funcionalidad principal solicitada en el archivo `ToDo.md`:
- **Autenticación**: Endpoint para generación de tokens JWT.
- **Gestión de Diagnósticos**: Endpoints protegidos para consultar y almacenar diagnósticos.
- **Filtrado**: Capacidad de filtrar diagnósticos por nombre del paciente y/o fecha.
- **Validaciones Extra**: Implementación de verificaciones robustas para DNI y Email.

### Calidad y Pruebas
Se han implementado **tests unitarios y de integración** para los módulos más críticos del sistema.
> [!NOTE]
> Debido al límite de tiempo, no se han realizado pruebas de carga ni de race condition, y me hubiera gustado extender la cobertura de tests y validaciones aún más.

Se han dedicado un total de **8 horas** a la resolución, priorizando la estructura de diagnóstico y pacientes.

## Stack Tecnológico

- **Lenguaje**: Go (Golang)
- **Base de Datos**: SQLite (GORM)
- **Logging**: `slog` (estándar) + `tint` (formato legible)
- **Configuración**: `cleanenv` (soporta `.yml` y `.env`)
- **Documentación**: `swag` (Swagger)
- **Contenerización**: Docker & Docker Compose

## Arquitectura y Decisiones

### Justificación de Herramientas
- **Librería HTTP estándar**: Se ha decidido usar la librería `net/http` nativa para demostrar conocimiento del lenguaje sin dependencias pesadas de frameworks externos.
- **Hexagonal Architecture**: El proyecto sigue principios de arquitectura limpia para separar la lógica de negocio de los adaptadores externos (HTTP, Persistencia).
- **Auto-migración**: El programa ejecuta `AutoMigrate` al inicio para asegurar la consistencia del esquema. *Nota: En un entorno real se optimizaría este proceso para evitar sobrecarga innecesaria en cada arranque.*

### Limitaciones Conocidas
- No se han implementado operaciones de `DELETE` o `UPDATE` por foco en la funcionalidad core.
- No se ha implementado capa de caché (considerado no crítico para esta prueba).
- Las respuestas de error podrían ser más granulares (ej. unicidad de usuarios).

---

## Configuración y Secretos

El proyecto utiliza variables de entorno y archivos YAML. A continuación se muestra un ejemplo de la estructura de configuración utilizada:

```yaml
logs:
  level: "debug"

database:
  DSN: "diagnostics.db"

api:
  port: "8050"
  jwt_secret: "docker_secret_key"
```

| Variable | Descripción | Valor por Defecto |
| :--- | :--- | :--- |
| `PORT` | Puerto del servidor HTTP | `8050` |
| `JWT_SECRET` | Clave secreta para tokens JWT | `secret` |

---

## API Documentation (Swagger)

La API cuenta con documentación automática accesible vía Swagger UI.

1. **Generar documentación**: `swag init -g cmd/api/main.go`
2. **Acceso UI**: Una vez iniciado el servidor, visita [http://localhost:8050/swagger/index.html](http://localhost:8050/swagger/index.html)

---

## Cómo ejecutar el proyecto

### Ejecución Nativa
1. **Instalar dependencias**: `go mod tidy`
2. **Ejecutar tests**: `go test ./...`
3. **Arrancar servidor**:
   ```bash
   go run ./cmd/api/main.go -config='configs/config.dev.yml'
   ```

### Ejecución con Docker
1. **Construir imagen**:
   ```bash
   docker build -t topdoctors-api .
   ```
2. **Ejecutar contenedor**:
   ```bash
   docker run -d -p 8050:8050 --name diagnostics-api topdoctors-api
   ```

---
*Prueba técnica realizada por Fran.*
