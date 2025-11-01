# 🍽️ Restaurante Concurrente - Productor-Consumidor en Go

[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![EbitenEngine](https://img.shields.io/badge/EbitenEngine-v2-00ADD8?style=flat)](https://ebitengine.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

> **Simulación visual de un restaurante usando el patrón Productor-Consumidor con Arquitectura Hexagonal en Go**

## 📖 Descripción

Este proyecto implementa una simulación interactiva de un restaurante donde:
- **Cocineros** (Productores) preparan platos concurrentemente
- **Meseros** (Consumidores) entregan los platos a los clientes
- **Barra** (Buffer) es un canal limitado que conecta productores y consumidores
- **Clientes** generan demanda que controla la producción

La aplicación demuestra conceptos avanzados de concurrencia en Go usando goroutines, canales, y mecanismos de sincronización, todo visualizado en una interfaz gráfica construida con EbitenEngine.

---

## 🎯 Patrón Implementado: Productor-Consumidor

```
┌─────────────┐         ┌─────────┐         ┌──────────┐
│  Cocinero 1 │────────▶│         │────────▶│ Mesero 1 │
│  Cocinero 2 │────────▶│  Barra  │────────▶│ Mesero 2 │
│     ...     │────────▶│ (Canal) │────────▶│ Mesero 3 │
└─────────────┘         └─────────┘         └──────────┘
  Productores            Buffer (5)          Consumidores
```

### Características del Patrón:
- ✅ **Múltiples productores**: 2+ cocineros trabajando en paralelo
- ✅ **Múltiples consumidores**: 3+ meseros sirviendo concurrentemente
- ✅ **Buffer limitado**: Canal buffered simula capacidad de la barra
- ✅ **Producción controlada**: Solo produce si hay demanda (clientes)
- ✅ **Sincronización automática**: Go channels manejan bloqueos
- ✅ **Control dinámico**: Pausar/reanudar, agregar/quitar clientes

---

## 🏗️ Arquitectura Hexagonal

El proyecto sigue el patrón de **Arquitectura Hexagonal** (Ports & Adapters):

```
├── cmd/
│   ├── app/          # Aplicación principal con UI
│   └── test/         # Test de concurrencia sin UI
├── internal/
│   ├── adapter/      # Adaptadores (implementaciones)
│   │   ├── primary/  # UI (Ebiten)
│   │   └── secondary/# Workers (Cocinero, Mesero)
│   ├── domain/       # Núcleo del negocio
│   │   ├── model/    # Entidades (Plato, Cliente)
│   │   ├── port/     # Interfaces (contratos)
│   │   └── service/  # Lógica de concurrencia
│   └── infraestructure/ # Logger, config
```

**Ventajas:**
- 🎯 Separación clara de responsabilidades
- 🧪 Fácil de testear (dominio independiente)
- 🔄 Fácil de extender y mantener
- 📦 El core no depende de detalles externos

---

## 🚀 Instalación y Ejecución

### Requisitos
- Go 1.20 o superior
- Sistema operativo: Windows, macOS, o Linux

### Clonar e Instalar Dependencias
```bash
git clone <tu-repo>
cd restaurant-concurrency
go mod download
```

### Ejecutar Aplicación con UI
```bash
go run ./cmd/app
```

### Ejecutar Test de Concurrencia (Sin UI)
```bash
go run ./cmd/test
```

### Verificar Race Conditions
```bash
go run -race ./cmd/app
```

---

## 🎮 Controles

| Tecla | Acción |
|-------|--------|
| `ESPACIO` | ⏸️  Pausar/Reanudar producción |
| `+` o `=` | 👥 Agregar un cliente |
| `-` | 👤 Quitar un cliente |
| `ESC` | ❌ Salir de la aplicación |

---

## 🧪 Conceptos de Concurrencia Implementados

### 1. Goroutines
```go
// Múltiples cocineros trabajando en paralelo
for i := 1; i <= numCocineros; i++ {
    go func(id int) {
        producer.Produce(ctx, barra, id)
    }(i)
}
```

### 2. Canales (Channels)
```go
// Canal buffered como buffer del productor-consumidor
barra := make(chan model.Plato, capacidadBarra)

// Envío
select {
case barra <- plato:
    fmt.Println("Plato en la barra")
case <-ctx.Done():
    return
}

// Recepción
plato := <-barra
```

### 3. Sincronización con Mutex
```go
func (s *RestaurantService) AgregarClientes(cantidad int) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.clientesActivos += cantidad
}
```

### 4. Context para Cancelación
```go
ctx, cancel := context.WithCancel(context.Background())

// Propagar cancelación
select {
case <-ctx.Done():
    return // Terminar goroutine
default:
    // Continuar
}
```

### 5. WaitGroup para Cierre Ordenado
```go
func (s *RestaurantService) Close() {
    s.cancel()      // Cancelar todas las goroutines
    s.wg.Wait()     // Esperar a que terminen
    close(s.barra)  // Cerrar canal
}
```

---

## 📊 Funcionamiento del Sistema

### Estados del Sistema

1. **Sin Clientes** → ❌ No produce platos
2. **Con Clientes** → ✅ Cocineros producen, meseros consumen
3. **Pausado** → ⏸️  Producción detenida (buffer se vacía)
4. **Barra Llena** → 🔴 Cocineros esperan (bloqueados)
5. **Barra Vacía** → 🟡 Meseros esperan (bloqueados)

### Flujo de Ejecución

```
1. Usuario agrega clientes (+)
   ↓
2. Cocineros empiezan a cocinar
   ↓
3. Platos se colocan en la barra (canal)
   ↓
4. Meseros toman platos de la barra
   ↓
5. Meseros entregan a clientes
   ↓
6. Proceso continúa mientras haya clientes
```

---

## 📈 Métricas en Tiempo Real

La UI muestra:
- 👥 **Clientes Activos**: Número actual de clientes
- 📊 **Barra**: Ocupación actual / capacidad total
- 📈 **Platos Producidos**: Total acumulado
- ✅ **Platos Servidos**: Total entregado
- ⏸️  **Estado**: Activo / Pausado

---

## 🔧 Configuración

Edita `config.json`:
```json
{
  "restaurant": {
    "capacidad_barra": 5,        // Tamaño del buffer
    "num_cocineros": 2,           // Productores
    "num_meseros": 3,             // Consumidores
    "clientes_inicial": 3,
    "tiempo_coccion_ms": 800,     // Velocidad de producción
    "tiempo_entrega_ms": 600      // Velocidad de consumo
  }
}
```

---

## 🧪 Pruebas

### Test Automático de Concurrencia
El archivo `cmd/test/main.go` ejecuta escenarios automáticos:

```bash
go run ./cmd/test
```

**Escenarios probados:**
1. ✅ Sistema sin clientes (no debe producir)
2. ✅ Con clientes (producción y consumo)
3. ✅ Pausar producción
4. ✅ Reanudar producción
5. ✅ Clientes se van (producción para)

### Resultados Esperados
```
✅ Sin race conditions
✅ Todas las goroutines terminan correctamente
✅ Canal se cierra sin bloqueos
✅ Contadores consistentes
```

---

## 📚 Decisiones de Diseño

### ¿Por qué Arquitectura Hexagonal?
- Permite testear la lógica sin UI
- Facilita cambiar implementaciones (UI, workers)
- Separa la lógica de negocio de detalles técnicos

### ¿Por qué Múltiples Productores/Consumidores?
- Demuestra concurrencia real (no simulación secuencial)
- Más realista (un restaurante tiene varios cocineros/meseros)
- Permite observar race conditions si no se sincronizan bien

### ¿Por qué Context en lugar de canales done?
- Más idiomático en Go moderno
- Permite cancelación en cascada
- Integración con bibliotecas estándar

### ¿Por qué RWMutex?
- Múltiples lectores pueden acceder al estado simultáneamente
- Solo escrituras se bloquean mutuamente
- Mejor performance que Mutex simple

---

## 🎓 Requisitos de la Actividad Cumplidos

| Criterio | Estado |
|----------|--------|
| ✅ Goroutines | Múltiples cocineros y meseros |
| ✅ Canales | Canal buffered (barra) |
| ✅ Sincronización | Mutex, WaitGroup, Context |
| ✅ Patrón de Concurrencia | Productor-Consumidor |
| ✅ Interfaz Gráfica | EbitenEngine con interacción |
| ✅ Sin Race Conditions | Verificado con `go run -race` |
| ✅ Documentación | README + análisis técnico |

---

## 🐛 Solución de Problemas

### La aplicación no inicia
```bash
# Verificar versión de Go
go version

# Reinstalar dependencias
go mod tidy
go mod download
```

### No se ven los sprites
- Verifica que `/internal/adapter/primary/ui/assets/` contenga las imágenes
- Los assets se embeden automáticamente con `//go:embed`

### Race conditions detectadas
```bash
# Ejecutar con race detector
go run -race ./cmd/app
```

---

## 📖 Referencias

- [Goroutines y Canales](https://go.dev/tour/concurrency/1)
- [Patrón Productor-Consumidor](https://en.wikipedia.org/wiki/Producer%E2%80%93consumer_problem)
- [EbitenEngine](https://ebitengine.org/)
- [Arquitectura Hexagonal](https://alistair.cockburn.us/hexagonal-architecture/)

---

## 👨‍💻 Autor

**Eduardo** - Programación Concurrente - Universidad

---

## 📄 Licencia

Este proyecto es de código abierto y está disponible bajo la Licencia MIT.

---

## 🙏 Agradecimientos

- EbitenEngine por la excelente biblioteca gráfica
- Comunidad de Go por las mejores prácticas de concurrencia
- Profesores y compañeros por el feedback
