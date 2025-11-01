package main

import (
	"fmt"
	"time"

	"restaurant-concurrency/internal/adapter/secondary/worker"
	"restaurant-concurrency/internal/domain/service"
)

// Test simple para verificar la lógica de concurrencia sin UI
func main() {
	fmt.Println("🧪 TEST DE CONCURRENCIA - PRODUCTOR-CONSUMIDOR")
	fmt.Println("==================================================")
	fmt.Println()

	// Configuración
	capacidadBarra := 3
	numCocineros := 2
	numMeseros := 2

	// Crear servicio
	restaurantService := service.NewRestaurantService(
		nil,
		nil,
		capacidadBarra,
		numCocineros,
		numMeseros,
	)

	// Crear workers
	cocinero := worker.NewCocinero(restaurantService.(*service.RestaurantService))
	mesero := worker.NewMesero(restaurantService.(*service.RestaurantService))

	// Inyectar dependencias
	restaurantService.(*service.RestaurantService).SetProducer(cocinero)
	restaurantService.(*service.RestaurantService).SetConsumer(mesero)

	// Iniciar sistema
	fmt.Println("▶️  Iniciando sistema...")
	restaurantService.Start()

	// Escenario 1: Sin clientes (no debería producir)
	fmt.Println("\n📋 ESCENARIO 1: Sin clientes")
	time.Sleep(2 * time.Second)
	mostrarEstado(restaurantService, "Sin clientes - No debería haber producción")

	// Escenario 2: Agregar 5 clientes
	fmt.Println("\n📋 ESCENARIO 2: Agregar 5 clientes")
	restaurantService.AgregarClientes(5)
	time.Sleep(5 * time.Second)
	mostrarEstado(restaurantService, "Con clientes - Debería producir y consumir")

	// Escenario 3: Pausar producción
	fmt.Println("\n📋 ESCENARIO 3: Pausar producción")
	restaurantService.TogglePausar()
	time.Sleep(3 * time.Second)
	mostrarEstado(restaurantService, "Pausado - No debería producir más")

	// Escenario 4: Reanudar
	fmt.Println("\n📋 ESCENARIO 4: Reanudar producción")
	restaurantService.TogglePausar()
	time.Sleep(3 * time.Second)
	mostrarEstado(restaurantService, "Reanudado - Debería producir de nuevo")

	// Escenario 5: Clientes se van
	fmt.Println("\n📋 ESCENARIO 5: Clientes se van")
	restaurantService.ClientesSeVan(5)
	time.Sleep(2 * time.Second)
	mostrarEstado(restaurantService, "Sin clientes - Debería dejar de producir")

	// Cerrar sistema
	fmt.Println("\n🛑 Cerrando sistema...")
	restaurantService.Close()

	fmt.Println("\n✅ Test completado exitosamente")
	fmt.Println("📊 Estadísticas finales:")
	estado := restaurantService.GetEstado()
	fmt.Printf("   - Total producido: %d platos\n", estado.PlatosTotales)
	fmt.Printf("   - Total servido: %d platos\n", estado.PlatosServidos)
	fmt.Printf("   - En barra: %d platos\n", estado.EnBarra)
}

func mostrarEstado(svc interface{}, titulo string) {
	rs := svc.(*service.RestaurantService)
	estado := rs.GetEstado()

	fmt.Println("  ", titulo)
	fmt.Printf("     👥 Clientes: %d\n", estado.ClientesActivos)
	fmt.Printf("     📊 Barra: %d/%d\n", estado.EnBarra, estado.CapacidadBarra)
	fmt.Printf("     📈 Producidos: %d\n", estado.PlatosTotales)
	fmt.Printf("     ✅ Servidos: %d\n", estado.PlatosServidos)
	fmt.Printf("     ⏸️  Pausado: %v\n", estado.Pausado)
}
