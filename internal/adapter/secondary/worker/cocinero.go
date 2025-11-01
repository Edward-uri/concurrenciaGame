package worker

import (
	"context"
	"fmt"
	"restaurant-concurrency/internal/domain/model"
	"restaurant-concurrency/internal/domain/service"
	"time"
)

// Cocinero implementa port.Producer
type Cocinero struct {
	service *service.RestaurantService
}

func NewCocinero(service *service.RestaurantService) *Cocinero {
	return &Cocinero{service: service}
}

func (c *Cocinero) Produce(ctx context.Context, output chan<- model.Plato, id int) {
	platoID := 0

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("👨‍🍳 Cocinero %d terminó\n", id)
			return
		default:
			// Verificar si debe producir
			if !c.service.DebeProducir() {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			// Simular cocción
			time.Sleep(800 * time.Millisecond)

			plato := model.NewPlato(platoID, id)

			select {
			case output <- plato:
				c.service.IncrementarPlatosTotales()
				fmt.Printf("👨‍🍳 Cocinero %d preparó plato #%d\n", id, platoID)
				platoID++
			case <-ctx.Done():
				return
			}
		}
	}
}
