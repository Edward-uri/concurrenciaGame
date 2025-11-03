package ui

import (
	"fmt"
	"image/color"
	"math"
	"restaurant-concurrency/internal/domain/model"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Renderer struct {
	assets *Assets
}

func NewRenderer() (*Renderer, error) {
	assets, err := LoadAssets()
	if err != nil {
		return nil, err
	}

	return &Renderer{
		assets: assets,
	}, nil
}

// DibujarPiso dibuja el piso repetido (tiled) en toda la pantalla
func (r *Renderer) DibujarPiso(screen *ebiten.Image, screenWidth, screenHeight int) {
	if r.assets.Piso == nil {
		// Fallback: color sólido si no hay imagen
		screen.Fill(color.RGBA{45, 35, 30, 255})
		return
	}

	// Obtener dimensiones del tile de piso
	tileWidth := r.assets.Piso.Bounds().Dx()
	tileHeight := r.assets.Piso.Bounds().Dy()

	// Dibujar tiles repetidos para cubrir toda la pantalla
	for y := 0; y < screenHeight; y += tileHeight {
		for x := 0; x < screenWidth; x += tileWidth {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(r.assets.Piso, op)
		}
	}
}

func (r *Renderer) DibujarMesa(screen *ebiten.Image, mesa model.MesaSnapshot) {
	x, y := float32(mesa.PosX), float32(mesa.PosY)

	// Dibujar sprite de la mesa
	if r.assets.Mesa != nil {
		op := &ebiten.DrawImageOptions{}
		scale := 2.0 // Escalar sprite 32x32 a 64x64
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(r.assets.Mesa, op)
	} else {
		// Fallback: rectángulo si no hay imagen
		colorMesa := color.RGBA{139, 69, 19, 255}
		vector.DrawFilledRect(screen, x, y, 64, 64, colorMesa, false)
		vector.StrokeRect(screen, x, y, 64, 64, 2, color.White, false)
	}

	if mesa.ClientesActivos > 0 {
		// Dibujar sprite de clientes
		if r.assets.ClienteFrames != nil && len(r.assets.ClienteFrames) > 0 {
			clienteSprite := r.assets.ClienteFrames[0] // Usar primer sprite de cliente
			op := &ebiten.DrawImageOptions{}
			scale := 1.5
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(float64(x-20), float64(y-30)) // Posicionar arriba de la mesa
			screen.DrawImage(clienteSprite, op)

			// Número de clientes
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("x%d", mesa.ClientesActivos),
				int(x-15), int(y-10))
		} else {
			// Fallback: texto simple
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Clientes: %d", mesa.ClientesActivos),
				int(x+5), int(y+5))
		}

		// Barra de paciencia
		paciencia := mesa.NivelPaciencia
		barWidth := float32(60)
		barHeight := float32(6)

		// Fondo de la barra
		vector.DrawFilledRect(screen, x, y+70, barWidth, barHeight,
			color.RGBA{50, 50, 50, 255}, false)

		// Barra de progreso (verde = paciente, rojo = impaciente)
		barraColor := interpolarColor(
			color.RGBA{0, 255, 0, 255}, // Verde
			color.RGBA{255, 0, 0, 255}, // Rojo
			paciencia,
		)
		vector.DrawFilledRect(screen, x, y+70, barWidth*(1-float32(paciencia)),
			barHeight, barraColor, false)

		// Plato entregado
		if mesa.TienePlato && r.assets.Plato != nil {
			op := &ebiten.DrawImageOptions{}
			scale := 1.2
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(float64(x+20), float64(y+20))
			screen.DrawImage(r.assets.Plato, op)
		}
	}
}

func (r *Renderer) DibujarBarra(screen *ebiten.Image, x, y float32, ocupado, capacidad int) {
	// Título de la barra - Buffer del patrón Productor-Consumidor
	ebitenutil.DebugPrintAt(screen, "BARRA", int(x-50), int(y-30))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Platos disponibles: %d/%d", ocupado, capacidad), int(x-50), int(y-15))

	slotWidth := float32(60)
	spacing := float32(15)

	for i := 0; i < capacidad; i++ {
		posX := x + float32(i)*(slotWidth+spacing)

		// Dibujar sprite de barra como fondo
		if r.assets.Barra != nil {
			op := &ebiten.DrawImageOptions{}
			scale := 2.0
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(float64(posX), float64(y))

			// Cambiar tono de color si está vacío (más oscuro)
			if i >= ocupado {
				op.ColorScale.Scale(0.6, 0.6, 0.6, 1.0) // Más oscuro pero no transparente
			}
			screen.DrawImage(r.assets.Barra, op)
		} else {
			// Fallback: rectángulo
			var col color.Color
			if i < ocupado {
				col = color.RGBA{255, 200, 0, 255}
			} else {
				col = color.RGBA{60, 60, 70, 255}
			}
			vector.DrawFilledRect(screen, posX, y, slotWidth, 50, col, false)
			vector.StrokeRect(screen, posX, y, slotWidth, 50, 2, color.White, false)
		}

		// Dibujar sprite de plato si el slot está ocupado
		if i < ocupado && r.assets.Plato != nil {
			op := &ebiten.DrawImageOptions{}
			scale := 1.8
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(float64(posX+10), float64(y+10))
			screen.DrawImage(r.assets.Plato, op)
		}
	}
}

func (r *Renderer) DibujarMesero(screen *ebiten.Image, mesero *model.Mesero) {
	x, y := float32(mesero.PosX), float32(mesero.PosY)

	// Dibujar sprite del mesero
	if r.assets.Mesero != nil {
		op := &ebiten.DrawImageOptions{}
		scale := 3.0 // Hacer más grande para que sea visible
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(x-48), float64(y-48)) // Centrar (32*3 = 96, 96/2 = 48)

		// Cambiar tono si tiene plato
		if mesero.TienePlato {
			op.ColorScale.ScaleWithColor(color.RGBA{255, 220, 180, 255}) // Tono dorado
		}

		screen.DrawImage(r.assets.Mesero, op)
	} else {
		// Fallback: círculo
		col := color.RGBA{100, 150, 255, 255}
		if mesero.TienePlato {
			col = color.RGBA{255, 150, 100, 255}
		}
		vector.DrawFilledCircle(screen, x, y, 16, col, false)
		vector.StrokeCircle(screen, x, y, 16, 2, color.White, false)
		ebitenutil.DebugPrintAt(screen, "M", int(x-8), int(y-8))
	}

	// Plato en mano (sprite)
	if mesero.TienePlato && r.assets.Plato != nil {
		op := &ebiten.DrawImageOptions{}
		scale := 2.0
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(x+20), float64(y-50)) // Arriba a la derecha
		screen.DrawImage(r.assets.Plato, op)
	}

	// Indicador de texto
	estado := "Libre"
	if mesero.TienePlato {
		estado = "Llevando plato"
	}
	ebitenutil.DebugPrintAt(screen, estado, int(x-30), int(y+55))
}

// DibujarCocinero dibuja al cocinero en la cocina
func (r *Renderer) DibujarCocinero(screen *ebiten.Image, x, y float32) {
	// Dibujar sprite del cocinero
	if r.assets.Cocinero != nil {
		op := &ebiten.DrawImageOptions{}
		scale := 3.5 // Más grande que el mesero
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(r.assets.Cocinero, op)
	} else {
		// Fallback: círculo naranja
		vector.DrawFilledCircle(screen, x+56, y+56, 20, color.RGBA{255, 140, 0, 255}, false)
		vector.StrokeCircle(screen, x+56, y+56, 20, 2, color.White, false)
		ebitenutil.DebugPrintAt(screen, "C", int(x+48), int(y+48))
	}

	// Etiqueta: CHEF es el productor en el patrón Productor-Consumidor
	ebitenutil.DebugPrintAt(screen, "CHEF", int(x-10), int(y+120))
}

// interpolarColor mezcla dos colores según un factor (0.0 a 1.0)
func interpolarColor(c1, c2 color.RGBA, factor float64) color.RGBA {
	f := float32(math.Max(0, math.Min(1, factor)))
	return color.RGBA{
		R: uint8(float32(c1.R)*(1-f) + float32(c2.R)*f),
		G: uint8(float32(c1.G)*(1-f) + float32(c2.G)*f),
		B: uint8(float32(c1.B)*(1-f) + float32(c2.B)*f),
		A: 255,
	}
}
