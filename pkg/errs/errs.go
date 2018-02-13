package errs

// ErrLogger interface has error logging functions
type ErrLogger interface {
	Error(args ...interface{})
	Fatal(args ...interface{})
}

var logger ErrLogger

// SetLogger sets errs package's global logger
func SetLogger(errLogger ErrLogger) {
	logger = errLogger
}

// Log checks if the given err is null and if not it logs it
func Log(err error) {
	if err != nil {
		logger.Error(err)
	}
}

// Fatal checks if the given err is null and if not it logs and exits the program
func Fatal(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}
