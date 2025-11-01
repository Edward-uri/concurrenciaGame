package infrastructure

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger es un wrapper sobre zerolog para logging estructurado
type Logger struct {
	logger zerolog.Logger
}

// NewLogger crea un nuevo logger basado en la configuración
func NewLogger(config LoggingConfig) (*Logger, error) {
	var output io.Writer = os.Stdout

	// Configurar output
	if config.Output == "file" {
		// Crear directorio si no existe
		if err := os.MkdirAll("logs", 0755); err != nil {
			return nil, err
		}

		file, err := os.OpenFile(
			config.FilePath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644,
		)
		if err != nil {
			return nil, err
		}
		output = file
	}

	// Configurar formato
	if !config.Structured {
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: time.RFC3339,
		}
	}

	// Crear logger
	logger := zerolog.New(output).With().Timestamp().Logger()

	// Configurar nivel
	switch config.Level {
	case "debug":
		logger = logger.Level(zerolog.DebugLevel)
	case "info":
		logger = logger.Level(zerolog.InfoLevel)
	case "warn":
		logger = logger.Level(zerolog.WarnLevel)
	case "error":
		logger = logger.Level(zerolog.ErrorLevel)
	default:
		logger = logger.Level(zerolog.InfoLevel)
	}

	// Establecer como logger global
	log.Logger = logger

	return &Logger{logger: logger}, nil
}

// Debug registra un mensaje de debug
func (l *Logger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

// Debugf registra un mensaje de debug formateado
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

// Info registra un mensaje informativo
func (l *Logger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

// Infof registra un mensaje informativo formateado
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

// Warn registra una advertencia
func (l *Logger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

// Warnf registra una advertencia formateada
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.Warn().Msgf(format, args...)
}

// Error registra un error
func (l *Logger) Error(msg string, err error) {
	l.logger.Error().Err(err).Msg(msg)
}

// Errorf registra un error formateado
func (l *Logger) Errorf(format string, err error, args ...interface{}) {
	l.logger.Error().Err(err).Msgf(format, args...)
}

// WithFields crea un logger con campos adicionales
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	event := l.logger.With()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	return &Logger{logger: event.Logger()}
}

// Cocinero registra eventos del cocinero
func (l *Logger) Cocinero(id int, platoID int, accion string) {
	l.logger.Info().
		Str("tipo", "cocinero").
		Int("cocinero_id", id).
		Int("plato_id", platoID).
		Str("accion", accion).
		Msg("Evento de cocinero")
}

// Mesero registra eventos del mesero
func (l *Logger) Mesero(id int, platoID int, accion string) {
	l.logger.Info().
		Str("tipo", "mesero").
		Int("mesero_id", id).
		Int("plato_id", platoID).
		Str("accion", accion).
		Msg("Evento de mesero")
}

// EstadoRestaurant registra el estado del restaurante
func (l *Logger) EstadoRestaurant(clientes, enBarra, capacidad int, pausado bool) {
	l.logger.Info().
		Str("tipo", "estado").
		Int("clientes", clientes).
		Int("en_barra", enBarra).
		Int("capacidad", capacidad).
		Bool("pausado", pausado).
		Msg("Estado del restaurante")
}
