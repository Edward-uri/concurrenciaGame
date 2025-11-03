package ui

import (
	"fmt"
	"image/color"
	"restaurant-concurrency/internal/domain/model"
	"restaurant-concurrency/internal/domain/service"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	service       *service.RestaurantService
	mesero        *model.Mesero
	inputHandler  *InputHandler
	renderer      *Renderer
	width, height int

	// Control de teclas para evitar repetición
	ePressedLastFrame     bool
	spacePressedLastFrame bool

	// Notificación temporal
	notificacion       string
	notificacionFrames int
}

func NewGame(service *service.RestaurantService, width, height int) (*Game, error) {
	renderer, err := NewRenderer()
	if err != nil {
		return nil, err
	}

	game := &Game{
		service:      service,
		mesero:       model.NewMesero(400, 200, 200), // Posición inicial
		inputHandler: NewInputHandler(),
		renderer:     renderer,
		width:        width,
		height:       height,
	}

	game.setupCallbacks()
	return game, nil
}

func (g *Game) setupCallbacks() {
	// Pasar métodos directamente en lugar de funciones anónimas
	g.inputHandler.SetCallbacks(
		g.service.TogglePausar, // Pausar - método directo
		nil,                    // Ya no agregamos clientes manualmente
		nil,                    // Ya no removemos clientes manualmente
		g.handleClose,          // Cerrar - método helper
		nil,
	)
}

// handleClose maneja el cierre del juego
func (g *Game) handleClose() {
	// Aquí se puede agregar lógica de cierre si es necesaria
}

func (g *Game) Update() error {
	// Procesar input
	g.inputHandler.Update()

	// Movimiento del mesero (WASD o flechas)
	dx, dy := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		dy = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		dy = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		dx = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		dx = 1
	}

	// Normalizar movimiento diagonal
	if dx != 0 && dy != 0 {
		factor := 0.707 // sqrt(2)/2
		dx *= factor
		dy *= factor
	}

	g.mesero.Mover(dx, dy, 1.0/60.0)

	// Recoger plato de la barra con E
	if g.inputHandler.IsKeyJustPressed(ebiten.KeyE) && !g.mesero.TienePlato {
		if g.meseroEnBarra() {
			if plato, ok := g.service.IntentarRecogerPlato(); ok {
				g.mesero.RecogerPlato(*plato)
				g.mostrarNotificacion(fmt.Sprintf("Plato #%d recogido", plato.ID))
			} else {
				g.mostrarNotificacion("No hay platos en la barra")
			}
		} else {
			g.mostrarNotificacion("Acercate a la barra (zona central superior)")
		}
	}

	// Entregar plato con ESPACIO
	if g.inputHandler.IsKeyJustPressed(ebiten.KeySpace) && g.mesero.TienePlato {
		if g.service.EntregarPlatoAMesa(g.mesero.PosX, g.mesero.PosY, 100) {
			delivered := g.mesero.EntregarPlato()
			if delivered != nil {
				g.mostrarNotificacion(fmt.Sprintf("Plato #%d entregado a la mesa", delivered.ID))
			} else {
				g.mostrarNotificacion("Plato entregado a la mesa")
			}
		} else {
			g.mostrarNotificacion("Acercate a una mesa con clientes")
		}
	}

	// Decrementar contador de notificación
	if g.notificacionFrames > 0 {
		g.notificacionFrames--
	}

	return nil
}

func (g *Game) meseroEnBarra() bool {
	// Verificar si el mesero está cerca de la barra (centro superior, coincidente con Renderer)
	barraX := float64(g.width/2 - 200)
	barraY := 80.0
	dx := g.mesero.PosX - barraX
	dy := g.mesero.PosY - barraY
	return dx*dx+dy*dy < 150*150 // Radio más grande para facilitar la interacción
}

func (g *Game) mostrarNotificacion(mensaje string) {
	g.notificacion = mensaje
	g.notificacionFrames = 120 // 2 segundos a 60 FPS
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Dibujar piso repetid
	g.renderer.DibujarPiso(screen, g.width, g.height)

	// Dibujar cocinero en la cocina (arriba a la izquierda)
	g.renderer.DibujarCocinero(screen, 50, 50)

	// Dibujar barra
	estadoBarra := g.service.GetEstadoBarra()
	capacidadBarra := g.service.GetCapacidadBarra()
	g.renderer.DibujarBarra(screen, float32(g.width/2-200), 80, estadoBarra, capacidadBarra)

	// Dibujar mesas con clientes (zona inferior)
	mesas := g.service.GetMesas()
	for _, mesa := range mesas {
		g.renderer.DibujarMesa(screen, mesa)
	}

	// Dibujar mesero (controlado por el jugador)
	g.renderer.DibujarMesero(screen, g.mesero)

	// Dibujar UI e información
	g.dibujarUI(screen)
}

func (g *Game) dibujarUI(screen *ebiten.Image) {
	totales, servidos, perdidos := g.service.GetMetricas()
	estadoBarra := g.service.GetEstadoBarra()
	capacidadBarra := g.service.GetCapacidadBarra()

	// PANEL IZQUIERDO - Información y controles
	panelX := 10
	y := 15

	// Título
	ebitenutil.DebugPrintAt(screen, "RESTAURANTE CONCURRENTE", panelX, y)
	y += 20
	ebitenutil.DebugPrintAt(screen, "Patron Productor-Consumidor", panelX, y)
	y += 30

	// Métricas del sistema
	ebitenutil.DebugPrintAt(screen, "===========================", panelX, y)
	y += 20
	ebitenutil.DebugPrintAt(screen, "METRICAS", panelX, y)
	y += 20
	ebitenutil.DebugPrintAt(screen, "===========================", panelX, y)
	y += 20
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Buffer: %d/%d", estadoBarra, capacidadBarra), panelX, y)
	y += 18
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Producidos: %d", totales), panelX, y)
	y += 18
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Servidos: %d", servidos), panelX, y)
	y += 18
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Perdidos: %d", perdidos), panelX, y)
	y += 30

	// Controles
	ebitenutil.DebugPrintAt(screen, "===========================", panelX, y)
	y += 20
	ebitenutil.DebugPrintAt(screen, "CONTROLES", panelX, y)
	y += 20
	ebitenutil.DebugPrintAt(screen, "===========================", panelX, y)
	y += 20
	ebitenutil.DebugPrintAt(screen, "[WASD] Mover mesero", panelX, y)
	y += 18
	ebitenutil.DebugPrintAt(screen, "[E] Recoger plato", panelX, y)
	y += 18
	ebitenutil.DebugPrintAt(screen, "[ESPACIO] Entregar", panelX, y)
	y += 30

	// Estado del mesero
	ebitenutil.DebugPrintAt(screen, "===========================", panelX, y)
	y += 20
	ebitenutil.DebugPrintAt(screen, "ESTADO DEL MESERO", panelX, y)
	y += 20
	ebitenutil.DebugPrintAt(screen, "===========================", panelX, y)
	y += 20
	if g.mesero.TienePlato {
		ebitenutil.DebugPrintAt(screen, "Con plato", panelX, y)
		y += 18
		ebitenutil.DebugPrintAt(screen, "-> Busca mesa", panelX, y)
	} else {
		ebitenutil.DebugPrintAt(screen, "Libre", panelX, y)
		y += 18
		ebitenutil.DebugPrintAt(screen, "-> Ve a barra", panelX, y)
	}

	// NOTIFICACIÓN CENTRAL (si existe)
	if g.notificacionFrames > 0 {
		// Calcular posición central
		notifX := g.width/2 - 200
		notifY := g.height - 150

		// Fondo semi-transparente para la notificación
		vector.DrawFilledRect(screen, float32(notifX-10), float32(notifY-10),
			400, 60, color.RGBA{0, 0, 0, 200}, false)

		// Borde de la notificación
		vector.StrokeRect(screen, float32(notifX-10), float32(notifY-10),
			400, 60, 3, color.RGBA{255, 255, 0, 255}, false)

		// Texto de la notificación (grande y centrado)
		ebitenutil.DebugPrintAt(screen, g.notificacion, notifX, notifY+10)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return g.width, g.height
}
