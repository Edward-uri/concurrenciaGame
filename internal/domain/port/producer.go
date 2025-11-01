package port

import (
	"context"
	"restaurant-concurrency/internal/domain/model"
)

// Producer define el contrato de un productor
type Producer interface {
	Produce(ctx context.Context, output chan<- model.Plato, id int)
}
