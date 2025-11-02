package model

import "time"

// EstadoMesero representa el estado del mesero controlable
type EstadoMesero int

const (
	MeseroIdle EstadoMesero = iota
	MeseroCaminando
	MeseroRecogiendo
	MeseroEntregando
)

// Mesero es el personaje controlable por el jugador
type Mesero struct {
	PosX, PosY       float64
	VelocidadX       float64
	VelocidadY       float64
	Speed            float64
	TienePlato       bool
	PlatoEnMano      *Plato
	Estado           EstadoMesero
	UltimoMovimiento time.Time
}

func NewMesero(x, y, speed float64) *Mesero {
	return &Mesero{
		PosX:   x,
		PosY:   y,
		Speed:  speed,
		Estado: MeseroIdle,
	}
}

// Mover actualiza la posición del mesero
func (m *Mesero) Mover(dx, dy float64, deltaTime float64) {
	m.PosX += dx * m.Speed * deltaTime
	m.PosY += dy * m.Speed * deltaTime

	if dx != 0 || dy != 0 {
		m.Estado = MeseroCaminando
		m.UltimoMovimiento = time.Now()
	} else {
		if time.Since(m.UltimoMovimiento) > 100*time.Millisecond {
			m.Estado = MeseroIdle
		}
	}
}

// RecogerPlato asigna un plato al mesero
func (m *Mesero) RecogerPlato(plato Plato) {
	m.TienePlato = true
	m.PlatoEnMano = &plato
	m.Estado = MeseroRecogiendo
}

// EntregarPlato entrega el plato y limpia las manos
func (m *Mesero) EntregarPlato() *Plato {
	plato := m.PlatoEnMano
	m.TienePlato = false
	m.PlatoEnMano = nil
	m.Estado = MeseroEntregando
	return plato
}

// GetBounds retorna el rectángulo de colisión
func (m *Mesero) GetBounds() (x, y, width, height float64) {
	return m.PosX, m.PosY, 32, 48 // Ajusta según tu sprite
}
