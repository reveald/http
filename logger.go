package http

type Logger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}

type noopLogger struct{}

func (nl *noopLogger) Debugf(_ string, _ ...interface{}) {}
func (nl *noopLogger) Infof(_ string, _ ...interface{})  {}
func (nl *noopLogger) Errorf(_ string, _ ...interface{}) {}
