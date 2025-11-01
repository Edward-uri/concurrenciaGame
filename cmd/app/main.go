package main

import (
	"fmt"
	"log"

	"restaurant-concurrency/internal/adapter/primary/ui"
	"restaurant-concurrency/internal/adapter/secondary/worker"
	"restaurant-concurrency/internal/domain/service"
	infrastructure "restaurant-concurrency/internal/infraestructure"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth    = 1920
	screenHeight   = 1080
	capacidadBarra = 5
	numCocineros   = 1 // Solo 1 chef produciendo
	numMeseros     = 1 // 1 mesero controlado por el jugador
)

func main() {
	fmt.Println("🏗️  Arquitectura Hexagonal")
	fmt.Println("🔷 Patrón: Productor-Consumidor")
	fmt.Println()

	// Crear logger
	logger, err := infrastructure.NewLogger(infrastructure.LoggingConfig{
		Level:      "info",
		Structured: false,
		Output:     "stdout",
	})
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	// Crear servicio (núcleo del dominio)
	restaurantService := service.NewRestaurantService(
		nil, // Se inyectan después para evitar dependencia circular
		nil,
		capacidadBarra,
		numCocineros,
		numMeseros,
	)

	// Crear adapters secundarios (workers)
	cocinero := worker.NewCocinero(restaurantService.(*service.RestaurantService))
	mesero := worker.NewMesero(restaurantService.(*service.RestaurantService))

	// Inyectar dependencias
	restaurantService.(*service.RestaurantService).SetProducer(cocinero)
	restaurantService.(*service.RestaurantService).SetConsumer(mesero)

	// Iniciar sistema
	restaurantService.Start()
	restaurantService.AgregarClientes(3)

	// Crear adapter primario (UI)
	game, err := ui.NewEbitenGame(restaurantService, logger, screenWidth, screenHeight)
	if err != nil {
		log.Fatalf("Error creating game: %v", err)
	}

	// Ejecutar
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hexagonal Architecture - Restaurant")

	if err := ebiten.RunGame(game); err != nil {
		log.Println("Cerrando:", err)
	}

	restaurantService.Close()
	fmt.Println("✅ Sistema cerrado correctamente")
}
