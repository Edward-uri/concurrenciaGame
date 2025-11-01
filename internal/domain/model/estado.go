package model

type EstadoRestaurant struct {
	ClientesActivos int
	PlatosTotales   int
	PlatosServidos  int
	EnBarra         int
	CapacidadBarra  int
	Pausado         bool
}
