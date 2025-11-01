package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// InputHandler maneja toda la entrada del usuario de manera centralizada
type InputHandler struct {
	// Estado de teclas (para detectar cambios)
	lastFrameKeys map[ebiten.Key]bool

	// Callbacks para acciones
	onPausar         func()
	onAgregarCliente func()
	onRemoverCliente func()
	onSalir          func()
	onReset          func()

	// Configuración
	enabled bool
}

// InputAction representa una acción que puede realizar el usuario
type InputAction int

const (
	ActionPausar InputAction = iota
	ActionAgregarCliente
	ActionRemoverCliente
	ActionSalir
	ActionReset
	ActionNone
)

// NewInputHandler crea un nuevo manejador de entrada
func NewInputHandler() *InputHandler {
	return &InputHandler{
		lastFrameKeys: make(map[ebiten.Key]bool),
		enabled:       true,
	}
}

// SetCallbacks configura las funciones callback para cada acción
func (h *InputHandler) SetCallbacks(
	onPausar func(),
	onAgregarCliente func(),
	onRemoverCliente func(),
	onSalir func(),
	onReset func(),
) {
	h.onPausar = onPausar
	h.onAgregarCliente = onAgregarCliente
	h.onRemoverCliente = onRemoverCliente
	h.onSalir = onSalir
	h.onReset = onReset
}

// Enable activa el manejador de entrada
func (h *InputHandler) Enable() {
	h.enabled = true
}

// Disable desactiva el manejador de entrada
func (h *InputHandler) Disable() {
	h.enabled = false
}

// IsEnabled retorna si el manejador está activo
func (h *InputHandler) IsEnabled() bool {
	return h.enabled
}

// Update procesa la entrada del usuario (llamar en cada frame)
func (h *InputHandler) Update() error {
	if !h.enabled {
		return nil
	}

	// Detectar teclas presionadas (solo al momento de presionar)
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if h.onPausar != nil {
			h.onPausar()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		if h.onAgregarCliente != nil {
			h.onAgregarCliente()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		if h.onRemoverCliente != nil {
			h.onRemoverCliente()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if h.onSalir != nil {
			h.onSalir()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		if h.onReset != nil {
			h.onReset()
		}
	}

	// Soporte para múltiples clientes a la vez
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			// Agregar 5 clientes con Shift+A
			for i := 0; i < 5; i++ {
				if h.onAgregarCliente != nil {
					h.onAgregarCliente()
				}
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			// Remover 5 clientes con Shift+R
			for i := 0; i < 5; i++ {
				if h.onRemoverCliente != nil {
					h.onRemoverCliente()
				}
			}
		}
	}

	return nil
}

// IsKeyPressed verifica si una tecla está presionada actualmente
func (h *InputHandler) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}

// IsKeyJustPressed verifica si una tecla acaba de ser presionada
func (h *InputHandler) IsKeyJustPressed(key ebiten.Key) bool {
	return inpututil.IsKeyJustPressed(key)
}

// IsKeyJustReleased verifica si una tecla acaba de ser soltada
func (h *InputHandler) IsKeyJustReleased(key ebiten.Key) bool {
	wasPressed := h.lastFrameKeys[key]
	isPressed := ebiten.IsKeyPressed(key)

	h.lastFrameKeys[key] = isPressed

	return wasPressed && !isPressed
}

// GetPressedAction retorna la acción correspondiente a la tecla presionada
func (h *InputHandler) GetPressedAction() InputAction {
	if !h.enabled {
		return ActionNone
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return ActionPausar
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		return ActionAgregarCliente
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		return ActionRemoverCliente
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ActionSalir
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		return ActionReset
	}

	return ActionNone
}

// GetHelpText retorna el texto de ayuda con todos los controles
func (h *InputHandler) GetHelpText() []string {
	return []string{
		"CONTROLES:",
		"[SPACE]     Pausar/Reanudar producción",
		"[A]         Agregar 1 cliente",
		"[Shift+A]   Agregar 5 clientes",
		"[R]         Remover 1 cliente",
		"[Shift+R]   Remover 5 clientes",
		"[F5]        Reiniciar restaurante",
		"[Q/ESC]     Salir",
	}
}

// MousePosition retorna la posición actual del mouse
func (h *InputHandler) MousePosition() (int, int) {
	return ebiten.CursorPosition()
}

// IsMouseButtonPressed verifica si un botón del mouse está presionado
func (h *InputHandler) IsMouseButtonPressed(button ebiten.MouseButton) bool {
	return ebiten.IsMouseButtonPressed(button)
}

// IsMouseButtonJustPressed verifica si un botón del mouse acaba de ser presionado
func (h *InputHandler) IsMouseButtonJustPressed(button ebiten.MouseButton) bool {
	return inpututil.IsMouseButtonJustPressed(button)
}
