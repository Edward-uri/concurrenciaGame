package ui

import (
	"time"
)

// EstadoMesero representa el estado de un mesero en la animación
type EstadoMesero string

const (
	MeseroEsperando        EstadoMesero = "esperando"   // En posición inicial
	MeseroYendoABarra      EstadoMesero = "yendo_barra" // Moviéndose hacia la barra
	MeseroTomandoPlato     EstadoMesero = "tomando"     // Tomando plato de la barra
	MeseroLlevandoACliente EstadoMesero = "llevando"    // Llevando plato al cliente
	MeseroRegresando       EstadoMesero = "regresando"  // Regresando a posición inicial
)

// MeseroAnimado representa un mesero con animación
type MeseroAnimado struct {
	ID                 int
	Estado             EstadoMesero
	X, Y               float64 // Posición actual
	InicioX, InicioY   float64 // Posición inicial
	DestinoX, DestinoY float64 // Posición destino
	Velocidad          float64 // Pixels por segundo
	TiempoInicio       time.Time
	TieneAlgo          bool // Si lleva un plato
	PlatoID            int  // ID del plato que lleva
}

// NewMeseroAnimado crea un nuevo mesero animado
func NewMeseroAnimado(id int, x, y float64) *MeseroAnimado {
	return &MeseroAnimado{
		ID:        id,
		Estado:    MeseroEsperando,
		X:         x,
		Y:         y,
		InicioX:   x,
		InicioY:   y,
		Velocidad: 200.0, // pixels por segundo
	}
}

// IrABarra hace que el mesero vaya a la barra
func (m *MeseroAnimado) IrABarra(barraX, barraY float64) {
	m.Estado = MeseroYendoABarra
	m.DestinoX = barraX
	m.DestinoY = barraY
	m.TiempoInicio = time.Now()
	m.TieneAlgo = false
}

// TomarPlato simula que el mesero toma un plato
func (m *MeseroAnimado) TomarPlato(platoID int) {
	m.Estado = MeseroTomandoPlato
	m.TieneAlgo = true
	m.PlatoID = platoID
	m.TiempoInicio = time.Now()
}

// LlevarACliente hace que el mesero lleve el plato al cliente
func (m *MeseroAnimado) LlevarACliente(clienteX, clienteY float64) {
	m.Estado = MeseroLlevandoACliente
	m.DestinoX = clienteX
	m.DestinoY = clienteY
	m.TiempoInicio = time.Now()
}

// Regresar hace que el mesero regrese a su posición inicial
func (m *MeseroAnimado) Regresar() {
	m.Estado = MeseroRegresando
	m.DestinoX = m.InicioX
	m.DestinoY = m.InicioY
	m.TiempoInicio = time.Now()
	m.TieneAlgo = false
}

// Actualizar actualiza la posición del mesero según su estado
func (m *MeseroAnimado) Actualizar(deltaTime float64) {
	switch m.Estado {
	case MeseroYendoABarra, MeseroLlevandoACliente, MeseroRegresando:
		// Calcular dirección
		dx := m.DestinoX - m.X
		dy := m.DestinoY - m.Y
		distancia := sqrt(dx*dx + dy*dy)

		if distancia < 2.0 {
			// Llegó al destino
			m.X = m.DestinoX
			m.Y = m.DestinoY

			// Cambiar estado según a dónde llegó
			if m.Estado == MeseroYendoABarra {
				m.Estado = MeseroEsperando // Esperar a tomar plato
			} else if m.Estado == MeseroLlevandoACliente {
				m.Estado = MeseroEsperando // Entregó el plato
			} else if m.Estado == MeseroRegresando {
				m.Estado = MeseroEsperando
			}
		} else {
			// Mover hacia el destino
			velocidadFrame := m.Velocidad * deltaTime
			factor := velocidadFrame / distancia
			if factor > 1.0 {
				factor = 1.0
			}
			m.X += dx * factor
			m.Y += dy * factor
		}

	case MeseroTomandoPlato:
		// Esperar un momento simulando que toma el plato
		if time.Since(m.TiempoInicio).Seconds() > 0.3 {
			m.Estado = MeseroEsperando
		}
	}
}

func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	// Aproximación simple de sqrt
	z := x
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}
