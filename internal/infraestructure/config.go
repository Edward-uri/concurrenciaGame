package infrastructure

import (
	"encoding/json"
	"os"
	"time"
)

// Config contiene la configuración de la aplicación
type Config struct {
	// Ventana
	Window WindowConfig `json:"window"`

	// Restaurante
	Restaurant RestaurantConfig `json:"restaurant"`

	// Rendimiento
	Performance PerformanceConfig `json:"performance"`

	// Logging
	Logging LoggingConfig `json:"logging"`
}

type WindowConfig struct {
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Title     string `json:"title"`
	Resizable bool   `json:"resizable"`
	VSync     bool   `json:"vsync"`
}

type RestaurantConfig struct {
	CapacidadBarra         int           `json:"capacidad_barra"`
	NumCocineros           int           `json:"num_cocineros"`
	NumMeseros             int           `json:"num_meseros"`
	ClientesInicial        int           `json:"clientes_inicial"`
	TiempoCoccion          time.Duration `json:"tiempo_coccion_ms"`
	TiempoEntrega          time.Duration `json:"tiempo_entrega_ms"`
	MaxClientesSpritesheet int           `json:"max_clientes_spritesheet"`
}

type PerformanceConfig struct {
	TargetFPS   int  `json:"target_fps"`
	EnableDebug bool `json:"enable_debug"`
}

type LoggingConfig struct {
	Level      string `json:"level"`  // debug, info, warn, error
	Output     string `json:"output"` // stdout, file
	FilePath   string `json:"file_path"`
	Structured bool   `json:"structured"` // JSON logs
}

// DefaultConfig retorna la configuración por defecto
func DefaultConfig() *Config {
	return &Config{
		Window: WindowConfig{
			Width:     800,
			Height:    600,
			Title:     "Restaurante Concurrente - Arquitectura Hexagonal",
			Resizable: true,
			VSync:     true,
		},
		Restaurant: RestaurantConfig{
			CapacidadBarra:         5,
			NumCocineros:           2,
			NumMeseros:             3,
			ClientesInicial:        3,
			TiempoCoccion:          800 * time.Millisecond,
			TiempoEntrega:          600 * time.Millisecond,
			MaxClientesSpritesheet: 8,
		},
		Performance: PerformanceConfig{
			TargetFPS:   60,
			EnableDebug: true,
		},
		Logging: LoggingConfig{
			Level:      "info",
			Output:     "stdout",
			FilePath:   "logs/restaurant.log",
			Structured: true,
		},
	}
}

// LoadConfig carga la configuración desde un archivo JSON
func LoadConfig(path string) (*Config, error) {
	config := DefaultConfig()

	file, err := os.Open(path)
	if err != nil {
		// Si no existe el archivo, usar configuración por defecto
		return config, nil
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

// SaveConfig guarda la configuración en un archivo JSON
func (c *Config) SaveConfig(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c)
}
