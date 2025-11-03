# ANÁLISIS Y CORRECCIÓN DE MEJORES PRÁCTICAS

## 📊 PROBLEMAS ENCONTRADOS Y CORREGIDOS

### ❌ ANTES: Funciones Anónimas Innecesarias

#### Problema 1: `restaurant.go` - Función anónima para goroutine
```go
// ❌ MAL - Función anónima innecesaria
go func(c *worker.Cocinero) {
    defer s.wg.Done()
    c.Producir(s.ctx, s.barra, s.hayDemanda)
}(cocinero)
```

**Por qué está mal:**
- Funciones anónimas solo deben usarse cuando realmente no se puede extraer a un método
- Dificulta el testing y debugging
- Hace el código menos legible
- Agrega complejidad innecesaria

**✅ CORRECCIÓN:**
```go
// Método principal
go s.ejecutarCocinero(cocinero)

// Método helper (añadido)
func (s *RestaurantService) ejecutarCocinero(cocinero *worker.Cocinero) {
    defer s.wg.Done()
    cocinero.Producir(s.ctx, s.barra, s.hayDemanda)
}
```

**Beneficios:**
- ✅ Más fácil de testear
- ✅ Más legible
- ✅ Reutilizable
- ✅ Mejor para debugging (tiene nombre propio)

---

#### Problema 2: `restaurant.go` - Función anónima para limpieza de mesa
```go
// ❌ MAL - Función anónima + time.Sleep
go func(m *model.Mesa) {
    time.Sleep(3 * time.Second)  // ← También mal
    s.mesasMu.Lock()
    m.ClientesSatisfechos()
    s.mesasMu.Unlock()
}(mesa)
```

**Por qué está mal:**
1. Función anónima innecesaria
2. `time.Sleep` no respeta cancelación de contexto
3. Si cancelamos el servicio, esta goroutine seguirá ejecutándose 3 segundos

**✅ CORRECCIÓN:**
```go
// Método principal
go s.limpiarMesaDespuesDeTiempo(mesa, 3*time.Second)

// Método helper (añadido)
func (s *RestaurantService) limpiarMesaDespuesDeTiempo(mesa *model.Mesa, duracion time.Duration) {
    select {
    case <-time.After(duracion):
        s.mesasMu.Lock()
        mesa.ClientesSatisfechos()
        s.mesasMu.Unlock()
    case <-s.ctx.Done():
        // Si se cancela el contexto, salir inmediatamente
        return
    }
}
```

**Beneficios:**
- ✅ Respeta cancelación de contexto
- ✅ Se puede testear independientemente
- ✅ No deja goroutines zombie
- ✅ Más legible y mantenible

---

#### Problema 3: `ebiten_game.go` - Funciones anónimas para callbacks
```go
// ❌ MAL - Funciones anónimas triviales
g.inputHandler.SetCallbacks(
    func() { g.service.TogglePausar() },  // Wrapper innecesario
    nil,
    nil,
    func() { /* Cerrar */ },  // Vacía
    nil,
)
```

**Por qué está mal:**
- Wrapping innecesario de un método que ya existe
- La función de cerrar está vacía (no hace nada)
- Agrega una capa de indirección sin valor

**✅ CORRECCIÓN:**
```go
// Pasar métodos directamente
g.inputHandler.SetCallbacks(
    g.service.TogglePausar,  // Método directo
    nil,
    nil,
    g.handleClose,  // Método helper
    nil,
)

// Método helper (añadido)
func (g *Game) handleClose() {
    // Lógica de cierre si es necesaria en el futuro
}
```

**Beneficios:**
- ✅ Más directo y claro
- ✅ Extensible (podemos agregar lógica a handleClose)
- ✅ Mejor para debugging

---

### ❌ ANTES: Uso Incorrecto de `time.Sleep`

#### Problema 1: `cocinero.go` - Sleep para espera sin demanda
```go
// ❌ MAL - time.Sleep no respeta cancelación
if !verificarDemanda() {
    time.Sleep(500 * time.Millisecond)
    continue
}
```

**Por qué está mal:**
- Si cancelamos el contexto, el cocinero seguirá dormido 500ms
- No es cancelable
- Hace que el cierre del programa sea lento

**✅ CORRECCIÓN:**
```go
if !verificarDemanda() {
    // Espera cancelable usando select con time.After
    select {
    case <-time.After(500 * time.Millisecond):
        continue
    case <-ctx.Done():
        return  // Cancelación inmediata
    }
}
```

**Beneficios:**
- ✅ Cancelación inmediata
- ✅ Cierre limpio y rápido
- ✅ Mejor práctica de Go

---

#### Problema 2: `cocinero.go` - Sleep para simular cocción
```go
// ❌ MAL - time.Sleep no respeta cancelación
tiempoCoccion := time.Duration(1500+rand.Intn(1000)) * time.Millisecond
time.Sleep(tiempoCoccion)
```

**Por qué está mal:**
- Si cancelamos mientras cocina, esperará hasta 2.5 segundos
- No es cancelable durante la cocción

**✅ CORRECCIÓN:**
```go
tiempoCoccion := time.Duration(1500+rand.Intn(1000)) * time.Millisecond

select {
case <-time.After(tiempoCoccion):
    // Continuar con la producción
case <-ctx.Done():
    return  // Cancelación durante cocción
}
```

**Beneficios:**
- ✅ Cancelable en cualquier momento
- ✅ No bloquea el cierre
- ✅ Respeta el contexto

---

#### Problema 3: `restaurant.go` - Sleep en goroutine de limpieza
```go
// ❌ MAL - Goroutine anónima + time.Sleep
go func(m *model.Mesa) {
    time.Sleep(3 * time.Second)
    // ... limpiar mesa
}(mesa)
```

**Ya corregido arriba** - Ahora usa `limpiarMesaDespuesDeTiempo` con `time.After` y `select`

---

## 📚 REGLAS Y MEJORES PRÁCTICAS APLICADAS

### 1. ⚠️ Funciones Anónimas

**Cuándo NO usarlas:**
- ❌ Para wrappear una simple llamada a método
- ❌ Cuando se puede extraer a un método con nombre
- ❌ Cuando se repite el patrón múltiples veces
- ❌ Cuando dificulta el testing

**Cuándo SÍ usarlas:**
- ✅ En callbacks one-shot donde extraer no agrega valor
- ✅ Cuando capturan variables locales que cambiarán
- ✅ En closures donde se necesita estado privado
- ✅ En operaciones muy triviales (1 línea, obvias)

**Ejemplo CORRECTO de uso:**
```go
// OK - Closure que captura estado
for i := 0; i < 10; i++ {
    i := i  // Capturar variable del loop
    go func() {
        fmt.Println(i)
    }()
}

// MEJOR - Extraer a método
for i := 0; i < 10; i++ {
    go procesarIndice(i)
}
```

---

### 2. ⏰ `time.Sleep` vs `time.After` + `select`

**❌ NUNCA usar `time.Sleep` en goroutines con contexto:**
```go
// MAL
go func() {
    time.Sleep(5 * time.Second)  // No cancelable
    // hacer algo
}()
```

**✅ SIEMPRE usar `time.After` con `select`:**
```go
// BIEN
go func() {
    select {
    case <-time.After(5 * time.Second):
        // hacer algo
    case <-ctx.Done():
        return  // Cancelable
    }
}()
```

**Razones:**
1. `time.Sleep` bloquea completamente, no es cancelable
2. `time.After` retorna un channel que se puede usar en `select`
3. `select` permite escuchar múltiples channels (timeout + cancelación)
4. Permite cierre limpio y rápido del programa

---

### 3. 🎯 Context Best Practices

**✅ Siempre escuchar `ctx.Done()` en goroutines:**
```go
for {
    select {
    case <-ctx.Done():
        return  // Salir inmediatamente
    default:
        // Trabajo normal
    }
}
```

**✅ Usar `select` para operaciones que pueden tardar:**
```go
select {
case <-time.After(duration):
    // Continuar
case result := <-ch:
    // Procesar
case <-ctx.Done():
    return  // Cancelar
}
```

**✅ Propagar contexto en llamadas:**
```go
func (s *Service) Start(ctx context.Context) {
    go s.worker1(ctx)  // Pasar contexto
    go s.worker2(ctx)  // Pasar contexto
}
```

---

### 4. 🧪 Testabilidad

**❌ Código difícil de testear:**
```go
func (s *Service) Start() {
    go func() {  // Función anónima
        time.Sleep(5 * time.Second)  // Sleep
        s.doSomething()
    }()
}
```

**✅ Código testeable:**
```go
func (s *Service) Start() {
    go s.workerWithTimeout(5 * time.Second)
}

func (s *Service) workerWithTimeout(timeout time.Duration) {
    select {
    case <-time.After(timeout):
        s.doSomething()
    case <-s.ctx.Done():
        return
    }
}

// Ahora se puede testear:
func TestWorkerWithTimeout(t *testing.T) {
    service.workerWithTimeout(100 * time.Millisecond)
    // Verificar comportamiento
}
```

---

## 📊 RESUMEN DE CAMBIOS

### Archivos Modificados:

1. **`internal/adapter/secondary/worker/cocinero.go`**
   - ✅ Reemplazado `time.Sleep(500ms)` por `select` con `time.After`
   - ✅ Reemplazado `time.Sleep(tiempoCoccion)` por `select` con `time.After`
   - ✅ Ambos ahora respetan cancelación de contexto

2. **`internal/domain/service/restaurant.go`**
   - ✅ Eliminada función anónima en `Start()` → `ejecutarCocinero()`
   - ✅ Eliminada función anónima en `EntregarPlatoAMesa()` → `limpiarMesaDespuesDeTiempo()`
   - ✅ Nuevo método helper usa `time.After` con `select` en lugar de `time.Sleep`

3. **`internal/adapter/primary/ui/ebiten_game.go`**
   - ✅ Eliminadas funciones anónimas triviales en callbacks
   - ✅ Se pasan métodos directamente
   - ✅ Nuevo método helper `handleClose()` para extensibilidad

### Impacto:

- ✅ **0 funciones anónimas innecesarias**
- ✅ **0 usos de `time.Sleep` en goroutines**
- ✅ **100% respeto a cancelación de contexto**
- ✅ **Mejor testabilidad**
- ✅ **Cierre más rápido y limpio**
- ✅ **Código más mantenible**

---

## 🎓 LECCIONES APRENDIDAS

### 1. Funciones Anónimas
**Principio:** Solo usar cuando realmente agregan valor (captura de estado, closures necesarios)

### 2. time.Sleep
**Principio:** Nunca en goroutines con contexto. Siempre `time.After` + `select`

### 3. Contexto
**Principio:** Todas las goroutines deben poder ser canceladas limpiamente

### 4. Métodos Helpers
**Principio:** Extraer lógica a métodos con nombre mejora legibilidad y testabilidad

### 5. Select Statement
**Principio:** Usar `select` para multiplexar channels (timeouts, cancelación, trabajo)

---

## ✅ VALIDACIÓN

### Compilación:
```bash
go build ./...
```
**Resultado:** ✅ Sin errores

### Race Detector:
```bash
go run -race cmd/app/main.go
```
**Resultado:** ✅ Sin race conditions

### Comportamiento:
- ✅ El programa sigue funcionando igual
- ✅ El cierre es más rápido (cancelación inmediata)
- ✅ No quedan goroutines zombie
- ✅ Todas las esperas respetan contexto

---

## 🎯 CONCLUSIÓN

El código ahora sigue las mejores prácticas de Go para concurrencia:

1. ✅ **No hay funciones anónimas innecesarias** - Solo métodos con nombre
2. ✅ **No hay `time.Sleep` en goroutines** - Solo `time.After` con `select`
3. ✅ **Todas las goroutines respetan el contexto** - Cancelación limpia
4. ✅ **Mejor testabilidad** - Métodos extraídos se pueden testear
5. ✅ **Código más mantenible** - Lógica con nombres descriptivos

**Calidad del código: Excelente** 🏆
