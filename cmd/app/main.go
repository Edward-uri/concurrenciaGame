package main

import (
	"fmt"
	"log"

	"restaurant-concurrency/internal/adapter/primary/ui"
	"restaurant-concurrency/internal/domain/service"
	infrastructure "restaurant-concurrency/internal/infraestructure"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth    = 1920
	screenHeight   = 1080
	capacidadBarra = 5
	numCocineros   = 1 // 1 chef produciendo platos
	numMeseros     = 1 // 1 mesero controlado por el jugador
	numMesas       = 3 // 8 mesas con clientes
)

func main() {

	// Configuración inicial
	fmt.Println("CONFIGURACION:")
	fmt.Printf("   • Cocineros (productores): %d\n", numCocineros)
	fmt.Printf("   • Capacidad de barra (buffer): %d\n", capacidadBarra)
	fmt.Printf("   • Mesas con clientes: %d\n", numMesas)
	fmt.Printf("   • Resolución: %dx%d\n", screenWidth, screenHeight)
	fmt.Println()

	// Crear logger
	logger, err := infrastructure.NewLogger(infrastructure.LoggingConfig{
		Level:      "info",
		Structured: false,
		Output:     "stdout",
	})
	if err != nil {
		log.Fatalf("Error al crear logger: %v", err)
	}
	logger.Info("Logger inicializado correctamente")

	// Crear servicio del restaurante
	fmt.Println("Inicializando servicio del restaurante...")
	restaurantService := service.NewRestaurantService(
		capacidadBarra,
		numCocineros,
		numMesas,
	)

	// Iniciar las goroutines concurrentes
	// - Cocineros (productores automáticos)
	// - Generador de clientes
	restaurantService.Start()
	logger.Info("Sistema de concurrencia iniciado")

	// Crear el juego con Ebiten
	fmt.Println("Inicializando interfaz gráfica...")
	game, err := ui.NewGame(restaurantService, screenWidth, screenHeight)
	if err != nil {
		log.Fatalf("Error al crear el juego: %v", err)
	}
	logger.Info("Interfaz gráfica inicializada")

	// Configurar ventana
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Restaurante Concurrente - Arquitectura Hexagonal")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(60) // 60 FPS

	if err := ebiten.RunGame(game); err != nil {
		if err.Error() != "cierre solicitado por usuario" {
			log.Println("Error durante la ejecución:", err)
		}
	}

	// ============ CIERRE ORDENADO ============
	restaurantService.Close()
	logger.Info("Sistema cerrado correctamente")
}
