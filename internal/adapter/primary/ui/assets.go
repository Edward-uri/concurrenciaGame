package ui

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var assetsFS embed.FS

// Assets gestiona todos los recursos gráficos
type Assets struct {
	// Sprites individuales
	Cocinero *ebiten.Image
	Mesero   *ebiten.Image
	Plato    *ebiten.Image
	Mesa     *ebiten.Image
	Barra    *ebiten.Image
	Piso     *ebiten.Image

	// Spritesheet de clientes (múltiples skins)
	ClientesSheet *ebiten.Image
	ClienteFrames []*ebiten.Image

	// Configuración del spritesheet
	ClienteWidth  int
	ClienteHeight int
	NumClientes   int
}

// LoadAssets carga todos los assets desde el sistema de archivos embebido
func LoadAssets() (*Assets, error) {
	assets := &Assets{
		ClienteWidth:  32, // Sprites de 32x32
		ClienteHeight: 32, // Sprites de 32x32
		NumClientes:   1,  // Al menos 1 para evitar división por cero
	}

	var err error

	// Cargar sprites individuales
	assets.Cocinero, err = loadImage("assets/cocinero.png")
	if err != nil {
		return nil, err
	}

	assets.Mesero, err = loadImage("assets/mesero.png")
	if err != nil {
		return nil, err
	}

	assets.Plato, err = loadImage("assets/plato.png")
	if err != nil {
		return nil, err
	}

	assets.Mesa, err = loadImage("assets/mesa.png")
	if err != nil {
		return nil, err
	}

	assets.Barra, err = loadImage("assets/barra.png")
	if err != nil {
		return nil, err
	}

	assets.Piso, err = loadImage("assets/piso.png")

	// Cargar spritesheet de clientes
	assets.ClientesSheet, err = loadImage("assets/clientes.png")
	if err != nil {
		return nil, err
	}

	// Extraer frames individuales del spritesheet
	assets.ClienteFrames, assets.NumClientes = extractSpritesheetFrames(
		assets.ClientesSheet,
		assets.ClienteWidth,
		assets.ClienteHeight,
	)

	return assets, nil
}

// loadImage carga una imagen desde el FS embebido
func loadImage(path string) (*ebiten.Image, error) {
	data, err := assetsFS.ReadFile(path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(img), nil
}

// extractSpritesheetFrames extrae frames individuales de un spritesheet horizontal
func extractSpritesheetFrames(sheet *ebiten.Image, frameWidth, frameHeight int) ([]*ebiten.Image, int) {
	if sheet == nil || frameWidth <= 0 || frameHeight <= 0 {
		// Retornar al menos un frame vacío para evitar errores
		return []*ebiten.Image{sheet}, 1
	}

	bounds := sheet.Bounds()
	sheetWidth := bounds.Dx()
	sheetHeight := bounds.Dy()

	// Si el frame es más grande que el sheet, usar el sheet completo
	if frameWidth > sheetWidth || frameHeight > sheetHeight {
		return []*ebiten.Image{sheet}, 1
	}

	// Calcular número de frames
	numFrames := sheetWidth / frameWidth

	if numFrames <= 0 {
		numFrames = 1
	}

	frames := make([]*ebiten.Image, numFrames)

	for i := 0; i < numFrames; i++ {
		x := i * frameWidth

		// Crear subimagen para este frame
		frame := sheet.SubImage(image.Rect(
			x, 0,
			x+frameWidth, frameHeight,
		)).(*ebiten.Image)

		frames[i] = frame
	}

	return frames, numFrames
}

// GetClienteSprite retorna el sprite de un cliente específico por índice
func (a *Assets) GetClienteSprite(index int) *ebiten.Image {
	if len(a.ClienteFrames) == 0 {
		return a.ClientesSheet // Retornar el sheet completo si no hay frames
	}

	if index < 0 || index >= len(a.ClienteFrames) {
		index = 0 // Default al primero si está fuera de rango
	}

	return a.ClienteFrames[index]
}
