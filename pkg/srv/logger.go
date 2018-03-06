package srv

type Logger interface {
	Infof(string, ...interface{})
	Error(...interface{})
}
