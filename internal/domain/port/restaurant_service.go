package port

import "restaurant-concurrency/internal/domain/model"

// RestaurantService define el contrato del servicio principal
type RestaurantService interface {
	// Control de clientes
	AgregarClientes(cantidad int)
	ClientesSeVan(cantidad int)

	// Control de producción
	TogglePausar()

	// Observabilidad
	GetEstado() model.EstadoRestaurant
	GetBarra() []model.Plato

	// Consumir plato (para UI manual)
	ConsumirPlato() *model.Plato
	EntregarPlato()

	// Ciclo de vida
	Start()
	Close()
}
