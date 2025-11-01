package worker

import (
	"context"
	"fmt"
	"restaurant-concurrency/internal/domain/model"
	"restaurant-concurrency/internal/domain/service"
	"time"
)

// Mesero implementa port.Consumer
type Mesero struct {
	service *service.RestaurantService
}

func NewMesero(service *service.RestaurantService) *Mesero {
	return &Mesero{service: service}
}

func (m *Mesero) Consume(ctx context.Context, input <-chan model.Plato, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("🍽️  Mesero %d terminó\n", id)
			return
		case plato, ok := <-input:
			if !ok {
				return
			}

			// Simular entrega
			time.Sleep(600 * time.Millisecond)

			m.service.IncrementarPlatosServidos()
			fmt.Printf("🍽️  Mesero %d entregó plato #%d\n", id, plato.ID)
		}
	}
}
