package worker

import (
	"context"
	"fmt"
	"math/rand"
	"restaurant-concurrency/internal/domain/model"
	"time"
)

// Cocinero es el worker que implementa el PRODUCTOR
// Este es un adapter secundario que ejecuta la lógica de producción
type Cocinero struct {
	id int
}

func NewCocinero(id int) *Cocinero {
	return &Cocinero{id: id}
}

// Producir ejecuta el loop de producción (goroutine)
// Recibe:
// - ctx: para cancelación
// - barra: canal donde deposita platos (buffer)
// - verificarDemanda: función que verifica si hay clientes esperando
func (c *Cocinero) Producir(
	ctx context.Context,
	barra chan<- model.Plato,
	verificarDemanda func() bool,
) {
	platoID := 0

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("👨‍🍳 Cocinero %d terminó su turno\n", c.id)
			return

		default:
			// Solo producir si hay demanda (clientes esperando)
			if !verificarDemanda() {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			// Simular tiempo de cocción (trabajo concurrente)
			tiempoCoccion := time.Duration(1500+rand.Intn(1000)) * time.Millisecond
			time.Sleep(tiempoCoccion)

			// Crear plato
			plato := model.NewPlato(platoID, c.id)

			// INTENTAR PONER EN LA BARRA
			// Si la barra está llena, SE BLOQUEA aquí (comportamiento del patrón)
			select {
			case barra <- plato:
				fmt.Printf("👨‍🍳 Cocinero %d preparó plato #%d (tiempo: %.1fs)\n",
					c.id, platoID, tiempoCoccion.Seconds())
				platoID++

			case <-ctx.Done():
				return
			}
		}
	}
}
