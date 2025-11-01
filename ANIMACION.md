# 🎬 Sistema de Animación de Meseros

## Descripción

El sistema de animación hace visible el patrón **Productor-Consumidor** mostrando cómo los meseros (consumidores) se mueven para:
1. **Ir a la barra** cuando hay platos disponibles
2. **Tomar un plato** del buffer (canal)
3. **Llevar el plato al cliente** 
4. **Regresar** a su posición inicial

## Estados del Mesero

Cada mesero animado puede estar en uno de estos estados:

- `MeseroEsperando` - En su posición inicial, esperando trabajo
- `MeseroYendoABarra` - Moviéndose hacia la barra para tomar un plato
- `MeseroTomandoPlato` - Tomando el plato de la barra (pausa breve)
- `MeseroLlevandoACliente` - Llevando el plato al cliente
- `MeseroRegresando` - Regresando a su posición inicial

## Cómo Funciona

### 1. Detección de Trabajo Disponible
```go
// Si hay platos en barra y el mesero está libre
if mesero.Estado == MeseroEsperando && platosEnBarra > 0 && hay_clientes {
    mesero.IrABarra(posicionBarra)
}
```

### 2. Movimiento Suave
Los meseros se mueven pixel por pixel hacia su destino:
- Velocidad: 200 pixels/segundo
- Movimiento interpolado para transición suave
- Detección automática al llegar al destino

### 3. Visual del Estado
- El mesero cambia de posición en tiempo real
- Muestra un plato flotante cuando lo lleva
- Etiqueta con el estado actual (→Barra, →Cliente, etc.)

## Integración con el Patrón Productor-Consumidor

```
COCINEROS (Productores)
    ↓
    Crean platos
    ↓
BARRA (Buffer/Canal) ← MESEROS detectan platos disponibles
    ↓
    Mesero se mueve a la barra
    ↓
    Toma plato del buffer
    ↓
CLIENTE ← Mesero lleva el plato
    ↓
    Entrega y regresa
```

## Visualización

- **Cocineros**: Estáticos en el lado izquierdo
- **Barra**: Centro, muestra platos disponibles
- **Meseros**: Se mueven dinámicamente entre barra y clientes
- **Plato en mano**: Visible cuando el mesero lo transporta

## Ejecución

```bash
go run ./cmd/app
```

### Controles
- `+` : Agregar cliente (más trabajo para meseros)
- `-` : Quitar cliente
- `ESPACIO` : Pausar producción (los meseros terminan su trabajo actual)
- `ESC` : Salir

## Mejoras Futuras Posibles

1. **Animación de cocineros cocinando**
2. **Clientes con indicador de satisfacción**
3. **Diferentes tipos de platos con colores**
4. **Trayectorias más complejas (evitar obstáculos)**
5. **Sonidos al tomar/entregar platos**
6. **Partículas/efectos visuales**
