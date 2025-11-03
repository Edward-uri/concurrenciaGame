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
			fmt.Printf("Cocinero %d terminó su turno\n", c.id)
			return

		default:
			// Solo producir si hay demanda (clientes esperando)
			if !verificarDemanda() {
				// Espera no bloqueante usando select con time.After
				select {
				case <-time.After(500 * time.Millisecond):
					continue
				case <-ctx.Done():
					return
				}
			}

			// Simular tiempo de cocción (trabajo concurrente) usando time.After
			tiempoCoccion := time.Duration(1500+rand.Intn(1000)) * time.Millisecond

			select {
			case <-time.After(tiempoCoccion):
				// Continuar con la producción
			case <-ctx.Done():
				return
			}

			// Crear plato
			plato := model.NewPlato(platoID, c.id)

			// INTENTAR PONER EN LA BARRA (canal buffered)
			// Si la barra está llena, SE BLOQUEA aquí hasta que haya espacio
			// Este es el comportamiento del patrón Productor-Consumidor
			select {
			case barra <- plato:
				fmt.Printf("Cocinero %d preparó plato #%d (tiempo: %.1fs)\n",
					c.id, platoID, tiempoCoccion.Seconds())
				platoID++

			case <-ctx.Done():
				return
			}
		}
	}
}
