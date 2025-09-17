package abstraction

// Interface that provides Logger.
// It contains Info(), Error(), Warn()
// That's enave for loging a small aplication
type Logger interface {
	Info(msg string, a ...any)  // Write info
	Error(msg string, a ...any) // Write error
	Warn(msg string, a ...any)  // write warning
	Debug(msg string, a ...any)
}
