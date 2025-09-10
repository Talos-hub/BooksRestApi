package adapters

type Logger interface {
	Info(msg string, a ...any)  // Write info
	Error(msg string, a ...any) // Write error
	Warn(msg string, a ...any)  // write warning
}
