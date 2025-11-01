package service

import (
	"context"
	"restaurant-concurrency/internal/domain/model"
	"restaurant-concurrency/internal/domain/port"
	"sync"
)

// RestaurantService implementa la lógica de concurrencia
type RestaurantService struct {
	// Dependencias (puertos)
	producer port.Producer
	consumer port.Consumer

	// Estado del dominio
	mu              sync.RWMutex
	clientesActivos int
	pausado         bool
	platosTotales   int
	platosServidos  int

	// Infraestructura de concurrencia
	barra  chan model.Plato
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Configuración
	capacidadBarra int
	numCocineros   int
	numMeseros     int
}

// NewRestaurantService crea el servicio con inyección de dependencias
func NewRestaurantService(
	producer port.Producer,
	consumer port.Consumer,
	capacidadBarra, numCocineros, numMeseros int,
) port.RestaurantService {
	ctx, cancel := context.WithCancel(context.Background())

	return &RestaurantService{
		producer:       producer,
		consumer:       consumer,
		barra:          make(chan model.Plato, capacidadBarra),
		ctx:            ctx,
		cancel:         cancel,
		capacidadBarra: capacidadBarra,
		numCocineros:   numCocineros,
		numMeseros:     numMeseros,
	}
}

func (s *RestaurantService) Start() {
	// Iniciar productores
	for i := 1; i <= s.numCocineros; i++ {
		s.wg.Add(1)
		go func(id int) {
			defer s.wg.Done()
			s.producer.Produce(s.ctx, s.barra, id)
		}(i)
	}

	// Iniciar consumidores
	for i := 1; i <= s.numMeseros; i++ {
		s.wg.Add(1)
		go func(id int) {
			defer s.wg.Done()
			s.consumer.Consume(s.ctx, s.barra, id)
		}(i)
	}
}

func (s *RestaurantService) AgregarClientes(cantidad int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clientesActivos += cantidad
}

func (s *RestaurantService) ClientesSeVan(cantidad int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cantidad > s.clientesActivos {
		cantidad = s.clientesActivos
	}
	s.clientesActivos -= cantidad
}

func (s *RestaurantService) TogglePausar() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pausado = !s.pausado
}

// SetProducer permite inyectar o reemplazar el productor en tiempo de ejecución
func (s *RestaurantService) SetProducer(p port.Producer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.producer = p
}

// SetConsumer permite inyectar o reemplazar el consumidor en tiempo de ejecución
func (s *RestaurantService) SetConsumer(c port.Consumer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.consumer = c
}

func (s *RestaurantService) GetEstado() model.EstadoRestaurant {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return model.EstadoRestaurant{
		ClientesActivos: s.clientesActivos,
		PlatosTotales:   s.platosTotales,
		PlatosServidos:  s.platosServidos,
		EnBarra:         len(s.barra),
		CapacidadBarra:  s.capacidadBarra,
		Pausado:         s.pausado,
	}
}

func (s *RestaurantService) GetBarra() []model.Plato {
	// Solo retorna la cantidad, no los platos reales
	// (para evitar race conditions con el canal)
	// La UI puede dibujar basándose en GetEstado().EnBarra
	return nil
}

func (s *RestaurantService) IncrementarPlatosTotales() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.platosTotales++
}

func (s *RestaurantService) IncrementarPlatosServidos() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.platosServidos++
}

func (s *RestaurantService) DebeProducir() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.clientesActivos > 0 && !s.pausado
}

// ConsumirPlato intenta tomar un plato de la barra (no bloqueante)
func (s *RestaurantService) ConsumirPlato() *model.Plato {
	select {
	case plato := <-s.barra:
		return &plato
	default:
		return nil // No hay platos disponibles
	}
}

// EntregarPlato incrementa el contador de platos servidos
func (s *RestaurantService) EntregarPlato() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.platosServidos++
}

func (s *RestaurantService) Close() {
	s.cancel()
	s.wg.Wait()
	close(s.barra)
}
