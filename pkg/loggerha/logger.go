package loggerha

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/labstack/echo/v4"
)

type ILogger interface {
	Infor(string)
	Error(string)
	Warn(string)
	Debug(string)
	SetWriter(io.Writer)
}
type Logger struct {
	writer io.Writer
	mu     sync.Mutex
}

func NewLogger() *Logger {
	return &Logger{
		writer: os.Stdout,
	}
}

func (l *Logger) Infor(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.writer, "INFO: %s\n", msg)
}

func (l *Logger) Error(msg string) {
	l.mu.Lock()
	defer l.mu.Lock()
	fmt.Fprintf(l.writer, "ERROR: %s\n", msg)
}

func (l *Logger) Debug(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.writer, "DEBUG: %s\n", msg)
}

func (l *Logger) SetWriter(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writer = w
}

func (l *Logger) Warn(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.writer, "WARN: %s\n", msg)
}

func LoggerMiddleware(logger ILogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			logger.Infor(fmt.Sprintf("Method: %s, URI: %s", req.Method, req.RequestURI))
			return next(c)
		}
	}
}
