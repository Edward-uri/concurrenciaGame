package port

import (
	"context"
	"restaurant-concurrency/internal/domain/model"
)

// Consumer define el contrato de un consumidor
type Consumer interface {
	Consume(ctx context.Context, input <-chan model.Plato, id int)
}
