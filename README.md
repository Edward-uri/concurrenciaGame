Desarrollo de una Aplicación Concurrente con
Interfaz Gráfica en Go
Descripción General
En esta actividad, los estudiantes deberán desarrollar una aplicación práctica en el lenguaje Go, aplicando
los principios de concurrencia y utilizando la librería EbitenEngine para construir una interfaz gráfica
interactiva.
El propósito es demostrar la comprensión de los mecanismos de sincronización, comunicación entre
goroutines y patrones de concurrencia a través de una simulación visual o juego interactivo.

Objetivos de Aprendizaje
Comprender y aplicar goroutines y canales para la ejecución concurrente.
Implementar mecanismos de sincronización como sync.Mutex , sync.WaitGroup o sync.Cond .
Utilizar patrones de concurrencia como Fan-out/Fan-in, Worker Pool, Productor-Consumidor, etc.
Integrar la lógica concurrente con una interfaz gráfica 2D utilizando EbitenEngine.
Desarrollar habilidades para diseñar, probar y documentar sistemas concurrentes en Go.

Instrucciones de la Actividad
1. Diseña una aplicación en Go con las siguientes características:
Utiliza goroutines para tareas concurrentes (ejemplo: animaciones, procesamiento, generación
de eventos).
Implementa canales para la comunicación entre las goroutines.
Emplea mecanismos de sincronización para evitar condiciones de carrera.
Usa al menos un patrón de concurrencia (por ejemplo: Productor-Consumidor, Worker Pool,
Fan-out/Fan-in, Pipeline).
Desarrolla una interfaz gráfica con EbitenEngine que visualice la ejecución concurrente.
2. Ejemplos de proyectos posibles:
Simulación de tráfico (autos moviéndose y controlados por semáforos concurrentes).
Juego de la serpiente con múltiples instancias controladas por goroutines.
Simulación de un sistema de producción y consumo (por ejemplo, una fábrica).
Procesamiento de datos o imágenes concurrente con visualización en pantalla.
3. Entrega los siguientes elementos:
Código fuente del proyecto (repositorio o archivo comprimido).
Documento técnico con la descripción de los patrones de concurrencia utilizados.
Capturas o video del funcionamiento de la aplicación.

Criterio Excelente (100%) Bueno (80%) Básico (60%)

Implementación
Concurrente

Se implementan múltiples
goroutines y canales de
forma eficiente y sin
bloqueos.

Se implementan
goroutines y canales,
con algunos errores
menores.

Implementación
limitada o con fallos
significativos.

Sincronización y
Seguridad

Usa correctamente
mecanismos como Mutex,
Cond o WaitGroup para
evitar condiciones de
carrera.

Usa parcialmente
mecanismos de
sincronización.

No se aplican
correctamente los
mecanismos de
sincronización.

Patrón de
Concurrencia

Se aplica al menos un
patrón concurrente

correctamente (Fan-
out/Fan-in, Worker Pool,

etc.).

Se intenta aplicar un
patrón, pero con
fallos en la lógica.

No se identifica el uso
de ningún patrón.

Interfaz Gráfica
(EbitenEngine)

Interfaz funcional, clara e
integrada con la lógica
concurrente. Buena
usabilidad y feedback
visual.

Interfaz funcional con
errores menores o
falta de integracion
clara.

Interfaz incompleta o
no funcional.

Documentación y
Presentación

Código documentado,
README.md completo y
clara explicación de los
patrones usados.

Documentación
parcial o poco clara.

Ausencia de
documentación o sin
explicación de los
patrones.

Creatividad y
Complejidad

Propone una aplicación
original con diseño bien
estructurado y aspectos
extra (métricas,
configuraciones).

Aplicación funcional
pero de complejidad
media.

Aplicación simple o
con poca originalidad.

Requisitos técnicos mínimos
Lenguaje: Go 1.20+ (o la versión estable más reciente que use la cátedra).
Librería gráfica: Ebiten (github.com/hajimehoshi/ebiten/v2).
El proyecto debe ejecutarse localmente (desktop).
Uso de go run -race para verificar ausencia de condiciones de carrera.

Rúbrica de Evaluación

Ponderación sugerida (opcional):
Implementación concurrente: 25%

Sincronización y seguridad: 20%
Patrones aplicados: 20%
Interfaz gráfica: 15%
Documentación: 10%
Creatividad: 10%

Criterios de Aprobación
Para aprobar la actividad, el proyecto deberá:
Ejecutarse correctamente sin condiciones de carrera detectadas con go run -race .
Mostrar uso real de concurrencia (no simulaciones secuenciales).
Cumplir con los elementos de diseño, funcionalidad y documentación requeridos.

Sugerencias de implementación (guía técnica)
Separación de responsabilidades: Mantén la lógica de concurrencia (workers, canales, coordinación)
separada de la capa de render (Ebiten). Evita bloquear el hilo de render en operaciones largas.
Comunicación segura con la UI: La UI (Ebiten's Update/Draw) debe interactuar con los datos
compartidos de forma segura — usa canales o copia inmutable del estado antes de dibujar. Nunca
realices operaciones bloqueantes directamente en Draw() o Update() ; usa buffers/colas y select con
default si es necesario.
Sincronización y cancelación: Usa context.Context para cancelar tareas cuando el usuario cierre la
aplicación o al cambiar de escena. Usa sync.WaitGroup para asegurar el cierre ordenado de
goroutines.
Evitar fugas de goroutines: Asegura que todas las goroutines puedan finalizar (propagación de
cancelación, cierre de canales o condiciones de salida).