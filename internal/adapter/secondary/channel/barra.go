package channel

import (
	"restaurant-concurrency/internal/domain/model"
	"sync"
)

// Barra representa el buffer (canal) con funcionalidad extendida
type Barra struct {
	canal     chan model.Plato
	capacidad int
	mu        sync.RWMutex
	platos    []model.Plato // Snapshot para visualización
}

// NewBarra crea una nueva barra con la capacidad especificada
func NewBarra(capacidad int) *Barra {
	return &Barra{
		canal:     make(chan model.Plato, capacidad),
		capacidad: capacidad,
		platos:    make([]model.Plato, 0),
	}
}

// Push agrega un plato al canal (puede bloquear si está llena)
func (b *Barra) Push(plato model.Plato) {
	b.canal <- plato
	b.actualizarSnapshot()
}

// Pop extrae un plato del canal (puede bloquear si está vacía)
func (b *Barra) Pop() (model.Plato, bool) {
	plato, ok := <-b.canal
	if ok {
		b.actualizarSnapshot()
	}
	return plato, ok
}

// TryPush intenta agregar sin bloquear
func (b *Barra) TryPush(plato model.Plato) bool {
	select {
	case b.canal <- plato:
		b.actualizarSnapshot()
		return true
	default:
		return false
	}
}

// TryPop intenta extraer sin bloquear
func (b *Barra) TryPop() (model.Plato, bool) {
	select {
	case plato := <-b.canal:
		b.actualizarSnapshot()
		return plato, true
	default:
		return model.Plato{}, false
	}
}

// GetSnapshot retorna una copia del estado actual (thread-safe)
func (b *Barra) GetSnapshot() []model.Plato {
	b.mu.RLock()
	defer b.mu.RUnlock()

	snapshot := make([]model.Plato, len(b.platos))
	copy(snapshot, b.platos)
	return snapshot
}

// actualizarSnapshot actualiza el estado interno para visualización
func (b *Barra) actualizarSnapshot() {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Crear snapshot no destructivo
	b.platos = make([]model.Plato, 0, len(b.canal))
	temp := make([]model.Plato, 0)

	// Extraer platos temporalmente
	done := false
	for !done {
		select {
		case p := <-b.canal:
			temp = append(temp, p)
		default:
			done = true
		}
	}

	// Devolver platos al canal
	for _, p := range temp {
		b.canal <- p
		b.platos = append(b.platos, p)
	}
}

// Len retorna la cantidad de platos en la barra
func (b *Barra) Len() int {
	return len(b.canal)
}

// Cap retorna la capacidad máxima de la barra
func (b *Barra) Cap() int {
	return b.capacidad
}

// IsFull indica si la barra está llena
func (b *Barra) IsFull() bool {
	return len(b.canal) == b.capacidad
}

// IsEmpty indica si la barra está vacía
func (b *Barra) IsEmpty() bool {
	return len(b.canal) == 0
}

// Close cierra el canal
func (b *Barra) Close() {
	close(b.canal)
}

// GetChannel retorna el canal interno (para uso en select statements)
func (b *Barra) GetChannel() <-chan model.Plato {
	return b.canal
}

// GetWriteChannel retorna el canal para escritura
func (b *Barra) GetWriteChannel() chan<- model.Plato {
	return b.canal
}
