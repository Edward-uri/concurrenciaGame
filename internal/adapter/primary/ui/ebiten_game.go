package ui

import (
	"fmt"
	"image/color"
	"restaurant-concurrency/internal/domain/port"
	infrastructure "restaurant-concurrency/internal/infraestructure"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// EbitenGame adapta el servicio a la interfaz de Ebiten
type EbitenGame struct {
	service      port.RestaurantService
	assets       *Assets
	inputHandler *InputHandler
	logger       *infrastructure.Logger
	width        int
	height       int
	shouldClose  bool

	// Mesero controlado por el jugador
	meseroX       float64
	meseroY       float64
	meseroVel     float64
	carryingPlato bool
	platoID       int

	// Control de tecla E para evitar múltiples pulsaciones
	ePressedLastFrame bool

	// Posiciones de zonas
	barraX, barraY     float64
	clienteX, clienteY float64
}

func NewEbitenGame(
	service port.RestaurantService,
	logger *infrastructure.Logger,
	width, height int,
) (*EbitenGame, error) {
	assets, err := LoadAssets()
	if err != nil {
		return nil, fmt.Errorf("error cargando assets: %w", err)
	}

	game := &EbitenGame{
		service:      service,
		assets:       assets,
		inputHandler: NewInputHandler(),
		logger:       logger,
		width:        width,
		height:       height,
		shouldClose:  false,

		// Posición inicial del mesero (centro-derecha)
		meseroX:       float64(width) - 200,
		meseroY:       float64(height) / 2,
		meseroVel:     5.0,
		carryingPlato: false,

		// Zona de la barra (centro)
		barraX: float64(width)/2 - 50,
		barraY: float64(height) / 2,

		// Zona de clientes (arriba a la derecha)
		clienteX: float64(width) - 150,
		clienteY: 100,
	}

	// Configurar callbacks del input handler
	game.setupInputCallbacks()

	return game, nil
}

// setupInputCallbacks configura las acciones del input handler
func (g *EbitenGame) setupInputCallbacks() {
	g.inputHandler.SetCallbacks(
		// onPausar
		func() {
			g.service.TogglePausar()
			estado := g.service.GetEstado()
			if estado.Pausado {
				g.logger.Info("Producción pausada por usuario")
			} else {
				g.logger.Info("Producción reanudada por usuario")
			}
		},
		// onAgregarCliente
		func() {
			g.service.AgregarClientes(1)
			g.logger.Infof("Cliente agregado por usuario")
		},
		// onRemoverCliente
		func() {
			g.service.ClientesSeVan(1)
			g.logger.Infof("Cliente removido por usuario")
		},
		// onSalir
		func() {
			g.shouldClose = true
			g.logger.Info("Cierre solicitado por usuario")
		},
		// onReset
		func() {
			g.logger.Info("Reset solicitado - funcionalidad pendiente")
			// Aquí podrías reiniciar el estado del restaurante
		},
	)
}

// Update maneja la lógica del juego (60 FPS)
func (g *EbitenGame) Update() error {
	// Procesar entrada
	if err := g.inputHandler.Update(); err != nil {
		return err
	}

	// Verificar si se solicitó cerrar
	if g.shouldClose {
		return fmt.Errorf("cierre solicitado por usuario")
	}

	// Controles del mesero con WASD
	g.moverMesero()

	// Acción con E (interactuar) - solo cuando se presiona, no cuando se mantiene
	ePressed := ebiten.IsKeyPressed(ebiten.KeyE)
	if ePressed && !g.ePressedLastFrame {
		g.accionMesero()
	}
	g.ePressedLastFrame = ePressed

	return nil
}

// moverMesero mueve al mesero con las teclas WASD
func (g *EbitenGame) moverMesero() {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.meseroY -= g.meseroVel
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.meseroY += g.meseroVel
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.meseroX -= g.meseroVel
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.meseroX += g.meseroVel
	}

	// Limitar al área de la pantalla
	if g.meseroX < 0 {
		g.meseroX = 0
	}
	if g.meseroX > float64(g.width)-100 {
		g.meseroX = float64(g.width) - 100
	}
	if g.meseroY < 0 {
		g.meseroY = 0
	}
	if g.meseroY > float64(g.height)-100 {
		g.meseroY = float64(g.height) - 100
	}
}

// accionMesero realiza la acción de recoger/entregar plato
func (g *EbitenGame) accionMesero() {
	estado := g.service.GetEstado()

	// Si está cerca de la barra y hay platos
	distBarra := distancia(g.meseroX, g.meseroY, g.barraX, g.barraY)
	if distBarra < 10000 && !g.carryingPlato && estado.EnBarra > 0 {
		// Consumir plato del canal (patrón Productor-Consumidor)
		plato := g.service.ConsumirPlato()
		if plato != nil {
			g.carryingPlato = true
			g.platoID = plato.ID
			g.logger.Infof("🍽️ Mesero tomó plato #%d de la barra (Buffer: %d/%d)", plato.ID, estado.EnBarra-1, estado.CapacidadBarra)
		} else {
			g.logger.Info("⚠️ No hay platos disponibles en la barra")
		}
		return
	}

	// Si está cerca del cliente y tiene plato
	distCliente := distancia(g.meseroX, g.meseroY, g.clienteX, g.clienteY)
	if distCliente < 10000 && g.carryingPlato && estado.ClientesActivos > 0 {
		// Entregar plato al cliente (SOLO AQUÍ se incrementa el contador de servidos)
		g.carryingPlato = false
		g.service.EntregarPlato()
		g.logger.Infof("✅ Plato #%d entregado al cliente! (Servidos: %d)", g.platoID, estado.PlatosServidos+1)
		return
	}

	// Mensajes si no está en posición correcta
	if !g.carryingPlato {
		g.logger.Info("⚠️ Acércate a la BARRA para tomar un plato")
	} else {
		g.logger.Info("⚠️ Acércate a los CLIENTES para entregar el plato")
	}
}

// distancia calcula la distancia euclidiana entre dos puntos
func distancia(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return dx*dx + dy*dy // Comparamos el cuadrado para evitar sqrt
}

// Draw renderiza la interfaz gráfica
func (g *EbitenGame) Draw(screen *ebiten.Image) {
	// Fondo con el piso
	screen.Fill(color.RGBA{60, 45, 35, 255})
	g.dibujarPiso(screen)

	estado := g.service.GetEstado()

	// Panel superior con información
	g.dibujarInfoSuperior(screen, estado)

	// === LAYOUT MEJORADO ===
	// Cocina (Izquierda): Chef produciendo
	cocinerosX := 50.0
	cocinerosY := 300.0

	// Barra (Centro): Buffer de platos
	barraX := float64(g.width)/2 - 50
	barraY := 300.0
	g.barraX = barraX
	g.barraY = barraY

	// Clientes (Derecha-Arriba): Esperando servicio
	clientesX := float64(g.width) - 200
	clientesY := 150.0
	g.clienteX = clientesX
	g.clienteY = clientesY

	// Etiquetas de áreas con explicación
	ebitenutil.DebugPrintAt(screen, "🍳 COCINA (PRODUCTOR)", int(cocinerosX), int(cocinerosY)-30)
	ebitenutil.DebugPrintAt(screen, "Chef cocinando...", int(cocinerosX), int(cocinerosY)-15)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("📦 BARRA (BUFFER %d/%d)", estado.EnBarra, estado.CapacidadBarra), int(barraX)-20, int(barraY)-30)
	ebitenutil.DebugPrintAt(screen, "Presiona E cerca para tomar", int(barraX)-60, int(barraY)-15)

	if estado.ClientesActivos > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("👥 CLIENTES (%d)", estado.ClientesActivos), int(clientesX)-20, int(clientesY)-30)
		ebitenutil.DebugPrintAt(screen, "Presiona E cerca para entregar", int(clientesX)-70, int(clientesY)-15)
	}

	// Dibujar elementos del restaurante
	g.dibujarCocineros(screen, cocinerosX, cocinerosY)
	g.dibujarBarra(screen, estado.EnBarra, estado.CapacidadBarra, barraX, barraY)
	g.dibujarClientes(screen, estado.ClientesActivos, clientesX, clientesY)
	g.dibujarMeseros(screen, 0, 0) // El mesero se dibuja en su posición actual

	// Controles en la parte inferior
	g.dibujarControlesSimple(screen)

	// Instrucciones del juego
	ebitenutil.DebugPrintAt(screen, "🎮 Usa W/A/S/D para mover al mesero", 20, g.height-80)
	ebitenutil.DebugPrintAt(screen, "🔄 Presiona E cerca de barra/cliente para tomar/entregar platos", 20, g.height-60)
}

// dibujarInfoSuperior dibuja información en la parte superior (simple)
func (g *EbitenGame) dibujarInfoSuperior(screen *ebiten.Image, estado interface{}) {
	st := g.service.GetEstado()

	// Título
	ebitenutil.DebugPrintAt(screen, "RESTAURANTE CONCURRENTE - Patron Productor-Consumidor", 20, 15)

	// Estadísticas en una sola línea
	status := "ACTIVO"
	if st.Pausado {
		status = "PAUSADO"
	}

	info := fmt.Sprintf("Estado: %s | Clientes: %d | Barra: %d/%d | Producidos: %d | Servidos: %d",
		status, st.ClientesActivos, st.EnBarra, st.CapacidadBarra, st.PlatosTotales, st.PlatosServidos)
	ebitenutil.DebugPrintAt(screen, info, 20, 35)
}

// dibujarControlesSimple dibuja los controles de manera simple
func (g *EbitenGame) dibujarControlesSimple(screen *ebiten.Image) {
	y := g.height - 110
	ebitenutil.DebugPrintAt(screen, "CONTROLES GENERALES:", 20, y)
	ebitenutil.DebugPrintAt(screen, "[P] Pausar producción", 20, y+15)
	ebitenutil.DebugPrintAt(screen, "[+] Agregar cliente", 200, y+15)
	ebitenutil.DebugPrintAt(screen, "[-] Quitar cliente", 360, y+15)
	ebitenutil.DebugPrintAt(screen, "[R] Reiniciar", 500, y+15)
	ebitenutil.DebugPrintAt(screen, "[ESC] Salir", 620, y+15)
}

// dibujarPiso dibuja el piso como fondo repetido
func (g *EbitenGame) dibujarPiso(screen *ebiten.Image) {
	if g.assets.Piso == nil {
		return
	}

	pisoW := g.assets.Piso.Bounds().Dx()
	pisoH := g.assets.Piso.Bounds().Dy()

	// Repetir el piso en toda la pantalla
	for y := 0; y < g.height; y += pisoH {
		for x := 0; x < g.width; x += pisoW {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			// Hacer el piso semi-transparente
			op.ColorScale.ScaleAlpha(0.3)
			screen.DrawImage(g.assets.Piso, op)
		}
	}
}

func (g *EbitenGame) dibujarCocineros(screen *ebiten.Image, x, y float64) {
	if g.assets.Cocinero == nil {
		ebitenutil.DebugPrintAt(screen, "[COCINERO]", int(x), int(y))
		return
	}

	// Dibujar 2 cocineros MUY GRANDES (sprites 32x32 -> 128x128)
	scale := 4.0
	for i := 0; i < 2; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(x+float64(i*140), y)
		screen.DrawImage(g.assets.Cocinero, op)

		// Etiqueta debajo
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Chef %d", i+1), int(x)+int(float64(i)*140), int(y)+135)
	}
}

func (g *EbitenGame) dibujarBarra(screen *ebiten.Image, ocupado, capacidad int, x, y float64) {
	if g.assets.Barra == nil || g.assets.Plato == nil {
		ebitenutil.DebugPrintAt(screen, "[BARRA]", int(x), int(y))
		return
	}

	// BARRA MUY GRANDE (32x32 -> 96x96)
	spacing := 100.0
	scale := 3.0

	for i := 0; i < capacidad; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(x+float64(i)*spacing, y)

		if i < ocupado {
			// Dibujar plato
			screen.DrawImage(g.assets.Plato, op)
		} else {
			// Dibujar espacio vacío de la barra (más transparente)
			op.ColorScale.ScaleAlpha(0.3)
			screen.DrawImage(g.assets.Barra, op)
		}
	}
}

func (g *EbitenGame) dibujarMeseros(screen *ebiten.Image, x, y float64) {
	if g.assets.Mesero == nil {
		ebitenutil.DebugPrintAt(screen, "[MESERO]", int(x), int(y))
		return
	}

	// Dibujar mesero CONTROLADO POR EL JUGADOR
	scale := 4.0
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(g.meseroX, g.meseroY)
	screen.DrawImage(g.assets.Mesero, op)

	// Si lleva un plato, dibujarlo
	if g.carryingPlato && g.assets.Plato != nil {
		platoOp := &ebiten.DrawImageOptions{}
		platoOp.GeoM.Scale(2.5, 2.5)
		platoOp.GeoM.Translate(g.meseroX+50, g.meseroY-30)
		screen.DrawImage(g.assets.Plato, platoOp)
	}

	// Etiqueta con estado del mesero
	estadoTexto := "Listo"
	if g.carryingPlato {
		estadoTexto = "Llevando plato"
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mesero: %s", estadoTexto), int(g.meseroX)-10, int(g.meseroY)+135)
}

func (g *EbitenGame) dibujarClientes(screen *ebiten.Image, cantidad int, x, y float64) {
	if cantidad == 0 {
		ebitenutil.DebugPrintAt(screen, "[ Sin clientes - Presiona + para agregar ]", int(x), int(y)+20)
		return
	}

	// CLIENTES MUY GRANDES (32x32 -> 112x112)
	scale := 3.5
	spacing := 130.0

	// Dibujar mesas primero (si existen)
	if g.assets.Mesa != nil {
		for i := 0; i < cantidad; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scale*0.9, scale*0.9)
			op.GeoM.Translate(x+float64(i)*spacing, y+40)
			op.ColorScale.ScaleAlpha(0.5)
			screen.DrawImage(g.assets.Mesa, op)
		}
	}

	// Dibujar clientes sobre las mesas
	for i := 0; i < cantidad; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(x+float64(i)*spacing+10, y)

		// Protección contra división por cero
		skinIndex := 0
		if g.assets.NumClientes > 0 {
			skinIndex = i % g.assets.NumClientes
		}

		clienteSprite := g.assets.GetClienteSprite(skinIndex)

		if clienteSprite != nil {
			screen.DrawImage(clienteSprite, op)
		} else {
			// Si no hay sprite, mostrar texto
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("C%d", i+1), int(x)+int(float64(i)*spacing), int(y))
		}
	}
}

func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}
