package logger

// NilLogger ...
type NilLogger struct {
}

// NewNilLogger is a constructor
func NewNilLogger() *NilLogger {
	return new(NilLogger)
}

// Info do nothing, just match the interface
func (l *NilLogger) Info(values ...interface{}) {
	return
}

// Error do nothing, just match the interface
func (l *NilLogger) Error(values ...interface{}) {
	return
}

// Warn do nothing, just match the interface
func (l *NilLogger) Warn(values ...interface{}) {
	return
}

// Close for NilLogger do nothing
func (l *NilLogger) Close() {
	return
}
