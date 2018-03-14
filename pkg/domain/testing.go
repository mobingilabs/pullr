package domain

type TestLogger struct{}

func (*TestLogger) Info(args ...interface{}) {
}

func (*TestLogger) Infof(format string, args ...interface{}) {
}

func (*TestLogger) Error(args ...interface{}) {
}

func (*TestLogger) Errorf(format string, args ...interface{}) {
}
