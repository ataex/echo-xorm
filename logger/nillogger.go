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

// Err do nothing, just match the interface
func (l *NilLogger) Err(values ...interface{}) {
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
