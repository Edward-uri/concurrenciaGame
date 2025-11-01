package model

import "time"

type Plato struct {
	ID         int
	Nombre     string
	CocineroID int
	Timestamp  time.Time
}

func NewPlato(id, cocineroID int) Plato {
	return Plato{
		ID:         id,
		Nombre:     "Plato especial",
		CocineroID: cocineroID,
		Timestamp:  time.Now(),
	}
}
