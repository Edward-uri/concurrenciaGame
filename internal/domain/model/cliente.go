package model

import (
	"time"
)

// Cliente representa un cliente en el restaurante
type Cliente struct {
	ID            int
	Nombre        string
	TiempoLlegada time.Time
	Satisfecho    bool
	SkinIndex     int // Índice del sprite en el spritesheet
}

// NewCliente crea un nuevo cliente
func NewCliente(id int, skinIndex int) Cliente {
	return Cliente{
		ID:            id,
		Nombre:        "Cliente",
		TiempoLlegada: time.Now(),
		Satisfecho:    false,
		SkinIndex:     skinIndex,
	}
}

// MarcarSatisfecho marca al cliente como satisfecho
func (c *Cliente) MarcarSatisfecho() {
	c.Satisfecho = true
}

// TiempoEspera retorna cuánto tiempo lleva esperando el cliente
func (c *Cliente) TiempoEspera() time.Duration {
	return time.Since(c.TiempoLlegada)
}
