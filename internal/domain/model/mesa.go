package model

import "time"

// Mesa representa una mesa con clientes esperando
type Mesa struct {
	ID              int
	PosX, PosY      float64
	ClientesActivos int
	TienePlato      bool
	TiempoEspera    time.Time
	Paciencia       time.Duration // Si tarda mucho, se van
}

func NewMesa(id int, x, y float64, paciencia time.Duration) *Mesa {
	return &Mesa{
		ID:         id,
		PosX:       x,
		PosY:       y,
		Paciencia:  paciencia,
		TienePlato: false,
	}
}

// AgregarClientes añade clientes a la mesa
func (m *Mesa) AgregarClientes(cantidad int) {
	if m.ClientesActivos == 0 {
		m.TiempoEspera = time.Now()
	}
	m.ClientesActivos += cantidad
}

// EntregarPlato marca que se entregó un plato
func (m *Mesa) EntregarPlato() {
	m.TienePlato = true
}

// ClientesSatisfechos limpia la mesa
func (m *Mesa) ClientesSatisfechos() {
	m.ClientesActivos = 0
	m.TienePlato = false
}

// EstaPaciente verifica si los clientes siguen esperando
func (m *Mesa) EstaPaciente() bool {
	if m.ClientesActivos == 0 {
		return true
	}
	return time.Since(m.TiempoEspera) < m.Paciencia
}

// GetNivelPaciencia retorna valor 0.0 a 1.0 (1.0 = muy impacientes)
func (m *Mesa) GetNivelPaciencia() float64 {
	if m.ClientesActivos == 0 {
		return 0
	}
	elapsed := time.Since(m.TiempoEspera)
	return float64(elapsed) / float64(m.Paciencia)
}

// MesaSnapshot es una copia inmutable de los datos de Mesa para renderizado thread-safe
type MesaSnapshot struct {
	ID              int
	PosX, PosY      float64
	ClientesActivos int
	TienePlato      bool
	NivelPaciencia  float64
}

// Snapshot crea una copia thread-safe de los datos de la mesa
// DEBE ser llamado mientras se tiene el lock de mesasMu
func (m *Mesa) Snapshot() MesaSnapshot {
	return MesaSnapshot{
		ID:              m.ID,
		PosX:            m.PosX,
		PosY:            m.PosY,
		ClientesActivos: m.ClientesActivos,
		TienePlato:      m.TienePlato,
		NivelPaciencia:  m.GetNivelPaciencia(),
	}
}
