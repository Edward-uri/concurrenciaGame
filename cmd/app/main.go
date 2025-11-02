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
	numMesas       = 8 // 8 mesas con clientes
)

func main() {

	// Configuración inicial
	fmt.Println("📋 Configuración:")
	fmt.Printf("   • Cocineros (productores): %d\n", numCocineros)
	fmt.Printf("   • Capacidad de barra (buffer): %d\n", capacidadBarra)
	fmt.Printf("   • Mesas con clientes: %d\n", numMesas)
	fmt.Printf("   • Resolución: %dx%d\n", screenWidth, screenHeight)
	fmt.Println()

	// ============ INFRAESTRUCTURA ============
	// Crear logger
	logger, err := infrastructure.NewLogger(infrastructure.LoggingConfig{
		Level:      "info",
		Structured: false,
		Output:     "stdout",
	})
	if err != nil {
		log.Fatalf("❌ Error al crear logger: %v", err)
	}
	logger.Info("Logger inicializado correctamente")

	// ============ DOMINIO ============
	// Crear servicio del restaurante (núcleo de la aplicación)
	fmt.Println("🍽️  Inicializando servicio del restaurante...")
	restaurantService := service.NewRestaurantService(
		capacidadBarra,
		numCocineros,
		numMesas,
	)

	// Iniciar las goroutines concurrentes
	// - Cocineros (productores automáticos)
	// - Generador de clientes
	// - Verificador de paciencia
	restaurantService.Start()
	logger.Info("Sistema de concurrencia iniciado")

	// ============ INTERFAZ GRÁFICA ============
	// Crear el juego con Ebiten
	fmt.Println("🎮 Inicializando interfaz gráfica...")
	game, err := ui.NewGame(restaurantService, screenWidth, screenHeight)
	if err != nil {
		log.Fatalf("❌ Error al crear el juego: %v", err)
	}
	logger.Info("Interfaz gráfica inicializada")

	fmt.Println()
	fmt.Println("✅ Sistema listo")
	fmt.Println()
	fmt.Println("CONTROLES:")
	fmt.Println("  [W/↑]         Mover mesero arriba")
	fmt.Println("  [S/↓]         Mover mesero abajo")
	fmt.Println("  [A/←]         Mover mesero izquierda")
	fmt.Println("  [D/→]         Mover mesero derecha")
	fmt.Println("  [E]           Recoger plato de la barra")
	fmt.Println("  [ESPACIO]     Entregar plato a mesa cercana")
	fmt.Println("  [ESPACIO]     Pausar/Reanudar (en barra)")
	fmt.Println("  [Q/ESC]       Salir")
	fmt.Println()
	fmt.Println("OBJETIVO:")
	fmt.Println("  1. Recoger platos de la BARRA (cocinero produce)")
	fmt.Println("  2. Entregar platos a las MESAS (clientes esperan)")
	fmt.Println("  3. Evitar que clientes pierdan paciencia")
	fmt.Println()
	fmt.Println()

	// ============ CONFIGURACIÓN DE EBITEN ============
	// Configurar ventana
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("🍽️ Restaurante Concurrente - Arquitectura Hexagonal")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(60) // 60 FPS

	// ============ EJECUTAR JUEGO ============
	// Ejecutar el loop del juego (bloquea hasta que se cierre)
	if err := ebiten.RunGame(game); err != nil {
		if err.Error() != "cierre solicitado por usuario" {
			log.Println("❌ Error durante la ejecución:", err)
		}
	}

	// ============ CIERRE ORDENADO ============
	fmt.Println()
	fmt.Println("🚪 Cerrando sistema...")
	restaurantService.Close()
	logger.Info("Sistema cerrado correctamente")

	fmt.Println("✅ Aplicación finalizada exitosamente")
	fmt.Println()
}

// Función auxiliar para repetir strings (si Go < 1.20)
func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
