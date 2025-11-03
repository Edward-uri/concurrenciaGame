package service

import (
	"context"
	"math/rand"
	"restaurant-concurrency/internal/adapter/secondary/worker"
	"restaurant-concurrency/internal/domain/model"
	"sync"
	"time"
)

type RestaurantService struct {
	// Canal productor-consumidor (BUFFER)
	barra          chan model.Plato
	capacidadBarra int

	// Mesas y clientes
	mesas   []*model.Mesa
	mesasMu sync.RWMutex

	// Métricas
	mu               sync.RWMutex
	platosTotales    int
	platosServidos   int
	clientesPerdidos int
	pausado          bool

	// Concurrencia
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Workers (adapters secundarios)
	cocineros    []*worker.Cocinero
	numCocineros int
}

func NewRestaurantService(capacidadBarra, numCocineros, numMesas int) *RestaurantService {
	ctx, cancel := context.WithCancel(context.Background())

	service := &RestaurantService{
		barra:          make(chan model.Plato, capacidadBarra),
		capacidadBarra: capacidadBarra,
		numCocineros:   numCocineros,
		ctx:            ctx,
		cancel:         cancel,
		mesas:          make([]*model.Mesa, 0, numMesas),
		cocineros:      make([]*worker.Cocinero, 0, numCocineros),
	}

	// Crear cocineros (productores en el patrón Productor-Consumidor)
	for i := 1; i <= numCocineros; i++ {
		service.cocineros = append(service.cocineros, worker.NewCocinero(i))
	}

	// Crear mesas
	positions := [][2]float64{
		{100, 300}, {300, 300}, {500, 300}, {700, 300},
		{100, 450}, {300, 450}, {500, 450}, {700, 450},
	}

	for i := 0; i < numMesas && i < len(positions); i++ {
		mesa := model.NewMesa(i, positions[i][0], positions[i][1], 30*time.Second)
		service.mesas = append(service.mesas, mesa)
	}

	return service
}

// Start inicia todas las goroutines
func (s *RestaurantService) Start() {
	// Iniciar cocineros
	for _, cocinero := range s.cocineros {
		s.wg.Add(1)
		go s.ejecutarCocinero(cocinero)
	}

	// Generador de clientes
	s.wg.Add(1)
	go s.generadorClientes()

	// Verificador de paciencia
	s.wg.Add(1)
	go s.verificadorPaciencia()
}

// ejecutarCocinero es un método helper para evitar función anónima
func (s *RestaurantService) ejecutarCocinero(cocinero *worker.Cocinero) {
	defer s.wg.Done()
	cocinero.Producir(s.ctx, s.barra, s.hayDemanda)
}

// hayDemanda verifica si hay clientes esperando (para que cocineros produzcan)
func (s *RestaurantService) hayDemanda() bool {
	s.mesasMu.RLock()
	defer s.mesasMu.RUnlock()

	s.mu.RLock()
	pausado := s.pausado
	s.mu.RUnlock()

	if pausado {
		return false
	}

	for _, mesa := range s.mesas {
		if mesa.ClientesActivos > 0 && !mesa.TienePlato {
			return true
		}
	}
	return false
}

func (s *RestaurantService) generadorClientes() {
	defer s.wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			// Agregar clientes aleatoriamente a mesas vacías
			s.mesasMu.Lock()
			for _, mesa := range s.mesas {
				if mesa.ClientesActivos == 0 && rand.Float64() < 0.4 {
					cantidadClientes := rand.Intn(3) + 1
					mesa.AgregarClientes(cantidadClientes)
				}
			}
			s.mesasMu.Unlock()
		}
	}
}

func (s *RestaurantService) verificadorPaciencia() {
	defer s.wg.Done()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.mesasMu.Lock()
			for _, mesa := range s.mesas {
				if mesa.ClientesActivos > 0 && !mesa.EstaPaciente() {
					// Clientes se fueron por falta de servicio
					s.mu.Lock()
					s.clientesPerdidos += mesa.ClientesActivos
					s.mu.Unlock()
					mesa.ClientesSatisfechos()
				}
			}
			s.mesasMu.Unlock()
		}
	}
}

// IntentarRecogerPlato permite al mesero (jugador) CONSUMIR de la barra
func (s *RestaurantService) IntentarRecogerPlato() (*model.Plato, bool) {
	select {
	case plato := <-s.barra:
		s.mu.Lock()
		s.platosTotales++
		s.mu.Unlock()
		return &plato, true
	default:
		return nil, false
	}
}

// EntregarPlatoAMesa entrega un plato a una mesa cercana
func (s *RestaurantService) EntregarPlatoAMesa(meseroX, meseroY float64, rango float64) bool {
	s.mesasMu.Lock()
	defer s.mesasMu.Unlock()

	for _, mesa := range s.mesas {
		if mesa.ClientesActivos > 0 && !mesa.TienePlato {
			// Verificar distancia
			dx := mesa.PosX - meseroX
			dy := mesa.PosY - meseroY
			distancia := dx*dx + dy*dy

			if distancia < rango*rango {
				mesa.EntregarPlato()
				s.mu.Lock()
				s.platosServidos++
				s.mu.Unlock()

				// Después de un tiempo, clientes se van satisfechos
				go s.limpiarMesaDespuesDeTiempo(mesa, 3*time.Second)

				return true
			}
		}
	}
	return false
}

// GetMesas retorna snapshots inmutables de las mesas (thread-safe para rendering)
func (s *RestaurantService) GetMesas() []model.MesaSnapshot {
	s.mesasMu.RLock()
	defer s.mesasMu.RUnlock()

	snapshots := make([]model.MesaSnapshot, len(s.mesas))
	for i, mesa := range s.mesas {
		snapshots[i] = mesa.Snapshot()
	}
	return snapshots
}

func (s *RestaurantService) GetEstadoBarra() int {
	return len(s.barra)
}

func (s *RestaurantService) GetCapacidadBarra() int {
	return s.capacidadBarra
}

func (s *RestaurantService) GetMetricas() (totales, servidos, perdidos int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.platosTotales, s.platosServidos, s.clientesPerdidos
}

func (s *RestaurantService) TogglePausar() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pausado = !s.pausado
}

// limpiarMesaDespuesDeTiempo limpia la mesa después de un tiempo especificado
// Usa time.After con select en lugar de time.Sleep para respetar cancelación
func (s *RestaurantService) limpiarMesaDespuesDeTiempo(mesa *model.Mesa, duracion time.Duration) {
	select {
	case <-time.After(duracion):
		s.mesasMu.Lock()
		mesa.ClientesSatisfechos()
		s.mesasMu.Unlock()
	case <-s.ctx.Done():
		// Si se cancela el contexto, no limpiar la mesa
		return
	}
}

func (s *RestaurantService) Close() {
	s.cancel()
	s.wg.Wait()
	close(s.barra)
}
